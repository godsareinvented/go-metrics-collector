package value_handler

import (
	"github.com/godsareinvented/go-metrics-collector/internal/general/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/general/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/server/buisness_logic/value_handler/handler"
	"github.com/godsareinvented/go-metrics-collector/internal/server/interfaces"
)

func GetValueHandler(metric dto.Metrics) interfaces.ValueHandlerInterface {
	switch metric.MType {
	case dictionary.GaugeMetricType:
		return &handler.GaugeValueHandler{}
	case dictionary.CounterMetricType:
		return &handler.CounterValueHandler{}
	default:
		panic("unknown metric type")
	}
}
