package repository

import (
	"context"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/interfaces"
)

type Repository struct {
	storage *interfaces.StorageInterface
}

func (repository *Repository) UpdateMetric(ctx context.Context, metric dto.Metrics) (string, error) {
	metricID, err := (*repository.storage).Save(ctx, metric)
	return metricID, err
}

func (repository *Repository) GetMetric(ctx context.Context, metric dto.Metrics) (dto.Metrics, bool, error) {
	var foundMetric dto.Metrics
	var isSet = false
	var err error

	if "" != metric.ID {
		foundMetric, isSet, err = (*repository.storage).GetByID(ctx, metric.ID, metric.MType)
		if isSet {
			return foundMetric, isSet, err
		}
	}
	if "" != metric.MName {
		foundMetric, isSet, err = (*repository.storage).GetByName(ctx, metric.MName, metric.MType)
		if isSet {
			return foundMetric, isSet, err
		}
	}

	return foundMetric, isSet, err
}

func (repository *Repository) GetMetricByID(ctx context.Context, metric dto.Metrics) (dto.Metrics, bool, error) {
	foundMetric, isSet, err := (*repository.storage).GetByID(ctx, metric.ID, metric.MType)
	return foundMetric, isSet, err
}

func (repository *Repository) GetMetricByName(ctx context.Context, metric dto.Metrics) (dto.Metrics, bool, error) {
	foundMetric, isSet, err := (*repository.storage).GetByName(ctx, metric.MName, metric.MType)
	return foundMetric, isSet, err
}

func (repository *Repository) GetAllMetrics(ctx context.Context) ([]dto.Metrics, error) {
	list, err := (*repository.storage).GetAll(ctx)
	return list, err
}

func (repository *Repository) CloseStorage() error {
	storageConnector := (*repository.storage).(interfaces.StorageConnectorInterface)
	return storageConnector.CloseConnect()
}

func (repository *Repository) PingStorage(ctx context.Context) (bool, error) {
	storageConnector := (*repository.storage).(interfaces.StorageConnectorInterface)
	return storageConnector.Ping(ctx)
}

func NewInstance(storageInterface *interfaces.StorageInterface) *Repository {
	return &Repository{storage: storageInterface}
}
