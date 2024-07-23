package interfaces

import (
	"github.com/oldhanasong/go-metrics-collector/internal/constraint"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
)

type ParsingStrategy[Num constraint.Numeric] interface {
	GetMetric(metricName string, metricData dto.CollectedMetricData) dto.Metric[Num]
}
