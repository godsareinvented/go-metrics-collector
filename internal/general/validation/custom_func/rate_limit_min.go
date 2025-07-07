package custom_func

import (
	"github.com/go-playground/validator/v10"
	"time"
)

func ValidateRateLimitMin() validator.Func {
	fn := func(fl validator.FieldLevel) bool {
		rateLimit := fl.Field().Int()
		pollInterval, ok := fl.Parent().FieldByName("PollInterval").Interface().(time.Duration)
		if !ok {
			return false
		}
		reportInterval, ok := fl.Parent().FieldByName("ReportInterval").Interface().(time.Duration)
		if !ok {
			return false
		}

		return float64(rateLimit) >= float64(reportInterval/pollInterval)
	}
	return fn
}
