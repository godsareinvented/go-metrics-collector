package interfaces

import (
	"github.com/godsareinvented/go-metrics-collector/internal/general/dto"
)

type ValueHandlerInterface interface {
	GetMutatedValueMetric(metric dto.Metrics, metricFromStorage dto.Metrics, isSetMetricIsStorage bool) dto.Metrics
}
