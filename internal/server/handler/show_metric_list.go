package handler

import (
	"fmt"
	"github.com/godsareinvented/go-metrics-collector/internal/dictionary"
	"github.com/godsareinvented/go-metrics-collector/internal/repository"
	"net/http"
)

func ShowMetricList(responseWriter http.ResponseWriter, request *http.Request) {
	repos := repository.GetInstance[int64](dictionary.GaugeMetricType)
	metricDTOList := repos.GetMetricList()

	fmt.Println(metricDTOList)
}
