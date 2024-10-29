package repository

import (
	"encoding/json"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/interfaces"
)

type Repository struct {
	storage *interfaces.Storage
}

func (repository *Repository) UpdateMetric(metric dto.Metrics) {
	key := getKey(metric)
	value, _ := json.Marshal(metric)
	(*repository.storage).Set(key, value)
}

func (repository *Repository) GetMetric(metric dto.Metrics) (dto.Metrics, bool) {
	key := getKey(metric)
	jsonMetric := (*repository.storage).Get(key)
	if jsonMetric == "" {
		return dto.Metrics{}, false
	}

	var metricDTO dto.Metrics
	err := json.Unmarshal(jsonMetric.([]byte), &metricDTO)

	if err != nil {
		panic("Cannot unmarshal metric")
	}

	return metricDTO, true
}

func (repository *Repository) GetAllMetrics() []dto.Metrics {
	var resultingList []dto.Metrics
	var metricDTO dto.Metrics

	metricJsonList := (*repository.storage).GetAll()
	for _, metricJson := range metricJsonList {
		_ = json.Unmarshal(metricJson.([]byte), &metricDTO)
		resultingList = append(resultingList, metricDTO)
	}

	return resultingList
}

func NewInstance(storage *interfaces.Storage) *Repository {
	return &Repository{storage: storage}
}

func getKey(metric dto.Metrics) string {
	return metric.MType + "/" + metric.ID
}
