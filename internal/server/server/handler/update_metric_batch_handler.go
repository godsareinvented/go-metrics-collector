package handler

import (
	"context"
	"github.com/go-playground/validator/v10"
	metricPackage "github.com/godsareinvented/go-metrics-collector/internal/server/service/metric"
	"github.com/godsareinvented/go-metrics-collector/internal/server/service/metric/parser"
	"net/http"
)

func UpdateMetricBatchMetric(ctx context.Context) http.HandlerFunc {
	fn := func(responseWriter http.ResponseWriter, request *http.Request) {
		// Комбинированный контекст, чтобы хендлер мог обработать завершение контекстов как приложения, так и запроса
		requestCtx, cancel := context.WithCancel(request.Context())
		defer cancel()

		// todo: Утечка памяти? Горутина будет вечно ожидать закрытия канала. Горутина завершиться при завершении main или при завершении и UpdateMetricBatchMetric тоже?
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

		// todo: Усилить валидацию сущностей метрик. Проверять, что задано нужное значение метрики, соответствующее типу.
		// Валидация корректности данных метрики, инъекций.
		v := validator.New(validator.WithRequiredStructEnabled())
		for _, metric := range metricBatch {
			err = v.Struct(metric)
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
