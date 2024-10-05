package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/562589540/Go-Midjourney/gclient"
	"github.com/562589540/Go-Midjourney/initialization"
	"github.com/562589540/Go-Midjourney/services"
	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

//type attachment struct {
//	Url      string `json:"url"`
//	ProxyUrl string `json:"proxy_url"`
//}

//var (
//	attachmentsMap = make(map[string]attachment)
//)

type ReqUploadFile struct {
	ImgData []byte `json:"imgData"`
	Name    string `json:"name"`
	Size    int64  `json:"size"`
}

func UploadFile(c *gin.Context) {
	var body ReqUploadFile
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	data, err := services.Attachments(body.Name, body.Size)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error saving image: %s", err)
		return
	}
	if len(data.Attachments) == 0 {
		c.String(http.StatusInternalServerError, "上传图片失败: %s", err)
		return
	}
	payload := bytes.NewReader(body.ImgData)
	client := &http.Client{}
	req, err := http.NewRequest("PUT", data.Attachments[0].UploadUrl, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "image/png")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	c.JSON(200, gin.H{"name": data.Attachments[0].UploadFilename})
}

// UploadFileWithLocalGetAttachment 上传图片发送到频道并获取上传后的信息
func UploadFileWithLocalGetAttachment(filePath string) (*discordgo.MessageAttachment, error) {
	attachment, err := UploadFileWithLocal(filePath)
	if err != nil {
		return nil, err
	}
	nonce, err := services.GetNextIDGenerator().Generate()
	if err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %v", err)
	}

	ctime := time.Now()
	//发送消息
	err = services.PutUploadMessages(fmt.Sprintf("%d", nonce), attachment.UploadFilename)
	if err != nil {
		return nil, err
	}

	time.Sleep(1 * time.Second)

	//在读取消息拿到url
	message, err := initialization.GetMessage(5)
	if err != nil {
		return nil, err
	}
	UploadName := filepath.Base(attachment.UploadFilename)
	for _, m := range message {
		if m.Timestamp.After(ctime) {
			if m.Attachments != nil && len(m.Attachments) > 0 {
				if m.Attachments[0].Filename == UploadName {
					return m.Attachments[0], nil
				}
			}
		}
	}
	return nil, errors.New("未获取到上传的图片")
}

// UploadFileWithLocal 仅上传图片
func UploadFileWithLocal(filePath string) (*services.ResFile, error) {
	// 打开文件获取信息
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// 获取文件状态信息
	info, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %v", err)
	}
	fileName := filepath.Base(filePath)
	fileSize := info.Size()

	data, err := services.Attachments(fileName, fileSize)
	if err != nil {
		return nil, fmt.Errorf("error saving image: %v", err)
	}
	if len(data.Attachments) == 0 {
		return nil, fmt.Errorf("上传图片失败: %v", err)
	}

	// 读取文件内容到缓冲区
	buffer := bytes.NewBuffer(nil)
	if _, err := io.Copy(buffer, file); err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	client := gclient.GetGclient().GetHTTPClient()

	// 创建 PUT 请求
	req, err := http.NewRequest(http.MethodPut, data.Attachments[0].UploadUrl, buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to create PUT request: %v", err)
	}
	req.Header.Set("Content-Type", "image/png") // 设置内容类型，根据实际情况修改

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send PUT request: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to upload file: received status code %d", resp.StatusCode)
	}
	return &data.Attachments[0], nil
}
