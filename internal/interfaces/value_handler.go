package interfaces

import (
	"github.com/godsareinvented/go-metrics-collector/internal/constraint"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
)

type ValueHandler[Num constraint.Numeric] interface {
	GetMutatedValueMetric(metric dto.Metric[Num]) dto.Metric[Num]
}
