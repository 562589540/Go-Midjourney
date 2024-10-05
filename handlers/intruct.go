package handlers

import (
	"github.com/562589540/Go-Midjourney/services"
)

func GenerateImage(prompt, nonce string) error {
	err := services.GenerateImage(prompt, nonce)
	return err
}

func ImageUpscale(index int64, discordMsgId string, msgHash string, messageFlags int64, nonce string) error {
	err := services.Upscale(index, discordMsgId, msgHash, messageFlags, nonce)
	return err
}

func ImageVariation(index int64, discordMsgId string, msgHash string, messageFlags int64, nonce string) error {
	err := services.Variate(index, discordMsgId, msgHash, messageFlags, nonce)
	return err
}

func ImageMaxUpscale(discordMsgId string, msgHash string) error {
	err := services.MaxUpscale(discordMsgId, msgHash)
	return err
}

func ImageReset(discordMsgId string, msgHash string) error {
	err := services.ReRoll(discordMsgId, msgHash)
	return err
}

func ImageDescribe(uploadName, nonce string) error {
	err := services.Describe(uploadName, nonce)
	return err
}

func UpscaleSubtle(discordMsgId string, msgHash string, messageFlags int64, nonce string) error {
	err := services.UpscaleSubtle(discordMsgId, msgHash, messageFlags, nonce)
	return err
}
