-- Increment click count
-- name: IncrementClickCount :exec
UPDATE click_counts
SET total_clicks = total_clicks + $1
WHERE url_id = $2;

-- Create click count record
-- name: InsertClickCount :exec
INSERT INTO click_counts (url_id)
VALUES ($1);