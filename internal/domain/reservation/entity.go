package reservation

import (
	"spotsync/internal/domain/user"
	"spotsync/internal/domain/zone"

	"gorm.io/gorm"
)

type Reservation struct {
	gorm.Model
	UserID       uint   `gorm:"not null;index"`
	ZoneID       uint   `gorm:"not null;index"`
	LicensePlate string `gorm:"type:varchar(15);not null"`
	Status       string `gorm:"type:varchar(20);not null;default:'active'"` // active, completed, cancelled

	User user.User        `gorm:"foreignKey:UserID"`
	Zone zone.ParkingZone `gorm:"foreignKey:ZoneID"`
}
