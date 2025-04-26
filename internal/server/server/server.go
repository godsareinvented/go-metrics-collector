package server

import (
	"context"
	"errors"
	"github.com/go-chi/chi"
	"github.com/godsareinvented/go-metrics-collector/internal/server/config"
	"github.com/godsareinvented/go-metrics-collector/internal/server/server/handler"
	"github.com/godsareinvented/go-metrics-collector/internal/server/server/middleware"
	"net"
	"net/http"
	"sync"
	"time"
)

type (
	OnStartCallback func(ctx context.Context, errCh chan<- error) error

	OnStopCallback func(ctx context.Context, errCh chan<- error)

	Server struct {
		OnStart    OnStartCallback
		OnStop     OnStopCallback
		server     *http.Server
		router     *chi.Mux
		stopCtx    context.Context
		stopCancel context.CancelFunc
		errCh      chan error
	}
)

func (s *Server) Start(ctx context.Context, cancel context.CancelFunc) (*sync.WaitGroup, chan error) {
	wg := &sync.WaitGroup{}
	s.stopCtx, s.stopCancel = context.WithCancel(context.Background())
	s.errCh = make(chan error, 3)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer cancel()
		defer close(s.errCh)
		defer wg.Done()

		s.createAndConfigureServer(ctx)
		runtimeCtx, runtimeCancel := context.WithCancel(context.Background())
		go func(cancel context.CancelFunc) {
			defer cancel()

			if err := s.run(ctx); nil != err && !errors.Is(err, http.ErrServerClosed) {
				s.errCh <- err
			}
		}(runtimeCancel)

		select {
		case <-ctx.Done():
		case <-s.stopCtx.Done():
			s.stop()
		case <-runtimeCtx.Done():
			s.onShutdownCallback()
			return
		}
	}(wg)

	return wg, s.errCh
}

func (s *Server) Stop() {
	if nil == s.stopCtx {
		return
	}

	s.stopCancel()
}

func (s *Server) stop() {
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := s.server.Shutdown(ctxWithTimeout)
	if nil != err {
		s.errCh <- err
	}
}

func (s *Server) run(ctx context.Context) error {
	l, err := net.Listen("tcp", config.Configuration.Endpoint)
	if nil != err {
		return err
	}

	defer func(l net.Listener) {
		err := l.Close()
		if nil != err {
			s.errCh <- err
		}
	}(l)

	err = s.OnStart(ctx, s.errCh)
	if nil != err {
		return err
	}

	if err = s.server.Serve(l); nil != err {
		return err
	}

	return nil
}

func (s *Server) createAndConfigureServer(ctx context.Context) {
	s.createAndConfigureRouter(ctx)
	s.server = &http.Server{
		Addr:    config.Configuration.Endpoint,
		Handler: s.router,
	}
	s.server.RegisterOnShutdown(s.onShutdownCallback)
}

func (s *Server) createAndConfigureRouter(ctx context.Context) {
	s.router = chi.NewRouter()

	s.router.Use(middleware.WithLogging)
	s.router.Use(middleware.GzipRequestDecompressing)
	s.router.Use(middleware.CheckRequestSign)
	s.router.Use(middleware.SigningResponse)
	s.router.Use(middleware.GzipResponseCompressing)

	s.router.Route("/", func(router chi.Router) {
		router.Get("/", handler.ShowMetricList(ctx))
		router.Route("/updates", func(router chi.Router) {
			router.Post("/", handler.UpdateMetricBatchMetric(ctx))
		})
		router.Route("/update", func(router chi.Router) {
			router.Post("/", handler.UpdateMetricJson(ctx))
			router.Post("/{type}/{name}/{value}", handler.UpdateMetric(ctx))
		})
		router.Route("/value", func(router chi.Router) {
			router.Post("/", handler.GetMetricJson(ctx))
			router.Get("/{type}/{name}", handler.GetMetric(ctx))
		})
		router.Route("/ping", func(router chi.Router) {
			router.Get("/", handler.DbPing(ctx))
		})
	})
}

func (s *Server) onShutdownCallback() {
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	s.OnStop(ctxWithTimeout, s.errCh)
}
