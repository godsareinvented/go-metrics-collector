package dto

type Metric struct {
	Type  string  `json:"type" validate:"required,contains=gauge|contains=counter"`
	Name  string  `json:"name" validate:"required,alpha"`
	Delta int64   `json:"delta" validate:"omitempty,required"` // todo: переписать валидацию.
	Value float64 `json:"value" validate:"omitempty,required"`
}
