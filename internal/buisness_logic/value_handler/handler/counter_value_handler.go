package handler

import (
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/repository"
)

type CounterValueHandler struct {
	Repository *repository.Repository
}

func (handler *CounterValueHandler) GetMutatedValueMetric(metric dto.Metrics) dto.Metrics {
	currentMetricFromStorage, isSet, _ := handler.Repository.GetMetric(metric)
	if isSet {
		*metric.Delta += *currentMetricFromStorage.Delta
	}
	return metric
}
