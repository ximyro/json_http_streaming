package entities

import "time"

type Event struct {
	T         string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}
