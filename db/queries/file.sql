-- name: CreateFiles :copyfrom
INSERT INTO t_files (attachable_item_id) VALUES ($1);

-- name: CreateFile :one
INSERT INTO t_files (attachable_item_id) VALUES ($1) RETURNING *;

-- name: DeleteFile :exec
DELETE FROM t_files WHERE file_id = $1;

-- name: FindFileByID :one
SELECT * FROM t_files WHERE file_id = $1;

-- name: FindFileByIDWithAttachableItem :one
SELECT sqlc.embed(t_files), sqlc.embed(t_attachable_items) FROM t_files
INNER JOIN t_attachable_items ON t_files.attachable_item_id = t_attachable_items.attachable_item_id
WHERE file_id = $1;

-- name: GetFiles :many
SELECT * FROM t_files
ORDER BY
	t_files_pkey DESC
LIMIT $1 OFFSET $2;

-- name: GetFilesWithAttachableItem :many
SELECT sqlc.embed(t_files), sqlc.embed(t_attachable_items) FROM t_files
INNER JOIN t_attachable_items ON t_files.attachable_item_id = t_attachable_items.attachable_item_id
ORDER BY
	t_files_pkey DESC
LIMIT $1 OFFSET $2;

-- name: CountFiles :one
SELECT COUNT(*) FROM t_files;
