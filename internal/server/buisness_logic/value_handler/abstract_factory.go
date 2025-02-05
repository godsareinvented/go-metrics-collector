package value_handler

import (
	"github.com/godsareinvented/go-metrics-collector/internal/general/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/general/dto"
	handler2 "github.com/godsareinvented/go-metrics-collector/internal/server/buisness_logic/value_handler/handler"
	"github.com/godsareinvented/go-metrics-collector/internal/server/interfaces"
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
