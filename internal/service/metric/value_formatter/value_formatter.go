package value_formatter

import (
	"github.com/godsareinvented/go-metrics-collector/internal/constraint"
	"github.com/godsareinvented/go-metrics-collector/internal/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
	"strconv"
)

func GetFormattedValue[Num constraint.Numeric](metricDTO dto.Metric[Num]) string {
	switch metricDTO.Type {
	case dictionary.GaugeMetricType:
		return strconv.FormatFloat(float64(metricDTO.Value), 'f', -1, 64)
	case dictionary.CounterMetricType:
		return strconv.FormatInt(int64(metricDTO.Value), 10)
	default:
		panic("Unknown metric type")
	}
}
