package retry

import (
	"context"
	"errors"
	"fmt"
	"github.com/oldhanasong/go-metrics-collector/internal/dto"
	"time"
)

type (
	UserDefinedCallback func() (error, bool)
)

var (
	ErrNonExecution   = errors.New("failed to execute user defined callback after N attempts")
	ErrInvalidOpts    = errors.New("invalid retry options. Need to pass DelayList or DelayCallback")
	ErrIncorrectDelay = errors.New("invalid retry options. Incorrect retry delay")
)

func DoWithRetry(ctx context.Context, retryOpts dto.RetryOptions, callback UserDefinedCallback) error {
	lastCallbackErr, needsRetry := callback()
	if !needsRetry {
		return lastCallbackErr
	}

	var iter uint = 0
	for {
		if iter >= retryOpts.Attempts {
			return fmt.Errorf("%w: %v", ErrNonExecution, lastCallbackErr)
		}

		delay, err := getNextDelay(iter, retryOpts)
		if nil != err {
			return err
		}

		if 0 == delay {
			if 0 == iter {
				return ErrIncorrectDelay
			}
			return fmt.Errorf("%w: %v", ErrNonExecution, lastCallbackErr)
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}

		lastCallbackErr, needsRetry = callback()
		if !needsRetry {
			return nil
		}

		iter += 1
	}
}

func getNextDelay(iteration uint, retryOpts dto.RetryOptions) (time.Duration, error) {
	if iteration >= retryOpts.Attempts {
		return 0, nil
	}

	if nil != retryOpts.DelayList {
		return getNextDelayFromList(iteration, retryOpts), nil
	}

	if nil != retryOpts.DelayCallback {
		return getNextDelayWithCallback(iteration, retryOpts), nil
	}

	return 0, ErrInvalidOpts
}

func getNextDelayFromList(iteration uint, retryOpts dto.RetryOptions) time.Duration {
	if iteration >= uint(len(retryOpts.DelayList)) {
		return 0
	}

	return retryOpts.DelayList[iteration]
}

func getNextDelayWithCallback(iteration uint, retryOpts dto.RetryOptions) time.Duration {
	return retryOpts.DelayCallback(iteration)
}
