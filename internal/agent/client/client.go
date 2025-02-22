package client

import (
	"github.com/go-resty/resty"
	requestPackage "github.com/godsareinvented/go-metrics-collector/internal/agent/client/request"
	"github.com/godsareinvented/go-metrics-collector/internal/agent/interfaces"
	"github.com/godsareinvented/go-metrics-collector/internal/general/dto"
	"time"
)

type Client struct {
	client resty.Client
}

func (s *Client) Send(metric dto.Metrics) error {
	return s.sendRequest(requestPackage.GetUpdateMetricJsonRequest(metric, &s.client))
}

func (s *Client) SendBatch(metrics *[]dto.Metrics) error {
	return s.sendRequest(requestPackage.GetUpdateMetricBatchRequest(metrics, &s.client))
}

func (s *Client) sendRequest(request *resty.Request) error {
	r := requestPackage.GetWrappedRequest(request)

	_, err := r.Execute(r.Method, r.URL)
	return err
}

// NewBaseClient unused
func _() interfaces.Client {
	client := resty.New().SetTimeout(2 * time.Second)

	return &Client{
		client: *client,
	}
}
