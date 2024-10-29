package handler

import (
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/repository"
)

type CounterValueHandler struct {
	Repository *repository.Repository
}

func (handler *CounterValueHandler) GetMutatedValueMetric(metricDTO dto.Metrics) dto.Metrics {
	currentMetricDTOFromDb, isSet := handler.Repository.GetMetric(metricDTO)
	if isSet {
		*metricDTO.Delta += *currentMetricDTOFromDb.Delta
	}
	return metricDTO
}
