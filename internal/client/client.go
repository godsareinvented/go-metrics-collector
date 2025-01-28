package client

import (
	"github.com/go-resty/resty"
	"github.com/godsareinvented/go-metrics-collector/internal/client/decorator"
	"github.com/godsareinvented/go-metrics-collector/internal/client/request"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/interfaces"
	"time"
)

type Client struct {
	client resty.Client
}

func (s *Client) Send(metricDTO dto.Metrics) error {
	return s.sendRequest(request.GetUpdateMetricJsonRequest(metricDTO, &s.client))
}

func (s *Client) SendBatch(metrics []dto.Metrics) error {
	return s.sendRequest(request.GetUpdateMetricBatchRequest(metrics, &s.client))
}

func (s *Client) sendRequest(request *resty.Request) error {
	r := decorator.GzipCompress(request)

	_, err := r.Execute(r.Method, r.URL)
	return err
}

func NewBaseClient() interfaces.Client {
	client := resty.New().SetTimeout(2 * time.Second)

	return &Client{
		client: *client,
	}
}
