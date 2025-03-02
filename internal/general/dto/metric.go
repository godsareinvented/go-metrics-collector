package dto

import (
	"fmt"
	"github.com/godsareinvented/go-metrics-collector/internal/general/dictionary"
	"strconv"
)

// Metrics todo: переписать валидацию. Как создавать группы констреинтов?
type Metrics struct {
	ID    string   `json:"id"               validate:"omitempty"`
	MType string   `json:"type"             validate:"required,oneof=gauge counter"`
	MName string   `json:"name"             validate:"required,alphanum"`
	Delta *int64   `json:"delta,omitempty"  validate:"omitempty"`
	Value *float64 `json:"value,omitempty"  validate:"omitempty"`
}

func (m Metrics) String() string {
	var valueString = m.GetFormattedValue()
	return fmt.Sprintf("%s\\%s\\%s: %s", m.MType, m.MName, m.ID, valueString)
}

func (m Metrics) GetFormattedValue() string {
	switch m.MType {
	case dictionary.GaugeMetricType:
		return strconv.FormatFloat(*m.Value, 'f', -1, 64)
	case dictionary.CounterMetricType:
		return strconv.FormatInt(*m.Delta, 10)
	default:
		panic("Unknown metric type")
	}
}
