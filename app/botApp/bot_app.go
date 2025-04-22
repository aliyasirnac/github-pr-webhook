package botapp

import (
	"context"
	"github.com/aliyasirnac/github-pr-webhook-bot/internal/apps/bot/botconfig"
	"github.com/aliyasirnac/github-pr-webhook-bot/pkg/pubsubinterface"
	"go.uber.org/zap"
)

type BotApp struct {
	PubSub pubsubinterface.PubSub
	Config *botconfig.BotConfig
}

func New(pubsub pubsubinterface.PubSub, config *botconfig.BotConfig) *BotApp {
	return &BotApp{PubSub: pubsub, Config: config}
}

func (s *BotApp) Start(ctx context.Context) error {
	err := s.PubSub.Subscribe("bot", func(message pubsubinterface.Message) {
		zap.L().Info("Received message", zap.Any("message", message))
	})
	if err != nil {
		zap.L().Error("failed to subscribe", zap.Error(err))
		return err
	}
	return nil
}

func (s *BotApp) Stop(ctx context.Context) error {
	err := s.PubSub.Close()
	if err != nil {
		zap.L().Error("Error closing PubSub", zap.Error(err))
		return err
	}
	return nil
}
