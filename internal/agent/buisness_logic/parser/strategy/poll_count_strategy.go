package strategy

import (
	agentdto "github.com/oldhanasong/go-metrics-collector/internal/agent/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/general/dictionary"
	generaldto "github.com/oldhanasong/go-metrics-collector/internal/general/dto"
)

type PollCountStrategy struct{}

func (strategy *PollCountStrategy) GetMetric(metricName string, metricData agentdto.CollectedMetricData) generaldto.Metrics {
	var value = metricData.PollCount
	return generaldto.Metrics{
		ID:    metricName,
		MType: dictionary.CounterMetricType,
		Delta: &value,
	}
}
