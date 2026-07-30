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
	ID          uuid.UUID
	ProjectID   uuid.UUID
	Title       string
	Description string
	Status      TaskStatus
	AssigneeID  *uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
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
