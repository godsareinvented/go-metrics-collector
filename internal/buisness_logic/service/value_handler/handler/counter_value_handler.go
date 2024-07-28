package handler

import (
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/repository"
)

type CounterValueHandler struct {
	Repository *repository.Repository
}

func (h *CounterValueHandler) GetMutatedValueMetric(metric dto.Metric) dto.Metric {
	currentMetricFromDb, isSet := h.Repository.GetMetric(metric)
	if isSet {
		metric.Delta += currentMetricFromDb.Delta
	}
	return metric
}
