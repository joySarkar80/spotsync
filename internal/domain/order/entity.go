package order

import (
	"spotsync/internal/domain/order/dto"

	"gorm.io/gorm"
)

const (
	OrderPending   = "pending"
	OrderConfirmed = "confirmed"
	OrderCancelled = "cancelled"
)

type Order struct {
	gorm.Model
	UserID     uint    `json:"user_id" gorm:"not null"`
	MangoID    uint    `json:"mango_id" gorm:"not null"`
	QuantityKg float64 `json:"quantity_kg" gorm:"not null"`
	TotalPrice float64 `json:"total_price" gorm:"not null"`
	Status     string  `json:"status" gorm:"type:varchar(50);not null"`
	OrderCode  string  `json:"order_code" gorm:"uniqueIndex;not null"`
}

func (o *Order) ToResponse() *dto.Response {
	return &dto.Response{
		ID:         o.ID,
		UserID:     o.UserID,
		MangoID:    o.MangoID,
		QuantityKg: o.QuantityKg,
		TotalPrice: o.TotalPrice,
		Status:     o.Status,
		OrderCode:  o.OrderCode,
		CreatedAt:  o.CreatedAt.String(),
	}
}
