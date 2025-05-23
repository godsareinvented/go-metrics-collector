package handler

import (
	"context"
	"github.com/oldhanasong/go-metrics-collector/internal/config"
	"github.com/oldhanasong/go-metrics-collector/internal/interfaces"
	"github.com/oldhanasong/go-metrics-collector/internal/storage/postgres"
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

		storage, err := postgres.NewInstance(config.Configuration.DatabaseDSN)
		if err != nil {
			http.Error(responseWriter, "error when accessing the storage", http.StatusInternalServerError)
			return
		}

		connector, ok := storage.(interfaces.StorageConnector)
		if !ok {
			http.Error(responseWriter, "storage doesn't import the StorageConnector interface", http.StatusInternalServerError)
			return
		}

		if ping, err := connector.Ping(requestCtx); !ping || err != nil {
			http.Error(responseWriter, "failed to ping db", http.StatusInternalServerError)
			return
		}
		responseWriter.WriteHeader(http.StatusOK)
	}
	return fn
}
