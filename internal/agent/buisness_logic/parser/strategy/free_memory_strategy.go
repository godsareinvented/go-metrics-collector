package strategy

import (
	agentdto "github.com/oldhanasong/go-metrics-collector/internal/agent/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/general/dictionary"
	generaldto "github.com/oldhanasong/go-metrics-collector/internal/general/dto"
)

type FreeMemoryStrategy struct {
	value float64
}

func (strategy *FreeMemoryStrategy) FillMetric(metric *generaldto.Metrics, metricData *agentdto.CollectedMetricData) {
	strategy.value = float64(metricData.VirtualMemoryStats.Free)

	metric.ID = dictionary.FreeMemoryMetricName
	metric.MType = dictionary.GaugeMetricType
	metric.Delta = nil
	metric.Value = &strategy.value
}
