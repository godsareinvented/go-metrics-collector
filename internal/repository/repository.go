package repository

import (
	"fmt"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/interfaces"
)

type Repository struct {
	storage *interfaces.Storage
}

func (repository *Repository) UpdateMetric(metric dto.Metrics) error {
	return (*repository.storage).Set(key(metric), metric)
}

func (repository *Repository) GetMetric(metric dto.Metrics) (dto.Metrics, bool, error) {
	return (*repository.storage).Get(key(metric))
}

func (repository *Repository) GetAllMetrics() ([]dto.Metrics, error) {
	return (*repository.storage).GetAll()
}

func key(metric dto.Metrics) string {
	return fmt.Sprintf("%s/%s", metric.MType, metric.ID)
}

func NewInstance(storageInterface *interfaces.Storage) *Repository {
	return &Repository{storage: storageInterface}
}
