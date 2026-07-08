package reservation

import "gorm.io/gorm"

// Full Reservation module (handler/service/repository) porer step-e banano hobe.
// Ekhon shudhu migration-er jonno struct ta declare kora holo, jate
// zone module-er available_spots query kaj kore.
type Reservation struct {
	gorm.Model
	UserID       uint   `gorm:"not null;index"`
	ZoneID       uint   `gorm:"not null;index"`
	LicensePlate string `gorm:"type:varchar(15);not null"`
	Status       string `gorm:"type:varchar(20);not null;default:'active'"` // active, completed, cancelled
}
