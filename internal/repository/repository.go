package repository

import (
	"context"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/interfaces"
)

type Repository struct {
	storage *interfaces.Storage
}

func (repository *Repository) UpdateMetric(metric dto.Metrics) error {
	return (*repository.storage).Set(metric)
}

func (repository *Repository) GetMetric(metric dto.Metrics) (dto.Metrics, bool, error) {
	return (*repository.storage).Get(metric)
}

func (repository *Repository) GetAllMetrics() ([]dto.Metrics, error) {
	return (*repository.storage).GetAll()
}

func (repository *Repository) CloseStorage() error {
	if connector, ok := (*repository.storage).(interfaces.StorageConnector); ok {
		return connector.Close()
	}

	return nil
}

func (repository *Repository) PingStorage(ctx context.Context) (bool, error) {
	if connector, ok := (*repository.storage).(interfaces.StorageConnector); ok {
		return connector.Ping(ctx)
	}

	return true, nil
}

func NewInstance(storageInterface *interfaces.Storage) *Repository {
	return &Repository{storage: storageInterface}
}
