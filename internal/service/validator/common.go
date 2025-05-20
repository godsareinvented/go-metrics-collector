package validator

import (
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
