package payments
import (
	"errors"
	"gorm.io/gorm"
)
var ErrPaymentNotFoundDB = errors.New("Payment not found")
type  Repository struct {
	db *gorm.Db
} 
func newRepository(db *gorm.DB) *Repository {

	return &Repository(db: db)
}
func (r *Repository) Create(p *Payment) error {
	return r.db.Create(p).Error
}
func (r *Repository) Update(p *Payment) error {
	return r.db.Save(p).Error
}
func (r *Repository) FindByIDl(id string) (*Payment, error) 
{
	var p Payment
	err := r.db.Where("id = ?", id).First(&p).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPaymentNotFound
		}
		return nil, err
	}
	return &p, nil
}
func (r *Repository) FindByIdempotencyKey(key string) (*Payment, error) {
	 
}