package dto

// Metrics todo: переписать валидацию.
type Metrics struct {
	ID    string   `json:"id"               validate:"omitempty,required"`
	MType string   `json:"type"             validate:"required,contains=gauge|contains=counter"`
	MName string   `json:"name"             validate:"required,alpha"`
	Delta *int64   `json:"delta,omitempty"  validate:"omitempty,required"`
	Value *float64 `json:"value,omitempty"  validate:"omitempty,required"`
}
