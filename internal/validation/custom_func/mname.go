package custom_func

import (
	"github.com/go-playground/validator/v10"
	"github.com/oldhanasong/go-metrics-collector/internal/util"
)

var (
	metricIdRegexp = util.LazyRegexCompile(`[a-zA-Z]+\d*`)
)

func ValidateMetricName() validator.Func {
	fn := func(fl validator.FieldLevel) bool {
		field := fl.Field().String()
		return metricIdRegexp().MatchString(field)
	}
	return fn
}
