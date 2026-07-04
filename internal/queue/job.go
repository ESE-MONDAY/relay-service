package queue

import "github.com/google/uuid"

type Job struct {
	ID      uuid.UUID
	EmailID uuid.UUID
	Type    string
}
