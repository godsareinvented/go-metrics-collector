package interfaces

import (
	"github.com/go-resty/resty/v2"
	"github.com/oldhanasong/go-metrics-collector/internal/general/dto"
)

type (
	Decorator func(request *resty.Request) *resty.Request

	Client interface {
		Send(metricDTO dto.Metrics) error
		SendBatch(metrics []dto.Metrics) error
		Use(decorator Decorator)
	}
)
