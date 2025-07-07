package metric

import (
	"github.com/oldhanasong/go-metrics-collector/internal/general/business_logic/dictionary"
	"github.com/oldhanasong/go-metrics-collector/internal/general/validation"
	"go.uber.org/multierr"
)

func ValidateMetricType(MType string) error {
	return validation.Validator().Var(MType, `oneof=`+dictionary.GaugeMetricType+` `+dictionary.CounterMetricType)
}

func ValidateMetricName(MName string) error {
	return validation.Validator().Var(MName, `required,mname`)
}

func ValidateMetricValue(MType, MValue string) error {
	return validation.Validator().Var(MValue, `required,mvalue_by=`+MType)
}

func ValidateMetricValues(MType, MName, MValue string) error {
	return multierr.Combine(ValidateMetricType(MType), ValidateMetricName(MName), ValidateMetricValue(MType, MValue))
}

func ValidateAbridgedMetricValues(MType, MName string) error {
	return multierr.Combine(ValidateMetricType(MType), ValidateMetricName(MName))
}
