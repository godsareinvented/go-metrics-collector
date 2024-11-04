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

	flag.StringVar(&Configuration.Endpoint, "a", ":8080", "Адрес эндпоинта HTTP-сервера")
	flag.IntVar(&Configuration.ReportInterval, "r", 10, "Частота отправки метрик на сервер")
	flag.IntVar(&Configuration.PollInterval, "p", 2, "Частота опроса метрик из пакета runtime")
	flag.IntVar(&Configuration.StoreInterval, "i", 300, "Интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск")
	flag.StringVar(&Configuration.FileStoragePath, "f", "metric_storage.txt", "Путь до файла, куда сохраняются текущие значения")
	flag.BoolVar(&Configuration.Restore, "e", true, "Булево значение, определяющее, загружать или нет ранее сохранённые значения из указанного файла при старте сервера")

	flag.Parse()
	// todo: Отрицательные значения?

	err := env.Parse(&Configuration)
	if err != nil {
		panic("Error parsing environment variables")
	}
}
