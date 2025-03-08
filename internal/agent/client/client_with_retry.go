package client

import (
	"context"
	"github.com/go-resty/resty"
	requestPackage "github.com/godsareinvented/go-metrics-collector/internal/agent/client/request"
	"github.com/godsareinvented/go-metrics-collector/internal/agent/interfaces"
	"github.com/godsareinvented/go-metrics-collector/internal/general/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/general/utils/service/retry"
	"time"
)

type ClientWithRetry struct {
	client    resty.Client
	retryOpts dto.RetryOptions
}

func (s *ClientWithRetry) Send(ctx context.Context, metric dto.Metrics) error {
	return s.sendRequest(ctx, requestPackage.GetUpdateMetricJsonRequest(metric, &s.client))
}

func (s *ClientWithRetry) SendBatch(ctx context.Context, metrics *[]dto.Metrics) error {
	return s.sendRequest(ctx, requestPackage.GetUpdateMetricBatchRequest(metrics, &s.client))
}

func (s *ClientWithRetry) sendRequest(ctx context.Context, request *resty.Request) error {
	r := requestPackage.GetWrappedRequest(request)
	callback := func() (error, bool) {
		_, err := r.SetContext(ctx).Execute(r.Method, r.URL)
		return err, nil != err
	}

	err := retry.DoWithRetry(ctx, s.retryOpts, callback)
	if nil != err {
		return err
	}

	return err
}

func NewClientWithRetry(retryOpts dto.RetryOptions) interfaces.ClientInterface {
	client := resty.New().SetTimeout(2 * time.Second)

	return &ClientWithRetry{
		client:    *client,
		retryOpts: retryOpts,
	}
}
