package mem_storage

import (
	"context"
	"fmt"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/interfaces"
	"maps"
	"slices"
)

type MemStorage struct {
	entityList map[string]dto.Metrics
}

var ()

func (memStorage *MemStorage) GetAll(_ context.Context) ([]dto.Metrics, error) {
	return slices.Collect(maps.Values(memStorage.entityList)), nil
}

func (memStorage *MemStorage) Get(_ context.Context, m dto.Metrics) (dto.Metrics, bool, error) {
	metrics, ok := memStorage.entityList[key(m)]
	return metrics, ok, nil
}

func (memStorage *MemStorage) Set(_ context.Context, m dto.Metrics) error {
	memStorage.entityList[key(m)] = m
	return nil
}

func key(metric dto.Metrics) string {
	return fmt.Sprintf("%s/%s", metric.MType, metric.ID)
}

func NewStorage() interfaces.Storage {
	return &MemStorage{
		entityList: make(map[string]dto.Metrics),
	}
}
