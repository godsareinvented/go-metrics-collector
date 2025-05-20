package handler

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/oldhanasong/go-metrics-collector/internal/config"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/repository"
	"github.com/oldhanasong/go-metrics-collector/internal/storage/mem_storage"
	"net/http"
)

var (
	v = validator.New()
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
