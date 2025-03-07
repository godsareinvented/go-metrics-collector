package interfaces

import (
	"github.com/godsareinvented/go-metrics-collector/internal/general/dto"
)

type ClientInterface interface {
	Send(metricDTO dto.Metrics) error
	SendBatch(metrics *[]dto.Metrics) error
}
