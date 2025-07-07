package value_handler

import (
	"fmt"
	"github.com/oldhanasong/go-metrics-collector/internal/general/business_logic/dictionary"
	"github.com/oldhanasong/go-metrics-collector/internal/general/dto"
	handler2 "github.com/oldhanasong/go-metrics-collector/internal/server/buisness_logic/metric/value_handler/handler"
	"github.com/oldhanasong/go-metrics-collector/internal/server/interfaces"
)

func ValueHandler(metric dto.Metrics) (interfaces.ValueHandler, error) {
	switch metric.MType {
	case dictionary.GaugeMetricType:
		return &handler2.GaugeValueHandler{}, nil
	case dictionary.CounterMetricType:
		return &handler2.CounterValueHandler{}, nil
	default:
		return nil, fmt.Errorf("unknown metric type: %s", metric.MType)
	}
}
