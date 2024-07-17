package handler

import (
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/repository"
	"reflect"
)

type CounterValuePreprocessor struct{}

func (preprocessor *CounterValuePreprocessor) GetMutatedValueMetric(metric dto.Metric) dto.Metric {
	currentMetricFromDb, isSet := repository.MetricRepository.GetMetric(metric)
	if isSet {
		metric.Value = sum(metric.Value, currentMetricFromDb.Value)
	}
	return metric
}

func sum(a interface{}, b interface{}) interface{} {
	if reflect.TypeOf(a).Name() == "int64" {
		return a.(int64) + b.(int64)
	}
	return a.(float64) + b.(float64)
}
