package dto

type Response struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Variety     string  `json:"variety"`
	PricePerKg  float64 `json:"price_per_kg"`
	StockKg     float64 `json:"stock_kg"`
	ImageURL    string  `json:"image_url,omitempty"`
	CreatedAt   string  `json:"created_at"`
}
