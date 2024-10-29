package strategy

import (
	"github.com/oldhanasong/go-metrics-collector/internal/dictionary"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
)

type PollCountStrategy struct{}

func (strategy *PollCountStrategy) GetMetric(metricName string, metricData dto.CollectedMetricData) dto.Metrics {
	var value = metricData.PollCount
	return dto.Metrics{
		ID:    metricName,
		MType: dictionary.CounterMetricType,
		Delta: &value,
	}
}
