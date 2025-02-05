package interfaces

import (
	"context"
	"github.com/oldhanasong/go-metrics-collector/internal/general/dto"
)

type Storage interface {
	GetAll(ctx context.Context) ([]dto.Metrics, error)
	Get(ctx context.Context, m dto.Metrics) (dto.Metrics, bool, error)
	Set(ctx context.Context, m dto.Metrics) error
	SetBatch(ctx context.Context, metricBatch []dto.Metrics) error
}

type StorageConnector interface {
	Ping(ctx context.Context) (bool, error)
	Close() error
}

type StorageConfigurator interface {
	Configure() error
}
