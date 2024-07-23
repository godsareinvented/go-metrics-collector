package parser

import (
	"github.com/go-chi/chi"
	"github.com/godsareinvented/go-metrics-collector/internal/constraint"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
	"net/http"
	"strconv"
)

type RequestParser[Num constraint.Numeric] struct{}

func (rp *RequestParser[Num]) GetMetricDTO(request *http.Request) dto.Metric[Num] {
	metricType := chi.URLParam(request, "type")
	metricName := chi.URLParam(request, "name")
	metricValue := chi.URLParam(request, "value")

	return dto.Metric[Num]{
		Type:  metricType,
		Name:  metricName,
		Value: getParsedValue[Num](metricValue),
	}
}

func getParsedValue[Num constraint.Numeric](metricValue string) Num {
	intVal, err := strconv.ParseInt(metricValue, 10, 64)
	if nil == err {
		return Num(intVal)
	}

	floatVal, _ := strconv.ParseFloat(metricValue, 64)
	return Num(floatVal)
}
