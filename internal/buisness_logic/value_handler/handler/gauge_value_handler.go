package handler

import (
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
)

type GaugeValueHandler struct{}

func (handler *GaugeValueHandler) GetMutatedValueMetric(metric dto.Metrics, _ dto.Metrics, _ bool) dto.Metrics {
	return metric
}
