package handler

import (
	"context"
	"encoding/json"
	"github.com/oldhanasong/go-metrics-collector/internal/general/util"
	"github.com/oldhanasong/go-metrics-collector/internal/server/config"
	"net/http"
)

func GetMetricJson(ctx context.Context) http.HandlerFunc {
	fn := func(responseWriter http.ResponseWriter, request *http.Request) {
		ctx, cancel := util.CombineContexts(ctx, request.Context())
		defer cancel()

		m, err := parsedJsonMetric(request)
		if err != nil {
			http.Error(responseWriter, "failed to decode request body", http.StatusBadRequest)
			return
		}

		if err = v.StructPartial(m, "ID", "MType"); err != nil {
			http.Error(responseWriter, "incorrect metric data", http.StatusBadRequest)
			return
		}

		resultingMetric, isSet, err := config.Configuration.Repository.GetMetric(ctx, m)
		if err != nil {
			http.Error(responseWriter, "failed to get the metric list", http.StatusInternalServerError)
			return
		}

		if !isSet {
			http.Error(responseWriter, "metric not found", http.StatusNotFound)
			return
		}

		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(responseWriter).Encode(resultingMetric); err != nil {
			http.Error(responseWriter, "failed to encode the metric", http.StatusInternalServerError)
		}
	}
	return fn
}
