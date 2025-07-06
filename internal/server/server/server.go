package server

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/oldhanasong/go-metrics-collector/internal/server/config"
	"github.com/oldhanasong/go-metrics-collector/internal/server/server/handler"
	"github.com/oldhanasong/go-metrics-collector/internal/server/server/middleware"
	"go.uber.org/multierr"
	"net"
	"net/http"
	"time"
)

type Server struct {
	server *http.Server
	router *chi.Mux

	OnStart func(ctx context.Context) error
	OnStop  func(ctx context.Context) error
}

func (s *Server) Start(ctx context.Context) {
	s.createServer(ctx)

	go func() {
		if err := s.startingServer(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
}

func (s *Server) Stop() error {
	if s.server == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return s.server.Shutdown(ctx)
}

func (s *Server) createAndConfigureRouter(ctx context.Context) {
	s.router = chi.NewRouter()

	s.router.Use(middleware.WithLogging)
	s.router.Use(middleware.GzipRequestDecompressing)
	s.router.Use(middleware.CheckRequestSign)
	s.router.Use(middleware.SigningResponse)
	s.router.Use(middleware.GzipResponseCompressing)

	s.router.Route("/", func(r chi.Router) {
		r.Get("/", handler.ShowMetricList(ctx))
		r.Route("/update", func(r chi.Router) {
			r.Post("/", handler.UpdateMetricJson(ctx))
			r.Route("/{type}/{name}/{value}", func(r chi.Router) {
				r.Post("/", handler.UpdateMetric(ctx))
			})
		})
		s.router.Route("/updates", func(router chi.Router) {
			router.Post("/", handler.UpdateMetricBatch(ctx))
		})
		r.Route("/value", func(r chi.Router) {
			r.Post("/", handler.GetMetricJson(ctx))
			r.Route("/{type}/{name}", func(r chi.Router) {
				r.Get("/", handler.GetMetric(ctx))
			})
		})
		r.Route("/ping", func(r chi.Router) {
			r.Get("/", handler.DbPing(ctx))
		})
	})
}

func (s *Server) createServer(ctx context.Context) {
	if s.server != nil {
		return
	}

	s.createAndConfigureRouter(ctx)
	s.server = &http.Server{
		Addr:    config.Configuration.Endpoint,
		Handler: s.router,
	}
	s.server.RegisterOnShutdown(func() {
		if s.OnStop != nil {
			if err := s.OnStop(ctx); err != nil && !errors.Is(err, context.Canceled) {
				panic(err)
			}
		}
	})
}

func (s *Server) startingServer(ctx context.Context) error {
	l, err := net.Listen("tcp", config.Configuration.Endpoint)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			errListenerClose := l.Close()
			errShutdown := s.Stop()
			if err = multierr.Combine(errListenerClose, errShutdown); err != nil {
				panic(err)
			}
		}
	}()

	if s.OnStart != nil {
		if err = s.OnStart(ctx); err != nil {
			errListenerClose := l.Close()
			return multierr.Combine(err, errListenerClose)
		}
	}

	if err = s.server.Serve(l); err != nil {
		errListenerClose := l.Close()
		return multierr.Combine(err, errListenerClose)
	}

	return nil
}
