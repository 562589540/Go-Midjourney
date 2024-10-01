package services

import (
	"fmt"
	discord "github.com/bwmarrin/discordgo"
	"regexp"
	"strconv"
	"strings"
)

// ExtractMessageDetails 提取 messageId 和 messageHash
func ExtractMessageDetails(m *discord.MessageCreate) (string, string) {
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

// ExtractImageDetails 提取绘画后的图片
func ExtractImageDetails(m *discord.MessageCreate) *discord.MessageAttachment {
	if len(m.Attachments) > 0 {
		return m.Attachments[0]
	}
	return nil
}

// ExtractHashFromFilename 从文件名中提取 messageHash（假设 messageHash 是 UUID 样式的字符串）
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

// ExtractProgress 提取字符串中的进度百分比数字
func ExtractProgress(content string) (int, error) {
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

// IsFirstGeneration 判断是否是首次生成的四宫格图片
func IsFirstGeneration(content string) bool {
	// 如果 content 中不包含 "Image #"，则为首次生成的四宫格图片
	return !strings.Contains(content, "Image #")
}

// IsUpscaledImage 判断是否是放大的单张图片
func IsUpscaledImage(content string) bool {
	// 如果 content 中包含 "Image #"，则为放大的图片
	return strings.Contains(content, "Image #")
}
