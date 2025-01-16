package request

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/oldhanasong/go-metrics-collector/internal/config"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
)

func GetUpdateMetricBatchRequest(metrics []dto.Metrics, client *resty.Client) *resty.Request {
	request := client.R()

	body, err := json.Marshal(metrics)
	if nil != err {
		panic(err)
	}

	request.URL = fmt.Sprintf("http://%s/updates/", config.Configuration.Endpoint)
	request.Method = resty.MethodPost
	request.SetBody(body)
	request.Header.Set("Content-Type", "application/json")

	return request
}
