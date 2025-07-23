package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
)

type Config struct {
	AppPort string `env:"APP_PORT" env-required:"true"`
	DBUser  string `env:"DB_USER" env-required:"true"`
	DBPass  string `env:"DB_PASSWORD" env-required:"true"`
	DBName  string `env:"DB_NAME" env-required:"true"`
}

func Load() *Config {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	cfg := &Config{}
	if err := cleanenv.ReadEnv(cfg); err != nil {
		log.Fatalf("Error reading env to config: %s", err)
	}
	return cfg
}
