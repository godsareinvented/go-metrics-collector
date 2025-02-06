package config

import (
	"context"
	"errors"
	"flag"
	"github.com/caarlos0/env"
	"github.com/godsareinvented/go-metrics-collector/internal/server/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/server/interfaces"
	"github.com/godsareinvented/go-metrics-collector/internal/server/logger"
	"github.com/godsareinvented/go-metrics-collector/internal/server/permanent_storage/file"
	"github.com/godsareinvented/go-metrics-collector/internal/server/repository"
	storagePackage "github.com/godsareinvented/go-metrics-collector/internal/server/storage"
	"os"
	"strings"
)

type ConfigConfigurator struct{}

var ErrNotStorageConnectorInterface = errors.New("the postgresql storage doesn't implement the StorageConnectorInterface")

func (c *ConfigConfigurator) ParseConfig() {
	Configuration = Config{
		GzipAcceptedContentTypes: []string{"application/json", "text/html"},
		GzipMinContentLength:     1400,
		Logger:                   logger.NewInstance(),
	}

	parseFlags()
	err := parseEnv()
	if nil != err {
		panic("Error parsing env: " + err.Error())
	}

	storage, storageConfigurator, err := getSuitableStorage()
	if nil != err {
		panic("Error configuring storage: " + err.Error())
	}
	err = storageConfigurator.Configure()
	if nil != err {
		panic("Error configuring storage: " + err.Error())
	}
	permanentStorage := file.NewInstance(Configuration.FileStoragePath)

	Configuration.Repository = repository.NewInstance(&storage)
	Configuration.PermanentStorage = &permanentStorage
}

func parseFlags() {
	flag.StringVar(&Configuration.Endpoint, "a", ":8080", "Адрес эндпоинта HTTP-сервера")
	flag.IntVar(&Configuration.StoreInterval, "i", 300, "Интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск")
	flag.StringVar(&Configuration.FileStoragePath, "f", getFileStoragePathDefaultValue(), "Путь до файла, куда сохраняются текущие значения")
	flag.BoolVar(&Configuration.Restore, "e", true, "Булево значение, определяющее, загружать или нет ранее сохранённые значения из указанного файла при старте сервера")
	flag.StringVar(&Configuration.DatabaseDSN, "d", "", "Адрес подключения к БД")
	flag.StringVar(&Configuration.HashKey, "k", "", "Ключ для вычисления хэша")

	flag.Parse()
	// todo: Отрицательные значения?
}

func parseEnv() error {
	return env.Parse(&Configuration)
}

func getFileStoragePathDefaultValue() string {
	filePathParts := []string{os.TempDir(), "metrics_snapshot.txt"}
	return strings.Join(filePathParts, string(os.PathSeparator))
}

func getSuitableStorage() (interfaces.StorageInterface, interfaces.StorageConfiguratorInterface, error) {
	if "" == Configuration.DatabaseDSN {
		return storagePackage.GetStorageAndConfigurator(storagePackage.StorageConfig{StorageType: dictionary.MemStorage})
	}

	storage, storageConfigurator, err := storagePackage.GetStorageAndConfigurator(storagePackage.StorageConfig{
		StorageType: dictionary.PostgresqlStorage,
		DatabaseDSN: Configuration.DatabaseDSN,
	})
	if nil != err {
		return nil, nil, err
	}
	storageConnector, ok := storage.(interfaces.StorageConnectorInterface)
	if !ok {
		return nil, nil, ErrNotStorageConnectorInterface
	}

	if res, err := storageConnector.Ping(context.Background()); nil == err && res {
		return storage, storageConfigurator, nil
	}

	return storagePackage.GetStorageAndConfigurator(storagePackage.StorageConfig{StorageType: dictionary.MemStorage})
}
