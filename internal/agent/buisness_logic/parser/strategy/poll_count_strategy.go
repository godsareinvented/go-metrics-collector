package strategy

import (
	agentDto "github.com/godsareinvented/go-metrics-collector/internal/agent/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/general/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/general/dto"
)

type PollCountStrategy struct{}

func (strategy *PollCountStrategy) ParseMetric(metric *dto.Metrics, metricData *agentDto.CollectedMetricData) error {
	var delta = metricData.PollCount
	metric.MType = dictionary.CounterMetricType
	metric.MName = dictionary.PollCountMetricName
	metric.Delta = &delta

	return nil
}
