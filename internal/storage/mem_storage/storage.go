package mem_storage

import (
	"github.com/godsareinvented/go-metrics-collector/internal/interfaces"
)

type MemStorage struct {
	storage map[string]interface{}
}

func (memStorage *MemStorage) Get(key string) interface{} {
	if value, ok := memStorage.storage[key]; ok {
		return value
	}
	return nil
}

func (memStorage *MemStorage) Set(key string, value interface{}) {
	memStorage.storage[key] = value
}

func NewInstance() interfaces.Storage {
	return &MemStorage{make(map[string]interface{})}
}
