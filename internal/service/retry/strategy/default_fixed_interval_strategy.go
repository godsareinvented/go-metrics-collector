package strategy

import (
	"github.com/godsareinvented/go-metrics-collector/internal/interfaces"
	"time"
)

// DefaultFixedIntervalStrategy todo: delay
type DefaultFixedIntervalStrategy struct {
	attempts              int
	timeIntervalList      [3]time.Duration
	timeIntervalGenerator func() time.Duration
	currentAttempt        int
}

func (s *DefaultFixedIntervalStrategy) NextDelayDuration() time.Duration {
	if s.currentAttempt >= s.attempts || s.currentAttempt >= len(s.timeIntervalList) {
		return 0
	}

	result := s.timeIntervalList[s.currentAttempt]
	s.currentAttempt += 1

	return result
}

func (s *DefaultFixedIntervalStrategy) Restore() {
	s.currentAttempt = 0
}

func NewDefaultFixedIntervalStrategy() interfaces.RetryStrategy {
	return &DefaultFixedIntervalStrategy{
		attempts:              3,
		timeIntervalList:      [3]time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second},
		timeIntervalGenerator: nil,
		currentAttempt:        0,
	}
}
