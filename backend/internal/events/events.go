// Package events publica y consume eventos de negocio a través de Kafka.
package events

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/segmentio/kafka-go"
)

// TopicOrderPaid es el canal donde se publica cada pedido que queda pagado.
const TopicOrderPaid = "order.paid"

// OrderPaid es el mensaje que viaja por Kafka cuando un pedido queda pagado.
type OrderPaid struct {
	OrderID     string    `json:"order_id"`
	UserID      string    `json:"user_id"`
	PaymentID   string    `json:"payment_id"`
	AmountCents int       `json:"amount_cents"`
	Currency    string    `json:"currency"`
	PaidAt      time.Time `json:"paid_at"`
}

// Publisher envía eventos a Kafka.
type Publisher struct {
	writer *kafka.Writer
}

// NewPublisher crea un publicador conectado a los brokers indicados (por ejemplo "localhost:9094").
func NewPublisher(brokers []string) *Publisher {
	return &Publisher{
		writer: &kafka.Writer{
			Addr:                   kafka.TCP(brokers...),
			Topic:                  TopicOrderPaid,
			Balancer:               &kafka.Hash{}, // el mismo pedido siempre va a la misma partición
			AllowAutoTopicCreation: true,
			RequiredAcks:           kafka.RequireOne,
			WriteTimeout:           5 * time.Second,
		},
	}
}

// PublishOrderPaid publica el evento usando el order_id como clave.
func (p *Publisher) PublishOrderPaid(ctx context.Context, ev OrderPaid) error {
	body, err := json.Marshal(ev)
	if err != nil {
		return err
	}
	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(ev.OrderID),
		Value: body,
	})
}

// Close cierra la conexión con Kafka.
func (p *Publisher) Close() error {
	return p.writer.Close()
}

// Handler procesa un evento consumido.
type Handler func(ctx context.Context, ev OrderPaid) error

// DecodeOrderPaid convierte el contenido de un mensaje en un OrderPaid.
func DecodeOrderPaid(body []byte) (OrderPaid, error) {
	var ev OrderPaid
	if err := json.Unmarshal(body, &ev); err != nil {
		return OrderPaid{}, err
	}
	if ev.OrderID == "" {
		return OrderPaid{}, errors.New("evento sin order_id")
	}
	return ev, nil
}

// ConsumeOrderPaid lee el topic dentro de un grupo de consumidores y llama a handle por cada evento.
// Corre hasta que ctx se cancela. Si Kafka no está disponible, reintenta solo.
func ConsumeOrderPaid(ctx context.Context, brokers []string, group string, handle Handler) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  group,
		Topic:    TopicOrderPaid,
		MinBytes: 1,
		MaxBytes: 10 << 20,
	})
	defer reader.Close()

	slog.Info("consumidor de Kafka iniciado", "topic", TopicOrderPaid, "group", group)

	for {
		msg, err := reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return // apagado normal del servidor
			}
			slog.Warn("no se pudo leer de Kafka, reintentando", "error", err)
			select {
			case <-ctx.Done():
				return
			case <-time.After(3 * time.Second):
			}
			continue
		}

		ev, err := DecodeOrderPaid(msg.Value)
		if err != nil {
			slog.Error("mensaje inválido, se descarta", "error", err, "offset", msg.Offset)
		} else if err := handle(ctx, ev); err != nil {
			slog.Error("error procesando evento", "order_id", ev.OrderID, "error", err)
		}

		if err := reader.CommitMessages(ctx, msg); err != nil && ctx.Err() == nil {
			slog.Error("no se pudo confirmar el mensaje", "error", err)
		}
	}
}
