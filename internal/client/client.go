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

func (s *MetricSender) Send(metric dto.Metric) error {
	_, err := resty.NewRequest().Post(getPreparedURL(metric))
	if err != nil {
		return err
	}
	return nil
}

func getPreparedURL(metric dto.Metric) string {
	if metric.Type == dictionary.GaugeMetricType {
		return fmt.Sprintf("http://%s/update/%s/%s/%.2f", config.Configuration.Endpoint, metric.Type, metric.Name, metric.Value)
	}
	return fmt.Sprintf("http://%s/update/%s/%s/%d", config.Configuration.Endpoint, metric.Type, metric.Name, metric.Delta)
}

func NewClient() *MetricSender {
	client := resty.New().SetTimeout(2 * time.Second)

	return &MetricSender{
		client: *client,
	}
}
