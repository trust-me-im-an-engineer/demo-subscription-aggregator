package main

import (
	"github.com/trust-me-im-an-engineer/demo-subscription-agregator/internal/config"
	"github.com/trust-me-im-an-engineer/demo-subscription-agregator/pkg/logger"
)

func main() {
	logger.Init()

	_ = config.Load()
}
