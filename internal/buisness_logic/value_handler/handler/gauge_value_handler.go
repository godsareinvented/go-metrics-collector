package handler

import (
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
)

type GaugeValueHandler struct{}

func (handler *GaugeValueHandler) GetMutatedValueMetric(metric dto.Metrics, _ dto.Metrics, _ bool) dto.Metrics {
	return metric
}
