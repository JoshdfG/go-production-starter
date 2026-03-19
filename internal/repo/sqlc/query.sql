-- name: GetAll :many
SELECT * FROM todos ORDER BY created_at DESC;

-- name: GetByID :one
SELECT * FROM todos WHERE id = $1;

-- name: Save :exec
INSERT INTO todos (id, title, done, created_at)
VALUES ($1, $2, $3, $4)
ON CONFLICT (id) DO UPDATE
SET title = EXCLUDED.title,
    done = EXCLUDED.done;

-- name: Delete :exec
DELETE FROM todos WHERE id = $1;

-- name: CreateUser :exec
INSERT INTO users (id, email, password, created_at)
VALUES ($1, $2, $3, $4);

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;
