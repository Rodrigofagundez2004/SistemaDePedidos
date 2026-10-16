package orders

import (
	"errors"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(order *Order, items []OrderItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		for i := range items {
			items[i].OrderID = order.ID
		}
		return tx.Create(&items).Error
	})
}

func (r *Repository) FindByID(id string) (*Order, error) {
	var o Order
	err := r.db.Where("id = ?", id).First(&o).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrOrderNotFound
	}
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *Repository) FindByIdempotencyKey(key string) (*Order, error) {
	var o Order
	err := r.db.Where("idempotency_key = ?", key).First(&o).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrOrderNotFound
	}
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *Repository) ListItems(orderID string) ([]OrderItem, error) {
	var items []OrderItem
	err := r.db.Where("order_id = ?", orderID).Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *Repository) ListByUser(userID string, q ListOrdersQuery) ([]Order, int64, error) {
	query := r.db.Model(&Order{}).Where("user_id = ?", userID)
	if q.Status != "" {
		query = query.Where("status = ?", q.Status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var orders []Order

	offset := (q.Page - 1) * q.Limit
	err := query.Order("created_at DESC").Limit(q.Limit).Offset(offset).Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}
	return orders, total, nil
}

func (r *Repository) UpdateStatus(orderID, status string) error {
	return r.db.Model(&Order{}).Where("id = ?", orderID).Update("status", status).Error
}
