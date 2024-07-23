package strategy

import (
	"github.com/godsareinvented/go-metrics-collector/internal/constraint"
	"github.com/godsareinvented/go-metrics-collector/internal/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
)

type PollCountStrategy[Num constraint.Numeric] struct{}

func (strategy *PollCountStrategy[Num]) GetMetric(metricName string, metricData dto.CollectedMetricData) dto.Metric[Num] {
	return dto.Metric[Num]{
		Type:  dictionary.CounterMetricType,
		Name:  metricName,
		Value: Num(metricData.PollCount), // int64
	}
}
