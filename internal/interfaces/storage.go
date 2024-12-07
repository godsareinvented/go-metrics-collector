package interfaces

import (
	"context"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
)

type Storage interface {
	GetAll() ([]dto.Metrics, error)
	Get(m dto.Metrics) (dto.Metrics, bool, error)
	Set(m dto.Metrics) error
}

type StorageConnector interface {
	Ping(ctx context.Context) (bool, error)
	Close() error
}

type StorageConfigurator interface {
	Configure() error
}
