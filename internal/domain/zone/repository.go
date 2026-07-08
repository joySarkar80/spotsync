package zone

import (
	"errors"

	"gorm.io/gorm"
)

var ErrZoneNotFound = errors.New("parking zone not found")

type Repository interface {
	CreateZone(zone *ParkingZone) error
	GetAllZones() ([]ParkingZone, error)
	GetZoneByID(id uint) (*ParkingZone, error)
	CountActiveReservations(zoneID uint) (int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateZone(zone *ParkingZone) error {
	return r.db.Create(zone).Error
}

func (r *repository) GetAllZones() ([]ParkingZone, error) {
	var zones []ParkingZone
	if err := r.db.Order("id asc").Find(&zones).Error; err != nil {
		return nil, err
	}
	return zones, nil
}

func (r *repository) GetZoneByID(id uint) (*ParkingZone, error) {
	var zone ParkingZone
	err := r.db.First(&zone, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrZoneNotFound
		}
		return nil, err
	}
	return &zone, nil
}

// CountActiveReservations counts currently active reservations for a zone.
// Uses the "reservations" table directly (raw table reference) so this
// package doesn't need to import the reservation domain package.
func (r *repository) CountActiveReservations(zoneID uint) (int64, error) {
	var count int64
	err := r.db.Table("reservations").
		Where("zone_id = ? AND status = ? AND deleted_at IS NULL", zoneID, "active").
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
