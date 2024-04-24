-- name: CreateMimeTypes :copyfrom
INSERT INTO m_mime_types (name, key) VALUES ($1, $2);

-- name: CreateMimeType :one
INSERT INTO m_mime_types (name, key) VALUES ($1, $2) RETURNING *;

-- name: DeleteMimeType :exec
DELETE FROM m_mime_types WHERE mime_type_id = $1;

-- name: DeleteMimeTypeByKey :exec
DELETE FROM m_mime_types WHERE key = $1;

-- name: FindMimeTypeByID :one
SELECT * FROM m_mime_types WHERE mime_type_id = $1;

-- name: FindMimeTypeByKey :one
SELECT * FROM m_mime_types WHERE key = $1;

-- name: GetMimeTypes :many
SELECT * FROM m_mime_types
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN name END ASC,
	m_mime_types_pkey DESC
LIMIT $1 OFFSET $2;

-- name: CountMimeTypes :one
SELECT COUNT(*) FROM m_mime_types;
