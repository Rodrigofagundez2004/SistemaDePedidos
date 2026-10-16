package orders

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID         uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Status         string    `gorm:"not null;default:pending" json:"status"`
	TotalCents     int       `gorm:"not null" json:"total_cents"`
	Currency       string    `gorm:"not null;default:USD" json:"currency"`
	IdempotencyKey string    `gorm:"uniqueIndex;not null" json:"idempotency_key"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (Order) TableName() string { return "orders" }

func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}

type OrderItem struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	OrderID        uuid.UUID `gorm:"type:uuid;not null" json:"order_id"`
	ProductID      uuid.UUID `gorm:"type:uuid;not null" json:"product_id"`
	ProductName    string    `gorm:"not null" json:"product_name"`
	UnitPriceCents int       `gorm:"not null" json:"unit_price_cents"`
	Quantity       int       `gorm:"not null" json:"quantity"`
	SubtotalCents  int       `gorm:"not null" json:"subtotal_cents"`
}

func (OrderItem) TableName() string { return "order_items" }

func (oi *OrderItem) BeforeCreate(tx *gorm.DB) error {
	if oi.ID == uuid.Nil {
		oi.ID = uuid.New()
	}
	return nil
}
