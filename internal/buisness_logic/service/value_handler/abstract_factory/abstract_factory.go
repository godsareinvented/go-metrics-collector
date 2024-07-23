package abstract_factory

import (
	"github.com/godsareinvented/go-metrics-collector/internal/buisness_logic/service/value_handler/handler"
	"github.com/godsareinvented/go-metrics-collector/internal/constraint"
	"github.com/godsareinvented/go-metrics-collector/internal/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/interfaces"
	"github.com/godsareinvented/go-metrics-collector/internal/repository"
)

func GetValueHandler[Num constraint.Numeric](metric dto.Metric[Num], repos repository.Repository[Num]) interfaces.ValueHandler[Num] {
	switch metric.Type {
	case dictionary.GaugeMetricType:
		return &handler.GaugeValueHandler[Num]{Repository: &repos}
	case dictionary.CounterMetricType:
		return &handler.CounterValueHandler[Num]{Repository: &repos}
	default:
		panic("unknown metric type")
	}
}
