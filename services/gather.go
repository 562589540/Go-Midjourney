package services

import (
	"fmt"
	"github.com/562589540/Go-Midjourney/initialization"
	"github.com/bwmarrin/discordgo"
)

// GatherResult 用于打包消息结果
type GatherResult struct {
	*discordgo.MessageAttachment        // 包含 Discord 附件
	MsgHash                      string `json:"msgHash"`
	MessageFlags                 int    `json:"messageFlags"`
	MessageId                    string `json:"messageId"`
}

// GatherDoneImage 收集已完成的图片
func GatherDoneImage(limit int) (map[string]GatherResult, error) {
	// 获取消息
	message, err := initialization.GetMessage(limit)
	if err != nil {
		return nil, err
	}

	// 检查 message 是否为 nil
	if message == nil {
		return nil, fmt.Errorf("no messages found")
	}

	//DebugDiscordMsg(message, "采集信息")

	// 创建结果 map
	resultMap := make(map[string]GatherResult)

	for _, d := range message {
		// 从内容中提取 taskId
		taskId := ExtractNonceFromContent(d.Content)
		if taskId == "" {
			continue // taskId 为空，跳过
		}

		// 遍历消息中的附件
		for _, attachment := range d.Attachments {
			// 只处理有宽度和高度的图片附件
			if attachment != nil && attachment.Width > 0 && attachment.Height > 0 {
				// 提取消息ID和哈希
				messageId, messageHash := ExtractMessageDetails(d)

				// 将信息打包到 GatherResult 结构体
				result := GatherResult{
					MessageAttachment: attachment,   // 包含图片附件
					MsgHash:           messageHash,  // 消息哈希
					MessageFlags:      int(d.Flags), // 消息标志
					MessageId:         messageId,    // 消息ID
				}

				// 存储到结果 map 中
				resultMap[taskId] = result
				break // 找到合适的附件后，跳出内部循环
			}
		}
	}
	return resultMap, nil
}
