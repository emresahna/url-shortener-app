-- Create a new shortened URL for a specific user
-- name: CreateURL :one
INSERT INTO urls (original_url, shortened_code, user_id)
VALUES ($1, $2, $3)
RETURNING id, original_url, shortened_code, user_id, created_at;

-- Get the original URL by shortened code
-- name: GetURLByCode :one
SELECT id, original_url, shortened_code, user_id, created_at
FROM urls
WHERE shortened_code = $1;

-- Get all URLs created by a specific user
-- name: GetURLsByUserID :many
SELECT id, original_url, shortened_code, user_id, created_at
FROM urls
WHERE user_id = $1;