package repository

import (
	"context"
	"errors"
	"github.com/oldhanasong/go-metrics-collector/internal/general/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/server/interfaces"
)

type Repository struct {
	storage *interfaces.Storage
}

var (
	ErrDontImplementConnectorInterface = errors.New("storage don't implement connector interface")
)

func (repository *Repository) UpdateMetric(ctx context.Context, metric dto.Metrics) error {
	return (*repository.storage).Set(ctx, metric)
}

func (repository *Repository) UpdateMetricBatch(ctx context.Context, metrics []dto.Metrics) error {
	return (*repository.storage).SetBatch(ctx, metrics)
}

func (repository *Repository) GetMetric(ctx context.Context, metric dto.Metrics) (dto.Metrics, bool, error) {
	return (*repository.storage).Get(ctx, metric)
}

func (repository *Repository) GetAllMetrics(ctx context.Context) ([]dto.Metrics, error) {
	return (*repository.storage).GetAll(ctx)
}

func (repository *Repository) CloseStorage() error {
	if connector, ok := (*repository.storage).(interfaces.StorageConnector); ok {
		return connector.Close()
	}

	return ErrDontImplementConnectorInterface
}

func (repository *Repository) PingStorage(ctx context.Context) (bool, error) {
	if connector, ok := (*repository.storage).(interfaces.StorageConnector); ok {
		return connector.Ping(ctx)
	}

	return false, ErrDontImplementConnectorInterface
}

func New(storageInterface interfaces.Storage) *Repository {
	return &Repository{storage: &storageInterface}
}
