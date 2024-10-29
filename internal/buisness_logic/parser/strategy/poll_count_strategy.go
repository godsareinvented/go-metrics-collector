package strategy

import (
	"github.com/godsareinvented/go-metrics-collector/internal/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
)

type PollCountStrategy struct{}

func (strategy *PollCountStrategy) GetMetric(metricName string, metricData dto.CollectedMetricData) dto.Metrics {
	var value = metricData.PollCount
	return dto.Metrics{
		MType: dictionary.CounterMetricType,
		MName: metricName,
		Delta: &value,
	}
}
