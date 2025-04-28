package custon_func

import (
	"github.com/go-playground/validator/v10"
	"sync"
)

var (
	v    *validator.Validate
	once sync.Once
)

func getValidator() *validator.Validate {
	once.Do(func() {
		v = validator.New()
	})

	return v
}
