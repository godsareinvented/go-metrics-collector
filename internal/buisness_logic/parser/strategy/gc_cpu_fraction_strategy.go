package strategy

import (
	"github.com/oldhanasong/go-metrics-collector/internal/dictionary"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
)

type GCCPUFractionStrategy struct{}

func (strategy *GCCPUFractionStrategy) GetMetric(metricName string, metricData dto.CollectedMetricData) dto.Metrics {
	var value = metricData.MemStats.GCCPUFraction
	return dto.Metrics{
		ID:    metricName,
		MType: dictionary.GaugeMetricType,
		Value: &value,
	}
}
