package handler

import (
	"context"
	"github.com/godsareinvented/go-metrics-collector/internal/server/config"
	"net/http"
)

func DbPing(ctx context.Context) http.HandlerFunc {
	fn := func(responseWriter http.ResponseWriter, request *http.Request) {
		requestCtx, cancel := GetCombinedContext(ctx, request.Context())
		defer cancel()

		_, err := config.Configuration.Repository.PingStorage(requestCtx)
		if nil != err {
			responseWriter.WriteHeader(http.StatusInternalServerError)
		}
		responseWriter.WriteHeader(http.StatusOK)
	}
	return fn
}
