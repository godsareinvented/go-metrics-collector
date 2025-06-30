package handler

import (
	"context"
	"github.com/oldhanasong/go-metrics-collector/internal/general/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/general/util"
	"github.com/oldhanasong/go-metrics-collector/internal/server/config"
	"github.com/oldhanasong/go-metrics-collector/internal/server/service/validator/metric"
	"net/http"
)

func GetMetric(ctx context.Context) http.HandlerFunc {
	fn := func(responseWriter http.ResponseWriter, request *http.Request) {
		ctx, cancel := util.CombineContexts(ctx, request.Context())
		defer cancel()

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

		resultingMetric, isSet, err := config.Configuration.Repository.GetMetric(ctx, m)
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
		if _, err = responseWriter.Write([]byte(preparedMetricValue)); err != nil {
			http.Error(responseWriter, "body record error", http.StatusInternalServerError)
		}
	}
	return fn
}
