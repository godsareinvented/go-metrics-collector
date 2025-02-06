package request

import (
	"github.com/go-resty/resty"
	"github.com/godsareinvented/go-metrics-collector/internal/agent/client/decorator"
)

func GetWrappedRequest(request *resty.Request) *resty.Request {
	return decorator.GzipCompress(decorator.HashCalculation(request))
}
