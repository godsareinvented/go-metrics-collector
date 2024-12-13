package handler

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/oldhanasong/go-metrics-collector/internal/config"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/repository"
	"github.com/oldhanasong/go-metrics-collector/internal/storage/mem_storage"
	"github.com/oldhanasong/go-metrics-collector/internal/validation/decorator"
	"go.uber.org/multierr"
	"net/http"
)

var (
	v, _ = decorator.GetRegisteredCustomFunctionsValidator(validator.New())
)

// parseAndCleanConfig For tests
func parseAndCleanConfig() func() {
	oldRepos := config.Configuration.Repository
	memStorage := mem_storage.NewStorage()
	config.Configuration.Repository = repository.NewInstance(memStorage)

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

// combineContext Получение комбинированного контекста, чтобы хендлер мог обработать завершение контекстов как приложения, так и запроса
func combineContext(serverCtx context.Context, requestCtx context.Context) (context.Context, context.CancelFunc) {
	combinedCtx, cancel := context.WithCancel(requestCtx)

	go func() {
		select {
		case <-serverCtx.Done():
			cancel()
		case <-combinedCtx.Done():
			return
		}
	}()

	return combinedCtx, cancel
}

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
