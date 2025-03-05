package dto

import (
	"fmt"
	"github.com/godsareinvented/go-metrics-collector/internal/general/dictionary"
	"strconv"
)

// Metrics todo: При добавлении новых типов метрик необходимо обновлять валидацию сущности.
type Metrics struct {
	ID    string   `json:"id"               validate:"omitempty"`
	MType string   `json:"type"             validate:"required,oneof=gauge counter"`
	MName string   `json:"name"             validate:"required,alphanum"`
	Delta *int64   `json:"delta,omitempty"  validate:"required_if=MType counter,excluded_if=MType gauge"`
	Value *float64 `json:"value,omitempty"  validate:"required_if=MType gauge,excluded_if=MType counter"`
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
