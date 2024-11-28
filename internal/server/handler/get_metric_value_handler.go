package handler

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/godsareinvented/go-metrics-collector/internal/config"
	"github.com/godsareinvented/go-metrics-collector/internal/service/metric/parser"
	"net/http"
)

func GetMetric(_ context.Context) http.HandlerFunc {
	fn := func(responseWriter http.ResponseWriter, request *http.Request) {
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

		resultingMetric, isSet, _ := config.Configuration.Repository.GetMetricByName(metricDTO)

		if isSet {
			preparedMetricValue := resultingMetric.GetFormattedValue()
			responseWriter.Write([]byte(preparedMetricValue))
			responseWriter.WriteHeader(http.StatusOK)
			return
		}

		http.NotFound(responseWriter, request)
	}
	return fn
}
