package handler

import (
	"github.com/oldhanasong/go-metrics-collector/internal/general/dto"
)

type GaugeValueHandler struct{}

func (handler *GaugeValueHandler) UpdateMetricValue(metric dto.Metrics, _ dto.Metrics, _ bool) dto.Metrics {
	return metric
}
