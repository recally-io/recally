-- name: CreateFile :one
INSERT INTO files (
    original_url,
    user_id,
    s3_key,
    s3_url,
    file_name,
    file_type,
    file_size,
    file_hash,
    metadata
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: GetFileByID :one
SELECT * FROM files
WHERE id = $1;

-- name: GetFileByOriginalURL :one
SELECT * FROM files
WHERE original_url = $1 
AND (user_id = $2 OR user_id = sqlc.narg('dummy_user_id'))  ;

-- name: GetFileByS3Key :one
SELECT * FROM files
WHERE s3_key = $1
AND (user_id = $2 OR user_id = sqlc.narg('dummy_user_id'))
;

-- name: ListFiles :many
SELECT * FROM files
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: SearchFilesByType :many
SELECT * FROM files
WHERE file_type = $1
AND user_id = $2
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- name: UpdateFile :one
UPDATE files
SET 
    s3_url = COALESCE($2, s3_url),
    file_name = COALESCE($3, file_name),
    file_type = COALESCE($4, file_type),
    file_size = COALESCE($5, file_size),
    metadata = COALESCE($6, metadata)
WHERE id = $1
RETURNING *;

-- name: DeleteFile :exec
DELETE FROM files
WHERE id = $1;

-- name: DeleteFileByOriginalURL :exec
DELETE FROM files
WHERE original_url = $1
AND (user_id = $2 OR user_id = sqlc.narg('dummy_user_id'));

-- name: DeleteFileByS3Key :exec
DELETE FROM files
WHERE s3_key = $1
AND (user_id = $2 OR user_id = sqlc.narg('dummy_user_id'));

