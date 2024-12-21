-- Increment click count
-- name: IncrementClickCount :exec
UPDATE click_counts
SET total_clicks = total_clicks + $1
    FROM urls
WHERE click_counts.url_id = urls.id AND urls.shortened_code = $2;

-- Create click count record
-- name: InsertClickCount :exec
INSERT INTO click_counts (url_id)
VALUES ($1);