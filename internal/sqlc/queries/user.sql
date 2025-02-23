-- Insert a new user
-- name: CreateUser :one
INSERT INTO users (username, password)
VALUES ($1, $2)
RETURNING *;

-- Get user by user id
-- name: GetUserByUserID :one
SELECT * FROM users
WHERE id = $1;

-- Get user by username
-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1;

-- Check username already taken
-- name: UserExists :one
SELECT EXISTS(
    SELECT 1
    FROM users
    WHERE username = $1
) AS exists;