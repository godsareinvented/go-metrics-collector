package strategy

import (
	"github.com/godsareinvented/go-metrics-collector/internal/constraint"
	"github.com/godsareinvented/go-metrics-collector/internal/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
)

type StackSysStrategy[Num constraint.Numeric] struct{}

func (strategy *StackSysStrategy[Num]) GetMetric(metricName string, metricData dto.CollectedMetricData) dto.Metric[Num] {
	return dto.Metric[Num]{
		Type:  dictionary.GaugeMetricType,
		Name:  metricName,
		Value: Num(metricData.MemStats.StackSys), // float64
	}
}
