package botconfig

import (
	"errors"
	"github.com/aliyasirnac/github-pr-webhook-bot/pkg/config"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/rawbytes"
	koanfyaml "github.com/knadh/koanf/v2"
	"go.uber.org/zap"
	"os"
)

type BotConfig struct {
	Telegram Telegram
	PubSub   config.PubSubConfig
}

type Telegram struct {
	ApiKey    string
	ChannelId string
}

func LoadConfig() (*BotConfig, error) {
	if _, err := os.Stat(".env"); err == nil {
		_ = godotenv.Load()
	} else {
		zap.L().Info(".env file not found, using system environment variables")
	}
	k := koanfyaml.New(".")

	// Load raw bytes from file
	raw, err := os.ReadFile("./config/bot/config.yaml")
	if err != nil {
		zap.L().Error("failed to read config file", zap.Error(err))
		return nil, err
	}

	yamlString := os.ExpandEnv(string(raw))

	if err := k.Load(rawbytes.Provider([]byte(yamlString)), yaml.Parser()); err != nil {
		zap.L().Error("failed to load config from raw bytes", zap.Error(err))
		return nil, err
	}

	var cfg BotConfig
	if err := k.Unmarshal("", &cfg); err != nil {
		zap.L().Error("failed to unmarshal config", zap.Error(err))
		return nil, err
	}

	if cfg.Telegram.ApiKey == "" {
		zap.L().Error("telegram api key is empty")
		return nil, errors.New("telegram api key is empty")
	}

	if cfg.Telegram.ChannelId == "" {
		zap.L().Error("telegram channel id is empty")
		return nil, errors.New("telegram channel id is empty")
	}

	if cfg.PubSub.ConnectionUrl == "" {
		zap.L().Error("pubsub connection url is empty")
		return nil, errors.New("pubsub connection url is empty")
	}

	zap.L().Info("loaded config", zap.Any("config", cfg))
	return &cfg, nil
}
