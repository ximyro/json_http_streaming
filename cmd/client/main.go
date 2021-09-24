package main

import (
	"context"
	"encoding/json"
	"fmt"
	"streaming/internal/entities"
	"streaming/pkg/adapters"
)

var ErrStopListening = fmt.Errorf("stop listening")

func handleFunc(data []byte) error {
	event := entities.Event{}
	err := json.Unmarshal(data, &event)
	if err != nil {
		return fmt.Errorf("can't unmarshal event: %w", err)
	}
	fmt.Printf("Received event: %s\n", string(data))

	// We received the event which we're waiting for
	if event.T == "order_finished" {
		fmt.Println("don't need to listen more")
		return ErrStopListening
	}

	return nil
}

func main() {
	ordersClient := adapters.NewOrdersClient("http://localhost:8080")
	client := adapters.NewEventsClient("localhost:8080")
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	body, err := ordersClient.CreateOrder()
	if err != nil {
		panic(err)
	}

	order := entities.Order{}
	err = json.Unmarshal(body, &order)
	if err != nil {
		panic(err)
	}

	err = client.ReadEvents(ctx, handleFunc)
	if err != nil {
		if err == ErrStopListening {
			fmt.Println("goodbye")
			return
		}
	}
}
