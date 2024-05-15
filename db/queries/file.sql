-- name: CreateFiles :copyfrom
INSERT INTO t_files (attachable_item_id) VALUES ($1);

-- name: CreateFile :one
INSERT INTO t_files (attachable_item_id) VALUES ($1) RETURNING *;

-- name: DeleteFile :execrows
DELETE FROM t_files WHERE file_id = $1;

-- name: PluralDeleteFiles :execrows
DELETE FROM t_files WHERE file_id = ANY($1::uuid[]);

-- name: FindFileByID :one
SELECT * FROM t_files WHERE file_id = $1;

-- name: FindFileByIDWithAttachableItem :one
SELECT sqlc.embed(t_files), sqlc.embed(t_attachable_items), sqlc.embed(m_mime_types) FROM t_files
LEFT JOIN t_attachable_items ON t_files.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
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
WHERE attachable_item_id = ANY(@attachable_item_ids::uuid[])
ORDER BY
	t_files_pkey ASC;

-- name: GetPluralFilesUseNumberedPaginate :many
SELECT * FROM t_files
WHERE attachable_item_id = ANY(@attachable_item_ids::uuid[])
ORDER BY
	t_files_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetFilesWithAttachableItem :many
SELECT sqlc.embed(t_files), sqlc.embed(t_attachable_items), sqlc.embed(m_mime_types) FROM t_files
LEFT JOIN t_attachable_items ON t_files.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
ORDER BY
	t_files_pkey ASC;

-- name: GetFilesWithAttachableItemUseNumberedPaginate :many
SELECT sqlc.embed(t_files), sqlc.embed(t_attachable_items), sqlc.embed(m_mime_types) FROM t_files
LEFT JOIN t_attachable_items ON t_files.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
ORDER BY
	t_files_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetFilesWithAttachableItemUseKeysetPaginate :many
SELECT sqlc.embed(t_files), sqlc.embed(t_attachable_items), sqlc.embed(m_mime_types) FROM t_files
LEFT JOIN t_attachable_items ON t_files.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
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
SELECT sqlc.embed(t_files), sqlc.embed(t_attachable_items), sqlc.embed(m_mime_types) FROM t_files
LEFT JOIN t_attachable_items ON t_files.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
WHERE attachable_item_id = ANY(@attachable_item_ids::uuid[])
ORDER BY
	t_files_pkey ASC;

-- name: GetPluralFilesWithAttachableItemUseNumberedPaginate :many
SELECT sqlc.embed(t_files), sqlc.embed(t_attachable_items), sqlc.embed(m_mime_types) FROM t_files
LEFT JOIN t_attachable_items ON t_files.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
WHERE attachable_item_id = ANY(@attachable_item_ids::uuid[])
ORDER BY
	t_files_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountFiles :one
SELECT COUNT(*) FROM t_files;
