package botconfig

import "github.com/aliyasirnac/github-pr-webhook-bot/pkg/config"

type BotConfig struct {
	Telegram Telegram
	Nats     config.NatsConfig
}

type Telegram struct {
	ApiKey    string
	ChannelId string //TODO to convert int
}

func LoadConfig() (*BotConfig, error) {
	var cfg BotConfig
	return &cfg, nil
}
