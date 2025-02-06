package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/oldhanasong/go-metrics-collector/internal/server/config"
	"net/http"
)

type srBufferResponseWriter struct {
	http.ResponseWriter
	buffer *bytes.Buffer
}

func (w srBufferResponseWriter) Write(b []byte) (int, error) {
	return w.buffer.Write(b)
}

func SigningResponse(handlerFunc http.Handler) http.Handler {
	fn := func(responseWriter http.ResponseWriter, request *http.Request) {
		if "" == config.Configuration.HashKey {
			handlerFunc.ServeHTTP(responseWriter, request)
			return
		}

		recorder := srBufferResponseWriter{
			ResponseWriter: responseWriter,
			buffer:         &bytes.Buffer{},
		}
		handlerFunc.ServeHTTP(&recorder, request)

		if len(recorder.buffer.Bytes()) == 0 {
			_, _ = responseWriter.Write(recorder.buffer.Bytes())
			return
		}

		h := hmac.New(sha256.New, []byte(config.Configuration.HashKey))
		h.Write(recorder.buffer.Bytes())
		dst := h.Sum(nil)

		responseWriter.Header().Set("HashSHA256", hex.EncodeToString(dst))
		_, _ = responseWriter.Write(recorder.buffer.Bytes())
	}
	return http.HandlerFunc(fn)
}
