package webhookConfig

import "context"

type ServerConfig struct {
	Postgres      Postgres
	OpenTelemetry OpenTelemetry
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
