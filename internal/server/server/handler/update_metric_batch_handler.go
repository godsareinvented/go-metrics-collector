package handler

import (
	"context"
	"github.com/godsareinvented/go-metrics-collector/internal/server/config"
	metricPackage "github.com/godsareinvented/go-metrics-collector/internal/server/service/metric"
	"github.com/godsareinvented/go-metrics-collector/internal/server/service/metric/parser"
	"net/http"
)

func UpdateMetricBatchMetric(ctx context.Context) http.HandlerFunc {
	fn := func(responseWriter http.ResponseWriter, request *http.Request) {
		requestCtx, cancel := GetCombinedContext(ctx, request.Context())
		defer cancel()

		jsonParser := parser.JsonParser{}
		metricBatch, err := jsonParser.GetMetricBatch(request)
		if nil != err {
			http.Error(responseWriter, err.Error(), http.StatusBadRequest)
			return
		}

		// Валидация корректности данных метрик.
		for _, metric := range metricBatch {
			err = config.Configuration.Validate.Struct(metric)
			if nil != err {
				message, statusCode := ProcessValidationError(err)
				http.Error(responseWriter, message, statusCode)
				return
			}
		}

		if nil != metricBatch {
			metricManager := metricPackage.MetricManager{}
			metricManager.UpdateMetrics(requestCtx, metricBatch)
		}

		responseWriter.WriteHeader(http.StatusOK)
	}
	return fn
}
