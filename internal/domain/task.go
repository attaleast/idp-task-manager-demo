package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrTaskNotFound  = errors.New("task not found")
	ErrInvalidStatus = errors.New("invalid task status")
)

type TaskStatus string

const (
	StatusTodo       TaskStatus = "todo"
	StatusInProgress TaskStatus = "in_progress"
	StatusDone       TaskStatus = "done"
)

type Task struct {
	ID          uuid.UUID  `json:"id"`
	ProjectID   uuid.UUID  `json:"project_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	AssigneeID  *uuid.UUID `json:"assignee_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (t *Task) ChangeStatus(newStatus TaskStatus) error {
	switch newStatus {
	case StatusDone, StatusInProgress, StatusTodo:
	default:
		return ErrInvalidStatus
	}

	t.Status = newStatus
	return nil
}
