package server

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/godsareinvented/go-metrics-collector/internal/config"
	"github.com/godsareinvented/go-metrics-collector/internal/server/handler"
	"github.com/godsareinvented/go-metrics-collector/internal/server/middleware"
	"net"
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
	s.createParentContext()
	defer (*s.cancel)()

	s.createAndConfigureRouter()
	s.createServer()

	go func() {
		if err := s.startingServer(); nil != err {
			panic(err)
		}
	}()

	select {
	case <-(*s.ctx).Done():
		err := s.server.Shutdown(*s.ctx)
		if nil != err {
			return
		}
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

	s.router.Route("/", func(router chi.Router) {
		s.router.Get("/", handler.ShowMetricList)
		s.router.Route("/update", func(router chi.Router) {
			router.Post("/", handler.UpdateMetricJson)
			router.Post("/{type}/{name}/{value}", handler.UpdateMetric)
		})
		s.router.Route("/value", func(router chi.Router) {
			router.Post("/", handler.GetMetricJson)
			router.Get("/{type}/{name}", handler.GetMetric)
		})
		s.router.Route("/ping", func(router chi.Router) {
			router.Get("/", handler.DbPing(*s.ctx))
		})
	})
}

func (s *Server) createServer() {
	s.server = &http.Server{
		Addr:    config.Configuration.Endpoint,
		Handler: s.router,
	}
}

func (s *Server) createParentContext() {
	ctx, cancel := context.WithCancel(context.Background())
	s.ctx = &ctx
	s.cancel = &cancel
}

func (s *Server) startingServer() error {
	l, err := net.Listen("tcp", config.Configuration.Endpoint)
	if nil != err {
		return err
	}

	err = s.executeOnStartCallback()
	if nil != err {
		_ = l.Close()
		return err
	}

	if err = s.server.Serve(l); err != nil {
		_ = l.Close()
		return err
	}

	return nil
}

func (s *Server) executeOnStartCallback() error {
	if nil == s.OnStart {
		return nil
	}
	if err := s.OnStart(); nil != err {
		return err
	}
	return nil
}
