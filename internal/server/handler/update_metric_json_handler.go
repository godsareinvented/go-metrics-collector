package handler

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/godsareinvented/go-metrics-collector/internal/service/metric/manager"
	"github.com/godsareinvented/go-metrics-collector/internal/service/metric/parser"
	"net/http"
)

func UpdateMetricJson(_ context.Context) http.HandlerFunc {
	fn := func(responseWriter http.ResponseWriter, request *http.Request) {
		requestParser := parser.JsonParser{}
		metricDTO, err := requestParser.GetMetricDTO(request)
		if nil != err {
			http.Error(responseWriter, err.Error(), http.StatusBadRequest)
			return
		}

		// Валидация корректности данных метрики, инъекций.
		err = validator.New(validator.WithRequiredStructEnabled()).Struct(metricDTO)
		if nil != err {
			message, statusCode := ProcessValidationError(err)
			http.Error(responseWriter, message, statusCode)
			return
		}

		metricManager := manager.MetricManager{}
		metricManager.UpdateMetric(metricDTO)

		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusOK)
	}
	return fn
}
