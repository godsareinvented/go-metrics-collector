package custom_func

import (
	"github.com/go-playground/validator/v10"
	"github.com/oldhanasong/go-metrics-collector/internal/general/util"
)

var (
	hashKeyRegexp = util.LazyRegexCompile(`^[a-zA-Z0-9\-]+$`) // Пример: kE5-y3R-Y3B-yBR
)

func ValidateHashkey() validator.Func {
	fn := func(fl validator.FieldLevel) bool {
		field := fl.Field().String()
		return hashKeyRegexp().MatchString(field)
	}
	return fn
}
