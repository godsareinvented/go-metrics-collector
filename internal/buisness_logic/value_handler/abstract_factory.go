package value_handler

import (
	"fmt"
	"github.com/oldhanasong/go-metrics-collector/internal/buisness_logic/value_handler/handler"
	"github.com/oldhanasong/go-metrics-collector/internal/dictionary"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/interfaces"
)

func GetValueHandler(metric dto.Metrics) (interfaces.ValueHandler, error) {
	switch metric.MType {
	case dictionary.GaugeMetricType:
		return &handler.GaugeValueHandler{}, nil
	case dictionary.CounterMetricType:
		return &handler.CounterValueHandler{}, nil
	default:
		return nil, fmt.Errorf("unknown metric type: %s", metric.MType)
	}
}
