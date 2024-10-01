package services

import (
	"fmt"
	"regexp"
)

// ExtractNonceFromContent 从消息内容中提取 `nonce`（假设 `nonce` 被包裹在 [] 中）
func ExtractNonceFromContent(content string) string {
	// 正则表达式提取 [nonce] 格式的内容
	//var nonceRegex = regexp.MustCompile(`\[(.*?)]`)

	//必然是数字
	var nonceRegex = regexp.MustCompile(`\[(\d+)\]`)

	// 使用正则表达式提取 [nonce] 格式的部分
	match := nonceRegex.FindStringSubmatch(content)
	if len(match) > 1 {
		return match[1] // 返回 [nonce] 中的内容
	}
	return ""
}

func AddPromptNonce(prompt string) (uint64, string, error) {
	// 生成唯一的 nonce 这个应该本地业务时候就生成在外部传入
	nonce, err := GetNextIDGenerator().Generate()
	if err != nil {
		return 0, "", err
	}
	// 在 prompt 中附加 nonce，确保它与结果消息相关联
	return nonce, fmt.Sprintf("[%d] %s", nonce, prompt), nil
}
