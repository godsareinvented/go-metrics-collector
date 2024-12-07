package storage

import (
	"database/sql"
	"errors"
	"github.com/godsareinvented/go-metrics-collector/internal/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/interfaces"
	"github.com/godsareinvented/go-metrics-collector/internal/storage/mem_storage"
	"github.com/godsareinvented/go-metrics-collector/internal/storage/postgressql"
	_ "github.com/golang-migrate/migrate/source/file" // todo: Для чего необходим импорт конкретно? Без него ошика unknown driver
	_ "github.com/jackc/pgx/v5/stdlib"
)

type StorageConfig struct {
	StorageType string
	DatabaseDSN string
}

func GetStorageAndConfigurator(sc StorageConfig) (interfaces.StorageInterface, interfaces.StorageConfiguratorInterface, error) {
	switch sc.StorageType {
	case dictionary.MemStorage:
		return mem_storage.NewInstance(), &mem_storage.MemStorageConfigurator{}, nil
	case dictionary.PostgresqlStorage:
		db, err := getPostgreSQLOpenedConnection(sc)
		if nil != err {
			return nil, nil, err
		}
		return postgressql.NewInstance(db), &postgressql.PostgreSQLConfigurator{Db: db}, nil
	default:
		return nil, nil, errors.New("unknown storage type")
	}
}

func getPostgreSQLOpenedConnection(sc StorageConfig) (*sql.DB, error) {
	if "" == sc.DatabaseDSN {
		return nil, errors.New("DATABASE_DSN is empty")
	}

	db, err := sql.Open("pgx", sc.DatabaseDSN)
	if nil != err {
		return nil, err
	}

	return db, nil
}
