package handler

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/godsareinvented/go-metrics-collector/internal/service/metric/manager"
	"github.com/godsareinvented/go-metrics-collector/internal/service/metric/parser"
	"net/http"
)

func UpdateMetricBatchMetric(ctx context.Context) http.HandlerFunc {
	fn := func(responseWriter http.ResponseWriter, request *http.Request) {
		// Комбинированный контекст, чтобы хендлер мог обработать завершение контекстов как приложения, так и запроса
		requestCtx, cancel := context.WithCancel(request.Context())
		defer cancel()

		go func() {
			<-ctx.Done()
			cancel()
		}()

		jsonParser := parser.JsonParser{}
		metricBatch, err := jsonParser.GetMetricBatch(request)
		if nil != err {
			http.Error(responseWriter, err.Error(), http.StatusBadRequest)
			return
		}

		// Валидация корректности данных метрики, инъекций.
		for _, metric := range metricBatch {
			err = validator.New(validator.WithRequiredStructEnabled()).Struct(metric)
			if nil != err {
				message, statusCode := ProcessValidationError(err)
				http.Error(responseWriter, message, statusCode)
				return
			}
		}

		if nil != metricBatch {
			metricManager := manager.MetricManager{}
			metricManager.UpdateMetrics(requestCtx, metricBatch)
		}

		responseWriter.WriteHeader(http.StatusOK)
	}
	return fn
}
