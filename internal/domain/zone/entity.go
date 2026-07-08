package zone

import "gorm.io/gorm"

type ParkingZone struct {
	gorm.Model
	Name          string  `gorm:"type:varchar(150);not null"`
	Type          string  `gorm:"type:varchar(20);not null"` // general, ev_charging, covered
	TotalCapacity int     `gorm:"not null"`
	PricePerHour  float64 `gorm:"not null"`
}
