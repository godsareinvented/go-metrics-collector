package server

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/oldhanasong/go-metrics-collector/internal/config"
	"github.com/oldhanasong/go-metrics-collector/internal/server/handler"
	"github.com/oldhanasong/go-metrics-collector/internal/server/middleware"
	"github.com/oldhanasong/go-metrics-collector/internal/util"
	"net"
	"net/http"
	"time"
)

type Server struct {
	server *http.Server
	router *chi.Mux

	OnStart func() error
	OnStop  func() error
}

func (s *Server) Start() {
	s.createServer()

	go func() {
		if err := s.startingServer(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
}

func (s *Server) Stop() error {
	if s.server == nil {
		return nil
	}

	var err error
	if s.OnStop != nil {
		err = s.OnStop()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	errShutdown := s.server.Shutdown(ctx)
	return util.WrappedErrs(err, errShutdown)
}

func (s *Server) createAndConfigureRouter() {
	s.router = chi.NewRouter()

	s.router.Use(middleware.WithLogging)
	s.router.Use(middleware.GzipRequestDecompressing)
	s.router.Use(middleware.GzipResponseCompressing)

	s.router.Route("/", func(r chi.Router) {
		s.router.Get("/", handler.ShowMetricList)
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
	})
}

func (s *Server) createServer() {
	if s.server != nil {
		return
	}

	s.createAndConfigureRouter()
	s.server = &http.Server{
		Addr:    config.Configuration.Endpoint,
		Handler: s.router,
	}
}

func (s *Server) startingServer() error {
	l, err := net.Listen("tcp", config.Configuration.Endpoint)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			errListenerClose := l.Close()
			errShutdown := s.Stop()
			if err = util.WrappedErrs(errListenerClose, errShutdown); err != nil {
				panic(err)
			}
		}
	}()

	if s.OnStart != nil {
		if err = s.OnStart(); err != nil {
			errListenerClose := l.Close()
			return util.WrappedErrs(err, errListenerClose)
		}
	}

	if err = s.server.Serve(l); err != nil {
		errListenerClose := l.Close()
		return util.WrappedErrs(err, errListenerClose)
	}

	return nil
}

func (s *Server) executeOnStartCallback() error {
	if s.OnStart == nil {
		return nil
	}
	return s.OnStart()
}
