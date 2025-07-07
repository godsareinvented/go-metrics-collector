package strategy

import (
	agentdto "github.com/oldhanasong/go-metrics-collector/internal/agent/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/general/business_logic/dictionary"
	generaldto "github.com/oldhanasong/go-metrics-collector/internal/general/dto"
)

type PollCountStrategy struct{}

func (strategy *PollCountStrategy) FillMetric(metric *generaldto.Metrics, metricData *agentdto.CollectedMetricData) {
	metric.ID = dictionary.PollCountMetricName
	metric.MType = dictionary.CounterMetricType
	metric.Delta = &metricData.PollCount
	metric.Value = nil
}
