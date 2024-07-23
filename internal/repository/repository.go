package repository

import (
	"encoding/json"
	"github.com/godsareinvented/go-metrics-collector/internal/constraint"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/interfaces"
)

type Repository[Num constraint.Numeric] struct {
	storage interfaces.Storage
}

func (repository *Repository[Num]) UpdateMetric(metricDTO dto.Metric[Num]) {
	key := getKey(metricDTO)
	value, _ := json.Marshal(metricDTO)
	repository.storage.Set(key, value)
}

func (repository *Repository[Num]) GetMetric(metric dto.Metric[Num]) (dto.Metric[Num], bool) {
	key := getKey(metric)
	jsonMetricDTO := repository.storage.Get(key)
	if nil == jsonMetricDTO {
		return dto.Metric[Num]{}, false
	}

	var metricDTO dto.Metric[Num]
	err := json.Unmarshal(jsonMetricDTO.([]uint8), &metricDTO)

	if nil != err {
		panic("Cannot unmarshal metric")
	}

	return metricDTO, true
}

func NewInstance[Num constraint.Numeric](storage interfaces.Storage) Repository[Num] {
	return Repository[Num]{storage: storage}
}

//func (repository *Repository) GetMetricList() []dto.Metric {
//	metricMap := repository.storage.GetAll()
//	for _, metricDTOJson := range metricMap {
//
//	}
//}

func getKey[Num constraint.Numeric](metric dto.Metric[Num]) string {
	return metric.Type + "/" + metric.Name
}
