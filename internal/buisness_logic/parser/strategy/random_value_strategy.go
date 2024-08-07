package strategy

import (
	"github.com/godsareinvented/go-metrics-collector/internal/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
)

type RandomValueStrategy struct{}

func (strategy *RandomValueStrategy) GetMetric(metricName string, metricData dto.CollectedMetricData) dto.Metric {
	return dto.Metric{
		Type:  dictionary.GaugeMetricType,
		Name:  metricName,
		Value: float64(metricData.RandomValue),
	}
}
