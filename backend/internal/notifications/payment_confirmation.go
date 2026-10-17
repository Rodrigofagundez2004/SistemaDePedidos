// Package notifications contiene los consumidores que reaccionan a eventos de negocio.
package notifications

import (
	"context"
	"log/slog"

	"github.com/RodrigoFagundez2004/SistemaDePedidos/backend/internal/events"
)

// SendPaymentConfirmation es un consumidor de ejemplo: cuando un pedido queda pagado,
// "envía" el email de confirmación. Acá iría la integración real (SMTP, SendGrid, etc.).
func SendPaymentConfirmation(ctx context.Context, ev events.OrderPaid) error {
	slog.Info("email de confirmación enviado (simulado)",
		"to_user_id", ev.UserID,
		"order_id", ev.OrderID,
		"amount_cents", ev.AmountCents,
		"currency", ev.Currency,
	)
	return nil
}
