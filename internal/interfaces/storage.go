package interfaces

import (
	"context"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
)

type StorageInterface interface {
	GetAll(ctx context.Context) ([]dto.Metrics, error)
	GetByID(ctx context.Context, ID string, mType string) (dto.Metrics, bool, error)
	GetByName(ctx context.Context, mName string, mType string) (dto.Metrics, bool, error)
	Save(ctx context.Context, metric dto.Metrics) (string, error)
	GetGeneratedID(ctx context.Context, metric dto.Metrics) (string, error)
}

type StorageConnectorInterface interface {
	CloseConnect() error
	Ping(ctx context.Context) (bool, error)
}

type StorageConfiguratorInterface interface {
	Configure() error
}
