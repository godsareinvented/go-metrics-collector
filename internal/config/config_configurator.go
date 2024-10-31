package config

import (
	"flag"
	"github.com/caarlos0/env"
	"github.com/godsareinvented/go-metrics-collector/internal/logger"
	"github.com/godsareinvented/go-metrics-collector/internal/repository"
	"github.com/godsareinvented/go-metrics-collector/internal/storage/mem_storage"
)

type ConfigConfigurator struct{}

func (c *ConfigConfigurator) ParseConfig() {
	// todo Для клиента не нужна инициализация MemStorage...
	memStorage := mem_storage.NewInstance()

	Configuration = Config{
		GzipAcceptedContentTypes: []string{"application/json", "text/html"},
		GzipMinContentLength:     1400,
		Repository:               repository.NewInstance(&memStorage),
		Logger:                   logger.NewInstance(),
	}

	flag.StringVar(&Configuration.Endpoint, "a", ":8080", "The endpoint of the collector")
	flag.IntVar(&Configuration.ReportInterval, "r", 10, "The interval of reporting metrics")
	flag.IntVar(&Configuration.PollInterval, "p", 2, "The interval of polling metrics")

	flag.Parse()
	// todo: Отрицательные значения?

	err := env.Parse(&Configuration)
	if err != nil {
		panic("Error parsing environment variables")
	}
}
