package storage

import (
	"errors"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/oldhanasong/go-metrics-collector/internal/server/dictionary"
	"github.com/oldhanasong/go-metrics-collector/internal/server/interfaces"
	"github.com/oldhanasong/go-metrics-collector/internal/server/storage/mem_storage"
	"github.com/oldhanasong/go-metrics-collector/internal/server/storage/postgres"
)

type Config struct {
	Type string
	DSN  string
}

func CreateStorageAndConfigurator(c Config) (interfaces.Storage, interfaces.StorageConfigurator, error) {
	switch c.Type {
	case dictionary.MemStorage:
		return mem_storage.NewStorage(), mem_storage.NewConfigurator(), nil
	case dictionary.PostgresqlStorage:
		db, err := postgres.OpenConnectionWithRetry(c.DSN)
		if err != nil {
			return nil, nil, err
		}
		return postgres.NewStorage(db), postgres.NewConfigurator(db), nil
	default:
		return nil, nil, errors.New("unknown storage type passed")
	}
}
