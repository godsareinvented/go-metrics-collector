package handler

import (
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
)

type CounterValueHandler struct{}

func (handler *CounterValueHandler) GetMutatedValueMetric(metric dto.Metrics, metricFromStorage dto.Metrics, isSetMetricIsStorage bool) dto.Metrics {
	if isSetMetricIsStorage {
		*metric.Delta += *metricFromStorage.Delta
	}
	return metric
}
