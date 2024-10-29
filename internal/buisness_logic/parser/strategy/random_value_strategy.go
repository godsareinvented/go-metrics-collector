package strategy

import (
	"github.com/oldhanasong/go-metrics-collector/internal/dictionary"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
)

type RandomValueStrategy struct{}

func (strategy *RandomValueStrategy) GetMetric(metricName string, metricData dto.CollectedMetricData) dto.Metrics {
	return dto.Metrics{
		ID:    metricName,
		MType: dictionary.GaugeMetricType,
		Value: &metricData.RandomValue,
	}
}
