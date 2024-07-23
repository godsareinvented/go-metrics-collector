package dto

import "runtime"

type CollectedMetricData struct {
	MemStats    runtime.MemStats
	PollCount   int64
	RandomValue float64
}
