package main

import (
	"context"
	"github.com/aliyasirnac/github-pr-webhook-bot/app"
	serverapp "github.com/aliyasirnac/github-pr-webhook-bot/app/serverApp"
	"github.com/aliyasirnac/github-pr-webhook-bot/internal/apps/webhook/webhookConfig"
	"go.uber.org/zap"
	"log"
)

func main() {
	ctx := context.Background()
	config, err := webhookConfig.LoadConfig(ctx)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	serverApp := serverapp.New(config)
	err = app.InitApp(serverApp)
	if err != nil {
		zap.L().Fatal("Error initializing app", zap.Error(err))
		panic(err)
	}
}
