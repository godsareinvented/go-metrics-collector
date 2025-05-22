package handler

import (
	"encoding/json"
	manager "github.com/oldhanasong/go-metrics-collector/internal/service/metric"
	"net/http"
)

func GetMetricJson(responseWriter http.ResponseWriter, request *http.Request) {
	m, err := parsedJsonMetric(request)
	if err != nil {
		http.Error(responseWriter, "failed to decode request body", http.StatusBadRequest)
		return
	}

	if err = v.StructPartial(m, "ID", "MType"); err != nil {
		http.Error(responseWriter, "incorrect metric data", http.StatusBadRequest)
		return
	}

	metricManager := manager.MetricManager{}
	resultingMetric, isSet, err := metricManager.Get(m)
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
