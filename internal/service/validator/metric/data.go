package metric

import (
	"github.com/oldhanasong/go-metrics-collector/internal/dictionary"
	"github.com/oldhanasong/go-metrics-collector/internal/service/validator"
	"github.com/oldhanasong/go-metrics-collector/internal/util"
)

func ValidateMetricType(MType string) error {
	return validator.GetValidator().Var(MType, `oneof=`+dictionary.GaugeMetricType+` `+dictionary.CounterMetricType)
}

func ValidateMetricName(MName string) error {
	return validator.GetValidator().Var(MName, `required,alphanum`)
}

func ValidateMetricValue(MType, MValue string) error {
	return validator.GetValidator().Var(MValue, `required,mvalue_by=`+MType)
}

func ValidateMetricValues(MType, MName, MValue string) error {
	return util.WrappedErrs(ValidateMetricType(MType), ValidateMetricName(MName), ValidateMetricValue(MType, MValue))
}

func ValidateAbridgedMetricValues(MType, MName string) error {
	return util.WrappedErrs(ValidateMetricType(MType), ValidateMetricName(MName))
}
