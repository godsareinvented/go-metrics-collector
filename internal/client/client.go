package client

import (
	"github.com/go-resty/resty"
	"github.com/godsareinvented/go-metrics-collector/internal/client/decorator"
	"github.com/godsareinvented/go-metrics-collector/internal/client/request"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
	"time"
)

type MetricSender struct {
	client resty.Client
}

func (s *MetricSender) Send(metricDTO dto.Metrics) error {
	r := decorator.GzipCompress(request.GetUpdateMetricJsonRequest(metricDTO, &s.client))

	_, err := r.Execute(r.Method, r.URL)
	return err
}

func (s *MetricSender) SendBatch(metrics []dto.Metrics) error {
	r := decorator.GzipCompress(request.GetUpdateMetricBatchRequest(metrics, &s.client))

	_, err := r.Execute(r.Method, r.URL)
	return err
}

func NewInstance() *MetricSender {
	client := resty.New().SetTimeout(2 * time.Second)

	return &MetricSender{
		client: *client,
	}
}
