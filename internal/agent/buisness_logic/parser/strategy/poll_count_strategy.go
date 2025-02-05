package strategy

import (
	dto2 "github.com/godsareinvented/go-metrics-collector/internal/agent/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/general/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/general/dto"
)

type PollCountStrategy struct{}

func (strategy *PollCountStrategy) GetMetric(metricName string, metricData dto2.CollectedMetricData) dto.Metrics {
	var value = metricData.PollCount
	return dto.Metrics{
		MType: dictionary.CounterMetricType,
		MName: metricName,
		Delta: &value,
	}
}
