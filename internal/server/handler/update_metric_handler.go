package handler

import (
	"context"
	"github.com/oldhanasong/go-metrics-collector/internal/dictionary"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	manager "github.com/oldhanasong/go-metrics-collector/internal/service/metric"
	"github.com/oldhanasong/go-metrics-collector/internal/service/validator/metric"
	"net/http"
	"strconv"
)

func UpdateMetric(ctx context.Context) http.HandlerFunc {
	fn := func(responseWriter http.ResponseWriter, request *http.Request) {
		combinedCtx, cancel := combineContext(ctx, request.Context())
		defer cancel()

		MType, MName, MValue := parsedMetricValues(request)
		if MType == "" || MName == "" || MValue == "" {
			http.Error(responseWriter, "empty metric data", http.StatusBadRequest)
			return
		}

		err := metric.ValidateMetricValues(MType, MName, MValue)
		if err != nil {
			http.Error(responseWriter, "incorrect metric data", http.StatusBadRequest)
			return
		}

		m := dto.Metrics{ID: MName, MType: MType}
		if MType == dictionary.GaugeMetricType {
			floatVal, _ := strconv.ParseFloat(MValue, 64)
			m.Value = &floatVal
		} else {
			intVal, _ := strconv.ParseInt(MValue, 10, 64)
			m.Delta = &intVal
		}

		metricManager := manager.MetricManager{}
		err = metricManager.UpdateMetrics(combinedCtx, m)
		if err != nil {
			http.Error(responseWriter, "failed to save the metric", http.StatusInternalServerError)
			return
		}

		responseWriter.WriteHeader(http.StatusOK)
	}
	return fn
}
