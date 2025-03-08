package client

import (
	"context"
	"github.com/go-resty/resty"
	requestPackage "github.com/godsareinvented/go-metrics-collector/internal/agent/client/request"
	"github.com/godsareinvented/go-metrics-collector/internal/agent/interfaces"
	"github.com/godsareinvented/go-metrics-collector/internal/general/dto"
	"time"
)

type Client struct {
	client resty.Client
}

func (s *Client) Send(ctx context.Context, metric dto.Metrics) error {
	return s.sendRequest(ctx, requestPackage.GetUpdateMetricJsonRequest(metric, &s.client))
}

func (s *Client) SendBatch(ctx context.Context, metrics *[]dto.Metrics) error {
	return s.sendRequest(ctx, requestPackage.GetUpdateMetricBatchRequest(metrics, &s.client))
}

func (s *Client) sendRequest(ctx context.Context, request *resty.Request) error {
	r := requestPackage.GetWrappedRequest(request)

	_, err := r.SetContext(ctx).Execute(r.Method, r.URL)
	return err
}

// NewBaseClient unused
func _() interfaces.ClientInterface {
	client := resty.New().SetTimeout(2 * time.Second)

	return &Client{
		client: *client,
	}
}
