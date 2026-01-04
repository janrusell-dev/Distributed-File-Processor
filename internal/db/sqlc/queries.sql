-- name: CreateFile :exec

INSERT INTO files (id, filename, size, mime_type, status)
VALUES ($1, $2, $3, $4, $5);

-- name: GetFile :one
SELECT * FROM files WHERE id = $1;

-- name: UpdateFileStatus :exec
UPDATE files
SET status = $2
WHERE id = $1;