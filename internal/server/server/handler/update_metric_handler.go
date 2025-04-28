package handler

import (
	"context"
	"github.com/godsareinvented/go-metrics-collector/internal/server/config"
	metricPackage "github.com/godsareinvented/go-metrics-collector/internal/server/service/metric"
	"github.com/godsareinvented/go-metrics-collector/internal/server/service/metric/parser"
	"net/http"
)

func UpdateMetric(ctx context.Context) http.HandlerFunc {
	fn := func(responseWriter http.ResponseWriter, request *http.Request) {
		requestCtx, cancel := GetCombinedContext(ctx, request.Context())
		defer cancel()

		requestParser := parser.RequestParser{}
		metric, err := requestParser.GetMetric(request, true)
		if nil != err {
			http.Error(responseWriter, err.Error(), http.StatusBadRequest)
			return
		}

		// Валидация корректности данных метрики, инъекций.
		err = config.Configuration.Validate.StructCtx(requestCtx, metric)
		if nil != err {
			if isContextError(err) {
				http.Error(responseWriter, "", http.StatusInternalServerError)
				return
			}
			message, statusCode := processValidationError(err)
			http.Error(responseWriter, message, statusCode)
			return
		}

		metricManager := metricPackage.MetricManager{}
		err = metricManager.UpdateMetric(requestCtx, metric)
		if nil != err {
			http.Error(responseWriter, "", http.StatusInternalServerError)
			return
		}

		responseWriter.WriteHeader(http.StatusOK)
	}
	return fn
}
