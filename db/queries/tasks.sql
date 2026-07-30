-- name: CreateTask :one
INSERT INTO tasks (project_id, title, description, assignee_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetTask :one
SELECT * FROM tasks WHERE id = $1;

-- name: ListTasksByProject :many
SELECT * FROM tasks WHERE project_id = $1 ORDER BY created_at DESC;

-- name: UpdateTaskStatus :one
UPDATE tasks SET status = $2, updated_at = NOW() WHERE id = $1 RETURNING *;
