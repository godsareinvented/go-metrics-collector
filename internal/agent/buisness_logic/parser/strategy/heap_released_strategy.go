package strategy

import (
	dto2 "github.com/godsareinvented/go-metrics-collector/internal/agent/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/general/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/general/dto"
)

type HeapReleasedStrategy struct{}

func (strategy *HeapReleasedStrategy) GetMetric(metricName string, metricData dto2.CollectedMetricData) dto.Metrics {
	var value = float64(metricData.MemStats.HeapReleased)
	return dto.Metrics{
		MType: dictionary.GaugeMetricType,
		MName: metricName,
		Value: &value,
	}
}
