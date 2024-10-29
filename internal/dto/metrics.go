package dto

import (
	"fmt"
	"github.com/oldhanasong/go-metrics-collector/internal/dictionary"
	"strconv"
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (m Metrics) String() string {
	var valueString = m.FormattedValue()
	return fmt.Sprintf("%s\\%s: %s", m.MType, m.ID, valueString)
}

func (m Metrics) FormattedValue() string {
	switch m.MType {
	case dictionary.GaugeMetricType:
		return strconv.FormatFloat(*m.Value, 'f', -1, 64)
	case dictionary.CounterMetricType:
		return strconv.FormatInt(*m.Delta, 10)
	default:
		panic("Unknown metric type")
	}
}
