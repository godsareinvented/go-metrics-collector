package repository

import (
	"encoding/json"
	"github.com/oldhanasong/go-metrics-collector/internal/constraint"
	"github.com/oldhanasong/go-metrics-collector/internal/dictionary"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/interfaces"
)

type Repository[Num constraint.Numeric] struct {
	storage interfaces.Storage
}

var (
	int64Repository   Repository[int64]
	float64Repository Repository[float64]
)

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

	// Необходимо для преобразования значения метрики к корректному (согласно типу метрики),
	// т.к. парсер json'а распознаёт любое значение как float64
	metricDTO.Value = Num(metricDTO.Value)

	if nil != err {
		panic("Cannot unmarshal metric")
	}

	return metricDTO, true
}

func NewInstance(storage interfaces.Storage) {
	int64Repository = Repository[int64]{storage: storage}
	float64Repository = Repository[float64]{storage: storage}
}

func GetInstance[Num constraint.Numeric](metricDTO dto.Metric[Num]) Repository[Num] {
	// todo: Тоже проблема. Надо передавать по ссылке.
	switch metricDTO.Type {
	case dictionary.GaugeMetricType:
		return Repository[Num](float64Repository)
	case dictionary.CounterMetricType:
		return Repository[Num](int64Repository)
	default:
		panic("Unknown metric type")
	}
}

func getKey[Num constraint.Numeric](metric dto.Metric[Num]) string {
	return metric.Type + "/" + metric.Name
}
