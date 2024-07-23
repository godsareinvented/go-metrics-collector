package handler

import (
	"github.com/oldhanasong/go-metrics-collector/internal/constraint"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/repository"
)

type CounterValuePreprocessor[Num constraint.Numeric] struct {
	Repository *repository.Repository[Num]
}

func (preprocessor *CounterValuePreprocessor[Num]) GetMutatedValueMetric(metric dto.Metric[Num]) dto.Metric[Num] {
	currentMetricFromDb, isSet := preprocessor.Repository.GetMetric(metric)
	if isSet {
		metric.Value += currentMetricFromDb.Value
	}
	return metric
}
