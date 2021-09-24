package services

import (
	"math/rand"
	"streaming/internal/entities"
	"time"
)

type EventsSender struct {
}

func NewEventsSender() *EventsSender {
	return &EventsSender{}
}

func (e *EventsSender) writeEvents(eventsChan chan entities.Event) {
	instantEvents := []entities.Event{
		{
			T:         "order_created",
			CreatedAt: time.Now(),
		},
		{
			T:         "order_checked",
			CreatedAt: time.Now(),
		},
		{
			T:         "order_in_routing",
			CreatedAt: time.Now(),
		},
	}
	for _, event := range instantEvents {
		eventsChan <- event
	}

	newEvents := []entities.Event{
		{
			T:         "order_started",
			CreatedAt: time.Now(),
		},
		{
			T:         "order_in_ride",
			CreatedAt: time.Now(),
		},
		{
			T:         "order_finished",
			CreatedAt: time.Now(),
		},
	}

	for i := 0; i < len(newEvents); i++ {
		time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
		eventsChan <- newEvents[i]
	}
	close(eventsChan)
}

func (e *EventsSender) GetEvents() chan entities.Event {
	eventsChan := make(chan entities.Event)
	go func() {
		e.writeEvents(eventsChan)
	}()

	return eventsChan
}
