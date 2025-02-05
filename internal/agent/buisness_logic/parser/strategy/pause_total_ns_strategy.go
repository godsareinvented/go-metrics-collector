package strategy

import (
	dto2 "github.com/godsareinvented/go-metrics-collector/internal/agent/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/general/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/general/dto"
)

type PauseTotalNsStrategy struct{}

func (strategy *PauseTotalNsStrategy) GetMetric(metricName string, metricData dto2.CollectedMetricData) dto.Metrics {
	var value = float64(metricData.MemStats.PauseTotalNs)
	return dto.Metrics{
		MType: dictionary.GaugeMetricType,
		MName: metricName,
		Value: &value,
	}
}
