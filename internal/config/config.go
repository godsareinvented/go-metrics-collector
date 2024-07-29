package config

import (
	"github.com/godsareinvented/go-metrics-collector/internal/repository"
)

type Config struct {
	Endpoint       string
	ReportInterval int
	PollInterval   int
	Repository     *repository.Repository
}

var Configuration Config
