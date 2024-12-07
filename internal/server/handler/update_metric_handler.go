package handler

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/godsareinvented/go-metrics-collector/internal/service/metric/manager"
	"github.com/godsareinvented/go-metrics-collector/internal/service/metric/parser"
	"net/http"
)

func UpdateMetric(ctx context.Context) http.HandlerFunc {
	fn := func(responseWriter http.ResponseWriter, request *http.Request) {
		// Комбинированный контекст, чтобы хендлер мог обработать завершение контекстов как приложения, так и запроса
		requestCtx, cancel := context.WithCancel(request.Context())
		defer cancel()

		go func() {
			<-ctx.Done()
			cancel()
		}()

		requestParser := parser.RequestParser{}
		metricDTO, err := requestParser.GetMetricDTO(request, true)
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
		metricManager.UpdateMetric(requestCtx, metricDTO)

		responseWriter.WriteHeader(http.StatusOK)
	}
	return fn
}
