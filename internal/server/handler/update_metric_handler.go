package handler

import (
	"github.com/oldhanasong/go-metrics-collector/internal/buisness_logic/manager"
	"github.com/oldhanasong/go-metrics-collector/internal/dictionary"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/repository"
	"github.com/oldhanasong/go-metrics-collector/internal/service/validator/metric"
	"net/http"
	"strconv"
)

func UpdateMetric(responseWriter http.ResponseWriter, request *http.Request) {
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

	m := dto.Metric{Type: MType, Name: MName}
	if MType == dictionary.GaugeMetricType {
		m.Value, _ = strconv.ParseFloat(MValue, 64)
	} else {
		m.Delta, _ = strconv.ParseInt(MValue, 10, 64)
	}

	metricManager := manager.MetricManager{Repository: repository.GetInstance()}
	metricManager.UpdateValue(m)
}
