package mem

import (
	agentdto "github.com/oldhanasong/go-metrics-collector/internal/agent/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/general/business_logic/dictionary"
	generaldto "github.com/oldhanasong/go-metrics-collector/internal/general/dto"
)

type MCacheSysStrategy struct {
	value float64
}

func (strategy *MCacheSysStrategy) FillMetric(metric *generaldto.Metrics, metricData *agentdto.CollectedMetricData) {
	strategy.value = float64(metricData.MemStats.MCacheSys)

	metric.ID = dictionary.MCacheSysMetricName
	metric.MType = dictionary.GaugeMetricType
	metric.Delta = nil
	metric.Value = &strategy.value
}
