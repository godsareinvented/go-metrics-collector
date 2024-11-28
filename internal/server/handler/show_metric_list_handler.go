package handler

import (
	"context"
	"github.com/godsareinvented/go-metrics-collector/internal/config"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
	"html/template"
	"net/http"
	"sort"
)

func ShowMetricList(_ context.Context) http.HandlerFunc {
	fn := func(responseWriter http.ResponseWriter, request *http.Request) {
		metricDTOList, _ := config.Configuration.Repository.GetAllMetrics()

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
	return fn
}
