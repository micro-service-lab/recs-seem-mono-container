-- name: CreateFiles :copyfrom
INSERT INTO t_files (attachable_item_id) VALUES ($1);

-- name: CreateFile :one
INSERT INTO t_files (attachable_item_id) VALUES ($1) RETURNING *;

-- name: DeleteFile :execrows
DELETE FROM t_files WHERE file_id = $1;

-- name: PluralDeleteFiles :execrows
DELETE FROM t_files WHERE file_id = ANY(@file_ids::uuid[]);

-- name: FindFileByID :one
SELECT * FROM t_files WHERE file_id = $1;

-- name: FindFileByIDWithAttachableItem :one
SELECT t_files.*, t_attachable_items.owner_id owner_id, t_attachable_items.from_outer from_outer,
t_attachable_items.url url, t_attachable_items.size size, t_attachable_items.mime_type_id mime_type_id
FROM t_files
LEFT JOIN t_attachable_items ON t_files.attachable_item_id = t_attachable_items.attachable_item_id
WHERE file_id = $1;

-- name: GetFiles :many
SELECT * FROM t_files
ORDER BY
	t_files_pkey ASC;

-- name: GetFilesUseNumberedPaginate :many
SELECT * FROM t_files
ORDER BY
	t_files_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetFilesUseKeysetPaginate :many
SELECT * FROM t_files
WHERE
	CASE @cursor_direction::text
		WHEN 'next' THEN
			t_files_pkey > @cursor::int
		WHEN 'prev' THEN
			t_files_pkey < @cursor::int
	END
ORDER BY
	CASE WHEN @cursor_direction::text = 'next' THEN t_files_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN t_files_pkey END DESC
LIMIT $1;

-- name: GetPluralFiles :many
SELECT * FROM t_files
WHERE file_id = ANY(@file_ids::uuid[])
ORDER BY
	t_files_pkey ASC;

-- name: GetPluralFilesUseNumberedPaginate :many
SELECT * FROM t_files
WHERE file_id = ANY(@file_ids::uuid[])
ORDER BY
	t_files_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetFilesWithAttachableItem :many
SELECT t_files.*, t_attachable_items.owner_id owner_id, t_attachable_items.from_outer from_outer,
t_attachable_items.url url, t_attachable_items.size size, t_attachable_items.mime_type_id mime_type_id
FROM t_files
LEFT JOIN t_attachable_items ON t_files.attachable_item_id = t_attachable_items.attachable_item_id
ORDER BY
	t_files_pkey ASC;

-- name: GetFilesWithAttachableItemUseNumberedPaginate :many
SELECT t_files.*, t_attachable_items.owner_id owner_id, t_attachable_items.from_outer from_outer,
t_attachable_items.url url, t_attachable_items.size size, t_attachable_items.mime_type_id mime_type_id
FROM t_files
LEFT JOIN t_attachable_items ON t_files.attachable_item_id = t_attachable_items.attachable_item_id
ORDER BY
	t_files_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetFilesWithAttachableItemUseKeysetPaginate :many
SELECT t_files.*, t_attachable_items.owner_id owner_id, t_attachable_items.from_outer from_outer,
t_attachable_items.url url, t_attachable_items.size size, t_attachable_items.mime_type_id mime_type_id
FROM t_files
LEFT JOIN t_attachable_items ON t_files.attachable_item_id = t_attachable_items.attachable_item_id
WHERE
	CASE @cursor_direction::text
		WHEN 'next' THEN
			t_files_pkey > @cursor::int
		WHEN 'prev' THEN
			t_files_pkey < @cursor::int
	END
ORDER BY
	CASE WHEN @cursor_direction::text = 'next' THEN t_files_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN t_files_pkey END DESC
LIMIT $1;

-- name: GetPluralFilesWithAttachableItem :many
SELECT t_files.*, t_attachable_items.owner_id owner_id, t_attachable_items.from_outer from_outer,
t_attachable_items.url url, t_attachable_items.size size, t_attachable_items.mime_type_id mime_type_id
FROM t_files
LEFT JOIN t_attachable_items ON t_files.attachable_item_id = t_attachable_items.attachable_item_id
WHERE file_id = ANY(@file_ids::uuid[])
ORDER BY
	t_files_pkey ASC;

-- name: GetPluralFilesWithAttachableItemUseNumberedPaginate :many
SELECT t_files.*, t_attachable_items.owner_id owner_id, t_attachable_items.from_outer from_outer,
t_attachable_items.url url, t_attachable_items.size size, t_attachable_items.mime_type_id mime_type_id
FROM t_files
LEFT JOIN t_attachable_items ON t_files.attachable_item_id = t_attachable_items.attachable_item_id
WHERE file_id = ANY(@file_ids::uuid[])
ORDER BY
	t_files_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountFiles :one
SELECT COUNT(*) FROM t_files;
