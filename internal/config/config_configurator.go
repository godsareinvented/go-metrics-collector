package config

import (
	"flag"
	"fmt"
	"github.com/oldhanasong/go-metrics-collector/internal/repository"
	"github.com/oldhanasong/go-metrics-collector/internal/storage/mem_storage"
	"sync"
)

type (
	ConfigConfigurator struct{}
)

var (
	once sync.Once
)

func (c *ConfigConfigurator) ParseConfig() {
	once.Do(func() {
		memStorage := mem_storage.NewInstance()

		Configuration = Config{
			Repository: repository.NewInstance(&memStorage),
		}

		var endpoint string
		flag.StringVar(&endpoint, "a", "localhost:8080", "The endpoint of the collector")
		flag.IntVar(&Configuration.ReportInterval, "r", 10, "The interval of reporting metrics")
		flag.IntVar(&Configuration.PollInterval, "p", 2, "The interval of polling metrics")

		flag.Parse()

		Configuration.Endpoint = fmt.Sprintf("http://%s", endpoint)
	})
}
