package handler

import (
	"github.com/oldhanasong/go-metrics-collector/internal/buisness_logic/manager"
	"github.com/oldhanasong/go-metrics-collector/internal/dictionary"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/service/validator/metric"
	"net/http"
	"strconv"
)

type UpdateMetricHandler struct{}

func (handler *UpdateMetricHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(responseWriter, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	MType, MName, MValue := parsedMetricValues(request)
	if MType == "" || MName == "" || MValue == "" {
		http.Error(responseWriter, "empty metric data", http.StatusNotFound)
		return
	}

	err := metric.ValidateMetricValues(MType, MName, MValue)
	if err != nil {
		http.Error(responseWriter, "incorrect metric data", http.StatusBadRequest)
		return
	}

	metrics := dto.Metric{Type: MType, Name: MName}
	if MType == dictionary.GaugeMetricType {
		metrics.Value, _ = strconv.ParseFloat(MValue, 64)
	} else {
		metrics.Value, _ = strconv.ParseInt(MValue, 10, 64)
	}

	metricManager := manager.MetricManager{}
	metricManager.UpdateValue(metrics)
}

func parsedMetricValues(r *http.Request) (string, string, string) {
	return r.PathValue("type"),
		r.PathValue("name"),
		r.PathValue("value")
}
