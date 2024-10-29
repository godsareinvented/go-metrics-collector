package strategy

import (
	"github.com/godsareinvented/go-metrics-collector/internal/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
)

type HeapInuseStrategy struct{}

func (strategy *HeapInuseStrategy) GetMetric(metricName string, metricData dto.CollectedMetricData) dto.Metrics {
	var value = float64(metricData.MemStats.HeapInuse)
	return dto.Metrics{
		MType: dictionary.GaugeMetricType,
		MName: metricName,
		Value: &value,
	}
}
