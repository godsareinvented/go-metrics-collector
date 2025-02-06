package client

import (
	"context"
	"github.com/go-resty/resty"
	requestPackage "github.com/godsareinvented/go-metrics-collector/internal/agent/client/request"
	"github.com/godsareinvented/go-metrics-collector/internal/agent/interfaces"
	"github.com/godsareinvented/go-metrics-collector/internal/general/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/server/service/retry"
	"github.com/godsareinvented/go-metrics-collector/internal/server/service/retry/prepared_option"
	"time"
)

type ClientWithRetry struct {
	client resty.Client
}

func (s *ClientWithRetry) Send(metric dto.Metrics) error {
	return s.sendRequest(requestPackage.GetUpdateMetricJsonRequest(metric, &s.client))
}

func (s *ClientWithRetry) SendBatch(metrics []dto.Metrics) error {
	return s.sendRequest(requestPackage.GetUpdateMetricBatchRequest(metrics, &s.client))
}

func (s *ClientWithRetry) sendRequest(request *resty.Request) error {
	r := requestPackage.GetWrappedRequest(request)
	callback := func() (error, bool) {
		_, err := r.Execute(r.Method, r.URL)
		return err, nil != err
	}

	err := retry.DoWithRetry(context.Background(), prepared_option.DefaultFixedDelayListOptions, callback)
	if nil != err {
		return err
	}

	return err
}

func NewClientWithRetry() interfaces.Client {
	client := resty.New().SetTimeout(20 * time.Second)

	return &ClientWithRetry{
		client: *client,
	}
}
