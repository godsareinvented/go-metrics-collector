package server

import (
	"github.com/godsareinvented/go-metrics-collector/internal/server/handler"
	"net/http"
)

type Server struct{}

func (server *Server) Start() {
	err := http.ListenAndServe("localhost:8080", getHandlers())
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func getHandlers() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/update/{type}/{name}/{value}", &handler.UpdateMetricHandler{})

	return mux
}
