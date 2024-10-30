package dto

import (
	"fmt"
	"github.com/godsareinvented/go-metrics-collector/internal/dictionary"
	"strconv"
)

// Metrics todo: переписать валидацию. Как создавать группы констреинтов?
type Metrics struct {
	ID    string   `json:"id"               validate:"omitempty,required"`
	MType string   `json:"type"             validate:"required,contains=gauge|contains=counter"`
	MName string   `json:"name"             validate:"required,alpha"`
	Delta *int64   `json:"delta,omitempty"  validate:"omitempty,required"`
	Value *float64 `json:"value,omitempty"  validate:"omitempty,required"`
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
