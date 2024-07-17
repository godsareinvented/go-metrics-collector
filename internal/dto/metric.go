package dto

type Metric struct {
	Type  string      `validate:"required,contains=gauge|contains=counter"`
	Name  string      `validate:"required,alpha"`
	Value interface{} `validate:"required,number,gt=0"`
}
