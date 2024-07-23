package handler

import (
	"github.com/oldhanasong/go-metrics-collector/internal/constraint"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/repository"
)

type GaugeValueHandler[Num constraint.Numeric] struct {
	Repository *repository.Repository[Num]
}

func (preprocessor *GaugeValueHandler[Num]) GetMutatedValueMetric(metric dto.Metric[Num]) dto.Metric[Num] {
	return metric
}
