package mango

import (
	"errors"

	"gorm.io/gorm"
)

var ErrMangoNotFound = errors.New("mango not found")

type Repository interface {
	Create(mango *Mango) error
	GetAll() ([]*Mango, error)
	GetByID(mangoId uint) (*Mango, error)
	Update(mango *Mango) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(mango *Mango) error {
	return r.db.Create(mango).Error
}

func (r *repository) GetAll() ([]*Mango, error) {
	var mangoes []*Mango
	if err := r.db.Find(&mangoes).Error; err != nil {
		return nil, err
	}
	return mangoes, nil
}

func (r *repository) GetByID(mangoId uint) (*Mango, error) {
	var mango Mango
	err := r.db.First(&mango, mangoId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMangoNotFound
		}
		return nil, err
	}
	return &mango, nil
}

func (r *repository) Update(mango *Mango) error {
	return r.db.Save(mango).Error
}
