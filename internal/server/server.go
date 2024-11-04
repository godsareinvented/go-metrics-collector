package server

import (
	"context"
	"errors"
	"github.com/go-chi/chi"
	"github.com/godsareinvented/go-metrics-collector/internal/config"
	"github.com/godsareinvented/go-metrics-collector/internal/server/handler"
	"github.com/godsareinvented/go-metrics-collector/internal/server/middleware"
	"net/http"
)

type Server struct {
	server *http.Server
	router *chi.Mux
	ctx    *context.Context
	cancel *context.CancelFunc

	OnStart func() error
	OnStop  func() error
}

func (s *Server) Start() {
	s.createAndConfigureRouter()

	s.server = &http.Server{
		Addr:    config.Configuration.Endpoint,
		Handler: s.router,
	}

	s.createContext()

	defer (*s.cancel)()

	serverIsRunning := make(chan bool)
	go func() {
		go func() {
			if err := s.server.ListenAndServe(); nil != err && !errors.Is(err, http.ErrServerClosed) {
				panic("ListenAndServe: " + err.Error())
			}
		}()

		serverIsRunning <- true
	}()

	select {
	case <-(*s.ctx).Done():
		err := s.server.Shutdown(*s.ctx)
		if nil != err {
			return
		}
	case <-serverIsRunning:
		s.executeOnStartCallback()
	}
}

func (s *Server) Stop() {
	if nil == s.OnStop {
		(*s.cancel)()
		return
	}

	err := s.OnStop()

	(*s.cancel)()

	if nil != err {
		panic("OnStop: " + err.Error())
	}
}

func (s *Server) createAndConfigureRouter() {
	s.router = chi.NewRouter()

	s.router.Use(middleware.WithLogging)
	s.router.Use(middleware.GzipRequestDecompressing)
	s.router.Use(middleware.GzipResponseCompressing)

	s.router.Route("/update", func(router chi.Router) {
		router.Post("/", handler.UpdateMetricJson)
		router.Route("/{type}/{name}/{value}", func(router chi.Router) {
			router.Post("/", handler.UpdateMetric)
		})
	})
	s.router.Route("/value", func(router chi.Router) {
		router.Post("/", handler.GetMetricJson)
		router.Route("/{type}/{name}", func(router chi.Router) {
			router.Get("/", handler.GetMetric)
		})
	})
	s.router.Get("/", handler.ShowMetricList)
}

func (s *Server) createContext() {
	ctx, cancel := context.WithCancel(context.Background())
	s.ctx = &ctx
	s.cancel = &cancel
}

func (s *Server) executeOnStartCallback() {
	if nil == s.OnStart {
		return
	}
	if err := s.OnStart(); nil != err {
		// todo: Не panic.
		panic("OnStart: " + err.Error())
	}
}
