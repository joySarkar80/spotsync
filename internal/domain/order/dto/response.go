package dto

type Response struct {
	ID         uint    `json:"id"`
	UserID     uint    `json:"user_id"`
	MangoID    uint    `json:"mango_id"`
	QuantityKg float64 `json:"quantity_kg"`
	TotalPrice float64 `json:"total_price"`
	Status     string  `json:"status"`
	OrderCode  string  `json:"order_code"`
	CreatedAt  string  `json:"created_at"`
}
