package strategy

import (
	"github.com/godsareinvented/go-metrics-collector/internal/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
)

type LookupsStrategy struct{}

func (strategy *LookupsStrategy) GetMetric(metricName string, metricData dto.CollectedMetricData) dto.Metrics {
	var value = float64(metricData.MemStats.Lookups)
	return dto.Metrics{
		MType: dictionary.GaugeMetricType,
		MName: metricName,
		Value: &value,
	}
}
