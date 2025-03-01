package strategy

import (
	agentdto "github.com/oldhanasong/go-metrics-collector/internal/agent/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/general/dictionary"
	generaldto "github.com/oldhanasong/go-metrics-collector/internal/general/dto"
)

type NextGCStrategy struct {
	value float64
}

func (strategy *NextGCStrategy) FillMetric(metric *generaldto.Metrics, metricData *agentdto.CollectedMetricData) {
	strategy.value = float64(metricData.MemStats.NextGC)

	metric.ID = dictionary.NextGCMetricName
	metric.MType = dictionary.GaugeMetricType
	metric.Delta = nil
	metric.Value = &strategy.value
}
