package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/562589540/Go-Midjourney/gclient"
	config "github.com/562589540/Go-Midjourney/initialization"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

var proxyURL = ""

const (
	url             string = "https://discord.com/api/v9/interactions"
	uploadUrlFormat string = "https://discord.com/api/v9/channels/%s/attachments"
	appId           string = "936929561302675456"
	SessionID       string = "b1bf42be072a0c8c706e153bf37585b6"
	Version         string = "1237876415471554623"
	Id              string = "938956540159881230"
)

func GenerateImage(prompt, nonce string) error {
	requestBody := ReqTriggerDiscord{
		Type:          2,
		GuildID:       config.GetConfig().DISCORD_SERVER_ID,
		ChannelID:     config.GetConfig().DISCORD_CHANNEL_ID,
		ApplicationId: appId,
		SessionId:     SessionID,
		Data: DSCommand{
			Version: Version,
			Id:      Id,
			Name:    "imagine",
			Type:    1,
			Options: []DSOption{{Type: 3, Name: "prompt", Value: prompt}},
			ApplicationCommand: DSApplicationCommand{
				Id:                       Id,
				ApplicationId:            appId,
				Version:                  Version,
				DefaultPermission:        true,
				DefaultMemberPermissions: nil,
				Type:                     1,
				Nsfw:                     false,
				Name:                     "imagine",
				Description:              "Lucky you!",
				DmPermission:             true,
				Options:                  []DSCommandOption{{Type: 3, Name: "prompt", Description: "The prompt to imagine", Required: true}},
			},
			Attachments: []ReqCommandAttachments{},
		},
		Nonce: nonce,
	}
	_, err := request(requestBody, url)
	return err
}

func Upscale(index int64, messageId string, messageHash string, messageFlags int64, nonce string) error {
	requestBody := ReqUpscaleDiscord{
		Type:          3,
		GuildId:       config.GetConfig().DISCORD_SERVER_ID,
		ChannelId:     config.GetConfig().DISCORD_CHANNEL_ID,
		MessageFlags:  messageFlags,
		MessageId:     messageId,
		ApplicationId: appId,
		SessionId:     SessionID,
		Data: UpscaleData{
			ComponentType: 2,
			CustomId:      fmt.Sprintf("MJ::JOB::upsample::%d::%s", index, messageHash),
		},
		Nonce: nonce,
	}
	_, err := request(requestBody, url)
	return err
}

// UpscaleSubtle 轻微放大：对图像进行轻微的放大处理，保持图像的细节相对不变，只是提升分辨率。
func UpscaleSubtle(messageId string, messageHash string, messageFlags int64, nonce string) error {
	requestBody := ReqUpscaleDiscord{
		Type:          3,
		GuildId:       config.GetConfig().DISCORD_SERVER_ID,
		ChannelId:     config.GetConfig().DISCORD_CHANNEL_ID,
		MessageFlags:  messageFlags,
		MessageId:     messageId,
		ApplicationId: appId,
		SessionId:     SessionID,
		Data: UpscaleData{
			ComponentType: 2,
			CustomId:      fmt.Sprintf("MJ::JOB::upsample_v6r1_2x_subtle::1::%s::SOLO", messageHash),
		},
		Nonce: nonce,
	}
	_, err := request(requestBody, url)
	return err
}

// MaxUpscale 未知
func MaxUpscale(messageId string, messageHash string) error {
	requestBody := ReqUpscaleDiscord{
		Type:          3,
		GuildId:       config.GetConfig().DISCORD_SERVER_ID,
		ChannelId:     config.GetConfig().DISCORD_CHANNEL_ID,
		MessageFlags:  0,
		MessageId:     messageId,
		ApplicationId: appId,
		SessionId:     SessionID,
		Data: UpscaleData{
			ComponentType: 2,
			CustomId:      fmt.Sprintf("MJ::JOB::variation::1::%s::SOLO", messageHash),
		},
	}
	_, err := request(requestBody, url)
	return err
}

// Variate 变体
func Variate(index int64, messageId string, messageHash string, messageFlags int64, nonce string) error {
	requestBody := ReqVariationDiscord{
		Type:          3,
		GuildId:       config.GetConfig().DISCORD_SERVER_ID,
		ChannelId:     config.GetConfig().DISCORD_CHANNEL_ID,
		MessageFlags:  messageFlags,
		MessageId:     messageId,
		ApplicationId: appId,
		SessionId:     SessionID,
		Data: UpscaleData{
			ComponentType: 2,
			CustomId:      fmt.Sprintf("MJ::JOB::variation::%d::%s", index, messageHash),
		},
		Nonce: nonce,
	}
	_, err := request(requestBody, url)
	return err
}

