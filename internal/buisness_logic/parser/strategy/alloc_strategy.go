package strategy

import (
	"github.com/godsareinvented/go-metrics-collector/internal/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
)

type AllocStrategy struct{}

func (strategy *AllocStrategy) GetMetric(metricName string, metricData dto.CollectedMetricData) dto.Metric {
	return dto.Metric{
		Type:  dictionary.GaugeMetricType,
		Name:  metricName,
		Value: float64(metricData.MemStats.Alloc),
	}
}
