package handler

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/oldhanasong/go-metrics-collector/internal/general/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/general/validation/decorator"
	"github.com/oldhanasong/go-metrics-collector/internal/server/buisness_logic/config"
	"github.com/oldhanasong/go-metrics-collector/internal/server/repository"
	"github.com/oldhanasong/go-metrics-collector/internal/server/storage/mem_storage"
	"go.uber.org/multierr"
	"io"
	"net/http"
	"net/http/httptest"
)

var (
	v, _ = decorator.RegisteredCustomFunctionsValidator(validator.New())
)

func parsedJsonMetric(r *http.Request) (dto.Metrics, error) {
	m := dto.Metrics{}
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		return m, err
	}
	return m, nil
}

func parsedJsonMetrics(r *http.Request) ([]dto.Metrics, error) {
	var m []dto.Metrics
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		return m, err
	}
	return m, nil
}

func parsedMetricValues(r *http.Request) (string, string, string) {
	return chi.URLParam(r, "type"),
		chi.URLParam(r, "name"),
		chi.URLParam(r, "value")
}

func parsedAbridgedMetricValues(r *http.Request) (string, string) {
	return chi.URLParam(r, "type"),
		chi.URLParam(r, "name")
}

// parseAndCleanConfig For tests
func parseAndCleanConfig() func() {
	oldRepos := config.Configuration.Repository
	memStorage := mem_storage.NewStorage()
	config.Configuration.Repository = repository.New(memStorage)

	oldStoreInterval := config.Configuration.StoreInterval
	config.Configuration.StoreInterval = 1

	return func() {
		config.Configuration.Repository = oldRepos
		config.Configuration.StoreInterval = oldStoreInterval
	}
}

// prepareStorage For tests
func prepareStorage() error {
	repos := config.Configuration.Repository
	err := repos.UpdateMetric(context.Background(), dto.Metrics{ID: "PollCount", MType: "counter", Delta: ptrInt(527)})
	err2 := repos.UpdateMetric(context.Background(), dto.Metrics{ID: "RandomValue", MType: "gauge", Value: ptrFloat(0.47)})
	return multierr.Combine(err2, err)
}

// ptrInt For tests
func ptrInt(val int64) *int64 {
	return &val
}

// ptrFloat For tests
func ptrFloat(val float64) *float64 {
	return &val
}

// sendRequest For tests
func sendRequest(router chi.Router, r *http.Request) (int, http.Header, io.ReadCloser) {
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	resp := w.Result()
	return resp.StatusCode, resp.Header, resp.Body
}
