package middleware

import (
	"bytes"
	"compress/gzip"
	"github.com/oldhanasong/go-metrics-collector/internal/server/config"
	"net/http"
	"slices"
	"strings"
)

type bufferResponseWriter struct {
	http.ResponseWriter
	buffer     bytes.Buffer
	statusCode int
}

func (w *bufferResponseWriter) Write(b []byte) (int, error) {
	return w.buffer.Write(b)
}

func (w *bufferResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

func GzipResponseCompressing(handlerFunc http.Handler) http.Handler {
	fn := func(responseWriter http.ResponseWriter, request *http.Request) {
		if !supportsGzip(request) {
			handlerFunc.ServeHTTP(responseWriter, request)
			return
		}

		recorder := bufferResponseWriter{
			ResponseWriter: responseWriter,
			buffer:         bytes.Buffer{},
			statusCode:     http.StatusOK,
		}
		handlerFunc.ServeHTTP(&recorder, request)

		if !isCompressionNeed(recorder) {
			responseWriter.WriteHeader(recorder.statusCode)
			if _, err := responseWriter.Write(recorder.buffer.Bytes()); err != nil {
				http.Error(responseWriter, "failed to write metric in the response", http.StatusInternalServerError)
			}
			return
		}

		responseWriter.Header().Set("Content-Encoding", "gzip")
		responseWriter.WriteHeader(recorder.statusCode)

		gzipWriter, err := gzip.NewWriterLevel(responseWriter, gzip.BestSpeed)
		if err != nil {
			http.Error(responseWriter, "failed to declare gzip writer", http.StatusInternalServerError)
			return
		}
		defer func(gzipWriter *gzip.Writer) {
			if err = gzipWriter.Close(); err != nil {
				http.Error(responseWriter, "failed to close gzip writer", http.StatusInternalServerError)
			}
		}(gzipWriter)

		if _, err = gzipWriter.Write(recorder.buffer.Bytes()); err != nil {
			http.Error(responseWriter, "failed to compress data via gzip writer", http.StatusInternalServerError)
		}
	}
	return http.HandlerFunc(fn)
}

func supportsGzip(request *http.Request) bool {
	return clientAcceptEncoding(request)
}

func isCompressionNeed(responseWriter bufferResponseWriter) bool {
	return appropriateContentType(responseWriter) && appropriateContentSize(responseWriter)
}

// Есть формат x-gzip, поэтому простая проверка на strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") не подойдёт (?).
func clientAcceptEncoding(request *http.Request) bool {
	formatList := strings.Split(request.Header.Get("Accept-Encoding"), ", ")
	return slices.Contains(formatList, "gzip")
}

func appropriateContentType(responseWriter bufferResponseWriter) bool {
	contentType := strings.Split(responseWriter.Header().Get("Content-Type"), ";")[0]
	return slices.Contains(config.Configuration.GzipAcceptedContentTypes, contentType)
}

func appropriateContentSize(responseWriter bufferResponseWriter) bool {
	return responseWriter.buffer.Len() >= config.Configuration.GzipMinContentLength
}
