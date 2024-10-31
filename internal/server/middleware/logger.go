package middleware

import (
	"github.com/godsareinvented/go-metrics-collector/internal/config"
	"net/http"
	"time"
)

type responseData struct {
	statusCode   *int
	responseSize *int
}

type loggingResponseWriter struct {
	http.ResponseWriter
	responseData responseData
}

func (lrw *loggingResponseWriter) Write(bytes []byte) (int, error) {
	size, err := lrw.ResponseWriter.Write(bytes)
	if nil == err {
		*lrw.responseData.responseSize += size
	}
	return size, err
}

func (lrw *loggingResponseWriter) WriteHeader(statusCode int) {
	*lrw.responseData.statusCode = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}

func WithLogging(handlerFunc http.Handler) http.Handler {
	fn := func(responseWriter http.ResponseWriter, request *http.Request) {
		start := time.Now()

		statusCode := http.StatusOK
		responseSize := 0
		lrw := loggingResponseWriter{
			ResponseWriter: responseWriter,
			responseData: responseData{
				responseSize: &responseSize,
				statusCode:   &statusCode,
			},
		}

		handlerFunc.ServeHTTP(&lrw, request)

		since := time.Since(start)

		config.Configuration.Logger.Sugar().Infoln(
			"uri", request.RequestURI,
			"method", request.Method,
			"status", *lrw.responseData.statusCode,
			"duration", since,
			"size", *lrw.responseData.responseSize,
		)
	}
	return http.HandlerFunc(fn)
}
