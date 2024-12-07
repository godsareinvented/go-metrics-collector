package interfaces

import (
	"context"
	"github.com/godsareinvented/go-metrics-collector/internal/dto"
)

type StorageInterface interface {
	GetAll() ([]dto.Metrics, error)
	GetByID(ID string, mType string) (dto.Metrics, bool, error)
	GetByName(mName string, mType string) (dto.Metrics, bool, error)
	Save(metric dto.Metrics) (string, error)
	GetGeneratedID(metric dto.Metrics) string
}

type StorageConnectorInterface interface {
	CloseConnect() error
	Ping(ctx context.Context) (bool, error)
}

type StorageConfiguratorInterface interface {
	Configure() error
}
