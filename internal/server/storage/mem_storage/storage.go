package mem_storage

import (
	"context"
	"fmt"
	"github.com/oldhanasong/go-metrics-collector/internal/general/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/server/interfaces"
	"maps"
	"slices"
	"sync"
)

type MemStorage struct {
	mu         sync.RWMutex
	entityList map[string]dto.Metrics
}

func (memStorage *MemStorage) GetAll(_ context.Context) ([]dto.Metrics, error) {
	memStorage.mu.RLock()
	defer memStorage.mu.RUnlock()

	return slices.Collect(maps.Values(memStorage.entityList)), nil
}

func (memStorage *MemStorage) Get(_ context.Context, m dto.Metrics) (dto.Metrics, bool, error) {
	memStorage.mu.RLock()
	defer memStorage.mu.RUnlock()

	metrics, ok := memStorage.entityList[key(m)]
	return metrics, ok, nil
}

func (memStorage *MemStorage) Set(_ context.Context, m dto.Metrics) error {
	memStorage.mu.Lock()
	defer memStorage.mu.Unlock()

	memStorage.entityList[key(m)] = m
	return nil
}

func (memStorage *MemStorage) SetBatch(_ context.Context, metricList []dto.Metrics) error {
	memStorage.mu.Lock()
	defer memStorage.mu.Unlock()

	for _, m := range metricList {
		memStorage.entityList[key(m)] = m
	}
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
