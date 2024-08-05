package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/oldhanasong/go-metrics-collector/internal/config"
	"github.com/oldhanasong/go-metrics-collector/internal/server/handler"
	"github.com/oldhanasong/go-metrics-collector/internal/server/middleware"
	"net/http"
)

type Server struct{}

func (server *Server) Start() {
	router := chi.NewRouter()

	router.Use(middleware.WithLogging)

	router.Post("/update/{type}/{name}/{value}", handler.UpdateMetric)
	router.Get("/value/{type}/{name}", handler.GetMetric)
	router.Get("/", handler.ShowMetricList)

	err := http.ListenAndServe(config.Configuration.Endpoint, router)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
