package initialization

import (
	"fmt"

	discord "github.com/bwmarrin/discordgo"
)

var discordClient *discord.Session

func LoadDiscordClient(create func(s *discord.Session, m *discord.MessageCreate), update func(s *discord.Session, m *discord.MessageUpdate)) error {
	// 如果客户端已经启动，避免重复启动
	if discordClient != nil {
		fmt.Println("Discord client is already running")
		return nil
	}

	var err error

	if config.ListenType == ListenBot {
		//使用机器人监听
		discordClient, err = discord.New("Bot " + config.DISCORD_BOT_TOKEN)
	} else {
		//使用用户监听 封禁风险较高
		discordClient, err = discord.New(config.DISCORD_USER_TOKEN)
	}

	if err != nil {
		return fmt.Errorf("error creating Discord session: %v", err)
	}

	err = discordClient.Open()
	if err != nil {
		return fmt.Errorf("error opening connection: %v", err)
	}

	discordClient.AddHandler(create)
	discordClient.AddHandler(update)
	return nil
}

// StopDiscordMonitor 停止 Discord 监控，关闭连接并释放资源
func StopDiscordMonitor() {
	if discordClient != nil {
		// 关闭 Discord 连接
		err := discordClient.Close()
		if err != nil {
			fmt.Println("error closing Discord session:", err)
		} else {
			fmt.Println("Discord client stopped successfully")
		}
		// 将 discordClient 设置为 nil，释放资源
		discordClient = nil
	} else {
		fmt.Println("Discord client is not running")
	}
}

func GetDiscordClient() *discord.Session {
	return discordClient
}
