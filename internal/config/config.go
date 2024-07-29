package config

import (
	"github.com/godsareinvented/go-metrics-collector/internal/repository"
)

type Config struct {
	Endpoint       string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	Repository     *repository.Repository
}

var Configuration Config
