package interfaces

import (
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
)

type ValueHandler interface {
	GetMutatedValueMetric(metric dto.Metric) dto.Metric
}
