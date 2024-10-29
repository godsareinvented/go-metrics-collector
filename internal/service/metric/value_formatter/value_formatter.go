package value_formatter

import (
	"github.com/oldhanasong/go-metrics-collector/internal/dictionary"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"strconv"
)

func GetFormattedValue(metric dto.Metrics) string {
	switch metric.MType {
	case dictionary.GaugeMetricType:
		return strconv.FormatFloat(*metric.Value, 'f', -1, 64)
	case dictionary.CounterMetricType:
		return strconv.FormatInt(*metric.Delta, 10)
	default:
		panic("Unknown metric type")
	}
}
