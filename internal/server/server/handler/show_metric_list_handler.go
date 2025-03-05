package handler

import (
	"context"
	"github.com/godsareinvented/go-metrics-collector/internal/general/dto"
	"github.com/godsareinvented/go-metrics-collector/internal/server/config"
	"html/template"
	"net/http"
	"sort"
)

func ShowMetricList(ctx context.Context) http.HandlerFunc {
	fn := func(responseWriter http.ResponseWriter, request *http.Request) {
		requestCtx, cancel := GetCombinedContext(ctx, request.Context())
		defer cancel()

		metricDTOList, _ := config.Configuration.Repository.GetAllMetrics(requestCtx)

		sort.Slice(metricDTOList, func(i, j int) bool {
			return metricDTOList[i].MName < metricDTOList[j].MName
		})

		tmpl := template.Must(template.ParseFiles("internal/server/template/main_page.html"))
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
