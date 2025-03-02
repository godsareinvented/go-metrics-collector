package custom_func

import "github.com/go-playground/validator/v10"

// ValidateMetricNameFunc Функция для валидации названия метрики
func ValidateMetricNameFunc() validator.Func {
	return getRegExpFunc(MetricName)
}
