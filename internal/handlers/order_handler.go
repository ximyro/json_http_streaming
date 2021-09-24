package handlers

import (
	"encoding/json"
	"net/http"
	"streaming/internal/entities"
	"time"
)

type OrderHandler struct{}

func NewOrderHandler() *OrderHandler {
	return &OrderHandler{}
}

func (o OrderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	order := entities.Order{
		ID:        "12345",
		Status:    "order_created",
		CreatedAt: time.Now(),
	}

	data, err := json.Marshal(order)
	if err != nil {
		panic(err)
	}
	_, err = w.Write(data)
	if err != nil {
		panic(err)
	}
}
