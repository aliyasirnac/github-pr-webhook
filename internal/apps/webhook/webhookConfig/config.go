package webhookConfig

import "context"

type ServerConfig struct {
	Postgres Postgres
}

type Postgres struct {
	Dsn string
}

func LoadConfig(ctx context.Context) (*ServerConfig, error) {
	//TODO: load config from yml viper or koanfyaml
	return &ServerConfig{}, nil
}
