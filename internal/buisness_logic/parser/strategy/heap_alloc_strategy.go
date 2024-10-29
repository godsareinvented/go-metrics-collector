package strategy

import (
	"github.com/godsareinvented/go-metrics-collector/internal/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
)

type HeapAllocStrategy struct{}

func (strategy *HeapAllocStrategy) GetMetric(metricName string, metricData dto.CollectedMetricData) dto.Metrics {
	var value = float64(metricData.MemStats.HeapAlloc)
	return dto.Metrics{
		MType: dictionary.GaugeMetricType,
		MName: metricName,
		Value: &value,
	}
}
