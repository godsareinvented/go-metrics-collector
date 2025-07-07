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

// ValidateHostnamePort Функция для валидации имени хоста, которое может включать порт
func ValidateHostnamePort() validator.Func {
	fn := func(fl validator.FieldLevel) bool {
		if ":" == fl.Field().String() {
			return false
		}
		hostnamePort := parsedHostnamePortStruct(fl.Field().String())
		err := getValidator().Struct(hostnamePort)
		return nil == err
	}
	return fn
}

func parsedHostnamePortStruct(field string) HostnamePort {
	parts := strings.Split(field, ":")
	hostnamePort := HostnamePort{
		Hostname: parts[0],
	}
	if len(parts) == 2 {
		hostnamePort.Port = parts[1]
	}
	return hostnamePort
}
