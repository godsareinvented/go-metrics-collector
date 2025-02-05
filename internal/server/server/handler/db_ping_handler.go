package handler

import (
	"context"
	"github.com/oldhanasong/go-metrics-collector/internal/server/config"
	"net/http"
)

func DbPing(ctx context.Context) http.HandlerFunc {
	fn := func(responseWriter http.ResponseWriter, request *http.Request) {
		combinedCtx, cancel := combineContext(ctx, request.Context())
		defer cancel()

		if ping, err := config.Configuration.Repository.PingStorage(combinedCtx); !ping || err != nil {
			http.Error(responseWriter, "failed to ping db", http.StatusInternalServerError)
			return
		}
		responseWriter.WriteHeader(http.StatusOK)
	}
	return fn
}
