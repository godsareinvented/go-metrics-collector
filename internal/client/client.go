package client

import (
	"github.com/go-resty/resty/v2"
	"github.com/oldhanasong/go-metrics-collector/internal/client/request"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/interfaces"
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
	for _, decorator := range s.decorators {
		r = decorator(r)
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
