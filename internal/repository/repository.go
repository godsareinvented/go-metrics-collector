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

	if nil != err {
		panic("Cannot unmarshal metric")
	}

	return metricDTO, true
}

func (repository *Repository) GetAllMetrics() []dto.Metric {
	var resultingList []dto.Metric
	var metricDTO dto.Metric

	metricJsonList := repository.storage.GetAll()
	for _, metricJson := range metricJsonList {
		_ = json.Unmarshal(metricJson.([]byte), &metricDTO)
		resultingList = append(resultingList, metricDTO)
	}

	return resultingList
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
