package decorator

import (
	"github.com/go-playground/validator/v10"
	"github.com/oldhanasong/go-metrics-collector/internal/general/business_logic/validation/custom_func"
)

var customFuncMap = map[string]validator.Func{
	"mname":         custom_func.ValidateMetricName(),
	"mvalue_by":     custom_func.ValidateMetricValue(),
	"hashkey":       custom_func.ValidateHashkey(),
	"ratelimmin":    custom_func.ValidateRateLimitMin(),
	"hostname_port": custom_func.ValidateHostnamePort(),
}

func RegisteredCustomFunctionsValidator(validate *validator.Validate) (*validator.Validate, error) {
	var err error
	for tag, validateFunc := range customFuncMap {
		if err = validate.RegisterValidation(tag, validateFunc); err != nil {
			return nil, err
		}
	}

	return validate, nil
}
