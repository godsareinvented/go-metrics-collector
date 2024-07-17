package handler

import "github.com/godsareinvented/go-metrics-collector/internal/dto"

type GaugeValuePreprocessor struct{}

func (preprocessor *GaugeValuePreprocessor) GetMutatedValueMetric(metric dto.Metric) dto.Metric {
	return metric
}
