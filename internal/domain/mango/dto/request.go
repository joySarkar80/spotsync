package dto

type CreateRequest struct {
	Name        string  `json:"name" validate:"required,min=2,max=100"`
	Description string  `json:"description" validate:"max=1000"`
	Variety     string  `json:"variety" validate:"required,min=2,max=100"`
	PricePerKg  float64 `json:"price_per_kg" validate:"required,gt=0"`
	StockKg     float64 `json:"stock_kg" validate:"required,gt=0"`
	ImageURL    string  `json:"image_url" validate:"omitempty,url"`
}

type UpdateRequest struct {
	Name        string  `json:"name" validate:"omitempty,min=2,max=100"`
	Description string  `json:"description" validate:"omitempty,max=1000"`
	Variety     string  `json:"variety" validate:"omitempty,min=2,max=100"`
	PricePerKg  float64 `json:"price_per_kg" validate:"omitempty,gt=0"`
	StockKg     float64 `json:"stock_kg" validate:"omitempty,gt=0"`
	ImageURL    string  `json:"image_url" validate:"omitempty,url"`
}
