package server

import (
	"github.com/go-chi/chi"
	"github.com/godsareinvented/go-metrics-collector/internal/server/handler"
	"github.com/godsareinvented/go-metrics-collector/internal/server/middleware"
	"net/http"
)

type Server struct{}

func (server *Server) Start() {
	router := chi.NewRouter()

	router.Use(middleware.WithLogging)

	router.Post("/update/{type}/{name}/{value}", handler.UpdateMetric)
	router.Get("/value/{type}/{name}", handler.GetMetric)
	router.Get("/", handler.ShowMetricList)

	err := http.ListenAndServe("localhost:8080", router)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
