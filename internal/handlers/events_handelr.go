package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"streaming/internal/entities"
)

type RequestSender interface {
	GetEvents() chan entities.Event
}

type EventsHandler struct {
	sender RequestSender
}

func NewEventsHandler(sender RequestSender) *EventsHandler {
	return &EventsHandler{
		sender: sender,
	}
}

func (e EventsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	events := e.sender.GetEvents()
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	for e := range events {
		if err := enc.Encode(e); err != nil {
			panic(err)
		}
		_, err := io.WriteString(w, buf.String()+"\n")
		if err != nil {
			panic(err)
		}
		w.(http.Flusher).Flush()
		buf.Reset()
	}
}
