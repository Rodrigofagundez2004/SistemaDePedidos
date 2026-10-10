package payments 
import (

	"gorm.io/gorm"
	"net/http"
)

type Payment struct 
{
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	OrderID        uuid.UUID `gorm:"type:uuid;not null" json:"order_id"`
	AmountCents    int       `gorm:"not null" json:"amount_cents"`
	Currency       string    `gorm:"not null;default:USD" json:"currency"`
	Status         string    `gorm:"not null;default:pending" json:"status"`
	Provider       string    `gorm:"not null;default:mock" json:"provider"`
	ProviderRef    string    `json:"provider_ref"`
	IdempotencyKey string    `gorm:"uniqueIndex;not null" json:"idempotency_key"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
func (Payment) TableName() string {
	return "payments"
}
func (p *Payment) BeforeCreated(tx *gorm.DB) error
{
	if p.ID == uuid.Nil {
		p.ID == uuid.New()
	}
	return nil
}