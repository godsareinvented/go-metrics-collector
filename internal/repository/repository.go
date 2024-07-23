package repository

import (
	"encoding/json"
	"github.com/oldhanasong/go-metrics-collector/internal/constraint"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/interfaces"
)

type Repository[Num constraint.Numeric] struct {
	storage interfaces.Storage
}

func (repository *Repository[Num]) UpdateMetric(metric dto.Metric[Num]) {
	key := getKey(metric)
	value, _ := json.Marshal(metric)
	repository.storage.Set(key, value)
}

func (repository *Repository[Num]) GetMetric(metric dto.Metric[Num]) (dto.Metric[Num], bool) {
	key := getKey(metric)
	jsonMetric := repository.storage.Get(key)
	if nil == jsonMetric {
		return dto.Metric[Num]{}, false
	}

	var metricDTO dto.Metric[Num]
	err := json.Unmarshal(jsonMetric.([]uint8), &metricDTO)

	if nil != err {
		panic("Cannot unmarshal metric")
	}

	return metricDTO, true
}

func NewInstance[Num constraint.Numeric](storage interfaces.Storage) Repository[Num] {
	return Repository[Num]{storage: storage}
}

func getKey[Num constraint.Numeric](metric dto.Metric[Num]) string {
	return metric.Type + "/" + metric.Name
}
