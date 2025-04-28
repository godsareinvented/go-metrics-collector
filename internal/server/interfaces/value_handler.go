package interfaces

import (
	"github.com/godsareinvented/go-metrics-collector/internal/general/dto"
)

type ValueHandlerInterface interface {
	// GetMutatedValueMetric todo: Стоит переназвать метод и структуры, т.к. это обычная логика прилоения
	GetMutatedValueMetric(metric dto.Metrics, metricFromStorage dto.Metrics, isSetMetricIsStorage bool) dto.Metrics
}
