package handler

import (
	"github.com/oldhanasong/go-metrics-collector/internal/buisness_logic/manager"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/service/metric/value_formatter"
	"github.com/oldhanasong/go-metrics-collector/internal/service/validator/metric"
	"net/http"
)

func GetMetric(responseWriter http.ResponseWriter, request *http.Request) {
	MType, MName := parsedAbridgedMetricValues(request)
	if MType == "" || MName == "" {
		http.Error(responseWriter, "empty metric data", http.StatusNotFound)
		return
	}

	err := metric.ValidateAbridgedMetricValues(MType, MName)
	if err != nil {
		http.Error(responseWriter, "incorrect metric data", http.StatusBadRequest)
		return
	}

	m := dto.Metric{Type: MType, Name: MName}

	metricManager := manager.MetricManager{}
	resultingMetric, isSet := metricManager.Get(m)

	if !isSet {
		http.NotFound(responseWriter, request)
	}

	preparedMetricValue := value_formatter.GetFormattedValue(resultingMetric)
	responseWriter.Write([]byte(preparedMetricValue))
}
