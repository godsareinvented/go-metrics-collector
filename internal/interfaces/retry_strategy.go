package interfaces

import "time"

type RetryStrategy interface {
	NextDelayDuration() time.Duration
	Restore()
}
