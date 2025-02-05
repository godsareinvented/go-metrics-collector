package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
)

func GzipRequestDecompressing(handlerFunc http.Handler) http.Handler {
	fn := func(responseWriter http.ResponseWriter, request *http.Request) {
		if !isRequestCompressed(request) {
			handlerFunc.ServeHTTP(responseWriter, request)
			return
		}

		gz, err := gzip.NewReader(request.Body)
		if err != nil {
			http.Error(responseWriter, "failed to declare gzip reader", http.StatusInternalServerError)
			return
		}
		defer func(gz *gzip.Reader) {
			if err = gz.Close(); err != nil {
				http.Error(responseWriter, "failed to close gzip reader", http.StatusInternalServerError)
			}
		}(gz)

		body, err := io.ReadAll(gz)
		if err != nil {
			http.Error(responseWriter, "failed to decompress data via gzip writer", http.StatusInternalServerError)
			return
		}

		request.Body = io.NopCloser(bytes.NewBuffer(body))
		request.ContentLength = int64(len(body))

		handlerFunc.ServeHTTP(responseWriter, request)
	}
	return http.HandlerFunc(fn)
}

func isRequestCompressed(request *http.Request) bool {
	contentEncoding := request.Header.Get("Content-Encoding")
	return contentEncoding == "gzip"
}
