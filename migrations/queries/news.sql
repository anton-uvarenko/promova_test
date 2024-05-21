-- name: AddNews :one
INSERT INTO news (
  title,
  content,
  created_at,
  updated_at
) VALUES (
  $1,
  $2,
  NOW(),
  NOW()
)
RETURNING id;

-- name: UpdateNews :exec
UPDATE news
SET 
  title = $2,
  content = $3,
  updated_at = NOW()
WHERE
  id = $1;

-- name: DeleteNews :exec
DELETE FROM news
WHERE id = $1;

-- name: GetAllNews :many
SELECT * FROM news;

-- name: GetNewsById :one
SELECT * FROM news
WHERE id = $1;

