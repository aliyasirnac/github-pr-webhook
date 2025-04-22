package serverapp

import (
	"context"
	"github.com/aliyasirnac/github-pr-webhook-bot/internal/apps/webhook/webhookConfig"
	"github.com/aliyasirnac/github-pr-webhook-bot/internal/apps/webhook/webhookHandler"
	"github.com/aliyasirnac/github-pr-webhook-bot/pkg/pubsubinterface"
	"github.com/aliyasirnac/github-pr-webhook-bot/pkg/tracing"
	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"time"
)

import (
	_ "github.com/aliyasirnac/github-pr-webhook-bot/internal/infra/postgres"
)

type ServerApp struct {
	Config *webhookConfig.ServerConfig
	app    *fiber.App
	PubSub pubsubinterface.PubSub
}

func New(config *webhookConfig.ServerConfig, sub pubsubinterface.PubSub) *ServerApp {
	return &ServerApp{
		Config: config,
		PubSub: sub,
	}
}

func (s *ServerApp) Start(ctx context.Context) error {
	zap.L().Info("Starting database")
	tr := tracing.New()
	_ = tr.InitTracer(s.Config)

	//_, err := postgres.New(s.Config.Postgres.Dsn)
	//if err != nil {
	//	zap.L().Fatal("Failed to connect to database", zap.Error(err))
	//	return err
	//}

	githubWebhookHandler := webhookHandler.NewGithubWebhookHandler(s.PubSub)
	zap.L().Info("Server start")

	errCh := make(chan error)
	app := fiber.New(fiber.Config{})
	s.app = app
	app.Use(otelfiber.Middleware())
	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))
	app.Post("/webhook", handle(githubWebhookHandler))

	go func() {
		err := app.Listen(":8080")
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
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := s.PubSub.Close(); err != nil {
		zap.L().Error("Error closing PubSub", zap.Error(err))
		return err
	}
	//close database
	if err := s.app.ShutdownWithTimeout(5 * time.Second); err != nil {
		zap.L().Error("Server shutdown failed", zap.Error(err))
		return err
	}

	zap.L().Info("Server shutdown")
	return nil
}
