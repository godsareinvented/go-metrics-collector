package strategy

import (
	dto2 "github.com/godsareinvented/go-metrics-collector/internal/agent/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/general/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/general/dto"
)

type HeapInuseStrategy struct{}

func (strategy *HeapInuseStrategy) GetMetric(metricName string, metricData dto2.CollectedMetricData) dto.Metrics {
	var value = float64(metricData.MemStats.HeapInuse)
	return dto.Metrics{
		MType: dictionary.GaugeMetricType,
		MName: metricName,
		Value: &value,
	}
}
