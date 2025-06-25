package value_handler

import (
	"fmt"
	"github.com/oldhanasong/go-metrics-collector/internal/general/dictionary"
	"github.com/oldhanasong/go-metrics-collector/internal/general/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/server/buisness_logic/value_handler/handler"
	"github.com/oldhanasong/go-metrics-collector/internal/server/interfaces"
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
