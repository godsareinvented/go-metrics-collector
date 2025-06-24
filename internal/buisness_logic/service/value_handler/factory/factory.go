package factory

import (
	"github.com/oldhanasong/go-metrics-collector/internal/buisness_logic/service/value_handler/handler"
	"github.com/oldhanasong/go-metrics-collector/internal/dictionary"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/interfaces"
)

func GetValueHandler(metric dto.Metric) interfaces.ValueHandler {
	switch metric.Type {
	case dictionary.GaugeMetricType:
		return &handler.GaugeValuePreprocessor{}
	case dictionary.CounterMetricType:
		return &handler.CounterValuePreprocessor{}
	default:
		panic("unknown metric type")
	}
}
