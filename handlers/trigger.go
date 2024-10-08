package handlers

import (
	"errors"

	"github.com/gin-gonic/gin"
)

type RequestTrigger struct {
	Type         string `json:"type"`
	DiscordMsgId string `json:"discordMsgId,omitempty"`
	MsgHash      string `json:"msgHash,omitempty"`
	Prompt       string `json:"prompt,omitempty"`
	Index        int64  `json:"index,omitempty"`
	Flags        int64  `json:"flags,omitempty"`
	Nonce        string `json:"nonce,omitempty"`
}

func MidjourneyBot(c *gin.Context) {
	var body RequestTrigger
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var err error
	switch body.Type {
	case "generate":
		err = GenerateImage(body.Prompt, body.Nonce)
	case "upscale":
		err = ImageUpscale(body.Index, body.DiscordMsgId, body.MsgHash, body.Flags, body.Nonce)
	case "variation":
		err = ImageVariation(body.Index, body.DiscordMsgId, body.MsgHash, body.Flags, body.Nonce)
	case "maxUpscale":
		err = ImageMaxUpscale(body.DiscordMsgId, body.MsgHash)
	case "reset":
		err = ImageReset(body.DiscordMsgId, body.MsgHash)
	case "describe":
		err = ImageDescribe(body.Prompt, body.Nonce)
	default:
		err = errors.New("invalid type")
	}

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "success"})
}
