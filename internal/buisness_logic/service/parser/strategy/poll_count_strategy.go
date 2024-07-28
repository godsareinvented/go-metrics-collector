package strategy

import (
	"github.com/oldhanasong/go-metrics-collector/internal/dictionary"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
)

type PollCountStrategy struct{}

func (strategy *PollCountStrategy) GetMetric(metricName string, metricData dto.CollectedMetricData) dto.Metric {
	return dto.Metric{
		Type:  dictionary.CounterMetricType,
		Name:  metricName,
		Delta: metricData.PollCount,
	}
}
