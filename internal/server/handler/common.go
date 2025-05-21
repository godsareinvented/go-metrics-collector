package handler

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/oldhanasong/go-metrics-collector/internal/config"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/repository"
	"github.com/oldhanasong/go-metrics-collector/internal/storage/mem_storage"
	"github.com/oldhanasong/go-metrics-collector/internal/util"
	"github.com/oldhanasong/go-metrics-collector/internal/validation/decorator"
	"net/http"
)

var (
	v, _ = decorator.GetRegisteredCustomFunctionsValidator(validator.New())
)

// parseAndCleanConfig For tests
func parseAndCleanConfig() *repository.Repository {
	configConfigurator := config.ConfigConfigurator{}
	configConfigurator.ParseConfig()

	oldRepos := config.Configuration.Repository
	memStorage := mem_storage.NewInstance()
	config.Configuration.Repository = repository.NewInstance(&memStorage)

	return oldRepos
}

// prepareStorage For tests
func prepareStorage() error {
	repos := config.Configuration.Repository
	delta := int64(527)
	value := 0.47
	err := repos.UpdateMetric(dto.Metrics{ID: "PollCount", MType: "counter", Delta: &delta})
	err2 := repos.UpdateMetric(dto.Metrics{ID: "RandomValue", MType: "gauge", Value: &value})
	return util.WrappedErrs(err2, err)
}

// ptrInt For tests
func ptrInt(val int64) *int64 {
	return &val
}

// ptrFloat For tests
func ptrFloat(val float64) *float64 {
	return &val
}

func parsedJsonMetric(r *http.Request) (dto.Metrics, error) {
	m := dto.Metrics{}
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
