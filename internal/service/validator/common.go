package validator

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/oldhanasong/go-metrics-collector/internal/validation/decorator"
	"sync"
)

var (
	v    *validator.Validate
	once sync.Once
)

func GetValidator() *validator.Validate {
	once.Do(func() {
		v, _ = decorator.GetRegisteredCustomFunctionsValidator(validator.New())
	})

	return v
}

func WrappedError(errors ...error) error {
	var err error
	var errWrapped error
	for _, err = range errors {
		if err == nil {
			continue
		}

		if errWrapped == nil {
			errWrapped = err
			continue
		}

		errWrapped = fmt.Errorf("%w: %w", err, errWrapped)
	}

	return errWrapped
}
