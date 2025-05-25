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

	var err error
	if s.OnStop != nil {
		err = s.OnStop(ctx)
	}

	errShutdown := s.server.Shutdown(ctx)
	return util.WrappedErrs(err, errShutdown)
}

func (s *Server) createAndConfigureRouter(ctx context.Context) {
	s.router = chi.NewRouter()

	s.router.Use(middleware.WithLogging)
	s.router.Use(middleware.GzipRequestDecompressing)
	s.router.Use(middleware.GzipResponseCompressing)

	s.router.Route("/", func(r chi.Router) {
		r.Get("/", handler.ShowMetricList(ctx))
		r.Route("/update", func(r chi.Router) {
			r.Post("/", handler.UpdateMetricJson(ctx))
			r.Route("/{type}/{name}/{value}", func(r chi.Router) {
				r.Post("/", handler.UpdateMetric(ctx))
			})
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
			if err = util.WrappedErrs(errListenerClose, errShutdown); err != nil {
				panic(err)
			}
		}
	}()

	if s.OnStart != nil {
		if err = s.OnStart(ctx); err != nil {
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
