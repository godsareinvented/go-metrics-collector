package request

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/oldhanasong/go-metrics-collector/internal/general/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/server/config"
)

func GetUpdateMetricJsonRequest(metric dto.Metrics, client *resty.Client) *resty.Request {
	request := client.R()

	body, err := json.Marshal(metric)
	if err != nil {
		panic(err)
	}

	request.URL = fmt.Sprintf("http://%s/update/", config.Configuration.Endpoint)
	request.Method = resty.MethodPost
	request.SetBody(body)
	request.Header.Set("Content-Type", "application/json")

	return request
}
