package event

import (
	"time"

	"github.com/google/uuid"
)

type EmailEvent struct {
	Version string    `json:"version"`
	EventID string    `json:"event_id"`
	Type    string    `json:"type"`
	EmailID uuid.UUID `json:"email_id"`
	Retry   int       `json:"retry"`
	Time    time.Time `json:"time"`
}

const (
	EmailEventTypeCreated = "email.created"
	EmailEventTypeRetry   = "email.retry"
	EmailEventTypeFailed  = "email.failed"
	EmailEventTypeDLQ     = "email.dlq"
)
