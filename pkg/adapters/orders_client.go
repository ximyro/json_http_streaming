package adapters

import (
	"fmt"
	"io"
	"net/http"
)

type OrdersClient struct {
	hostWithPort string
}

const createOrderEndpoint = "/api/v1/orders"

func NewOrdersClient(hostWithPort string) *OrdersClient {
	return &OrdersClient{
		hostWithPort: hostWithPort,
	}
}

func (o *OrdersClient) CreateOrder() ([]byte, error) {
	client := http.Client{
		Transport: http.DefaultTransport,
	}

	req, err := http.NewRequest("POST", o.hostWithPort+createOrderEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("can't initialize request for creating an order: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("can't make a request: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("can't read body: %w", err)
	}

	return body, nil
}
