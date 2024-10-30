package interfaces

import "github.com/oldhanasong/go-metrics-collector/internal/dto"

type Storage interface {
	GetAll() ([]dto.Metrics, error)
	Get(key string) (dto.Metrics, bool, error)
	Set(key string, metric dto.Metrics) error
}
