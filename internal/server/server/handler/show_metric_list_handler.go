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

		metricList, err := config.Configuration.Repository.GetAllMetrics(requestCtx)
		if isContextError(err) {
			http.Error(responseWriter, "", http.StatusInternalServerError)
			return
		}

		sort.Slice(metricList, func(i, j int) bool {
			return metricList[i].MName < metricList[j].MName
		})

		responseWriter.Header().Set("Content-Type", "text/html")
		responseWriter.WriteHeader(http.StatusOK)

		tmpl := template.Must(template.ParseFiles("internal/server/template/main_page.html"))
		data := struct {
			Items []dto.Metrics
		}{
			Items: metricList,
		}

		err = tmpl.Execute(responseWriter, data)
		if nil != err {
			http.Error(responseWriter, "", http.StatusInternalServerError)
		}
	}
	return fn
}
