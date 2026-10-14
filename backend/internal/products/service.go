package products

import (
	"errors"

	"github.com/google/uuid"
)

var ErrInvalidStock = errors.New("Stock is invalid")
var ErrPriceInvalid = errors.New("Price is Invalid")

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}
func (s *Service) CreateProduct(req *CreateProductRequest) (*Product, error) {
	if req.PriceCents < 0 {
		return nil, ErrPriceInvalid
	}
	if req.Stock < 0 {
		return nil, ErrInvalidStock
	}
	// verifico tambien que la categoria existe si es dada
	if req.CategoryID != "" {
		_, err := s.repo.FindCategoryByID(req.CategoryID)
		if err != nil {
			return nil, err
		}
	}
	// aplicar defaults
	currency := req.Currency
	if currency == "" {
		currency = "USD"
	}
	//convierto el dto a model
	product := &Product{
		Name:        req.Name,
		Description: req.Description,
		PriceCents:  req.PriceCents,
		Currency:    currency,
		Stock:       req.Stock,
		ImageURL:    req.ImageURL,
		IsActive:    true,
	}
	// 6 asignar categoria si viene
	if req.CategoryID != "" {
		catID, err := uuid.Parse(req.CategoryID)
		if err != nil {
			return nil, ErrCategoryNotFound
		}
		product.CategoryID = &catID
	}
	//guardo en la db
	if err := s.repo.CreateProduct(product); err != nil {
		return nil, err
	}
	return product, nil
}
func (s *Service) GetProductById(id string) (*Product, error) {
	return s.repo.FindByID(id)
}
func (s *Service) ListProducts(query ListProductsQuery) ([]Product, int64, error) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 || query.Limit > 100 {
		query.Limit = 20
	}

	return s.repo.ListProducts(query)
}
