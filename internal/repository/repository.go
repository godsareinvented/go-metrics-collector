package repository

import (
	"context"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/interfaces"
)

type Repository struct {
	storage *interfaces.StorageInterface
}

func (repository *Repository) UpdateMetric(metric dto.Metrics) (string, error) {
	metricID, err := (*repository.storage).Save(metric)
	return metricID, err
}

func (repository *Repository) GetMetric(metric dto.Metrics) (dto.Metrics, bool, error) {
	var foundMetric dto.Metrics
	var isSet = false
	var err error

	if "" != metric.ID {
		foundMetric, isSet, err = (*repository.storage).GetByID(metric.ID, metric.MType)
		if isSet {
			return foundMetric, isSet, err
		}
	}
	if "" != metric.MName {
		foundMetric, isSet, err = (*repository.storage).GetByName(metric.MName, metric.MType)
		if isSet {
			return foundMetric, isSet, err
		}
	}

	return foundMetric, isSet, err
}

func (repository *Repository) GetMetricByID(metric dto.Metrics) (dto.Metrics, bool, error) {
	foundMetric, isSet, err := (*repository.storage).GetByID(metric.ID, metric.MType)
	return foundMetric, isSet, err
}

func (repository *Repository) GetMetricByName(metric dto.Metrics) (dto.Metrics, bool, error) {
	foundMetric, isSet, err := (*repository.storage).GetByName(metric.MName, metric.MType)
	return foundMetric, isSet, err
}

func (repository *Repository) GetAllMetrics() ([]dto.Metrics, error) {
	list, err := (*repository.storage).GetAll()
	return list, err
}

func (repository *Repository) CloseStorage() error {
	return (*repository.storage).Close()
}

func (repository *Repository) PingStorage(ctx context.Context) (bool, error) {
	return (*repository.storage).Ping(ctx)
}

func NewInstance(storageInterface *interfaces.StorageInterface) *Repository {
	return &Repository{storage: storageInterface}
}
