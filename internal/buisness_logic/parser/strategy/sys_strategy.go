package strategy

import (
	"github.com/godsareinvented/go-metrics-collector/internal/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
)

type SysStrategy struct{}

func (strategy *SysStrategy) GetMetric(metricName string, metricData dto.CollectedMetricData) dto.Metrics {
	var value = float64(metricData.MemStats.Sys)
	return dto.Metrics{
		MType: dictionary.GaugeMetricType,
		MName: metricName,
		Value: &value,
	}
}
