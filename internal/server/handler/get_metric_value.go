package handler

import (
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"github.com/godsareinvented/go-metrics-collector/internal/buisness_logic/manager"
	"github.com/godsareinvented/go-metrics-collector/internal/constraint"
	"github.com/godsareinvented/go-metrics-collector/internal/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/service/metric/parser"
	"github.com/godsareinvented/go-metrics-collector/internal/service/metric/value_formatter"
	"net/http"
	"strconv"
)

func GetMetric(responseWriter http.ResponseWriter, request *http.Request) {
	if chi.URLParam(request, "type") == dictionary.GaugeMetricType {
		handleGettingMetric[float64](responseWriter, request)
		return
	}
	handleGettingMetric[int64](responseWriter, request)
}

func handleGettingMetric[Num constraint.Numeric](responseWriter http.ResponseWriter, request *http.Request) {
	requestParser := parser.RequestParser[Num]{}
	metricDTO := requestParser.GetMetricDTO(request)

	metricDTO.Type = strconv.Quote(metricDTO.Type)
	metricDTO.Name = strconv.Quote(metricDTO.Name)

	err := validator.New().Struct(metricDTO)

	if nil != err {
		message, statusCode := ProcessValidationError(err)
		http.Error(responseWriter, message, statusCode)
		return
	}

	metricManager := manager.MetricManager[Num]{}
	resultingMetric, isSet := metricManager.Get(metricDTO)

	if isSet {
		preparedMetricValue := value_formatter.GetFormattedValue(resultingMetric)
		responseWriter.Write([]byte(preparedMetricValue))
		return
	}

	http.NotFound(responseWriter, request)
}
