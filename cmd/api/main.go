package main

import (
	"context"
	"github.com/aliyasirnac/github-pr-webhook-bot/app"
	serverapp "github.com/aliyasirnac/github-pr-webhook-bot/app/serverApp"
	"github.com/aliyasirnac/github-pr-webhook-bot/internal/apps/webhook/webhookConfig"
	_ "github.com/aliyasirnac/github-pr-webhook-bot/internal/logger"
	"github.com/aliyasirnac/github-pr-webhook-bot/internal/pubsub"
	"go.uber.org/zap"
	"log"
)

func main() {
	ctx := context.Background()
	config, err := webhookConfig.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	rabbitmq := pubsub.NewRabbitMQ(&config.PubSub)
	serverApp := serverapp.New(config, rabbitmq)
	err = app.InitApp(serverApp, ctx)
	if err != nil {
		zap.L().Fatal("Error initializing app", zap.Error(err))
		panic(err)
	}
}
