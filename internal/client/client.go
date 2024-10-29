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

func (s *MetricSender) Send(metricDTO dto.Metrics) error {
	request := s.client.R()
	_, err := request.Post(getPreparedURL(metricDTO))
	return err
}

// todo: Перенести шаблоны урлов в конфиг
func getPreparedURL(metricDTO dto.Metrics) string {
	if dictionary.GaugeMetricType == metricDTO.MType {
		return fmt.Sprintf(
			"http://%s/update/%s/%s/%.2f",
			config.Configuration.Endpoint,
			metricDTO.MType,
			metricDTO.MName,
			*metricDTO.Value,
		)
	}
	return fmt.Sprintf(
		"http://%s/update/%s/%s/%d",
		config.Configuration.Endpoint,
		metricDTO.MType,
		metricDTO.MName,
		*metricDTO.Delta,
	)
}

func NewInstance() *MetricSender {
	client := resty.New().SetTimeout(2 * time.Second)

	return &MetricSender{
		client: *client,
	}
}
