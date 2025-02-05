package interfaces

import (
	"github.com/godsareinvented/go-metrics-collector/internal/general/dto"
)

type PermanentStorage interface {
	Import() ([]dto.Metrics, error)
	Export([]dto.Metrics) error
	Close()
}
