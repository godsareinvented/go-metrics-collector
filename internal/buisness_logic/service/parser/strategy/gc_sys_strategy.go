package strategy

import (
	"github.com/godsareinvented/go-metrics-collector/internal/constraint"
	"github.com/godsareinvented/go-metrics-collector/internal/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
)

type GCSysStrategy[Num constraint.Numeric] struct{}

func (strategy *GCSysStrategy[Num]) GetMetric(metricName string, metricData dto.CollectedMetricData) dto.Metric[Num] {
	return dto.Metric[Num]{
		Type:  dictionary.GaugeMetricType,
		Name:  metricName,
		Value: Num(metricData.MemStats.GCSys), // float64
	}
}
