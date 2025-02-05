package strategy

import (
	agentdto "github.com/oldhanasong/go-metrics-collector/internal/agent/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/general/dictionary"
	generaldto "github.com/oldhanasong/go-metrics-collector/internal/general/dto"
)

type PauseTotalNsStrategy struct{}

func (strategy *PauseTotalNsStrategy) GetMetric(metricName string, metricData agentdto.CollectedMetricData) generaldto.Metrics {
	var value = float64(metricData.MemStats.PauseTotalNs)
	return generaldto.Metrics{
		ID:    metricName,
		MType: dictionary.GaugeMetricType,
		Value: &value,
	}
}
