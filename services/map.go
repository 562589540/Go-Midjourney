package services

import (
	"regexp"
)

var requestMap = make(map[string]string)

func saveRequest(nonce, prompt string) {
	requestMap[nonce] = prompt
}

func FindNonce(nonce string) string {
	if prompt, ok := requestMap[nonce]; ok {
		return prompt
	}
	return ""
}

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
