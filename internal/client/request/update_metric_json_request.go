package request

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty"
	"github.com/godsareinvented/go-metrics-collector/internal/config"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
)

func GetUpdateMetricJsonRequest(metric dto.Metrics, client *resty.Client) *resty.Request {
	request := client.R()

	body, err := json.Marshal(metric)
	if nil != err {
		panic(err)
	}

	request.URL = getURL()
	request.Method = resty.MethodPost
	request.SetBody(body)
	request.Header.Set("Content-Type", "application/json")

	return request
}

func getURL() string {
	return fmt.Sprintf("http://%s/update/", config.Configuration.Endpoint)
}
