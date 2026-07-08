package models

import (
	"time"

	"github.com/google/uuid"
)

type Email struct {
	ID uuid.UUID `json:"id"`

	From string `json:"from"`
	To   string `json:"to"`

	Subject string `json:"subject"`
	Text    string `json:"text"`
	HTML    string `json:"html"`

	Status EmailStatus `json:"status"`

	RetryCount int
	LastError  string

	IdempotencyKey *string `json:"-"`

	CreatedAt time.Time `json:"created_at"`
}
