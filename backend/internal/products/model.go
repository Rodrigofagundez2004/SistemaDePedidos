package products

import 
(
	"net/http"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)
type Category struct
{
	ID uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name string    `gorm:"uniqueIndex;not null" json:"name"`
	Slug string     `gorm:"uniqueIndex;not null" json:"slug"`



}
func (Category) TableName() string {
	return "categories"
}
type Product struct{
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	CategoryID  *uuid.UUID `gorm:"type:uuid" json:"category_id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	PriceCents  int       `gorm:"not null" json:"price_cents"`
	Currency    string    `gorm:"not null;default:USD" json:"currency"`
	Stock       int       `gorm:"not null;default:0" json:"stock"`
	ImageURL    string    `json:"image_url"`
	IsActive    bool      `gorm:"not null;default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
func (Product) TableName() string{
	return "products"
}
func (p *Product) BeforeCreate(tx *gorm.DB) (err error)
{
	if p.ID == uuid.Nil{
		p.ID = uuid.New()
	}
	return nil
}