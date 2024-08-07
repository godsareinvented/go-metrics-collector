package client

import (
	"fmt"
	"github.com/go-resty/resty"
	"github.com/godsareinvented/go-metrics-collector/internal/config"
	"github.com/godsareinvented/go-metrics-collector/internal/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
	"time"
)

type MetricSender struct {
	client resty.Client
}

func (s *MetricSender) Send(metricDTO dto.Metric) error {
	request := s.client.R()
	_, err := request.Post(getPreparedURL(metricDTO))
	return err
}

func getPreparedURL(metricDTO dto.Metric) string {
	if dictionary.GaugeMetricType == metricDTO.Type {
		return fmt.Sprintf(
			"http://%s/update/%s/%s/%.2f",
			config.Configuration.Endpoint,
			metricDTO.Type,
			metricDTO.Name,
			metricDTO.Value,
		)
	}
	return fmt.Sprintf(
		"http://%s/update/%s/%s/%d",
		config.Configuration.Endpoint,
		metricDTO.Type,
		metricDTO.Name,
		metricDTO.Delta,
	)
}

func NewInstance() *MetricSender {
	client := resty.New().SetTimeout(2 * time.Second)

	return &MetricSender{
		client: *client,
	}
}
