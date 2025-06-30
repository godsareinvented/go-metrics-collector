package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/oldhanasong/go-metrics-collector/internal/general/util"
	manager "github.com/oldhanasong/go-metrics-collector/internal/server/service/metric"
	"net/http"
)

func UpdateMetricBatch(ctx context.Context) http.HandlerFunc {
	fn := func(responseWriter http.ResponseWriter, request *http.Request) {
		ctx, cancel := util.CombineContexts(ctx, request.Context())
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
		resMetrics, err := metricManager.UpdateMetrics(ctx, metricList)
		if err != nil {
			http.Error(responseWriter, "failed to update metrics", http.StatusInternalServerError)
			return
		}

		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(responseWriter).Encode(resMetrics); err != nil {
			http.Error(responseWriter, "failed to encode the metrics", http.StatusInternalServerError)
		}
	}
	return fn
}
