package client

import (
	"fmt"
	"github.com/go-resty/resty"
	"github.com/oldhanasong/go-metrics-collector/internal/config"
	"github.com/oldhanasong/go-metrics-collector/internal/dictionary"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"time"
)

type MetricSender struct {
	client resty.Client
}

func (s *MetricSender) Send(metric dto.Metrics) error {
	_, err := s.client.NewRequest().Post(getPreparedURL(metric))
	if err != nil {
		return err
	}
	return nil
}

func getPreparedURL(metric dto.Metrics) string {
	if metric.MType == dictionary.GaugeMetricType {
		return fmt.Sprintf("http://%s/update/%s/%s/%.2f", config.Configuration.Endpoint, metric.MType, metric.ID, *metric.Value)
	}
	return fmt.Sprintf("http://%s/update/%s/%s/%d", config.Configuration.Endpoint, metric.MType, metric.ID, *metric.Delta)
}

func NewClient() *MetricSender {
	client := resty.New().SetTimeout(2 * time.Second)

	return &MetricSender{
		client: *client,
	}
}
