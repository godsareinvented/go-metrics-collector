package middleware

import (
	"github.com/godsareinvented/go-metrics-collector/internal/config"
	"net/http"
	"time"
)

type responseData struct {
	statusCode   int
	responseSize int
}

type loggingResponseWriter struct {
	responseWriter http.ResponseWriter
	responseData   responseData
}

func (lrw *loggingResponseWriter) WriteHeader(statusCode int) {
	lrw.responseData.statusCode = statusCode
	lrw.responseWriter.WriteHeader(statusCode)
}

func (lrw *loggingResponseWriter) Write(bytes []byte) (int, error) {
	lrw.responseData.responseSize += len(bytes)
	return lrw.responseWriter.Write(bytes)
}

func (lrw *loggingResponseWriter) Header() http.Header {
	return lrw.responseWriter.Header()
}

func WithLogging(handlerFunc http.Handler) http.Handler {
	fn := func(responseWriter http.ResponseWriter, request *http.Request) {
		start := time.Now()

		lrw := loggingResponseWriter{
			responseWriter: responseWriter,
			responseData: responseData{
				responseSize: 0,
			},
		}

		handlerFunc.ServeHTTP(&lrw, request)

		since := time.Since(start)

		config.Configuration.Logger.Sugar().Infoln(
			"uri", request.RequestURI,
			"method", request.Method,
			"status", lrw.responseData.statusCode,
			"duration", since,
			"size", lrw.responseData.responseSize,
		)
	}
	return http.HandlerFunc(fn)
}
