package handlers

import (
	"fmt"
	"github.com/562589540/Go-Midjourney/initialization"
	"github.com/562589540/Go-Midjourney/services"
	discord "github.com/bwmarrin/discordgo"
	"strings"
)

func DiscordMsgCreate(s *discord.Session, m *discord.MessageCreate) {
	// 过滤频道
	if m.ChannelID != initialization.GetConfig().DISCORD_CHANNEL_ID {
		return
	}
	//获取自己上传的图片 暂时弃用
	//if m.Author.ID == s.State.User.ID && m.Content == "" && m.Nonce != "" {
	//	if m.Attachments != nil && len(m.Attachments) > 0 {
	//		attachmentsMap[m.Nonce] = attachment{
	//			Url:      m.Attachments[0].URL,
	//			ProxyUrl: m.Attachments[0].ProxyURL,
	//		}
	//	}
	//}
	// 过滤掉自己发送的消息
	if m.Author.ID == s.State.User.ID {
		return
	}

	services.DebugDiscordMsg(m, "消息创建")

	if strings.Contains(m.Content, "(Waiting to start)") && !strings.Contains(m.Content, "Rerolling **") {
		//开始工作
		notice(m.Message, 0, WaitingToStart, "")
		return
	}

	if m.Nonce != "" {
		notice(m.Message, 0, BindMessageId, "")
		return
	}

	//创建了反推提示词
	if m.Interaction != nil && m.Interaction.Name == "describe" {
		notice(m.Message, 0, FirstDescribe, "")
	}

	//提示词获取开始
	//有绘画结果
	for _, attachment := range m.Attachments {
		if attachment.Width > 0 && attachment.Height > 0 {
			//绘画结束
			notice(m.Message, 1, GenerateEnd, "")
			return
		}
	}

	//一些错误处理
	if len(m.Embeds) > 0 {
		embeds := m.Embeds[0]
		switch embeds.Color {
		case 16711680:
			if embeds.Title == "Action needed to continue" {
				services.DebugDiscordMsg(embeds, "Action needed to continue")
				return
			} else if embeds.Title == "Pending mod message" {
				services.DebugDiscordMsg(embeds, "Pending mod message")
				return
			}
			notice(m.Message, 0, GenerateEditError, embeds.Description)
			break
		case 16776960:
			fmt.Println("警告", embeds.Description)
			break
		default:
			if strings.Contains(embeds.Title, "continue") && strings.Contains(embeds.Description, "verify you're human") {
				notice(m.Message, 0, GenerateEditError, "人机验证："+embeds.Description)
				return
			}

			if strings.Contains(embeds.Title, "Invalid") {
				notice(m.Message, 0, GenerateEditError, "无效的"+embeds.Description)
				return
			}
		}
		return
	}
}

func DiscordMsgUpdate(s *discord.Session, m *discord.MessageUpdate) {
	// 过滤频道
	if m.ChannelID != initialization.GetConfig().DISCORD_CHANNEL_ID {
		return
	}

	if m.Author == nil {
		return
	}

	// 过滤掉自己发送的消息
	if m.Author.ID == s.State.User.ID {
		return
	}

	services.DebugDiscordMsg(m, "消息更新")

	if strings.Contains(m.Content, "(Waiting to start)") && !strings.Contains(m.Content, "Rerolling **") {
		//开始工作
		notice(m.Message, 0, WaitingToStart, "")
		return
	}

	//提取到了进度
	if progress, err := services.ExtractProgress(m.Content); err == nil {
		notice(m.Message, progress, GenerateProgress, "")
		return
	}

	//有错误？？？？
	if strings.Contains(m.Content, "(Stopped)") {
		notice(m.Message, 0, GenerateEditError, "")
		return
	}

	//反推的
	if m.Interaction != nil && m.Interaction.Name == "describe" {
		if m.Embeds != nil && len(m.Embeds) > 0 && m.Embeds[0].Description != "" && m.Embeds[0].Image != nil {
			if m.Embeds[0].Image.Width > 0 && m.Embeds[0].Image.Height > 0 {
				notice(m.Message, 0, Describe, "")
				return
			}
		}
		notice(m.Message, 0, DescribeGet, "")
		return
	}
}
