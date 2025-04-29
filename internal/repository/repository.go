package repository

import (
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/interfaces"
	"github.com/oldhanasong/go-metrics-collector/internal/storage/mem_storage"
)

var MetricRepository = Repository{storage: mem_storage.NewInstance()}

type Repository struct {
	storage interfaces.Storage
}

func (repository *Repository) UpdateMetric(metric dto.Metric) {
	key := getKey(metric)
	value := metric.Value
	repository.storage.Set(key, value)
}

func (repository *Repository) GetMetric(metric dto.Metric) (dto.Metric, bool) {
	key := getKey(metric)
	value := repository.storage.Get(key)
	if value == nil {
		return dto.Metric{}, false
	}
	return dto.Metric{Type: metric.Type, Name: metric.Name, Value: value}, true
}

func getKey(metric dto.Metric) string {
	return metric.Type + "/" + metric.Name
}
