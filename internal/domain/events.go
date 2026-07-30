package domain

import "github.com/google/uuid"

type TaskCreatedEvent struct {
	TaskID    uuid.UUID `json:"task_id"`
	ProjectID uuid.UUID `json:"project_id"`
	Title     string    `json:"title"`
}
