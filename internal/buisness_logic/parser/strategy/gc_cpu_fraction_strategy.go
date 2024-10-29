package strategy

import (
	"github.com/godsareinvented/go-metrics-collector/internal/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
)

type GCCPUFractionStrategy struct{}

func (strategy *GCCPUFractionStrategy) GetMetric(metricName string, metricData dto.CollectedMetricData) dto.Metrics {
	var value = metricData.MemStats.GCCPUFraction
	return dto.Metrics{
		MType: dictionary.GaugeMetricType,
		MName: metricName,
		Value: &value,
	}
}
