package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/oldhanasong/go-metrics-collector/internal/server/config"
	"io"
	"net/http"
)

func CheckRequestSign(handlerFunc http.Handler) http.Handler {
	fn := func(responseWriter http.ResponseWriter, request *http.Request) {
		hash := request.Header.Get("HashSHA256")
		if "" == hash {
			handlerFunc.ServeHTTP(responseWriter, request)
			return
		}

		body, err := io.ReadAll(request.Body)
		if nil != err {
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
			return
		}

		request.Body = io.NopCloser(bytes.NewReader(body))

		h := hmac.New(sha256.New, []byte(config.Configuration.HashKey))
		h.Write(body)
		calculatedDst := h.Sum(nil)

		requestDst, err := hex.DecodeString(hash)
		if nil != err {
			http.Error(responseWriter, "Failed to decode request dst: "+err.Error(), http.StatusBadRequest)
			return
		}

		if !hmac.Equal(requestDst, calculatedDst) {
			http.Error(responseWriter, "The signatures don't match", http.StatusBadRequest)
			return
		}

		handlerFunc.ServeHTTP(responseWriter, request)
	}
	return http.HandlerFunc(fn)
}
