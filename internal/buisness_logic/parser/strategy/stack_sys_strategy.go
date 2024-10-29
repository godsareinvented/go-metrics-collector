package strategy

import (
	"github.com/godsareinvented/go-metrics-collector/internal/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
)

type StackSysStrategy struct{}

func (strategy *StackSysStrategy) GetMetric(metricName string, metricData dto.CollectedMetricData) dto.Metrics {
	var value = float64(metricData.MemStats.StackSys)
	return dto.Metrics{
		MType: dictionary.GaugeMetricType,
		MName: metricName,
		Value: &value,
	}
}
