package application

import (
	"context"

	"github.com/attaleast/idp-sdk/messaging"
	"github.com/attaleast/idp-task-manager-demo/internal/domain"
	"github.com/google/uuid"
)

type TaskRepository interface {
	Create(ctx context.Context, task *domain.Task) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Task, error)
	Update(ctx context.Context, task *domain.Task) error
	ListByProject(ctx context.Context, projectID uuid.UUID) ([]*domain.Task, error)
}

type TaskUseCase struct {
	repo  TaskRepository
	queue messaging.Publisher
}

func NewTaskUseCase(repo TaskRepository, queue messaging.Publisher) *TaskUseCase {
	return &TaskUseCase{repo: repo, queue: queue}
}

type CreateTaskInput struct {
	ProjectID   uuid.UUID
	Title       string
	Description string
}

func (uc *TaskUseCase) CreateTask(ctx context.Context, in CreateTaskInput) (*domain.Task, error) {
	task := &domain.Task{
		ProjectID:   in.ProjectID,
		Title:       in.Title,
		Description: in.Description,
		Status:      domain.StatusTodo,
	}

	if err := uc.repo.Create(ctx, task); err != nil {
		return nil, err
	}

	if err := uc.queue.Publish(ctx, "tasks.created", &domain.TaskCreatedEvent{
		Title:     task.Title,
		ProjectID: task.ProjectID,
		TaskID:    task.ID,
	}); err != nil {
		// TODO: add outbox
		return nil, err
	}

	return task, nil
}

func (uc *TaskUseCase) ListTasks(ctx context.Context, projectID uuid.UUID) ([]*domain.Task, error) {
	return uc.repo.ListByProject(ctx, projectID)
}
