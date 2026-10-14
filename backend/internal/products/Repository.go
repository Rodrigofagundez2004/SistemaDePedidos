package products

import (
	"errors"

	"gorm.io/gorm"
)

var ErrProductNotFound = errors.New("Product not found")
var ErrCategoryNotFound = errors.New("Category not found ")
var ErrPrdouctNameTaken = errors.New("Name of product already exists")
var ErrCategorySlugTaken = errors.New("The slug of category already exists")
var ErrCategoryNameTaken = errors.New("The name of category already exists")
var ErrInvalidPrice = errors.New("The price is invalid")
var ErrInvalidStock = errors.New("The stock cannot be negativo ")
var ErrInsufficientStock = errors.New("Stock Insufficient")

type Repository struct {
	db *gorm.DB
}

// categories
func (r *Repository) CreateCategory(c *Category) error {
	//verificar que la categoria no exista
	var count int64
	if err := r.db.Model(&Category{}).Where("name = ?", c.Name).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return ErrCategoryNameTaken
	}
	count = 0
	if err := r.db.Model(&Category{}).Where("slug = ?", c.Slug).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return ErrCategorySlugTaken
	}
	return r.db.Create(c).Error

}
func (r *Repository) FindCategoryByID(id string) (*Category, error) {
	var c Category
	if err := r.db.Where("id = ?", id).First(&c).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	return &c, nil
}
func (r *Repository) FindCategoryBySlug(slug string) (*Category, error) {
	var c Category
	if err := r.db.Where("slug = ?", slug).First(&c).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err

	}
	return &c, nil
}
func (r *Repository) ListCategories() ([]Category, error) {
	var categories []Category
	if err := r.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// Products
func (r *Repository) CreateProduct(p *Product) error {
	var count int64
	if err := r.db.Model(&Product{}).Where("name = ?", p.Name).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return ErrPrdouctNameTaken
	}
	return r.db.Create(p).Error

}
func (r *Repository) UpdateProduct(p *Product) error {
	result := r.db.Save(p)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrProductNotFound
	}
	return nil

}

func (r *Repository) ListProducts(q ListProductsQuery) ([]Product, int64, error) {
	query := r.db.Model(&Product{})

	if q.Category != "" {
		query = query.Where("category_id", q.Category)
	}
	if q.Search != "" {
		query = query.Where("name ILIKE ?", "%"+q.Search+"%")
	}
	if q.MinPrice > 0 {
		query = query.Where("price_cents >= ?", q.MinPrice)
	}
	if q.MaxPrice > 0 {
		query = query.Where("price_cents <= ?", q.MaxPrice)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	switch q.Sort {
	case "price_asc":
		query = query.Order("price_cents ASC")
	case "price_desc":
		query = query.Order("price_cents DESC")
	case "newest":
		query = query.Order("created_at DESC")
	default:
		query = query.Order("created_at ASC")
	}
	// Paginación
	offset := (q.Page - 1) * q.Limit
	query = query.Limit(q.Limit).Offset(offset)

	var products []Product
	if err := query.Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *Repository) DeleteProduct(id string) error {
	result := r.db.Delete(&Product{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrProductNotFound
	}
	return nil

}

func (r *Repository) FindByID(id string) (*Product, error) {
	var p Product
	err := r.db.Where("id = ?", id).First(&p).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrProductNotFound
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}
