package custom_func

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

const (
	HashKeyRegexp = `^[a-zA-Z0-9\-]+$`    // Пример: kE5-y3R-Y3B-yBR
	MetricName    = `^[A-Z][a-zA-Z]*\d*$` // Пример: RandomValue, CPUutilization0
)

var (
	validate = validator.New(validator.WithRequiredStructEnabled())
)

func getRegExpFunc(regex string) validator.Func {
	fn := func(fl validator.FieldLevel) bool {
		matched, err := regexp.MatchString(regex, fl.Field().String())
		if nil != err {
			return false
		}
		return matched
	}
	return fn
}
