package handler

import (
	manager "github.com/oldhanasong/go-metrics-collector/internal/service/metric"
	"net/http"
)

func UpdateMetricJson(responseWriter http.ResponseWriter, request *http.Request) {
	m, err := parsedJsonMetric(request)
	if err != nil {
		http.Error(responseWriter, "failed to get the metric list", http.StatusBadRequest)
		return
	}

	err = v.Struct(m)
	if err != nil {
		http.Error(responseWriter, "incorrect metric data", http.StatusBadRequest)
		return
	}

	metricManager := manager.MetricManager{}
	err = metricManager.UpdateValue(m)
	if err != nil {
		http.Error(responseWriter, "failed to save the metric", http.StatusInternalServerError)
		return
	}

	responseWriter.WriteHeader(http.StatusOK)
}
