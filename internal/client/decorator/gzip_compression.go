package decorator

import (
	"bytes"
	"compress/gzip"
	"github.com/go-resty/resty"
	"github.com/godsareinvented/go-metrics-collector/internal/config"
	"strconv"
)

func GzipCompress(request *resty.Request) *resty.Request {
	body := request.Body.([]byte)
	if len(body) >= config.Configuration.GzipMinContentLength {
		buf := getCompressBodyBuffer(body)
		compressedBody := buf.Bytes()

		request.SetBody(compressedBody)
		request.Header.Set("Content-Length", strconv.Itoa(len(compressedBody)))
		request.Header.Set("Accept-Encoding", "gzip")
		request.Header.Set("Content-Encoding", "gzip")
	}

	return request
}

func getCompressBodyBuffer(body []byte) *bytes.Buffer {
	buffer := new(bytes.Buffer)

	gzipWriter, err := gzip.NewWriterLevel(buffer, gzip.BestSpeed)
	if nil != err {
		panic(err)
	}
	defer gzipWriter.Close()

	_, _ = gzipWriter.Write(body)

	return buffer
}
