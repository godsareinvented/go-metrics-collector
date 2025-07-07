package config

import (
	"context"
	"errors"
	"github.com/oldhanasong/go-metrics-collector/internal/server/dictionary"
	"github.com/oldhanasong/go-metrics-collector/internal/server/interfaces"
	"github.com/oldhanasong/go-metrics-collector/internal/server/storage"
	"time"
)

type handler func() (interfaces.Storage, interfaces.StorageConfigurator, error)

var creators = []handler{createPostgresStorage, createMemStorage}

func CreateSuitableStorageAndConfigurator() (interfaces.Storage, interfaces.StorageConfigurator) {
	for _, f := range creators {
		if s, sc, err := f(); err == nil {
			return s, sc
		}
	}

	s, sc, _ := createMemStorage()
	return s, sc
}

func createPostgresStorage() (interfaces.Storage, interfaces.StorageConfigurator, error) {
	if Configuration.DatabaseDSN == "" {
		return nil, nil, errors.New("dsn flag not set")
	}
	conf := storage.Config{Type: dictionary.PostgresqlStorage, DSN: Configuration.DatabaseDSN}
	s, sc, err := storage.CreateStorageAndConfigurator(conf)
	if err != nil {
		return nil, nil, err
	}
	c, ok := s.(interfaces.StorageConnector)
	if !ok {
		return nil, nil, errors.New("storage doesn't implement the interface StorageConnector")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ping, err := c.Ping(ctx)
	if err != nil {
		return nil, nil, err
	}
	if !ping {
		return nil, nil, errors.New("failed to ping the storage")
	}

	return s, sc, nil
}

func createMemStorage() (interfaces.Storage, interfaces.StorageConfigurator, error) {
	return storage.CreateStorageAndConfigurator(storage.Config{Type: dictionary.MemStorage})
}
