package handlers

import (
	"encoding/json"
	"fmt"
	"midjourney/initialization"
	"midjourney/services"
	"regexp"
	"strconv"
	"strings"
	"time"

	discord "github.com/bwmarrin/discordgo"
	"github.com/k0kubun/pp/v3"
)

type Scene string

const (
	// FirstTrigger /** 首次触发生成 */
	FirstTrigger Scene = "FirstTrigger"
	// GenerateEnd /** 生成图片结束 */
	GenerateEnd Scene = "GenerateEnd"
	// GenerateEditError /** 发送的指令midjourney生成过程中发现错误 */
	GenerateEditError Scene = "GenerateEditError"
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

	debugDiscordMsg(m, "消息创建")

	/******** *********/
	pp.Println(m.Content)
	pp.Println(m.Attachments)
	/******** *********/

	nonce := services.ExtractNonceFromContent(m.Content)
	prompt := services.FindNonce(nonce)
	pp.Println("获取到的数据", prompt, nonce)

	if strings.Contains(m.Content, "(Waiting to start)") && !strings.Contains(m.Content, "Rerolling **") {
		//开始工作
		trigger(m.Content, FirstTrigger)
		return
	}
	for _, attachment := range m.Attachments {
		if attachment.Width > 0 && attachment.Height > 0 {
			//绘画结束
			replay(m)
			return
		}
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

	debugDiscordMsg(m, "消息更新")

	if strings.Contains(m.Content, "(Stopped)") {
		trigger(m.Content, GenerateEditError)
		return
	}
	if len(m.Embeds) > 0 {
		send(m.Embeds)
		return
	}
}

type ReqCb struct {
	Embeds  []*discord.MessageEmbed `json:"embeds,omitempty"`
	Discord *discord.MessageCreate  `json:"discord,omitempty"`
	Content string                  `json:"content,omitempty"`
	Type    Scene                   `json:"type"`
}

func replay(m *discord.MessageCreate) {
	body := ReqCb{
		Discord: m,
		Type:    GenerateEnd,
	}
	// 提取 messageId 和 messageHash
	messageId, messageHash := extractMessageDetails(m)
	fmt.Printf("Extracted messageId: %s, messageHash: %s\n", messageId, messageHash)

	// 等待1秒
	time.Sleep(2 * time.Second)

	err := ImageUpscale(2, messageId, messageHash)
	if err != nil {
		return
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

func trigger(content string, t Scene) {
	body := ReqCb{
		Content: content,
		Type:    t,
	}
	request(body)
}

// 通知客户服务器
func request(params interface{}) {
	//data, err := json.Marshal(params)
	//
	//// fmt.Println("请求回调接口", string(data))
	//
	//if err != nil {
	//	fmt.Println("json marshal error: ", err)
	//	return
	//}
	//req, err := http.NewRequest("POST", initialization.GetConfig().CB_URL, strings.NewReader(string(data)))
	//if err != nil {
	//	fmt.Println("http request error: ", err)
	//	return
	//}
	//req.Header.Set("Content-Type", "application/json")
	//client := &http.Client{}
	//resp, err := client.Do(req)
	//if err != nil {
	//	fmt.Println("http request error: ", err)
	//	return
	//}
	//defer resp.Body.Close()
}

func debugDiscordMsg(m any, t string) {
	// 序列化 m 结构体为 JSON 格式
	jsonData, err := json.MarshalIndent(m, "", "  ") // 将结构体格式化为漂亮的 JSON
	if err != nil {
		fmt.Println("Error serializing message:", err)
		return
	}
	// 打印序列化后的 JSON 字符串
	fmt.Println(t)
	fmt.Println(string(jsonData))
}

// 提取 messageId 和 messageHash
func extractMessageDetails(m *discord.MessageCreate) (string, string) {
	// 提取 messageId
	messageId := m.ID

	// 提取 messageHash，从第一个附件的文件名中获取
	var messageHash string
	if len(m.Attachments) > 0 {
		attachment := m.Attachments[0]
		// 假设 messageHash 是文件名中 UUID 样式的一部分，如 "e474232b-c979-420b-9d59-d7b207ebc33f"
		messageHash = extractHashFromFilename(attachment.Filename)
	}

	return messageId, messageHash
}

// 从文件名中提取 messageHash（假设 messageHash 是 UUID 样式的字符串）
func extractHashFromFilename(filename string) string {
	// 定义一个正则表达式来匹配 UUID（忽略文件扩展名）
	// 正则匹配形如: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx 的 UUID
	re := regexp.MustCompile(`[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}`)

	// 去掉文件扩展名
	filenameWithoutExt := strings.TrimSuffix(filename, ".png")

	// 使用正则表达式查找文件名中的 UUID
	match := re.FindString(filenameWithoutExt)
	if match != "" {
		return match // 如果找到 UUID，返回它
	}
	return ""
}

// 提取字符串中的进度百分比数字
func extractProgress(content string) (int, error) {
	// 定义正则表达式匹配百分比数值，例如 "78%"
	re := regexp.MustCompile(`\((\d+)%\)`)

	// 查找匹配的字符串
	match := re.FindStringSubmatch(content)

	// 检查是否找到匹配项
	if len(match) > 1 {
		// 将匹配的百分比部分转为整数
		progress, err := strconv.Atoi(match[1])
		if err != nil {
			return 0, err
		}
		return progress, nil
	}

	// 如果未找到匹配的进度，则返回错误
	return 0, fmt.Errorf("no progress found in content")
}
