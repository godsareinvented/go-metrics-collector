package interfaces

import (
	"context"
	"github.com/godsareinvented/go-metrics-collector/internal/general/dto"
)

type ClientInterface interface {
	Send(ctx context.Context, metricDTO dto.Metrics) error
	SendBatch(ctx context.Context, metrics *[]dto.Metrics) error
}
