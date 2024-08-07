package handler

import (
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	manager "github.com/oldhanasong/go-metrics-collector/internal/service/metric"
	"html/template"
	"net/http"
	"sort"
)

var (
	mainPageTplPath = "template/main_page.html"
)

func ShowMetricList(responseWriter http.ResponseWriter, _ *http.Request) {
	metricManager := manager.MetricManager{}
	metricDTOList := metricManager.GetList()

	sort.Slice(metricDTOList, func(i, j int) bool {
		return metricDTOList[i].Name < metricDTOList[j].Name
	})

	tmpl := template.Must(template.ParseFiles(mainPageTplPath))
	data := struct {
		Items []dto.Metric
	}{
		Items: metricDTOList,
	}

	err := tmpl.Execute(responseWriter, data)
	if err != nil {
		return
	}

	responseWriter.WriteHeader(http.StatusOK)
}
