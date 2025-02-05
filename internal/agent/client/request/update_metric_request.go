package request

import (
	"fmt"
	"github.com/go-resty/resty"
	"github.com/godsareinvented/go-metrics-collector/internal/agent/config"
	"github.com/godsareinvented/go-metrics-collector/internal/general/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/general/dto"
)

// GetUpdateMetricRequest unused
func _(metric dto.Metrics, client *resty.Client) *resty.Request {
	request := client.R()

	request.URL = getPreparedURL(metric)
	request.Method = resty.MethodPost

	return request
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
