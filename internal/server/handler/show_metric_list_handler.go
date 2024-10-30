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
	metricList, err := metricManager.GetList()

	if err != nil {
		http.Error(responseWriter, "failed to get the metric list", http.StatusInternalServerError)
		return
	}

	sort.Slice(metricList, func(i, j int) bool {
		return metricList[i].ID < metricList[j].ID
	})

	tmpl := template.Must(template.ParseFiles(mainPageTplPath))
	data := struct {
		Items []dto.Metrics
	}{
		Items: metricList,
	}

	err = tmpl.Execute(responseWriter, data)
	if err != nil {
		return
	}

	responseWriter.WriteHeader(http.StatusOK)
}
