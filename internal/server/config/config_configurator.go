package config

import (
	"errors"
	"flag"
	"github.com/caarlos0/env"
	"github.com/oldhanasong/go-metrics-collector/internal/server/logger"
	"github.com/oldhanasong/go-metrics-collector/internal/server/permanent_storage/file"
	"github.com/oldhanasong/go-metrics-collector/internal/server/repository"
	"os"
	"strings"
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
		createPermanentStorage()
		if err := createRepository(); err != nil {
			panic(err)
		}
	})
}

func parseFlags() {
	var storeInterval int

	flag.StringVar(&Configuration.Endpoint, "a", "localhost:8080", "Адрес эндпоинта HTTP-сервера")
	flag.IntVar(&storeInterval, "i", 300, "Интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск")
	flag.StringVar(&Configuration.FileStoragePath, "f", getFileStoragePathDefaultValue(), "Путь до файла, куда сохраняются текущие значения")
	flag.BoolVar(&Configuration.Restore, "e", true, "Булево значение, определяющее, загружать или нет ранее сохранённые значения из указанного файла при старте сервера")
	flag.StringVar(&Configuration.DatabaseDSN, "d", "", "Адрес подключения к БД")

	flag.Parse()

	Configuration.StoreInterval = time.Duration(storeInterval)
}

func parseEnv() error {
	return env.Parse(&Configuration)
}

func createPermanentStorage() {
	permanentStorage := file.New(Configuration.FileStoragePath)
	Configuration.PermanentStorage = &permanentStorage
}

func createRepository() error {
	storage, configurator := createSuitableStorageAndConfigurator()
	if storage == nil || configurator == nil {
		return errors.New("storage or configurator is not set")
	}
	if err := configurator.Configure(); err != nil {
		return err
	}
	Configuration.Repository = repository.New(storage)

	return nil
}

func getFileStoragePathDefaultValue() string {
	filePathParts := []string{os.TempDir(), "metrics_snapshot.txt"}
	return strings.Join(filePathParts, string(os.PathSeparator))
}
