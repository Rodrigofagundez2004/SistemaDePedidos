package orders

import (
	"github.com/google/uuid"
)

type Service struct {
	repo           *Repository
	productsClient ProductsClient
}

func NewService(repo *Repository, productsClient ProductsClient) *Service {
	return &Service{
		repo:           repo,
		productsClient: productsClient,
	}
}

func (s *Service) Create(userID string, req *CreateOrderRequest) (*Order, error) {
	// 1. Idempotencia
	if existing, err := s.repo.FindByIdempotencyKey(req.IdempotencyKey); err == nil {
		return existing, nil
	}

	if len(req.Items) == 0 {
		return nil, ErrEmptyOrder
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	order := &Order{
		UserID:         userUUID,
		Status:         "pending",
		Currency:       "USD",
		IdempotencyKey: req.IdempotencyKey,
	}

	var items []OrderItem
	totalCents := 0

	for _, item := range req.Items {
		if item.Quantity <= 0 {
			return nil, ErrInvalidQuantity
		}

		product, err := s.productsClient.GetProduct(item.ProductID)
		if err != nil {
			return nil, err
		}

		subtotal := product.PriceCents * item.Quantity
		totalCents += subtotal

		items = append(items, OrderItem{
			ProductID:      uuid.MustParse(item.ProductID),
			ProductName:    product.Name,
			UnitPriceCents: product.PriceCents,
			Quantity:       item.Quantity,
			SubtotalCents:  subtotal,
		})

		order.Currency = product.Currency
	}

	order.TotalCents = totalCents

	if err := s.repo.Create(order, items); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *Service) GetByID(id string) (*Order, []OrderItem, error) {
	order, err := s.repo.FindByID(id)
	if err != nil {
		return nil, nil, err
	}
	items, err := s.repo.ListItems(id)
	if err != nil {
		return nil, nil, err
	}
	return order, items, nil
}

func (s *Service) ListByUser(userID string, q ListOrdersQuery) ([]Order, int64, error) {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 || q.Limit > 100 {
		q.Limit = 20
	}
	return s.repo.ListByUser(userID, q)
}

func (s *Service) UpdateStatus(orderID, status string) error {
	return s.repo.UpdateStatus(orderID, status)
}
