package handler

import (
	"encoding/json"
	manager "github.com/oldhanasong/go-metrics-collector/internal/service/metric"
	"net/http"
)

func GetMetricJson(responseWriter http.ResponseWriter, request *http.Request) {
	m, err := parsedJsonMetric(request)
	if err != nil {
		http.Error(responseWriter, "failed to get the metric list", http.StatusBadRequest)
		return
	}

	err = v.StructPartial(m, "ID", "MType")
	if err != nil {
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
	_, err = responseWriter.Write(metricJson)
	if err != nil {
		http.Error(responseWriter, "body record error", http.StatusInternalServerError)
	}
}
