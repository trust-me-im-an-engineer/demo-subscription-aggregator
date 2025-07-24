package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
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

func Load() *Config {
	if err := godotenv.Load(".env"); err != nil {
		slog.Error("Error loading .env file", "error", err)
		os.Exit(1)
	}

	cfg := &Config{}
	if err := cleanenv.ReadEnv(cfg); err != nil {
		slog.Error("Error reading env to config", err)
		os.Exit(1)
	}

	slog.Info("Environment variables loaded")
	return cfg
}
