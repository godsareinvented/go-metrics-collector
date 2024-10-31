package server

import (
	"github.com/go-chi/chi"
	"github.com/godsareinvented/go-metrics-collector/internal/config"
	"github.com/godsareinvented/go-metrics-collector/internal/server/handler"
	"github.com/godsareinvented/go-metrics-collector/internal/server/middleware"
	"net/http"
)

type Server struct{}

func (server *Server) Start() {
	router := chi.NewRouter()

	router.Use(middleware.WithLogging)
	router.Use(middleware.GzipResponseCompressing)

	router.Route("/update", func(router chi.Router) {
		router.Post("/", handler.UpdateMetricJson)
		router.Route("/{type}/{name}/{value}", func(router chi.Router) {
			router.Post("/", handler.UpdateMetric)
		})
	})
	router.Route("/value", func(router chi.Router) {
		router.Post("/", handler.GetMetricJson)
		router.Route("/{type}/{name}", func(router chi.Router) {
			router.Get("/", handler.GetMetric)
		})
	})
	router.Get("/", handler.ShowMetricList)

	err := http.ListenAndServe(config.Configuration.Endpoint, router)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
