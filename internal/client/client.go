package client

import (
	"github.com/go-resty/resty/v2"
	"github.com/oldhanasong/go-metrics-collector/internal/client/decorator"
	"github.com/oldhanasong/go-metrics-collector/internal/client/request"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"time"
)

type MetricSender struct {
	client resty.Client
}

func (s *MetricSender) Send(metric dto.Metrics) error {
	r := decorator.GzipCompress(request.GetUpdateMetricJsonRequest(metric, &s.client))

	_, err := r.Execute(r.Method, r.URL)
	return err
}

func NewClient() *MetricSender {
	client := resty.New().SetTimeout(2 * time.Second)

	return &MetricSender{
		client: *client,
	}
}
