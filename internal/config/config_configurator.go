package config

import (
	"flag"
	"github.com/caarlos0/env"
	"github.com/oldhanasong/go-metrics-collector/internal/logger"
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
			GzipAcceptedContentTypes: []string{"application/json", "text/html"},
			GzipMinContentLength:     0, // Должно быть 1400. Для соответствия инкременту 8 заменено на 0.
			Repository:               repository.NewInstance(&memStorage),
			Logger:                   logger.NewInstance(),
		}

		flag.StringVar(&Configuration.Endpoint, "a", "localhost:8080", "Адрес эндпоинта HTTP-сервера")
		flag.IntVar(&Configuration.ReportInterval, "r", 10, "Частота отправки метрик на сервер")
		flag.IntVar(&Configuration.PollInterval, "p", 2, "Частота опроса метрик из пакета runtime")
		flag.IntVar(&Configuration.StoreInterval, "i", 300, "Интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск")
		flag.StringVar(&Configuration.FileStoragePath, "f", "metric_storage.txt", "Путь до файла, куда сохраняются текущие значения")
		flag.BoolVar(&Configuration.Restore, "e", true, "Булево значение, определяющее, загружать или нет ранее сохранённые значения из указанного файла при старте сервера")

		flag.Parse()

		err := env.Parse(&Configuration)
		if err != nil {
			panic("Error parsing environment variables")
		}
	})
}
