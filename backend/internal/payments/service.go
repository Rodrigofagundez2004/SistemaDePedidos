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

func (s *Service) CreatePayment(req *CreatePaymentRequest) (*PaymentResponse, error) {
	existing, err := s.repo.FindByIdempotencyKey(req.IdempotencyKey)
	if err == nil {
		// Si existe una orden con id distinto que
		if existing.OrderID.String() != req.OrderID {
			return nil, ErrIdempotencyKeyConflict
		}
		return ToPaymentResponse(existing), nil
	}
	//consultar el orden
	order, err := s.ordersClient.GetOrder(req.OrderID)
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
		AmountCents:    order.TotalCents, //  ¡Sácalo de la orden que consultaste antes! (ajusta el nombre si tu struct de orden lo llama diferente, ej: order.Amount, order.AmountCents, etc.)
		Currency:       order.Currency,   //  Igualmente la moneda viene de la orden
		Status:         "pending",
		Provider:       "mock",
		IdempotencyKey: req.IdempotencyKey, // Clave de idempotencia para evitar duplicados
	}
	if err := s.repo.Create(payment); err != nil {
		return nil, err
	}
	// 6. Simular cobro
	if simulateSuccess() {
		payment.Status = "succeeded"
		payment.ProviderRef = "mock_" + uuid.New().String()
		if err := s.ordersClient.UpdateOrderStatus(req.OrderID, "paid"); err != nil {
			slog.Error("pago exitoso pero orden no actualizada",
				"order_id", req.OrderID,
				"payment_id", payment.ID,
				"error", err,
			)
		}
	} else {
		payment.Status = "failed"
		s.ordersClient.UpdateOrderStatus(req.OrderID, "failed")
	}

	// 7. Actualizar
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
