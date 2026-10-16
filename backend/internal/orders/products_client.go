package orders

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ProductsClient interface {
	GetProduct(productID string) (*ProductInfo, error)
}

type ProductInfo struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	PriceCents int    `json:"price_cents"`
	Currency   string `json:"currency"`
	Stock      int    `json:"stock"`
	IsActive   bool   `json:"is_active"`
}

type HTTPProductsClient struct {
	baseURL string
	client  *http.Client
}

func NewHTTPProductsClient(baseURL string) *HTTPProductsClient {
	return &HTTPProductsClient{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *HTTPProductsClient) GetProduct(productID string) (*ProductInfo, error) {
	url := fmt.Sprintf("%s/products/%s", c.baseURL, productID)
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrProductNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var p ProductInfo
	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		return nil, err
	}
	return &p, nil
}
