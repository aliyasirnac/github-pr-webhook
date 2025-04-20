package webhookConfig

import (
	"context"
	"github.com/aliyasirnac/github-pr-webhook-bot/pkg/config"
)

type ServerConfig struct {
	Postgres      Postgres
	OpenTelemetry OpenTelemetry
	Nats          config.NatsConfig
}

type Postgres struct {
	Dsn string
}

type OpenTelemetry struct {
	OtelTraceEndpoint string
}

func LoadConfig(ctx context.Context) (*ServerConfig, error) {
	//TODO: load config from yml viper or koanfyaml
	return &ServerConfig{}, nil
}
