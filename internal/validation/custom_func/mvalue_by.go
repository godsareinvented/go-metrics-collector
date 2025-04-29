package custom_func

import (
	"github.com/go-playground/validator/v10"
	"github.com/oldhanasong/go-metrics-collector/internal/dictionary"
)

// ValidateMetricValue Валидация значения метрики согласованно с её типом
// Синтаксис: mvalue_by=<metric type>
//
//	Пример: mvalue_by=gauge
func ValidateMetricValue() validator.Func {
	fn := func(fl validator.FieldLevel) bool {
		field := fl.Field().String()
		param := fl.Param()

		var errs error
		switch param {
		case dictionary.GaugeMetricType:
			errs = getValidator().Var(field, `required,numeric`)
		case dictionary.CounterMetricType:
			errs = getValidator().Var(field, `required,number`)
		default:
			return false
		}

		if errs != nil {
			return false
		}
		return true
	}
	return fn
}
