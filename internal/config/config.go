package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
	"time"
)

type Config struct {
	App AppConfig `yaml:"app"`
	DB  DBConfig  `yaml:"db"`
}

type AppConfig struct {
	Port            string        `yaml:"port"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

type DBConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	User string `yaml:"user"`
	Pass string `yaml:"password"`
	Name string `yaml:"name"`
}

func Load() (Config, error) {
	cfg := Config{}

	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return Config{}, fmt.Errorf("failed to read config.yaml: %w", err)
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("failed to parse config.yaml: %w", err)
	}

	slog.Info("configuration loaded from config.yaml")
	return cfg, nil
}
