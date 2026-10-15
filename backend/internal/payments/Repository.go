package payments

import (
	"errors"

	"gorm.io/gorm"
)

var ErrPaymentNotFoundDB = errors.New("Payment not found")

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {

	return &Repository{db: db}
}
func (r *Repository) Create(p *Payment) error {
	return r.db.Create(p).Error
}
func (r *Repository) Update(p *Payment) error {
	return r.db.Save(p).Error
}
func (r *Repository) FindByID(id string) (*Payment, error) {
	var p Payment
	if err := r.db.Where("id = ?", id).First(&p).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPaymentNotFound
		}
		return nil, err
	}
	return &p, nil
}
func (r *Repository) FindByIdempotencyKey(key string) (*Payment, error) {
	var p Payment
	if err := r.db.Where("idempotency_key = ?", key).First(&p).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPaymentNotFoundDB
		}
		return nil, err
	}
	return &p, nil

}
func (r *Repository) ListByOrder(orderID string) ([]Payment, error) {
	var payments []Payment
	err := r.db.Where("order_id = ?", orderID).Order("created_at DESC").Find(&payments).Error
	if err != nil {
		return nil, err
	}
	return payments, nil
}
func (r *Repository) ListByOrderAsc(orderID string) ([]Payment, error) {
	var payments []Payment
	err := r.db.Where("order_id = ?", orderID).Order("created_at ASC").Find(&payments).Error
	if err != nil {
		return nil, err
	}
	return payments, nil
}
