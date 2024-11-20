package handler

import (
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
)

type GaugeValueHandler struct{}

func (handler *GaugeValueHandler) GetMutatedValueMetric(metric dto.Metrics, metricFromStorage dto.Metrics, isSetMetricIsStorage bool) dto.Metrics {
	return metric
}
