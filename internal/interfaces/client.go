package interfaces

import (
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
)

type Client interface {
	Send(metricDTO dto.Metrics) error
	SendBatch(metrics []dto.Metrics) error
}
