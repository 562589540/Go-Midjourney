package initialization

import (
	"fmt"

	"github.com/spf13/viper"
)

var (
	Debug   = false
	OnProxy = false
)

// 定义回调类型为字符串类型
type CallBackType string

// 定义允许的回调类型常量
const (
	CallBackPost   CallBackType = "post"
	CallBackDirect CallBackType = "direct"
)

// ListenType 监听类型
type ListenType string

// 定义允许的回调类型常量
const (
	ListenBot  ListenType = "bot"
	ListenUser ListenType = "user"
)

type Config struct {
	DISCORD_USER_TOKEN string
	DISCORD_BOT_TOKEN  string
	DISCORD_SERVER_ID  string
	DISCORD_CHANNEL_ID string
	CB_URL             string
	MJ_PORT            string
	ListenType         ListenType
	CallBackType       CallBackType // 使用自定义类型
}

var config *Config

func LoadConfig(cfg string) *Config {
	viper.SetConfigFile(cfg)
	viper.ReadInConfig()
	viper.AutomaticEnv()
	config = &Config{
		DISCORD_USER_TOKEN: getViperStringValue("DISCORD_USER_TOKEN"),
		DISCORD_BOT_TOKEN:  getViperStringValue("DISCORD_BOT_TOKEN"),
		DISCORD_SERVER_ID:  getViperStringValue("DISCORD_SERVER_ID"),
		DISCORD_CHANNEL_ID: getViperStringValue("DISCORD_CHANNEL_ID"),
		CB_URL:             getViperStringValue("CB_URL"),
		MJ_PORT:            getDefaultValue("MJ_PORT", "16007"),
		CallBackType:       CallBackType(getDefaultValue("CallBackType", "post")),
		ListenType:         ListenType(getDefaultValue("ListenType", "bot")),
	}
	return config
}

func GetConfig() *Config {
	return config
}

func getViperStringValue(key string) string {
	value := viper.GetString(key)
	if value == "" {
		panic(fmt.Errorf("%s MUST be provided in environment or config.yaml file", key))
	}
	return value
}

func getDefaultValue(key string, defaultValue string) string {
	value := viper.GetString(key)
	if value == "" {
		return defaultValue
	} else {
		return value
	}
}

func NewConfig() {
	config = &Config{}
}

// SetDiscordUserToken 手动设置用户token
func SetDiscordUserToken(discordUserToken string) {
	config.DISCORD_USER_TOKEN = discordUserToken
}

// SetDiscordBotToken 手动设置机器人token
func SetDiscordBotToken(discordBotToken string) {
	config.DISCORD_BOT_TOKEN = discordBotToken
}

// SetDiscordServerId 手动设置服务器id
func SetDiscordServerId(discordServerId string) {
	config.DISCORD_SERVER_ID = discordServerId
}

// SetDiscordChannelId 手动设置频道id
func SetDiscordChannelId(discordChannelId string) {
	config.DISCORD_CHANNEL_ID = discordChannelId
}

// SetCallBackType 手动设置回掉设置
func SetCallBackType(callBackType CallBackType) {
	config.CallBackType = callBackType
}

// SetListenType 手动设置监控类型
func SetListenType(ListenType ListenType) {
	config.ListenType = ListenType
}

func SetDebug(debug bool) {
	Debug = debug
}
