package orders

type CreateOrderItemRequest struct {
	ProductID string `json:"product_id" binding:"required,uuid"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
}

type CreateOrderRequest struct {
	Items          []CreateOrderItemRequest `json:"items" binding:"required,min=1,dive"`
	IdempotencyKey string                   `json:"idempotency_key" binding:"required,uuid"`
}

type OrderResponse struct {
	ID             string              `json:"id"`
	UserID         string              `json:"user_id"`
	Status         string              `json:"status"`
	TotalCents     int                 `json:"total_cents"`
	Currency       string              `json:"currency"`
	IdempotencyKey string              `json:"idempotency_key"`
	Items          []OrderItemResponse `json:"items,omitempty"`
	CreatedAt      string              `json:"created_at"`
	UpdatedAt      string              `json:"updated_at"`
}

type OrderItemResponse struct {
	ID             string `json:"id"`
	ProductID      string `json:"product_id"`
	ProductName    string `json:"product_name"`
	UnitPriceCents int    `json:"unit_price_cents"`
	Quantity       int    `json:"quantity"`
	SubtotalCents  int    `json:"subtotal_cents"`
}

type ListOrdersQuery struct {
	Page   int    `form:"page"`
	Limit  int    `form:"limit"`
	Status string `form:"status"`
}

type OrderInfo struct {
	ID         string `json:"id"`
	Status     string `json:"status"`
	TotalCents int    `json:"total_cents"`
	Currency   string `json:"currency"`
}
