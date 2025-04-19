package serverapp

import (
	"context"
	"github.com/aliyasirnac/github-pr-webhook-bot/internal/apps/webhook/webhookConfig"
	"github.com/aliyasirnac/github-pr-webhook-bot/internal/apps/webhook/webhookHandler"
	_ "github.com/aliyasirnac/github-pr-webhook-bot/internal/infra/postgres"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"time"
)

type ServerApp struct {
	Config *webhookConfig.ServerConfig
	app    *fiber.App
}

func New(config *webhookConfig.ServerConfig) *ServerApp {
	return &ServerApp{
		Config: config,
	}
}

func (s *ServerApp) Start(ctx context.Context) error {
	zap.L().Info("Starting database")
	//_, err := postgres.New(s.Config.Postgres.Dsn)
	//if err != nil {
	//	zap.L().Fatal("Failed to connect to database", zap.Error(err))
	//	return err
	//}

	githubWebhookHandler := webhookHandler.NewGithubWebhookHandler()
	zap.L().Info("Server start")

	errCh := make(chan error)
	app := fiber.New(fiber.Config{})
	s.app = app
	app.Post("/webhook", handle(githubWebhookHandler))

	go func() {
		err := app.Listen(":8000")
		if err != nil {
			zap.L().Error("Error starting server", zap.Error(err))
			errCh <- err
		}
	}()
	err := <-errCh
	if err != nil {
		zap.L().Error("Server start failed", zap.Error(err))
		return err
	}

	return nil
}

func (s *ServerApp) Stop(ctx context.Context) error {
	//close database
	if err := s.app.ShutdownWithTimeout(5 * time.Second); err != nil {
		zap.L().Error("Server shutdown failed", zap.Error(err))
		return err
	}

	zap.L().Info("Server shutdown")
	return nil
}
