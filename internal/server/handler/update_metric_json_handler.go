package handler

import (
	"encoding/json"
	manager "github.com/oldhanasong/go-metrics-collector/internal/service/metric"
	"net/http"
)

func UpdateMetricJson(responseWriter http.ResponseWriter, request *http.Request) {
	m, err := parsedJsonMetric(request)
	if err != nil {
		http.Error(responseWriter, "failed to get the metric list", http.StatusBadRequest)
		return
	}

	if err = v.Struct(m); err != nil {
		http.Error(responseWriter, "incorrect metric data", http.StatusBadRequest)
		return
	}

	metricManager := manager.MetricManager{}
	if err = metricManager.UpdateMetrics(m); err != nil {
		http.Error(responseWriter, "failed to write metric in the response", http.StatusInternalServerError)
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(responseWriter).Encode(m); err != nil {
		http.Error(responseWriter, "failed to encode the metric", http.StatusInternalServerError)
	}
}
