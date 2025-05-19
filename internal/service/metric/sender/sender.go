package sender

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

func (s *MetricSender) Send(metric dto.Metric) {
	_, _ = resty.NewRequest().Post(getPreparedURL(metric))
}

func getPreparedURL(metric dto.Metric) string {
	if metric.Type == dictionary.GaugeMetricType {
		return fmt.Sprintf("http://%s/update/%s/%s/%.2f", config.Configuration.Endpoint, metric.Type, metric.Name, metric.Value)
	}
	return fmt.Sprintf("http://%s/update/%s/%s/%d", config.Configuration.Endpoint, metric.Type, metric.Name, metric.Delta)
}

func NewSender() *MetricSender {
	client := resty.New().SetTimeout(2 * time.Second)

	return &MetricSender{
		client: *client,
	}
}
