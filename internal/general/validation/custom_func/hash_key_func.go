package custom_func

import "github.com/go-playground/validator/v10"

// ValidateHashKeyFunc Функция для валидации строки, содержащей: буквы, числа, знаки
func ValidateHashKeyFunc() validator.Func {
	return getRegExpFunc(HashKeyRegexp)
}
