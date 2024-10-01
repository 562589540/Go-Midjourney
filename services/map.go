package services

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
