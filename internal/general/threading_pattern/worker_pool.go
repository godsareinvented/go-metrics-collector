package threading_pattern

import "errors"

var (
	ErrNoArguments = errors.New("no arguments")
)

func InitWorkerPool[T interface{}](workerCount int, taskCh <-chan T, callback func(workerId int, task T)) error {
	if nil == taskCh || nil == callback {
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
