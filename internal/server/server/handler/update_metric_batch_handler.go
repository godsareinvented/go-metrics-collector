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
		// metricBatch todo: Добавить валидацию на уникальность метрик (unique=field)?
		for _, metric := range metricBatch {
			err = config.Configuration.Validate.StructCtx(requestCtx, metric)
			if nil == err {
				continue
			}
			if isContextError(err) {
				http.Error(responseWriter, "", http.StatusInternalServerError)
				return
			}

			message, statusCode := processValidationError(err)
			http.Error(responseWriter, message, statusCode)
			return
		}

		if nil != metricBatch {
			metricManager := metricPackage.MetricManager{}
			err = metricManager.UpdateMetrics(requestCtx, metricBatch)
			if nil != err {
				http.Error(responseWriter, "", http.StatusInternalServerError)
				return
			}
		}

		responseWriter.WriteHeader(http.StatusOK)
	}
	return fn
}
