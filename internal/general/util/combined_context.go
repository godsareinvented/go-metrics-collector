package util

import "context"

// CombineContexts Устанавливает зависимость побочных контекстов от состояния главного контекста
// Важно: после выполнения работы побочные контексты должны быть отменены. Иначе произойдёт утечка памяти из-за зависания горутин!
func CombineContexts(mainCtx, sideCtx context.Context) (context.Context, context.CancelFunc) {
	sideCtx, cancel := context.WithCancel(sideCtx)

	go func(mainCtx, ctx context.Context) {
		select {
		case <-mainCtx.Done():
			cancel()
		case <-ctx.Done():
			return
		}
	}(mainCtx, sideCtx)

	return sideCtx, cancel
}
