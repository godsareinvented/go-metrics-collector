package dto

import "time"

type (
	DelayCallback func(iteration uint) time.Duration

	RetryOptions struct {
		Attempts      uint
		DelayList     []time.Duration
		DelayCallback DelayCallback
	}
)
