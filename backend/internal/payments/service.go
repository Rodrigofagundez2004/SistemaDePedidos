package payments 

import (
	"log/slog"
	"math/rand"

	"github.com/google/uuid"
)

type Service struct {
	repo         *Repository
	ordersClient OrdersClient
}

func NewService(repo *Repository, ordersClient OrdersClient) *Service {
	return &Service{
		repo:         repo,
		ordersClient: ordersClient,
	}
}

func (s *Service) CreatePayment (req *CreatePaymentRequest) (*PaymentResponse, error) {


}