package config

import (
	"context"
	"errors"
	"github.com/oldhanasong/go-metrics-collector/internal/dictionary"
	"github.com/oldhanasong/go-metrics-collector/internal/interfaces"
	"github.com/oldhanasong/go-metrics-collector/internal/storage/mem_storage"
	"github.com/oldhanasong/go-metrics-collector/internal/storage/postgres"
	"time"
)

func suitableStorage() interfaces.Storage {
	for _, storageType := range []string{dictionary.PostgresqlStorage, dictionary.MemStorage} {
		if s, err := createStorage(storageType); err == nil {
			return s
		}
	}

	return nil
}

func createStorage(storageType string) (interfaces.Storage, error) {
	switch storageType {
	case dictionary.PostgresqlStorage:
		return createPostgreSQLStorage()
	case dictionary.MemStorage:
	default:
		return createMemStorage()
	}

	return nil, errors.New("unknown storage type passed")
}

func createPostgreSQLStorage() (interfaces.Storage, error) {
	if Configuration.DatabaseDSN == "" {
		return nil, errors.New("dsn flag not set")
	}
	s, err := postgres.NewInstance(Configuration.DatabaseDSN)
	if err != nil {
		return nil, err
	}
	c, ok := s.(interfaces.StorageConnector)
	if !ok {
		return nil, errors.New("storage doesn't implement the interface StorageConnector")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var ping bool
	ping, err = c.Ping(ctx)
	if err != nil {
		return nil, err
	}

	if !ping {
		return nil, errors.New("failed to ping the storage")
	}

	return s, nil
}

func createMemStorage() (interfaces.Storage, error) {
	return mem_storage.NewInstance(), nil
}
