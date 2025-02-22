package interfaces

import (
	"github.com/godsareinvented/go-metrics-collector/internal/general/dto"
)

type Client interface {
	Send(metricDTO dto.Metrics) error
	SendBatch(metrics *[]dto.Metrics) error
}
