package request

import (
	"fmt"
	"github.com/go-resty/resty"
	"github.com/oldhanasong/go-metrics-collector/internal/config"
	"github.com/oldhanasong/go-metrics-collector/internal/dictionary"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
)

func GetUpdateMetricRequest(metric dto.Metrics, client *resty.Client) *resty.Request {
	request := client.R()

	request.URL = getPreparedURL(metric)
	request.Method = resty.MethodPost

	return request
}

func getPreparedURL(metric dto.Metrics) string {
	if metric.MType == dictionary.GaugeMetricType {
		return fmt.Sprintf("http://%s/update/%s/%s/%.2f", config.Configuration.Endpoint, metric.MType, metric.ID, *metric.Value)
	}
	return fmt.Sprintf("http://%s/update/%s/%s/%d", config.Configuration.Endpoint, metric.MType, metric.ID, *metric.Delta)
}
