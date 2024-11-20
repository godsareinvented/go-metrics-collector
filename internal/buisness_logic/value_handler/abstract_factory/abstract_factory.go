package abstract_factory

import (
	handler2 "github.com/godsareinvented/go-metrics-collector/internal/buisness_logic/value_handler/handler"
	"github.com/godsareinvented/go-metrics-collector/internal/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/interfaces"
)

func GetValueHandler(metric dto.Metrics) interfaces.ValueHandlerInterface {
	switch metric.MType {
	case dictionary.GaugeMetricType:
		return &handler2.GaugeValueHandler{}
	case dictionary.CounterMetricType:
		return &handler2.CounterValueHandler{}
	default:
		panic("unknown metric type")
	}
}
