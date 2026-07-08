package mango

import (
	"spotsync/internal/domain/mango/dto"

	"gorm.io/gorm"
)

type Mango struct {
	gorm.Model
	Name        string  `json:"name" gorm:"type:varchar(100);not null"`
	Description string  `json:"description" gorm:"type:text"`
	Variety     string  `json:"variety" gorm:"type:varchar(100);not null"`
	PricePerKg  float64 `json:"price_per_kg" gorm:"not null"`
	StockKg     float64 `json:"stock_kg" gorm:"not null"`
	ImageURL    string  `json:"image_url" gorm:"type:varchar(500)"`
}

func (m *Mango) ToResponse() *dto.Response {
	return &dto.Response{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		Variety:     m.Variety,
		PricePerKg:  m.PricePerKg,
		StockKg:     m.StockKg,
		ImageURL:    m.ImageURL,
		CreatedAt:   m.CreatedAt.String(),
	}
}
