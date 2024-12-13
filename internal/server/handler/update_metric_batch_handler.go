package handler

import (
	"context"
	"encoding/json"
	"fmt"
	manager "github.com/oldhanasong/go-metrics-collector/internal/service/metric"
	"net/http"
)

func UpdateMetricBatch(ctx context.Context) http.HandlerFunc {
	fn := func(responseWriter http.ResponseWriter, request *http.Request) {
		combinedCtx, cancel := combineContext(ctx, request.Context())
		defer cancel()

		metricList, err := parsedJsonMetrics(request)
		if err != nil {
			http.Error(responseWriter, "failed to get the metric list", http.StatusBadRequest)
			return
		}

		for _, m := range metricList {
			if err = v.Struct(m); err != nil {
				http.Error(responseWriter, fmt.Sprintf("incorrect metric data (%s)", m), http.StatusBadRequest)
				return
			}
		}

		metricManager := manager.MetricManager{}
		if err = metricManager.UpdateMetrics(combinedCtx, metricList); err != nil {
			http.Error(responseWriter, "failed to write metrics in the response", http.StatusInternalServerError)
			return
		}

		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(responseWriter).Encode(metricList); err != nil {
			http.Error(responseWriter, "failed to encode the metrics", http.StatusInternalServerError)
		}
	}
	return fn
}
