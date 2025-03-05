package handler

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/godsareinvented/go-metrics-collector/internal/server/config"
	"github.com/godsareinvented/go-metrics-collector/internal/server/service/metric/parser"
	"net/http"
)

func GetMetric(ctx context.Context) http.HandlerFunc {
	fn := func(responseWriter http.ResponseWriter, request *http.Request) {
		requestCtx, cancel := GetCombinedContext(ctx, request.Context())
		defer cancel()

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

		resultingMetric, isSet, _ := config.Configuration.Repository.GetMetricByName(requestCtx, metricDTO)

		if isSet {
			preparedMetricValue := resultingMetric.GetFormattedValue()
			responseWriter.WriteHeader(http.StatusOK)
			_, _ = responseWriter.Write([]byte(preparedMetricValue))
			return
		}

		http.NotFound(responseWriter, request)
	}
	return fn
}
