package middleware

import (
	"bytes"
	"compress/gzip"
	"github.com/godsareinvented/go-metrics-collector/internal/config"
	"io"
	"net/http"
	"slices"
	"strings"
)

type bufferResponseWriter struct {
	http.ResponseWriter
	buffer     *bytes.Buffer
	statusCode *int
}

func (w bufferResponseWriter) Write(b []byte) (int, error) {
	return w.buffer.Write(b)
}

func (w bufferResponseWriter) WriteHeader(statusCode int) {
	*w.statusCode = statusCode
}

func GzipResponseCompressing(handlerFunc http.Handler) http.Handler {
	fn := func(responseWriter http.ResponseWriter, request *http.Request) {
		if !supportsGzip(request) {
			handlerFunc.ServeHTTP(responseWriter, request)
			return
		}

		statusCode := http.StatusOK
		recorder := bufferResponseWriter{
			ResponseWriter: responseWriter,
			buffer:         &bytes.Buffer{},
			statusCode:     &statusCode,
		}
		handlerFunc.ServeHTTP(&recorder, request)

		if !isCompressionNeed(recorder) {
			_, _ = responseWriter.Write(recorder.buffer.Bytes())
			return
		}

		responseWriter.Header().Set("Content-Encoding", "gzip")
		responseWriter.WriteHeader(*recorder.statusCode)

		gzipWriter, err := gzip.NewWriterLevel(responseWriter, gzip.BestSpeed)
		if err != nil {
			_, _ = io.WriteString(responseWriter, err.Error())
			return
		}
		defer gzipWriter.Close()

		_, _ = gzipWriter.Write(recorder.buffer.Bytes())
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
