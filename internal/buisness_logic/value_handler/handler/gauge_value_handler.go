package handler

import (
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/repository"
)

type GaugeValueHandler struct {
	Repository *repository.Repository
}

func (handler *GaugeValueHandler) GetMutatedValueMetric(metric dto.Metrics) dto.Metrics {
	return metric
}
