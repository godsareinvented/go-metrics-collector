package client

import (
	"context"
	"github.com/go-resty/resty"
	"github.com/godsareinvented/go-metrics-collector/internal/client/decorator"
	"github.com/godsareinvented/go-metrics-collector/internal/client/request"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/interfaces"
	"github.com/godsareinvented/go-metrics-collector/internal/service/retry"
	"github.com/godsareinvented/go-metrics-collector/internal/service/retry/strategy"
	"time"
)

type ClientWithRetry struct {
	client       resty.Client
	retryService retry.RetryService
}

func (s *ClientWithRetry) Send(metricDTO dto.Metrics) error {
	return s.sendRequest(request.GetUpdateMetricJsonRequest(metricDTO, &s.client))
}

func (s *ClientWithRetry) SendBatch(metrics []dto.Metrics) error {
	return s.sendRequest(request.GetUpdateMetricBatchRequest(metrics, &s.client))
}

func (s *ClientWithRetry) sendRequest(request *resty.Request) error {
	r := decorator.GzipCompress(request)
	callback := func() (error, bool) {
		_, err := r.Execute(r.Method, r.URL)
		return err, nil != err
	}

	err := s.retryService.DoWithRetry(context.Background(), callback)
	if nil != err {
		return err
	}

	return err
}

func NewClientWithRetry() interfaces.Client {
	client := resty.New().SetTimeout(2 * time.Second)
	retryService := retry.NewInstance(strategy.NewDefaultFixedIntervalStrategy())

	return &ClientWithRetry{
		client:       *client,
		retryService: retryService,
	}
}
