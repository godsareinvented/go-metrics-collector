package handler

import (
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/repository"
)

type CounterValueHandler struct {
	Repository *repository.Repository
}

func (h *CounterValueHandler) GetMutatedValueMetric(metric dto.Metrics, metricFromStorage dto.Metrics, isSetMetricIsStorage bool) dto.Metrics {
	if isSetMetricIsStorage {
		*metric.Delta += *metricFromStorage.Delta
	}
	return metric
}
