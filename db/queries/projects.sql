-- name: CreateProject :one
INSERT INTO projects (name, description)
VALUES ($1, $2)
RETURNING *;

-- name: GetProject :one
SELECT * FROM projects WHERE id = $1;
