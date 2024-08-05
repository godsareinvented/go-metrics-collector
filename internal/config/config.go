package config

import (
	"github.com/godsareinvented/go-metrics-collector/internal/repository"
	"go.uber.org/zap"
)

type Config struct {
	Endpoint       string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	Repository     *repository.Repository
	Logger         *zap.Logger
}

var Configuration Config
