package payments

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type OrdersClient interface {
	GetOrder(orderID string) (*OrderInfo, error)
	UpdateOrderStatus(orderID, status string) error
}

type OrderInfo struct {
	ID         string `json:"id"`
	Status     string `json:"status"`
	TotalCents int    `json:"total_cents"`
	Currency   string `json:"currency"`
}

type HTTPOrdersClient struct {
	baseURL string
	client  *http.Client
}

func NewHTTPOrdersClient(baseURL string) *HTTPOrdersClient {
	return &HTTPOrdersClient{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *HTTPOrdersClient) GetOrder(orderID string) (*OrderInfo, error) {
	url := fmt.Sprintf("%s/orders/%s", c.baseURL, orderID)
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrOrderNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	// La respuesta viene como {order: {...}, items: [...]}
	var wrapper struct {
		Order OrderInfo `json:"order"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, err
	}
	return &wrapper.Order, nil
}

func (c *HTTPOrdersClient) UpdateOrderStatus(orderID, status string) error {
	url := fmt.Sprintf("%s/orders/%s/status", c.baseURL, orderID)
	body, _ := json.Marshal(map[string]string{"status": status})

	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
	return nil
}
