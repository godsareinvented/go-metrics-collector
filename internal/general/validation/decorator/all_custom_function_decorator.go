package decorator

import (
	"github.com/go-playground/validator/v10"
	"github.com/godsareinvented/go-metrics-collector/internal/general/validation/custom_func"
)

var customFuncMap = map[string]validator.Func{
	"hash_key":      custom_func.ValidateHashKeyFunc(),
	"hostname_port": custom_func.ValidateHostnamePortFunc(),
	"metric_name":   custom_func.ValidateMetricNameFunc(),
}

func GetRegisteredCustomFunctionsValidator(validate *validator.Validate) (*validator.Validate, error) {
	var err error
	for tag, validateFunc := range customFuncMap {
		err = validate.RegisterValidation(tag, validateFunc)
		if err != nil {
			return nil, err
		}
	}

	return validate, nil
}
