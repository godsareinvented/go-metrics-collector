package strategy

import (
	"errors"
	"runtime"
)

var (
	ErrEmptyCollectedData = errors.New("empty collected data")

	emptyMemStatsStruct = runtime.MemStats{}
)

func isMemStatEmpty(memStats *runtime.MemStats) bool {
	return *memStats == emptyMemStatsStruct
}
