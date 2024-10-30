package mem_storage

import (
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"github.com/oldhanasong/go-metrics-collector/internal/interfaces"
	"maps"
	"slices"
)

type MemStorage struct {
	entityList map[string]dto.Metrics
}

func (memStorage *MemStorage) GetAll() ([]dto.Metrics, error) {
	return slices.Collect(maps.Values(memStorage.entityList)), nil
}

func (memStorage *MemStorage) Get(key string) (dto.Metrics, bool, error) {
	metrics, ok := memStorage.entityList[key]
	return metrics, ok, nil
}

func (memStorage *MemStorage) Set(key string, metric dto.Metrics) error {
	memStorage.entityList[key] = metric
	return nil
}

func NewInstance() interfaces.Storage {
	return &MemStorage{
		entityList: make(map[string]dto.Metrics),
	}
}
