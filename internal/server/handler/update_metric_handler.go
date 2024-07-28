package handler

import (
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"github.com/godsareinvented/go-metrics-collector/internal/buisness_logic/manager"
	"github.com/godsareinvented/go-metrics-collector/internal/constraint"
	"github.com/godsareinvented/go-metrics-collector/internal/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/service/metric/parser"
	"html"
	"net/http"
)

func UpdateMetric(responseWriter http.ResponseWriter, request *http.Request) {
	if chi.URLParam(request, "type") == dictionary.GaugeMetricType {
		handleMetricUpdates[float64](responseWriter, request)
		return
	}
	handleMetricUpdates[int64](responseWriter, request)
}

func handleMetricUpdates[Num constraint.Numeric](responseWriter http.ResponseWriter, request *http.Request) {
	requestParser := parser.RequestParser[Num]{}
	metricDTO := requestParser.GetMetricDTO(request)

	// todo: Защита точно нужна? chi не находит хэндлер со спец-символами в URL.
	metricDTO.Type = html.EscapeString(metricDTO.Type)
	metricDTO.Name = html.EscapeString(metricDTO.Name)

	err := validator.New().Struct(metricDTO)

	if nil != err {
		message, statusCode := ProcessValidationError(err)
		http.Error(responseWriter, message, statusCode)
		return
	}

	metricManager := manager.MetricManager[Num]{}
	metricManager.UpdateValue(metricDTO)
}
