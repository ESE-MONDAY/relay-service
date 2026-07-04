package models

type EmailStatus string

const (
	EmailQueued     EmailStatus = "queued"
	EmailProcessing EmailStatus = "processing"
	EmailSent       EmailStatus = "sent"
	EmailFailed     EmailStatus = "failed"
)
