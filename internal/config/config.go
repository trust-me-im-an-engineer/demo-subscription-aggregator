package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log/slog"
)

type Config struct {
	AppPort string `env:"APP_PORT" env-required:"true"`
	DB      DBConfig
}
type DBConfig struct {
	Host string `env:"DB_HOST" env-default:"localhost"`
	Port int    `env:"DB_PORT" env-default:"5432"`
	User string `env:"DB_USER" env-required:"true"`
	Pass string `env:"DB_PASSWORD" env-required:"true"`
	Name string `env:"DB_NAME" env-required:"true"`
}

func Load() (Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		return Config{}, fmt.Errorf("failed to load .env: %w", err)
	}

	cfg := Config{}
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return Config{}, fmt.Errorf("failed to read env: %w", err)
	}

	slog.Info("Environment variables loaded")
	return cfg, nil
}