// VariatePrompt 带提示词的放大 暂时不需要
func VariatePrompt(index int64, messageId string, messageHash string, prompt string) error {
	requestBody := ReqVariationVariatePromptDiscord{
		Type:          5,
		ApplicationId: appId,
		ChannelId:     config.GetConfig().DISCORD_CHANNEL_ID,
		GuildId:       config.GetConfig().DISCORD_SERVER_ID,
		Data: VariatePromptData{
			Id:       messageId,
			CustomId: fmt.Sprintf("MJ::RemixModal::%s::%d::1", messageHash, index),
			Components: []Component{
				{
					Type: 1,
					Components: []Component{
						{
							Type:     4,
							CustomId: "MJ::RemixModal::new_prompt",
							Value:    prompt,
						},
					},
				},
			},
		},
		SessionId: SessionID,
	}

	// 将结构体转换为 JSON 字符串并格式化
	jsonData, err := json.MarshalIndent(requestBody, "", "  ")
	if err != nil {
		return err
	}

	// 打印 JSON 格式的 requestBody
	fmt.Println(string(jsonData))

	_, err = request(requestBody, url)
	return err
}

// ReRoll 重绘
func ReRoll(messageId string, messageHash string) error {
	requestBody := ReqResetDiscord{
		Type:          3,
		GuildId:       config.GetConfig().DISCORD_SERVER_ID,
		ChannelId:     config.GetConfig().DISCORD_CHANNEL_ID,
		MessageFlags:  0,
		MessageId:     messageId,
		ApplicationId: appId,
		SessionId:     SessionID,
		Data: UpscaleData{
			ComponentType: 2,
			CustomId:      fmt.Sprintf("MJ::JOB::reroll::0::%s::SOLO", messageHash),
		},
	}
	_, err := request(requestBody, url)
	return err
}

func Describe(uploadName string) error {
	requestBody := ReqTriggerDiscord{
		Type:          2,
		GuildID:       config.GetConfig().DISCORD_SERVER_ID,
		ChannelID:     config.GetConfig().DISCORD_CHANNEL_ID,
		ApplicationId: appId,
		SessionId:     SessionID,
		Data: DSCommand{
			Version: Version,
			Id:      Id,
			Name:    "describe",
			Type:    1,
			Options: []DSOption{{Type: 11, Name: "image", Value: 0}},
			ApplicationCommand: DSApplicationCommand{
				Id:                       Id,
				ApplicationId:            appId,
				Version:                  Version,
				DefaultPermission:        true,
				DefaultMemberPermissions: nil,
				Type:                     1,
				Nsfw:                     false,
				Name:                     "describe",
				Description:              "Writes a prompt based on your image.",
				DmPermission:             true,
				Options:                  []DSCommandOption{{Type: 11, Name: "image", Description: "The image to describe", Required: true}},
			},
			Attachments: []ReqCommandAttachments{{
				Id:             "0",
				Filename:       filepath.Base(uploadName),
				UploadFilename: uploadName,
			}},
		},
	}
	_, err := request(requestBody, url)
	return err
}

func Attachments(name string, size int64) (ResAttachments, error) {
	requestBody := ReqAttachments{
		Files: []ReqFile{{
			Filename: name,
			FileSize: size,
			Id:       "1",
		}},
	}
	uploadUrl := fmt.Sprintf(uploadUrlFormat, config.GetConfig().DISCORD_CHANNEL_ID)
	body, err := request(requestBody, uploadUrl)
	var data ResAttachments

	if err = json.Unmarshal(body, &data); err != nil {
		return ResAttachments{}, err
	}

	return data, err
}

// 生成U指令结构体
func genUpscaleDiscordReq(messageId, nonce string, data UpscaleData) ReqUpscaleDiscord {
	return ReqUpscaleDiscord{
		Type:          3,
		GuildId:       config.GetConfig().DISCORD_SERVER_ID,
		ChannelId:     config.GetConfig().DISCORD_CHANNEL_ID,
		MessageFlags:  0,
		MessageId:     messageId,
		ApplicationId: appId,
		SessionId:     SessionID,
		Data:          data,
		Nonce:         nonce,
	}
}

func request(params interface{}, url string) ([]byte, error) {
	requestData, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", config.GetConfig().DISCORD_USER_TOKEN)

	client := gclient.GetGclient().GetHTTPClient()

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	bod, respErr := ioutil.ReadAll(response.Body)
	fmt.Println("response:", string(bod), respErr, response.Status, url)
	return bod, respErr
}
