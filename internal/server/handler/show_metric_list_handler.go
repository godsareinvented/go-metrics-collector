package handler

import (
	"context"
	"github.com/oldhanasong/go-metrics-collector/internal/config"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"html/template"
	"net/http"
	"sort"
)

var (
	mainPageTplPath = "template/main_page.html"
)

func ShowMetricList(_ context.Context) http.HandlerFunc {
	fn := func(responseWriter http.ResponseWriter, request *http.Request) {
		metricList, err := config.Configuration.Repository.GetAllMetrics()
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

		responseWriter.Header().Set("Content-Type", "text/html")
		responseWriter.WriteHeader(http.StatusOK)
		if err = tmpl.Execute(responseWriter, data); err != nil {
			http.Error(responseWriter, "failed to write html in the response", http.StatusInternalServerError)
		}
	}
	return fn
}
