package handler

import (
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/repository"
)

type GaugeValueHandler struct {
	Repository *repository.Repository
}

func (_ *GaugeValueHandler) GetMutatedValueMetric(metric dto.Metric) dto.Metric {
	return metric
}
