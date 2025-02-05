package mem_storage

import (
	"github.com/oldhanasong/go-metrics-collector/internal/server/interfaces"
)

type MemStorageConfigurator struct{}

func (c *MemStorageConfigurator) Configure() error {
	return nil
}

func NewConfigurator() interfaces.StorageConfigurator {
	return &MemStorageConfigurator{}
}
