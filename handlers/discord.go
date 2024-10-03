package handlers

import (
	"fmt"
	"github.com/562589540/Go-Midjourney/initialization"
	"github.com/562589540/Go-Midjourney/services"
	discord "github.com/bwmarrin/discordgo"
	"github.com/k0kubun/pp/v3"
	"strings"
)

type Scene string

const (
	// FirstTrigger /** 首次触发生成 */
	FirstTrigger Scene = "FirstTrigger"
	// GenerateEnd /** 生成图片结束 */
	GenerateEnd Scene = "GenerateEnd"
	// GenerateEditError /** 发送的指令midjourney生成过程中发现错误 */
	GenerateEditError Scene = "GenerateEditError"
	// GenerateProgress /** 生成图片的进度 */
	GenerateProgress Scene = "GenerateProgress"
	/**
	 * 富文本
	 */
	RichText Scene = "RichText"
	/**
	 * 发送的指令midjourney直接报错或排队阻塞不在该项目中处理 在业务服务中处理
	 * 例如：首次触发生成多少秒后没有回调业务服务判定会指令错误或者排队阻塞
	 */
)

func DiscordMsgCreate(s *discord.Session, m *discord.MessageCreate) {
	// 过滤频道
	if m.ChannelID != initialization.GetConfig().DISCORD_CHANNEL_ID {
		return
	}

	// 过滤掉自己发送的消息
	if m.Author.ID == s.State.User.ID {
		return
	}

	services.DebugDiscordMsg(m, "消息创建")

	/******** *********/
	pp.Println(m.Content)
	pp.Println(m.Attachments)
	/******** *********/

	//nonce := services.ExtractNonceFromContent(m.Content)
	//prompt := services.FindNonce(nonce)
	//pp.Println("获取到的数据", prompt, nonce)

	if strings.Contains(m.Content, "(Waiting to start)") && !strings.Contains(m.Content, "Rerolling **") {
		//开始工作
		triggerCreate(m.Content, m, FirstTrigger)
		return
	}

	//快速出图用尽

	for _, attachment := range m.Attachments {
		if attachment.Width > 0 && attachment.Height > 0 {
			//绘画结束
			replay(m)
			return
		}
	}

	if len(m.Embeds) > 0 {
		send(m.Embeds)
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

	//进度更新
	progress, err := services.ExtractProgress(m.Content)
	if err == nil {
		progressCd(m, progress)
	}

	if strings.Contains(m.Content, "(Stopped)") {
		triggerUpdate(m.Content, m, GenerateEditError)
		return
	}
	if len(m.Embeds) > 0 {
		send(m.Embeds)
		return
	}
}

type ReqCb struct {
	Embeds        []*discord.MessageEmbed `json:"embeds,omitempty"`
	Discord       *discord.MessageCreate  `json:"discord,omitempty"`
	DiscordUpdate *discord.MessageUpdate  `json:"discordUpdate,omitempty"`
	Content       string                  `json:"content,omitempty"`
	Progress      int                     `json:"progress,omitempty"`
	Type          Scene                   `json:"type"`
}

func replay(m *discord.MessageCreate) {
	body := ReqCb{
		Discord: m,
		Type:    GenerateEnd,
	}
	request(body)
}

func send(embeds []*discord.MessageEmbed) {
	body := ReqCb{
		Embeds: embeds,
		Type:   RichText,
	}
	request(body)
}

func triggerCreate(content string, m *discord.MessageCreate, t Scene) {
	body := ReqCb{
		Discord: m,
		Content: content,
		Type:    t,
	}
	request(body)
}

func triggerUpdate(content string, m *discord.MessageUpdate, t Scene) {
	body := ReqCb{
		DiscordUpdate: m,
		Content:       content,
		Type:          t,
	}
	request(body)
}

func progressCd(message *discord.MessageUpdate, progress int) {
	body := ReqCb{
		Progress:      progress,
		Content:       message.Content,
		DiscordUpdate: message,
		Type:          GenerateProgress,
	}
	request(body)
}

// 通知客户
func request(params ReqCb) {
	handler := getCallbackHandler()
	if err := handler.HandleCallback(params); err != nil {
		fmt.Printf("Error during callback: %v\n", err)
	}
}
