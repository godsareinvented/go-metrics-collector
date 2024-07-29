package config

import (
	"github.com/oldhanasong/go-metrics-collector/internal/repository"
)

type Config struct {
	Endpoint       string
	ReportInterval int
	PollInterval   int
	Repository     *repository.Repository
}

var Configuration Config
