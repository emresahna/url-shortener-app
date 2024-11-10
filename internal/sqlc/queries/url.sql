-- Create a new shortened URL for a specific user
-- name: CreateURL :one
INSERT INTO urls (original_url, shortened_code, user_id, expire_time)
VALUES ($1, $2, $3, $4)
RETURNING id, original_url, shortened_code, user_id, expire_time;

-- Get the original URL by shortened code
-- name: GetURLByCode :one
SELECT *
FROM urls
WHERE shortened_code = $1;

-- Get all URLs created by a specific user
-- name: GetURLsByUserID :many
SELECT *
FROM urls
WHERE user_id = $1;

-- Delete URL by url ID
-- name: DeleteURLByID :exec
UPDATE urls
SET deleted_at = NOW(), is_deleted = True, is_active = False
WHERE id = $1;