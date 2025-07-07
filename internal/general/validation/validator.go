package validation

import (
	"github.com/go-playground/validator/v10"
	"github.com/oldhanasong/go-metrics-collector/internal/general/validation/decorator"
	"sync"
)

var (
	v    *validator.Validate
	once sync.Once
)

func Validator() *validator.Validate {
	once.Do(func() {
		v, _ = decorator.RegisteredCustomFunctionsValidator(validator.New())
	})

	return v
}
