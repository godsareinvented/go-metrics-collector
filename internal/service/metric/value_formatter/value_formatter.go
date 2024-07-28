package value_formatter

import (
	"github.com/godsareinvented/go-metrics-collector/internal/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
	"strconv"
)

func GetFormattedValue(metricDTO dto.Metric) string {
	switch metricDTO.Type {
	case dictionary.GaugeMetricType:
		return strconv.FormatFloat(metricDTO.Value, 'f', -1, 64)
	case dictionary.CounterMetricType:
		return strconv.FormatInt(metricDTO.Delta, 10)
	default:
		panic("Unknown metric type")
	}
}
