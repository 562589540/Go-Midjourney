package services

import (
	"fmt"
	"github.com/562589540/Go-Midjourney/initialization"
	"github.com/bwmarrin/discordgo"
)

func GatherDoneImage(limit int) (map[string]*discordgo.MessageAttachment, error) {
	message, err := initialization.GetMessage(limit)
	if err != nil {
		return nil, err
	}

	// 检查 message 是否为 nil
	if message == nil {
		return nil, fmt.Errorf("no messages found")
	}

	resultMap := make(map[string]*discordgo.MessageAttachment)
	for _, d := range message {
		taskId := ExtractNonceFromContent(d.Content)
		if taskId == "" {
			continue
		}
		for _, attachment := range d.Attachments {
			if attachment.Width > 0 && attachment.Height > 0 {
				resultMap[taskId] = attachment
				break
			}
		}
	}
	return resultMap, nil
}
