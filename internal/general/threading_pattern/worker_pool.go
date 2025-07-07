package threading_pattern

import "errors"

var (
	ErrNoArguments = errors.New("no arguments")
)

func InitWorkerPool[T interface{}](workerCount int, taskCh <-chan T, callback func(workerId int, task T)) error {
	if taskCh == nil || callback == nil {
		return ErrNoArguments
	}

	for i := 0; i < workerCount; i++ {
		go func(i int) {
			for task := range taskCh {
				callback(i, task)
			}
		}(i)
	}

	return nil
}
