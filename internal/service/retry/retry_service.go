package retry

import (
	"context"
	"errors"
	"fmt"
	"github.com/oldhanasong/go-metrics-collector/internal/interfaces"
	"time"
)

type (
	RetryService struct {
		strategy interfaces.RetryStrategy
	}
	UserDefinedCallback func() (error, bool)
)

var (
	ErrNonExecution = errors.New("failed to execute user defined callback after N attempts")
)

func (s *RetryService) DoWithRetry(ctx context.Context, callback UserDefinedCallback) error {
	err, needsRetry := callback()
	if !needsRetry {
		return err
	}

	s.strategy.Restore()
	for {
		delayDuration := s.strategy.NextDelayDuration()
		if 0 == delayDuration {
			return fmt.Errorf("%w: %v", ErrNonExecution, err)
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delayDuration):
		}

		err, needsRetry = callback()
		if !needsRetry {
			return nil
		}
	}
}

func NewInstance(strategy interfaces.RetryStrategy) RetryService {
	return RetryService{strategy: strategy}
}
