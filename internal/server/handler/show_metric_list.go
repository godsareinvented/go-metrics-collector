package handler

import (
	"github.com/godsareinvented/go-metrics-collector/internal/buisness_logic/manager"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
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
		return
	}
}
