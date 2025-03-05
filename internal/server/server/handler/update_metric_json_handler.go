package handler

import (
	"context"
	"github.com/godsareinvented/go-metrics-collector/internal/server/config"
	metricPackage "github.com/godsareinvented/go-metrics-collector/internal/server/service/metric"
	"github.com/godsareinvented/go-metrics-collector/internal/server/service/metric/parser"
	"net/http"
)

func UpdateMetricJson(ctx context.Context) http.HandlerFunc {
	fn := func(responseWriter http.ResponseWriter, request *http.Request) {
		requestCtx, cancel := GetCombinedContext(ctx, request.Context())
		defer cancel()

		requestParser := parser.JsonParser{}
		metric, err := requestParser.GetMetricDTO(request)
		if nil != err {
			http.Error(responseWriter, err.Error(), http.StatusBadRequest)
			return
		}

		// Валидация корректности данных метрики, инъекций.
		err = config.Configuration.Validate.Struct(metric)
		if nil != err {
			message, statusCode := ProcessValidationError(err)
			http.Error(responseWriter, message, statusCode)
			return
		}

		metricManager := metricPackage.MetricManager{}
		metricManager.UpdateMetric(requestCtx, metric)

		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusOK)
	}
	return fn
}
