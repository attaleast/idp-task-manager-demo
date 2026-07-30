package postgres

import (
	"context"

	"github.com/attaleast/idp-task-manager-demo/internal/domain"
	"github.com/attaleast/idp-task-manager-demo/internal/infrastructure/postgres/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepository struct {
	db *pgxpool.Pool
	q  *sqlc.Queries
}

func NewTaskRepository(db *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{
		db: db,
		q:  sqlc.New(db),
	}
}

func (r *TaskRepository) Create(ctx context.Context, t *domain.Task) error {
	args := sqlc.CreateTaskParams{
		Title:       t.Title,
		Description: &t.Description,
		ProjectID:   t.ProjectID,
		AssigneeID:  t.AssigneeID,
	}

	_, err := r.q.CreateTask(ctx, args)
	return err
}

func (r *TaskRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	row, err := r.q.GetTask(ctx, id)
	if err != nil {
		return nil, err
	}

	t := &domain.Task{
		ID:          row.ID,
		Title:       row.Title,
		Description: *row.Description,
		ProjectID:   row.ProjectID,
		AssigneeID:  row.AssigneeID,
		Status:      domain.TaskStatus(row.Status),
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}

	return t, nil
}

func (r *TaskRepository) Update(ctx context.Context, t *domain.Task) error {
	args := sqlc.UpdateTaskStatusParams{
		ID:     t.ID,
		Status: string(t.Status),
	}

	_, err := r.q.UpdateTaskStatus(ctx, args)
	return err
}
