package strategy

import (
	"github.com/godsareinvented/go-metrics-collector/internal/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
)

type LookupsStrategy struct{}

func (strategy *LookupsStrategy) GetMetric(metricName string, metricData dto.CollectedMetricData) dto.Metric {
	return dto.Metric{
		Type:  dictionary.GaugeMetricType,
		Name:  metricName,
		Value: float64(metricData.MemStats.Lookups),
	}
}
