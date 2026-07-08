package reservation

import (
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"spotsync/internal/domain/zone"
)

var (
	ErrZoneNotFound        = errors.New("parking zone not found")
	ErrZoneFull            = errors.New("parking zone is full")
	ErrReservationNotFound = errors.New("reservation not found")
	ErrForbidden           = errors.New("you are not allowed to cancel this reservation")
)

type Repository interface {
	CreateReservation(userID, zoneID uint, licensePlate string) (*Reservation, error)
	GetMyReservations(userID uint) ([]Reservation, error)
	GetReservationByID(id uint) (*Reservation, error)
	CancelReservation(id uint) error
	GetAllReservations() ([]Reservation, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// CreateReservation implements the "EV Spot Bottleneck" fix:
// row-level lock (FOR UPDATE) on the zone + capacity check + insert,
// all inside a single DB transaction, so two concurrent requests for
// the last spot can never both succeed.
func (r *repository) CreateReservation(userID, zoneID uint, licensePlate string) (*Reservation, error) {
	var reservation Reservation

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var z zone.ParkingZone

		// 1. Lock the zone row — any other transaction trying to lock the
		// same row will block here until this transaction commits/rolls back.
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&z, zoneID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrZoneNotFound
			}
			return err
		}

		// 2. Count current active reservations for this zone (safe now,
		// since the zone row is locked — no other txn can insert concurrently
		// and commit before we do).
		var activeCount int64
		if err := tx.Model(&Reservation{}).
			Where("zone_id = ? AND status = ?", zoneID, "active").
			Count(&activeCount).Error; err != nil {
			return err
		}

		// 3. Capacity check
		if int(activeCount) >= z.TotalCapacity {
			return ErrZoneFull
		}

		// 4. Create the reservation
		reservation = Reservation{
			UserID:       userID,
			ZoneID:       zoneID,
			LicensePlate: licensePlate,
			Status:       "active",
		}
		return tx.Create(&reservation).Error
	})

	if err != nil {
		return nil, err
	}
	return &reservation, nil
}

func (r *repository) GetMyReservations(userID uint) ([]Reservation, error) {
	var reservations []Reservation
	err := r.db.Preload("Zone").
		Where("user_id = ?", userID).
		Order("id desc").
		Find(&reservations).Error
	if err != nil {
		return nil, err
	}
	return reservations, nil
}

func (r *repository) GetReservationByID(id uint) (*Reservation, error) {
	var reservation Reservation
	err := r.db.First(&reservation, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrReservationNotFound
		}
		return nil, err
	}
	return &reservation, nil
}

func (r *repository) CancelReservation(id uint) error {
	result := r.db.Model(&Reservation{}).
		Where("id = ?", id).
		Update("status", "cancelled")
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrReservationNotFound
	}
	return nil
}

func (r *repository) GetAllReservations() ([]Reservation, error) {
	var reservations []Reservation
	err := r.db.Preload("User").Preload("Zone").
		Order("id desc").
		Find(&reservations).Error
	if err != nil {
		return nil, err
	}
	return reservations, nil
}
