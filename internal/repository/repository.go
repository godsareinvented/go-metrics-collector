package repository

import (
	"encoding/json"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/interfaces"
)

type Repository struct {
	storage *interfaces.Storage
}

func (repository *Repository) UpdateMetric(metricDTO dto.Metrics) {
	key := getKey(metricDTO)
	value, _ := json.Marshal(metricDTO)
	(*repository.storage).Set(key, value)
}

func (repository *Repository) GetMetric(metric dto.Metrics) (dto.Metrics, bool) {
	key := getKey(metric)
	jsonMetricDTO := (*repository.storage).Get(key)
	if "" == jsonMetricDTO {
		return dto.Metrics{}, false
	}

	var metricDTO dto.Metrics
	err := json.Unmarshal(jsonMetricDTO.([]byte), &metricDTO)

	if nil != err {
		panic("Cannot unmarshal metric")
	}

	return metricDTO, true
}

func (repository *Repository) GetAllMetrics() []dto.Metrics {
	var resultingList []dto.Metrics

	metricJsonList := (*repository.storage).GetAll()
	for _, metricJson := range metricJsonList {
		var metricDTO dto.Metrics
		_ = json.Unmarshal(metricJson.([]byte), &metricDTO)
		resultingList = append(resultingList, metricDTO)
	}

	return resultingList
}

func NewInstance(storage *interfaces.Storage) *Repository {
	return &Repository{storage: storage}
}

func getKey(metric dto.Metrics) string {
	return metric.MType + "/" + metric.MName
}
