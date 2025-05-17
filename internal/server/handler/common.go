package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/oldhanasong/go-metrics-collector/internal/config"
	"github.com/oldhanasong/go-metrics-collector/internal/repository"
	"github.com/oldhanasong/go-metrics-collector/internal/storage/mem_storage"
	"net/http"
)

// parseAndCleanConfig For tests
func parseAndCleanConfig() {
	configConfigurator := config.ConfigConfigurator{}
	configConfigurator.ParseConfig()

	memStorage := mem_storage.NewInstance()
	config.Configuration.Repository = repository.NewInstance(&memStorage)
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
