package handler

import (
	"github.com/oldhanasong/go-metrics-collector/internal/constraint"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/repository"
)

type GaugeValuePreprocessor[Num constraint.Numeric] struct {
	Repository *repository.Repository[Num]
}

func (preprocessor *GaugeValuePreprocessor[Num]) GetMutatedValueMetric(metric dto.Metric[Num]) dto.Metric[Num] {
	return metric
}
