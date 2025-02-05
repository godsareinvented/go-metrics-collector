package handler

import (
	"github.com/godsareinvented/go-metrics-collector/internal/general/dto"
)

type CounterValueHandler struct{}

func (handler *CounterValueHandler) GetMutatedValueMetric(metric dto.Metrics, metricFromStorage dto.Metrics, isSetMetricInStorage bool) dto.Metrics {
	if isSetMetricInStorage {
		*metric.Delta += *metricFromStorage.Delta
	}
	return metric
}
