package dto

type CreateRequest struct {
	MangoID    uint    `json:"mango_id" validate:"required"`
	QuantityKg float64 `json:"quantity_kg" validate:"required,gt=0"`
}
