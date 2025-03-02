package custom_func

import (
	"github.com/go-playground/validator/v10"
	"strings"
)

type (
	HostnamePort struct {
		Hostname string `validate:"omitempty,hostname"`
		Port     string `validate:"omitempty,numeric,min=1,max=65535"`
	}
)

// ValidateHostnamePortFunc Функция для валидации имени хоста, которое может включать порт
func ValidateHostnamePortFunc() validator.Func {
	fn := func(fl validator.FieldLevel) bool {
		if ":" == fl.Field().String() {
			return false
		}
		hostnamePort := getParsedHostnamePortStruct(fl.Field().String())
		err := validate.Struct(hostnamePort)
		return nil == err
	}
	return fn
}

func getParsedHostnamePortStruct(field string) HostnamePort {
	parts := strings.Split(field, ":")
	hostnamePort := HostnamePort{
		Hostname: parts[0],
	}
	if len(parts) == 2 {
		hostnamePort.Port = parts[1]
	}
	return hostnamePort
}
