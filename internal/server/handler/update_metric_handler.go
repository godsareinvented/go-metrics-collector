package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/godsareinvented/go-metrics-collector/internal/buisness_logic/action"
	"github.com/godsareinvented/go-metrics-collector/internal/service/metric/parser"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type UpdateMetricHandler struct{}

func (handler *UpdateMetricHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(responseWriter, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	requestParser := parser.RequestParser{}
	metricDTO := requestParser.GetParsedMetric(request)

	err := validator.New().Struct(metricDTO)

	if nil != err {
		message, statusCode := processError(err)
		http.Error(responseWriter, message, statusCode)
		return
	}

	if reflect.TypeOf(metricDTO.Value).Name() == "int64" {
		metricDTO.Value, _ = strconv.ParseInt(metricDTO.Value.(string), 10, 64)
	} else {
		metricDTO.Value, _ = strconv.ParseFloat(metricDTO.Value.(string), 64)
	}

	action.UpdateValue(metricDTO)
}

func processError(error error) (string, int) {
	errors := error.(validator.ValidationErrors)
	if strings.Contains(errors[0].Field(), "Name") {
		return "Metric name not passed on or incorrect", http.StatusNotFound
	}
	return "incorrect metric data", http.StatusBadRequest
}
