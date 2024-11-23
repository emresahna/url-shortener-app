-- Create a new shortened URL for a specific user
-- name: CreateURL :one
INSERT INTO urls (original_url, shortened_code, user_id)
VALUES ($1, $2, $3)
RETURNING id, original_url, shortened_code, user_id;

-- Get the original URL by shortened code
-- name: GetURLByCode :one
SELECT original_url
FROM urls
WHERE shortened_code = $1;

-- Get url ID by short code
-- name: GetIDByShortCode :one
SELECT id
FROM urls
where shortened_code = $1;

-- Delete expired url by short code
-- name: DeleteExpiredUrlByShortCode :exec
UPDATE urls
SET is_deleted = true, is_active = false, deleted_at = $1
WHERE shortened_code = $2;