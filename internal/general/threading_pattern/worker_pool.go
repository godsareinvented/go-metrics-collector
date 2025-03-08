package threading_pattern

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrNoArguments = errors.New("no arguments")
)

func InitWorkerPool[T interface{}](
	ctx context.Context,
	workerCount int,
	taskCh <-chan T,
	callback func(workerId int, task T) error,
) (*sync.WaitGroup, error) {
	wg := sync.WaitGroup{}
	if nil == taskCh || nil == callback {
		return &wg, ErrNoArguments
	}

	for i := 0; i < workerCount; i++ {
		wg.Add(1)

		go func(ctx context.Context, i int, wg *sync.WaitGroup) {
			defer wg.Done()

			var err error
			for task := range taskCh {
				select {
				case <-ctx.Done():
					return
				default:
					err = callback(i, task)
					if nil != err {
						return
					}
				}
			}
		}(ctx, i, &wg)
	}

	return &wg, nil
}
