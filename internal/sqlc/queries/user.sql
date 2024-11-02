-- Insert a new user
-- name: CreateUser :one
INSERT INTO users (username)
VALUES ($1)
RETURNING *;

-- Get user by user id
-- name: GetUserByUserID :one
SELECT * FROM users
WHERE id = $1;