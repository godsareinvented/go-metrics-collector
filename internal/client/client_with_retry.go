package client

import (
	"context"
	"github.com/go-resty/resty/v2"
	"github.com/oldhanasong/go-metrics-collector/internal/client/request"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/interfaces"
	"github.com/oldhanasong/go-metrics-collector/internal/service/retry"
	"github.com/oldhanasong/go-metrics-collector/internal/service/retry/strategy"
	"time"
)

type ClientWithRetry struct {
	client       resty.Client
	retryService retry.RetryService
	decorators   []interfaces.Decorator
}

func (s *ClientWithRetry) Send(metrics dto.Metrics) error {
	return s.sendRequest(request.GetUpdateMetricJsonRequest(metrics, &s.client))
}

func (s *ClientWithRetry) SendBatch(metrics []dto.Metrics) error {
	return s.sendRequest(request.GetUpdateMetricBatchRequest(metrics, &s.client))
}

func (s *ClientWithRetry) Use(decorator interfaces.Decorator) {
	if decorator != nil {
		s.decorators = append(s.decorators, decorator)
	}
}

func (s *ClientWithRetry) sendRequest(r *resty.Request) error {
	r = s.prepareRequest(r)

	callback := func() (error, bool) {
		_, err := r.Execute(r.Method, r.URL)
		return err, err != nil
	}
	if err := s.retryService.DoWithRetry(context.Background(), callback); err != nil {
		return err
	}

	return nil
}

func (s *ClientWithRetry) prepareRequest(r *resty.Request) *resty.Request {
	for _, decorator := range s.decorators {
		r = decorator(r)
	}
	return r
}

func NewClientWithRetry() interfaces.Client {
	client := resty.New().SetTimeout(2 * time.Second)
	retryService := retry.NewInstance(strategy.NewDefaultFixedIntervalStrategy())

	return &ClientWithRetry{
		client:       *client,
		decorators:   make([]interfaces.Decorator, 0),
		retryService: retryService,
	}
}
