package decorator

import (
	"github.com/go-playground/validator/v10"
	"github.com/oldhanasong/go-metrics-collector/internal/general/validation/custom_func"
)

var customFuncMap = map[string]validator.Func{
	"mname":     custom_func.ValidateMetricName(),
	"mvalue_by": custom_func.ValidateMetricValue(),
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
