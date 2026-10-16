package payments

import (
	"errors"
	"log/slog"
	"math/rand"

	"github.com/google/uuid"
)

var (
	ErrPaymentNotFound          = errors.New("pago no encontrado")
	ErrOrderNotFound            = errors.New("orden no encontrada")
	ErrOrderAlreadyPaid         = errors.New("la orden ya está pagada")
	ErrOrderNotPending          = errors.New("la orden no está pendiente de pago")
	ErrIdempotencyKeyConflict   = errors.New("idempotency key ya usada con otra orden")
	ErrInvalidAmount            = errors.New("monto inválido")
	ErrPaymentNotRefundable     = errors.New("el pago no se puede reembolsar")
	ErrOrdersServiceUnavailable = errors.New("servicio de órdenes no disponible")
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

// CreatePayment cobra una orden. authHeader es el token del usuario que hace la compra;
// se reenvía a Orders para que acepte la consulta y el cambio de estado.
func (s *Service) CreatePayment(req *CreatePaymentRequest, authHeader string) (*PaymentResponse, error) {
	existing, err := s.repo.FindByIdempotencyKey(req.IdempotencyKey)
	if err == nil {
		// Si existe una orden con id distinto que
		if existing.OrderID.String() != req.OrderID {
			return nil, ErrIdempotencyKeyConflict
		}
		return ToPaymentResponse(existing), nil
	}

	// consultar la orden
	order, err := s.ordersClient.GetOrder(req.OrderID, authHeader)
	if err != nil {
		if err == ErrOrderNotFound {
			return nil, ErrOrderNotFound
		}
		return nil, ErrOrdersServiceUnavailable
	}

	// validar estado
	if order.Status == "paid" {
		return nil, ErrOrderAlreadyPaid
	}
	if order.Status != "pending" {
		return nil, ErrOrderNotPending
	}

	payment := &Payment{
		ID:             uuid.New(),
		OrderID:        uuid.MustParse(req.OrderID),
		AmountCents:    order.TotalCents,
		Currency:       order.Currency,
		Status:         "pending",
		Provider:       "mock",
		IdempotencyKey: req.IdempotencyKey,
	}
	if err := s.repo.Create(payment); err != nil {
		return nil, err
	}

	// simular cobro
	if simulateSuccess() {
		payment.Status = "succeeded"
		payment.ProviderRef = "mock_" + uuid.New().String()
		if err := s.ordersClient.UpdateOrderStatus(req.OrderID, "paid", authHeader); err != nil {
			slog.Error("pago exitoso pero orden no actualizada",
				"order_id", req.OrderID,
				"payment_id", payment.ID,
				"error", err,
			)
		}
	} else {
		payment.Status = "failed"
		s.ordersClient.UpdateOrderStatus(req.OrderID, "failed", authHeader)
	}

	// actualizar el pago con el resultado
	if err := s.repo.Update(payment); err != nil {
		return nil, err
	}

	return ToPaymentResponse(payment), nil
}

func (s *Service) GetByID(id string) (*Payment, error) {
	return s.repo.FindByID(id)
}

func (s *Service) ListByOrder(orderID string) ([]Payment, error) {
	return s.repo.ListByOrder(orderID)
}

func simulateSuccess() bool {
	return rand.Intn(100) < 90
}
