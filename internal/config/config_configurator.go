package config

import (
	"flag"
	"github.com/caarlos0/env"
	"github.com/godsareinvented/go-metrics-collector/internal/repository"
	"github.com/godsareinvented/go-metrics-collector/internal/storage/mem_storage"
)

type ConfigConfigurator struct{}

func (c *ConfigConfigurator) ParseConfig() {
	memStorage := mem_storage.NewInstance()

	Configuration = Config{
		Repository: repository.NewInstance(&memStorage),
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
