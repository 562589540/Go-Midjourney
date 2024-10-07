package handlers

import (
	"fmt"
	discord "github.com/bwmarrin/discordgo"
)

type Scene string

const (
	// BindMessageId /** 绑定消息Id */
	BindMessageId Scene = "BindMessageId"
	// WaitingToStart /** 等待开始 */
	WaitingToStart Scene = "WaitingToStart "
	// FirstTrigger /** 首次触发生成 */
	FirstTrigger Scene = "FirstTrigger"
	// FirstDescribe /** 反推首次触发生成 */
	FirstDescribe Scene = "FirstDescribe"
	// GenerateEnd /** 生成图片结束 */
	GenerateEnd Scene = "GenerateEnd"
	// GenerateEditError /** 发送的指令midjourney生成过程中发现错误 */
	GenerateEditError Scene = "GenerateEditError"
	// GenerateProgress /** 生成图片的进度 */
	GenerateProgress Scene = "GenerateProgress"
	// Describe 反推描述词获取成功
	Describe Scene = "Describe"
	// DescribeGet 请执行get消息获取描述词
	DescribeGet Scene = "DescribeGet"
	/**
	 * 富文本
	 */
	RichText Scene = "RichText"
	/**
	 * 发送的指令midjourney直接报错或排队阻塞不在该项目中处理 在业务服务中处理
	 * 例如：首次触发生成多少秒后没有回调业务服务判定会指令错误或者排队阻塞
	 */
)

type ReqMessage struct {
	Message  *discord.Message `json:"message,omitempty"`
	Content  string           `json:"content,omitempty"`
	Error    string           `json:"error,omitempty"`
	Progress int              `json:"progress,omitempty"`
	Type     Scene            `json:"type"`
}

func notice(m *discord.Message, progress int, scene Scene, err string) {
	request(ReqMessage{
		Message:  m,         //消息体
		Content:  m.Content, //消息文本
		Type:     scene,     //场景
		Progress: progress,  //进度
		Error:    err,       //进度
	})
}

// 通知客户
func request(params ReqMessage) {
	handler := getCallbackHandler()
	if err := handler.HandleCallback(params); err != nil {
		fmt.Printf("Error during callback: %v\n", err)
	}
}
