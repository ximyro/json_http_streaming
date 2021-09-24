package main

import (
	"fmt"
	"net/http"
	"streaming/internal/handlers"
	"streaming/internal/services"
	"time"
)

func main() {
	mux := http.NewServeMux()
	eventsSender := services.NewEventsSender()
	eventsHandler := handlers.NewEventsHandler(eventsSender)
	orderHandler := handlers.NewOrderHandler()
	mux.Handle("/api/v1/events", eventsHandler)
	mux.Handle("/api/v1/orders", orderHandler)

	s := &http.Server{
		Addr:           ":8080",
		Handler:        mux,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Println("Server is started")
	err := s.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
