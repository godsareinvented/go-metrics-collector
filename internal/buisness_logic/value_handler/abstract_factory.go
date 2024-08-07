package value_handler

import (
	handler2 "github.com/oldhanasong/go-metrics-collector/internal/buisness_logic/value_handler/handler"
	"github.com/oldhanasong/go-metrics-collector/internal/dictionary"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/interfaces"
	"github.com/oldhanasong/go-metrics-collector/internal/repository"
)

func GetValueHandler(metric dto.Metric, repos *repository.Repository) interfaces.ValueHandler {
	switch metric.Type {
	case dictionary.GaugeMetricType:
		return &handler2.GaugeValueHandler{Repository: repos}
	case dictionary.CounterMetricType:
		return &handler2.CounterValueHandler{Repository: repos}
	default:
		panic("unknown metric type")
	}
}
