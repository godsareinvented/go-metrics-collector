package config

import (
	"context"
	"flag"
	"github.com/caarlos0/env"
	"github.com/godsareinvented/go-metrics-collector/internal/interfaces"
	"github.com/godsareinvented/go-metrics-collector/internal/logger"
	"github.com/godsareinvented/go-metrics-collector/internal/permanent_storage/file"
	"github.com/godsareinvented/go-metrics-collector/internal/repository"
	"github.com/godsareinvented/go-metrics-collector/internal/storage/mem_storage"
	"github.com/godsareinvented/go-metrics-collector/internal/storage/postgres"
	"os"
	"strings"
)

type ConfigConfigurator struct{}

func (c *ConfigConfigurator) ParseConfig() {
	Configuration = Config{
		GzipAcceptedContentTypes: []string{"application/json", "text/html"},
		GzipMinContentLength:     1400,
		Logger:                   logger.NewInstance(),
	}

	flag.StringVar(&Configuration.Endpoint, "a", ":8080", "Адрес эндпоинта HTTP-сервера")
	flag.IntVar(&Configuration.ReportInterval, "r", 10, "Частота отправки метрик на сервер")
	flag.IntVar(&Configuration.PollInterval, "p", 2, "Частота опроса метрик из пакета runtime")
	flag.IntVar(&Configuration.StoreInterval, "i", 300, "Интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск")
	flag.StringVar(&Configuration.FileStoragePath, "f", getFileStoragePathDefaultValue(), "Путь до файла, куда сохраняются текущие значения")
	flag.BoolVar(&Configuration.Restore, "e", true, "Булево значение, определяющее, загружать или нет ранее сохранённые значения из указанного файла при старте сервера")
	flag.StringVar(&Configuration.DatabaseDSN, "d", "", "Адрес подключения к БД")

	flag.Parse()
	// todo: Отрицательные значения?

	err := env.Parse(&Configuration)
	if err != nil {
		panic("Error parsing environment variables")
	}

	// todo Для клиента не нужна инициализация хранилищ...
	storage := getSuitableStorage()
	permanentStorage := file.NewInstance(Configuration.FileStoragePath)

	Configuration.Repository = repository.NewInstance(&storage)
	Configuration.PermanentStorage = &permanentStorage
}

func getFileStoragePathDefaultValue() string {
	filePathParts := []string{os.TempDir(), "metrics_snapshot.txt"}
	return strings.Join(filePathParts, string(os.PathSeparator))
}

func getSuitableStorage() interfaces.StorageInterface {
	if "" == Configuration.DatabaseDSN {
		return mem_storage.NewInstance()
	}

	storage := postgres.NewInstance(Configuration.DatabaseDSN)
	if res, err := storage.Ping(context.Background()); nil == err && res {
		return storage
	}

	return mem_storage.NewInstance()
}
