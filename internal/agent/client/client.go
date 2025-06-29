package client

import (
	"github.com/go-resty/resty/v2"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/client/request"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/interfaces"
	"github.com/oldhanasong/go-metrics-collector/internal/general/dto"
	"time"
)

type Client struct {
	client     resty.Client
	decorators []interfaces.Decorator
}

func (s *Client) Send(metrics dto.Metrics) error {
	return s.sendRequest(request.GetUpdateMetricJsonRequest(metrics, &s.client))
}

func (s *Client) SendBatch(metrics []dto.Metrics) error {
	return s.sendRequest(request.GetUpdateMetricBatchRequest(metrics, &s.client))
}

func (s *Client) Use(decorator interfaces.Decorator) {
	if decorator != nil {
		s.decorators = append(s.decorators, decorator)
	}
}

func (s *Client) sendRequest(r *resty.Request) error {
	r = s.prepareRequest(r)
	_, err := r.Execute(r.Method, r.URL)
	return err
}

func (s *Client) prepareRequest(r *resty.Request) *resty.Request {
	for i := len(s.decorators) - 1; i >= 0; i-- {
		r = s.decorators[i](r)
	}
	return r
}

func NewClient() interfaces.Client {
	client := resty.New().SetTimeout(2 * time.Second)

	return &Client{
		client:     *client,
		decorators: make([]interfaces.Decorator, 0),
	}
}
