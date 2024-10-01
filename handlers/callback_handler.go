package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/562589540/Go-Midjourney/initialization"
	"net/http"
	"strings"
)

var callBackFun func(params ReqCb)

// CallbackHandler defines a contract for callback handling
type CallbackHandler interface {
	HandleCallback(params ReqCb) error
}

// HTTPCallbackHandler handles the callback via HTTP POST request
type HTTPCallbackHandler struct{}

// LocalCallbackHandler handles the callback via local method invocation
type LocalCallbackHandler struct{}

// HandleCallback for HTTP
func (h *HTTPCallbackHandler) HandleCallback(params ReqCb) error {
	data, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("json marshal error: %v", err)
	}

	req, err := http.NewRequest("POST", initialization.GetConfig().CB_URL, strings.NewReader(string(data)))
	if err != nil {
		return fmt.Errorf("http request creation error: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("http request error: %v", err)
	}
	defer resp.Body.Close()

	return nil
}

// HandleCallback for Local
func (l *LocalCallbackHandler) HandleCallback(params ReqCb) error {
	if callBackFun != nil {
		callBackFun(params)
	}
	return nil
}

// Select callback handler based on the callback type
func getCallbackHandler() CallbackHandler {
	switch initialization.GetConfig().CallBackType {
	case initialization.CallBackPost:
		return &HTTPCallbackHandler{}
	case initialization.CallBackDirect:
		return &LocalCallbackHandler{}
	default:
		return &LocalCallbackHandler{} // fallback to local callback
	}
}

// BindCallBackFun 绑定回掉方法 对外暴露方便调用本库使用者使用监控回调
func BindCallBackFun(cd func(params ReqCb)) {
	callBackFun = cd
}
