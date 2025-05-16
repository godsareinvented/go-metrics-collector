package repository

import (
	"encoding/json"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/interfaces"
)

type Repository struct {
	storage interfaces.Storage
}

var (
	repository Repository
)

func (repository *Repository) UpdateMetric(metric dto.Metric) {
	key := getKey(metric)
	value, _ := json.Marshal(metric)
	repository.storage.Set(key, value)
}

func (repository *Repository) GetMetric(metric dto.Metric) (dto.Metric, bool) {
	key := getKey(metric)
	jsonMetric := repository.storage.Get(key)
	if jsonMetric == "" {
		return dto.Metric{}, false
	}

	var metricDTO dto.Metric
	err := json.Unmarshal(jsonMetric.([]byte), &metricDTO)

	if err != nil {
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
