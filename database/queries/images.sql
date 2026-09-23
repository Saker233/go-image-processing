-- name: CreateImage :one
INSERT INTO images (
    user_id,
    original_filename,
    original_s3_key,
    mime_type,
    original_size,
    width,
    height
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
)
RETURNING *;


-- name: GetImageByID :one
SELECT *
FROM images
WHERE id = $1
LIMIT 1;


-- name: GetImagesByUser :many
SELECT *
FROM images
WHERE user_id = $1
ORDER BY created_at DESC;


-- name: UpdateProcessedImage :one
UPDATE images
SET
    processed_s3_key = $2,
    processed_size = $3
WHERE id = $1
RETURNING *;


-- name: DeleteImage :exec
DELETE FROM images
WHERE id = $1;