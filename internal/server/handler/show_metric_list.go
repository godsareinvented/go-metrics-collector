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
		return metricDTOList[i].MName < metricDTOList[j].MName
	})

	tmpl := template.Must(template.ParseFiles("internal/template/main_page.html"))
	data := struct {
		Items []dto.Metrics
	}{
		Items: metricDTOList,
	}

	err := tmpl.Execute(responseWriter, data)
	if err != nil {
		panic(err)
	}

	responseWriter.Header().Set("Content-Type", "text/html")
	responseWriter.WriteHeader(http.StatusOK)
}
