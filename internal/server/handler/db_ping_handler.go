package handler

import (
	"context"
	"github.com/godsareinvented/go-metrics-collector/internal/config"
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

		_, err := config.Configuration.Repository.PingStorage(requestCtx)
		if nil != err {
			responseWriter.WriteHeader(http.StatusInternalServerError)
		}
		responseWriter.WriteHeader(http.StatusOK)
	}
	return fn
}
