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

	switch MType {
	case dictionary.GaugeMetricType:
		m := dto.Metric[float64]{Type: MType, Name: MName}
		m.Value, _ = strconv.ParseFloat(MValue, 64)

		metricManager := manager.MetricManager[float64]{Repository: repository.GetInstance[float64](m.Type)}
		metricManager.UpdateValue(m)
	case dictionary.CounterMetricType:
		m := dto.Metric[int64]{Type: MType, Name: MName}
		m.Value, _ = strconv.ParseInt(MValue, 10, 64)

		metricManager := manager.MetricManager[int64]{Repository: repository.GetInstance[int64](m.Type)}
		metricManager.UpdateValue(m)
	default:
		http.Error(responseWriter, "invalid metric type", http.StatusBadRequest)
	}
}
