package handler

import (
	"context"
	"github.com/oldhanasong/go-metrics-collector/internal/config"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/service/validator/metric"
	"net/http"
)

func GetMetric(_ context.Context) http.HandlerFunc {
	fn := func(responseWriter http.ResponseWriter, request *http.Request) {
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

		resultingMetric, isSet, err := config.Configuration.Repository.GetMetric(m)
		if err != nil {
			http.Error(responseWriter, "failed to get the metric list", http.StatusInternalServerError)
			return
		}

		if !isSet {
			http.NotFound(responseWriter, request)
			return
		}

		responseWriter.WriteHeader(http.StatusOK)
		preparedMetricValue := resultingMetric.FormattedValue()
		_, err = responseWriter.Write([]byte(preparedMetricValue))
		if err != nil {
			http.Error(responseWriter, "body record error", http.StatusInternalServerError)
		}
	}
	return fn
}
