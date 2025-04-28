package handler

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strings"
)

func isContextError(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
}

func processValidationError(error error) (string, int) {
	errors := error.(validator.ValidationErrors)
	if strings.Contains(errors[0].Field(), "Name") {
		return "Metric name not passed on or incorrect", http.StatusNotFound
	}
	return "incorrect metric data", http.StatusBadRequest
}

// GetCombinedContext Получение комбинированного контекста, чтобы хендлер мог обработать завершение контекстов как приложения, так и запроса
func GetCombinedContext(serverCtx context.Context, requestCtx context.Context) (context.Context, context.CancelFunc) {
	wrappedRequestCtx, cancel := context.WithCancel(requestCtx)

	go func() {
		select {
		case <-serverCtx.Done():
			cancel()
		case <-wrappedRequestCtx.Done():
			return
		}
	}()

	return wrappedRequestCtx, cancel
}
