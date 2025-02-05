package client

import (
	"context"
	"github.com/go-resty/resty/v2"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/client/request"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/interfaces"
	"github.com/oldhanasong/go-metrics-collector/internal/general/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/server/service/retry"
	"github.com/oldhanasong/go-metrics-collector/internal/server/service/retry/prepared_option"
	"time"
)

type ClientWithRetry struct {
	client     resty.Client
	decorators []interfaces.Decorator
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
	if err := retry.DoWithRetry(context.Background(), prepared_option.DefaultFixedDelayListOptions, callback); err != nil {
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

	return &ClientWithRetry{
		client:     *client,
		decorators: make([]interfaces.Decorator, 0),
	}
}
