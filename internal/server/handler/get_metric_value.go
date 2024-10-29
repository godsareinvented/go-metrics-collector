package handler

import (
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	manager "github.com/oldhanasong/go-metrics-collector/internal/service/metric"
	"github.com/oldhanasong/go-metrics-collector/internal/service/validator/metric"
	"net/http"
)

func GetMetric(responseWriter http.ResponseWriter, request *http.Request) {
	MType, MName := parsedAbridgedMetricValues(request)
	if MType == "" || MName == "" {
		http.Error(responseWriter, "empty metric data", http.StatusBadRequest)
		return
	}

	err := metric.ValidateAbridgedMetricValues(MType, MName)
	if err != nil {
		http.Error(responseWriter, "incorrect metric data", http.StatusBadRequest)
		return
	}

	m := dto.Metrics{ID: MName, MType: MType}

	metricManager := manager.MetricManager{}
	resultingMetric, isSet := metricManager.Get(m)

	if !isSet {
		http.NotFound(responseWriter, request)
	}

	responseWriter.WriteHeader(http.StatusOK)
	preparedMetricValue := resultingMetric.FormattedValue()
	responseWriter.Write([]byte(preparedMetricValue))
}
