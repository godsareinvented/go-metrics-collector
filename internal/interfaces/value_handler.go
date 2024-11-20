package interfaces

import (
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
)

type ValueHandlerInterface interface {
	GetMutatedValueMetric(metric dto.Metrics, metricFromStorage dto.Metrics, isSetMetricIsStorage bool) dto.Metrics
}
