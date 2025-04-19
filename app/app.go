package app

import (
	"context"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

type App interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

func InitApp(app App) error {
	parentCtx := context.Background()
	closeChan := make(chan os.Signal, 2)
	signal.Notify(closeChan, syscall.SIGTERM, syscall.SIGINT)

	zap.L().Info("Server starting")
	go func() {
		err := app.Start(parentCtx)
		if err != nil {
			zap.L().Error("Server start failed", zap.Error(err))
			os.Exit(1)
		}
	}()

	sig := <-closeChan
	zap.L().Info("Caught signal %s: shutting down.", zap.String("signal", sig.String()))
	if err := app.Stop(parentCtx); err != nil {
		zap.L().Fatal("Failed to stop app", zap.Error(err))
		return err
	}
	return nil
}
