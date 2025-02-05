package storage

import (
	"errors"
	"github.com/godsareinvented/go-metrics-collector/internal/server/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/server/interfaces"
	"github.com/godsareinvented/go-metrics-collector/internal/server/storage/mem_storage"
	"github.com/godsareinvented/go-metrics-collector/internal/server/storage/postgressql"
	_ "github.com/golang-migrate/migrate/source/file" // todo: Для чего необходим импорт конкретно? Без него ошика unknown driver
	_ "github.com/jackc/pgx/v5/stdlib"
)

type StorageConfig struct {
	StorageType string
	DatabaseDSN string
}

var (
	ErrUnknownStorageType = errors.New("unknown storage type")
)

func GetStorageAndConfigurator(sc StorageConfig) (interfaces.StorageInterface, interfaces.StorageConfiguratorInterface, error) {
	switch sc.StorageType {
	case dictionary.MemStorage:
		idIndex := make(map[string]int)
		nameIndex := make(map[string]int)
		return mem_storage.NewInstance(idIndex, nameIndex), &mem_storage.MemStorageConfigurator{}, nil
	case dictionary.PostgresqlStorage:
		db, err := postgressql.GetOpenedConnectionWithRetry(sc.DatabaseDSN)
		if nil != err {
			return nil, nil, err
		}
		return postgressql.NewInstance(db), &postgressql.PostgreSQLConfigurator{Db: db}, nil
	default:
		return nil, nil, ErrUnknownStorageType
	}
}
