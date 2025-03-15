package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"time"

	"github.com/caarlos0/env/v9"
)

type Config struct {
	Env            string        `env:"ENV" envDefault:"development"`
	DBURL          string        `env:"DB_URL" envDefault:"postgres://postgres@postgres:5432/usdt_rates?sslmode=disable"`
	GRPCPort       int           `env:"GRPC_PORT" envDefault:"50051"`
	MetricsPort    int           `env:"METRICS_PORT" envDefault:"9090"`
	JaegerEndpoint string        `env:"JAEGER_ENDPOINT" envDefault:"http://localhost:14268/api/traces"`
	MexcTimeout    time.Duration `env:"MEXC_TIMEOUT" envDefault:"5s"`
	ServiceName    string        `env:"SERVICE_NAME" envDefault:"usdt-rate-service"`
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
