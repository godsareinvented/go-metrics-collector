package handler

import (
	"github.com/oldhanasong/go-metrics-collector/internal/general/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/server/repository"
)

type CounterValueHandler struct {
	Repository *repository.Repository
}

func (h *CounterValueHandler) UpdateMetricValue(metric dto.Metrics, metricFromStorage dto.Metrics, isSetInStorage bool) dto.Metrics {
	if isSetInStorage {
		*metric.Delta += *metricFromStorage.Delta
	}
	return metric
}
