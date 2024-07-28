package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/oldhanasong/go-metrics-collector/internal/buisness_logic/manager"
	"github.com/oldhanasong/go-metrics-collector/internal/constraint"
	"github.com/oldhanasong/go-metrics-collector/internal/dictionary"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/repository"
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

	switch chi.URLParam(request, "type") {
	case dictionary.GaugeMetricType:
		getMetric[float64](responseWriter, request, MType, MName)
	case dictionary.CounterMetricType:
		getMetric[int64](responseWriter, request, MType, MName)
	default:
		http.Error(responseWriter, "invalid metric type", http.StatusBadRequest)
	}
}

func getMetric[Num constraint.Numeric](responseWriter http.ResponseWriter, request *http.Request, MType, MName string) {
	m := dto.Metric[Num]{Type: MType, Name: MName}

	metricManager := manager.MetricManager[Num]{Repository: repository.GetInstance[Num](m.Type)}
	resultingMetric, isSet := metricManager.Get(m)

	if !isSet {
		http.NotFound(responseWriter, request)
	}

	preparedMetricValue := value_formatter.GetFormattedValue(resultingMetric)
	responseWriter.Write([]byte(preparedMetricValue))
}
