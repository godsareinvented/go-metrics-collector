package value_handler

import (
	"fmt"
	"github.com/oldhanasong/go-metrics-collector/internal/buisness_logic/value_handler/handler"
	"github.com/oldhanasong/go-metrics-collector/internal/dictionary"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/interfaces"
	"github.com/oldhanasong/go-metrics-collector/internal/repository"
)

func GetValueHandler(metric dto.Metrics, repos *repository.Repository) (interfaces.ValueHandler, error) {
	switch metric.MType {
	case dictionary.GaugeMetricType:
		return &handler.GaugeValueHandler{Repository: repos}, nil
	case dictionary.CounterMetricType:
		return &handler.CounterValueHandler{Repository: repos}, nil
	default:
		return nil, fmt.Errorf("unknown metric type: %s", metric.MType)
	}
}
