package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"net/http"
)

func GzipRequestDecompressing(handlerFunc http.Handler) http.Handler {
	fn := func(responseWriter http.ResponseWriter, request *http.Request) {
		if !isRequestCompressed(request) {
			handlerFunc.ServeHTTP(responseWriter, request)
			return
		}

		gz, err := gzip.NewReader(request.Body)
		if nil != err {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
			return
		}
		defer gz.Close()

		body, err := io.ReadAll(gz)
		if nil != err {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
			return
		}

		request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		request.ContentLength = int64(len(body))

		handlerFunc.ServeHTTP(responseWriter, request)
	}
	return http.HandlerFunc(fn)
}

func isRequestCompressed(request *http.Request) bool {
	contentEncoding := request.Header.Get("Content-Encoding")
	return "gzip" == contentEncoding
}
