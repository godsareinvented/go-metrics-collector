package handler

import (
	"context"
	"errors"
	"github.com/oldhanasong/go-metrics-collector/internal/general/util"
	"github.com/oldhanasong/go-metrics-collector/internal/server/config"
	"github.com/oldhanasong/go-metrics-collector/internal/server/repository"
	"net/http"
)

func DbPing(ctx context.Context) http.HandlerFunc {
	fn := func(responseWriter http.ResponseWriter, request *http.Request) {
		ctx, cancel := util.CombineContexts(ctx, request.Context())
		defer cancel()

		ping, err := config.Configuration.Repository.PingStorage(ctx)
		if errors.Is(err, repository.ErrDontImplementConnectorInterface) {
			http.Error(responseWriter, "storage unavailable", http.StatusServiceUnavailable)
			return
		}
		if err != nil || !ping {
			http.Error(responseWriter, "failed to ping db", http.StatusInternalServerError)
			return
		}
		responseWriter.WriteHeader(http.StatusOK)
	}
	return fn
}
