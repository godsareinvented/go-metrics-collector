package dto

import "github.com/godsareinvented/go-metrics-collector/internal/constraint"

type Metric[Num constraint.Numeric] struct {
	Type  string `json:"type" validate:"required,contains=gauge|contains=counter"`
	Name  string `json:"name" validate:"required,alpha"`
	Value Num    `json:"value" validate:"omitnil,required,numeric,gt=0"` // todo: Как будет работать?
}
