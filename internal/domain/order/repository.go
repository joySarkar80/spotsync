package order

import (
	"errors"
	"haddibanga/internal/domain/mango"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrOrderNotFound        = errors.New("order not found")
	ErrNotEnoughStock       = errors.New("not enough stock available")
	ErrOrderAlreadyCancelled = errors.New("order already cancelled")
	ErrForbiddenOrderAccess = errors.New("you do not own this order")
)

type Repository interface {
	Create(order *Order) error
	GetByID(orderId uint) (*Order, error)
	GetByUserID(userId uint) ([]*Order, error)
	Update(order *Order) error
	CreateWithStockUpdate(userId, mangoId uint, quantityKg float64) (*Order, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(order *Order) error {
	return r.db.Create(order).Error
}

func (r *repository) GetByID(orderId uint) (*Order, error) {
	var order Order
	err := r.db.First(&order, orderId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}
	return &order, nil
}

func (r *repository) GetByUserID(userId uint) ([]*Order, error) {
	var orders []*Order
	if err := r.db.Where("user_id = ?", userId).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *repository) Update(order *Order) error {
	return r.db.Save(order).Error
}

func (r *repository) CreateWithStockUpdate(userId, mangoId uint, quantityKg float64) (*Order, error) {
	var order Order

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var mangoData mango.Mango

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&mangoData, mangoId).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return mango.ErrMangoNotFound
			}
			return err
		}

		if mangoData.StockKg < quantityKg {
			return ErrNotEnoughStock
		}

		order = Order{
			UserID:     userId,
			MangoID:    mangoData.ID,
			QuantityKg: quantityKg,
			TotalPrice: quantityKg * mangoData.PricePerKg,
			Status:     OrderConfirmed,
			OrderCode:  generateOrderCode(),
		}

		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		mangoData.StockKg -= quantityKg
		if err := tx.Save(&mangoData).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &order, nil
}
