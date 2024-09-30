package main

import (
	"fmt"
	"github.com/562589540/Go-Midjourney/handlers"
	"github.com/562589540/Go-Midjourney/initialization"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
)

func main() {
	os.Setenv("HTTP_PROXY", "http://127.0.0.1:7890")
	os.Setenv("HTTPS_PROXY", "http://127.0.0.1:7890")

	cfg := pflag.StringP("config", "c", "./config.yaml", "api server config file path.")

	pflag.Parse()

	initialization.LoadConfig(*cfg)

	if err := initialization.LoadDiscordClient(handlers.DiscordMsgCreate, handlers.DiscordMsgUpdate); err != nil {
		fmt.Println("Error starting Discord monitor:", err)
		return
	}
	r := gin.Default()

	r.POST("/v1/trigger/midjourney-bot", handlers.MidjourneyBot)
	r.POST("/v1/trigger/upload", handlers.UploadFile)

	r.Run(fmt.Sprintf(":%s", initialization.GetConfig().MJ_PORT))
}

func localWork() {
	os.Setenv("HTTP_PROXY", "http://127.0.0.1:7890")
	os.Setenv("HTTPS_PROXY", "http://127.0.0.1:7890")
	//任务队列，默认队列10，并发3。
}
