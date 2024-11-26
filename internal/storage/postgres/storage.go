package postgres

import (
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/interfaces"
)

type PostgreSQLStorage struct{}

func (memStorage *PostgreSQLStorage) GetAll() ([]dto.Metrics, error) {
	// ...
	return nil, nil
}

func (memStorage *PostgreSQLStorage) Get(key string) (dto.Metrics, bool, error) {
	// ...
	return dto.Metrics{}, false, nil
}

func (memStorage *PostgreSQLStorage) Set(key string, metric dto.Metrics) error {
	// ...
	return nil
}

func NewInstance(dbDsn string) interfaces.Storage {
	return &PostgreSQLStorage{}
}
