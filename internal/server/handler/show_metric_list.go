package handler

import (
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/service/metric/manager"
	"html/template"
	"net/http"
	"sort"
)

func ShowMetricList(responseWriter http.ResponseWriter, _ *http.Request) {
	metricManager := manager.MetricManager{}
	metricDTOList := metricManager.GetList()

	sort.Slice(metricDTOList, func(i, j int) bool {
		return metricDTOList[i].Name < metricDTOList[j].Name
	})

	tmpl := template.Must(template.ParseFiles("internal/template/main_page.html"))
	data := struct {
		Items []dto.Metric
	}{
		Items: metricDTOList,
	}

	err := tmpl.Execute(responseWriter, data)
	if err != nil {
		panic(err)
	}

	responseWriter.WriteHeader(http.StatusOK)
}
