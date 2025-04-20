package main

import (
	"context"
	"github.com/aliyasirnac/github-pr-webhook-bot/app"
	botapp "github.com/aliyasirnac/github-pr-webhook-bot/app/botApp"
	"github.com/aliyasirnac/github-pr-webhook-bot/internal/apps/bot/botconfig"
	_ "github.com/aliyasirnac/github-pr-webhook-bot/internal/logger"
	"github.com/aliyasirnac/github-pr-webhook-bot/internal/pubsub"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()
	config, err := botconfig.LoadConfig()
	if err != nil {
		panic(err)
	}
	natsPubSub := pubsub.NewRabbitMQ()
	botApp := botapp.New(natsPubSub, config)
	err = app.InitApp(botApp, ctx)
	if err != nil {
		zap.L().Fatal("Error initializing app", zap.Error(err))
		panic(err)
	}
}
