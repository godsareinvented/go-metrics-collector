package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/godsareinvented/go-metrics-collector/internal/buisness_logic/manager"
	"github.com/godsareinvented/go-metrics-collector/internal/service/metric/parser"
	"github.com/godsareinvented/go-metrics-collector/internal/service/metric/value_formatter"
	"net/http"
)

func GetMetric(responseWriter http.ResponseWriter, request *http.Request) {
	requestParser := parser.RequestParser{}
	metricDTO, err := requestParser.GetMetricDTO(request, false)
	if nil != err {
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	err = validator.New().Struct(metricDTO)
	if nil != err {
		message, statusCode := ProcessValidationError(err)
		http.Error(responseWriter, message, statusCode)
		return
	}

	metricManager := manager.MetricManager{}
	resultingMetric, isSet := metricManager.Get(metricDTO)

	if isSet {
		preparedMetricValue := value_formatter.GetFormattedValue(resultingMetric)
		responseWriter.Write([]byte(preparedMetricValue))
		return
	}

	http.NotFound(responseWriter, request)
}
