package services

import (
	"fmt"
	"github.com/562589540/Go-Midjourney/initialization"
	"github.com/bwmarrin/discordgo"
	"regexp"
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

	DebugDiscordMsg(message, "采集信息")

	//- Image #1 - Image #3 //符合这种类型的就是 放大指令的结果
	// Making variations for
	//- Variations 带有这种的是 v指令进行中 - Variations (Strong) //v指令结束
	//我们需要的是过滤 因为采集 我们只要4宫格

	//重绘开始 Making variations for image #1 with prompt [148390667393957889] 1 女孩, 日本动漫风格，HD，8k --niji 6 --ar 1:1 - @wong (Waiting to start)

	// 创建结果 map
	resultMap := make(map[string]GatherResult)

	for _, d := range message {
		// 从内容中提取 taskId
		taskId := ExtractNonceFromContent(d.Content)
		if taskId == "" {
			continue // taskId 为空，跳过
		}

		//过滤uv指令的数据
		if !isZoomInstructionResult(d.Content) {
			continue
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

// 判断字符串是否是符合放大指令的结果
func isZoomInstructionResult(s string) bool {
	// 匹配 "- Image #X" 的模式，X 为不定的数字
	zoomPattern := regexp.MustCompile(`- Image #\d+`)

	// 排除包含 "- Variations" 或 "- Variations (Strong)"
	variationPattern := regexp.MustCompile(`- Variations| - Variations \(Strong\)`)

	// 先检查是否匹配放大指令
	if zoomPattern.MatchString(s) {
		// 再检查是否包含不需要的 Variations 格式
		if !variationPattern.MatchString(s) {
			return true // 符合要求的字符串
		}
	}

	return false // 不符合要求
}
