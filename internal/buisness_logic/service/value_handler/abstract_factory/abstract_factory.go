package abstract_factory

import (
	"github.com/godsareinvented/go-metrics-collector/internal/buisness_logic/service/value_handler/handler"
	"github.com/godsareinvented/go-metrics-collector/internal/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/interfaces"
	"github.com/godsareinvented/go-metrics-collector/internal/repository"
)

func GetValueHandler(metric dto.Metric, repos repository.Repository) interfaces.ValueHandler {
	switch metric.Type {
	case dictionary.GaugeMetricType:
		return &handler.GaugeValueHandler{Repository: &repos}
	case dictionary.CounterMetricType:
		return &handler.CounterValueHandler{Repository: &repos}
	default:
		panic("unknown metric type")
	}
}
