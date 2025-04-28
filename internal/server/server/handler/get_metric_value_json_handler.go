package handler

import (
	"context"
	"encoding/json"
	"github.com/godsareinvented/go-metrics-collector/internal/server/config"
	"github.com/godsareinvented/go-metrics-collector/internal/server/service/metric/parser"
	"net/http"
)

type InputMetrics struct {
	ID    string `json:"id"   validate:"required"`
	MType string `json:"type" validate:"required,oneof=gauge counter"`
}

func GetMetricJson(ctx context.Context) http.HandlerFunc {
	fn := func(responseWriter http.ResponseWriter, request *http.Request) {
		requestCtx, cancel := GetCombinedContext(ctx, request.Context())
		defer cancel()

		requestParser := parser.JsonParser{}
		metric, err := requestParser.GetMetric(request)
		if nil != err {
			http.Error(responseWriter, err.Error(), http.StatusBadRequest)
			return
		}

		inputMetric := InputMetrics{
			ID:    metric.ID,
			MType: metric.MType,
		}
		err = config.Configuration.Validate.StructCtx(requestCtx, inputMetric)
		if nil != err {
			if isContextError(err) {
				http.Error(responseWriter, "", http.StatusInternalServerError)
				return
			}
			message, statusCode := processValidationError(err)
			http.Error(responseWriter, message, statusCode)
			return
		}

		resultingMetric, isSet, err := config.Configuration.Repository.GetMetricByID(requestCtx, metric)
		if nil != err {
			http.Error(responseWriter, "", http.StatusInternalServerError)
		}

		if !isSet {
			http.NotFound(responseWriter, request)
			return
		}

		metricJson, err := json.Marshal(resultingMetric)
		if nil != err {
			http.Error(responseWriter, err.Error(), http.StatusBadRequest)
			return
		}

		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusOK)
		responseWriter.Write(metricJson)
	}
	return fn
}
