package decorator

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/go-resty/resty/v2"
	"github.com/oldhanasong/go-metrics-collector/internal/agent/business_logic/config"
)

func HashCalculation(request *resty.Request) *resty.Request {
	if config.Configuration.HashKey == "" {
		return request
	}

	body := request.Body.([]byte)

	h := hmac.New(sha256.New, []byte(config.Configuration.HashKey))
	h.Write(body)
	dst := h.Sum(nil)

	request.SetHeader("HashSHA256", hex.EncodeToString(dst))

	return request
}
