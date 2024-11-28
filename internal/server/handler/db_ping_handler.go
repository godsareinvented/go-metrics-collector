package handler

import (
	"context"
	"github.com/oldhanasong/go-metrics-collector/internal/config"
	"net/http"
)

func DbPing(ctx context.Context) http.HandlerFunc {
	fn := func(responseWriter http.ResponseWriter, request *http.Request) {
		// Комбинированный контекст, чтобы хендлер мог обработать завершение контекстов как приложения, так и запроса
		requestCtx, cancel := context.WithCancel(request.Context())
		defer cancel()

		go func() {
			<-ctx.Done()
			cancel()
		}()

		if ping, err := config.Configuration.Repository.PingStorage(requestCtx); !ping || err != nil {
			http.Error(responseWriter, "failed to ping db", http.StatusInternalServerError)
			return
		}
		responseWriter.WriteHeader(http.StatusOK)
	}
	return fn
}
