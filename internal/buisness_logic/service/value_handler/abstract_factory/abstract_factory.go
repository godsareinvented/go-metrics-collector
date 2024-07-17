package abstract_factory

import (
	handler2 "github.com/godsareinvented/go-metrics-collector/internal/buisness_logic/service/value_handler/handler"
	"github.com/godsareinvented/go-metrics-collector/internal/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/interfaces"
)

func GetValueHandler(metric dto.Metric) interfaces.ValueHandler {
	switch metric.Type {
	case dictionary.GaugeMetricType:
		return &handler2.GaugeValuePreprocessor{}
	case dictionary.CounterMetricType:
		return &handler2.CounterValuePreprocessor{}
	default:
		panic("unknown metric type")
	}
}
