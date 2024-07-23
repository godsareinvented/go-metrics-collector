package interfaces

import (
	"github.com/godsareinvented/go-metrics-collector/internal/constraint"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
)

type ParsingStrategy[Num constraint.Numeric] interface {
	GetMetric(metricName string, metricData dto.CollectedMetricData) dto.Metric[Num]
}
