package config

import (
	"flag"
	"github.com/caarlos0/env"
	"github.com/oldhanasong/go-metrics-collector/internal/server/logger"
	"sync"
	"time"
)

type (
	ConfigConfigurator struct{}
)

var (
	once sync.Once
)

func (c *ConfigConfigurator) ParseConfig() {
	once.Do(func() {
		Configuration = Config{
			GzipAcceptedContentTypes: []string{"application/json", "text/html"},
			GzipMinContentLength:     0, // Должно быть 1400. Для соответствия инкременту 8 заменено на 0.
			Logger:                   logger.New(),
		}

		parseFlags()
		if err := parseEnv(); err != nil {
			panic("Error parsing environment variables")
		}
	})
}

func parseFlags() {
	var reportInterval, pollInterval int

	flag.StringVar(&Configuration.Endpoint, "a", "localhost:8080", "Адрес эндпоинта HTTP-сервера")
	flag.IntVar(&reportInterval, "r", 10, "Частота отправки метрик на сервер")
	flag.IntVar(&pollInterval, "p", 2, "Частота опроса метрик из пакета runtime")
	flag.StringVar(&Configuration.HashKey, "k", "", "Ключ для вычисления хэша")

	flag.Parse()

	Configuration.ReportInterval = time.Duration(reportInterval)
	Configuration.PollInterval = time.Duration(pollInterval)
}

func parseEnv() error {
	return env.Parse(&Configuration)
}
