package payments

type OrdersClient interface {
	GetOrder(orderID string) (*OrderInfo, error)
	UpdateOrderStatus(orderID, status string) error
}
