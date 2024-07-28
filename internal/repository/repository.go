package repository

import (
	"encoding/json"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/interfaces"
)

type Repository struct {
	storage interfaces.Storage
}

var repository Repository

func (repository *Repository) UpdateMetric(metricDTO dto.Metric) {
	key := getKey(metricDTO)
	value, _ := json.Marshal(metricDTO)
	repository.storage.Set(key, value)
}

func (repository *Repository) GetMetric(metric dto.Metric) (dto.Metric, bool) {
	key := getKey(metric)
	jsonMetricDTO := repository.storage.Get(key)
	if "" == jsonMetricDTO {
		return dto.Metric{}, false
	}

	var metricDTO dto.Metric
	err := json.Unmarshal(jsonMetricDTO.([]byte), &metricDTO)

	// Необходимо для преобразования значения метрики к корректному (согласно типу метрики),
	// т.к. парсер json'а распознаёт любое значение как float64
	// todo: Позже изменить.
	metricDTO.Value = float64(metricDTO.Value)

	if nil != err {
		panic("Cannot unmarshal metric")
	}

	return metricDTO, true
}

func NewInstance(storage interfaces.Storage) {
	repository = Repository{storage: storage}
}

func GetInstance() Repository {
	return repository
}

func getKey(metric dto.Metric) string {
	return metric.Type + "/" + metric.Name
}
