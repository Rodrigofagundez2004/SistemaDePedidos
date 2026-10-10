package payments

import "time"

// ============================================
// REQUEST DTOs (lo que el cliente manda)
// ============================================

type CreatePaymentRequest struct {
	OrderID        string `json:"order_id" binding:"required,uuid"`
	IdempotencyKey string `json:"idempotency_key" binding:"required,uuid"`
}

// ============================================
// RESPONSE DTOs (lo que el backend devuelve)
// ============================================

type PaymentResponse struct {
	ID             string `json:"id"`
	OrderID        string `json:"order_id"`
	AmountCents    int    `json:"amount_cents"`
	Currency       string `json:"currency"`
	Status         string `json:"status"`
	Provider       string `json:"provider"`
	ProviderRef    string `json:"provider_ref,omitempty"`
	IdempotencyKey string `json:"idempotency_key"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

// ============================================
// CONVERSOR modelo → DTO
// ============================================

func ToPaymentResponse(p *Payment) *PaymentResponse {
	return &PaymentResponse{
		ID:             p.ID.String(),
		OrderID:        p.OrderID.String(),
		AmountCents:    p.AmountCents,
		Currency:       p.Currency,
		Status:         p.Status,
		Provider:       p.Provider,
		ProviderRef:    p.ProviderRef,
		IdempotencyKey: p.IdempotencyKey,
		CreatedAt:      p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      p.UpdatedAt.Format(time.RFC3339),
	}
}

// Convierte una lista de modelos a una lista de DTOs
func ToPaymentResponseList(payments []Payment) []PaymentResponse {
	result := make([]PaymentResponse, len(payments))
	for i, p := range payments {
		result[i] = *ToPaymentResponse(&p)
	}
	return result
}