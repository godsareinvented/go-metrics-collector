package mem

import (
	agentdto "github.com/oldhanasong/go-metrics-collector/internal/agent/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/general/business_logic/dictionary"
	generaldto "github.com/oldhanasong/go-metrics-collector/internal/general/dto"
)

type GCCPUFractionStrategy struct{}

func (strategy *GCCPUFractionStrategy) FillMetric(metric *generaldto.Metrics, metricData *agentdto.CollectedMetricData) {
	metric.ID = dictionary.GCCPUFractionMetricName
	metric.MType = dictionary.GaugeMetricType
	metric.Delta = nil
	metric.Value = &metricData.MemStats.GCCPUFraction
}
