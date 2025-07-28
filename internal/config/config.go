package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"log/slog"
	"time"
)

type Config struct {
	App AppConfig
	DB  DBConfig
}

type AppConfig struct {
	Address         string        `env:"APP_ADDRESS" envDefault:"0.0.0.0:8080"`
	ShutdownTimeout time.Duration `env:"APP_SHUTDOWN_TIMEOUT" envDefault:"10s"`
}

type DBConfig struct {
	Host string `env:"DB_HOST"`
	Port int    `env:"DB_PORT"`
	User string `env:"DB_USER"`
	Pass string `env:"DB_PASSWORD"`
	Name string `env:"DB_NAME"`
}

func Load() (Config, error) {
	cfg := Config{}

	if err := env.Parse(&cfg); err != nil {
		return Config{}, fmt.Errorf("failed to parse config from environment: %w", err)
	}

	slog.Debug("configuration loaded", "config", cfg)

	return cfg, nil
}
