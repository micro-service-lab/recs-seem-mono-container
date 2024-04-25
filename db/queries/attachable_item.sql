-- name: CreateAttachableItems :copyfrom
INSERT INTO t_attachable_items (url, size, mime_type_id) VALUES ($1, $2, $3);

-- name: CreateAttachableItem :one
INSERT INTO t_attachable_items (url, size, mime_type_id) VALUES ($1, $2, $3) RETURNING *;

-- name: DeleteAttachableItem :exec
DELETE FROM t_attachable_items WHERE attachable_item_id = $1;

-- name: FindAttachableItemByID :one
SELECT sqlc.embed(t_attachable_items), sqlc.embed(t_images), sqlc.embed(t_files) FROM t_attachable_items
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE t_attachable_items.attachable_item_id = $1;

-- name: FindAttachableItemByIDWithMimeType :one
SELECT sqlc.embed(t_attachable_items), sqlc.embed(m_mime_types), sqlc.embed(t_images), sqlc.embed(t_files) FROM t_attachable_items
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE t_attachable_items.attachable_item_id = $1;

-- name: GetAttachableItems :many
SELECT sqlc.embed(t_attachable_items), sqlc.embed(t_images), sqlc.embed(t_files) FROM t_attachable_items
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE
	CASE WHEN @where_mime_type_id::boolean = true THEN mime_type_id = @mime_type_id ELSE TRUE END
ORDER BY
	t_attachable_items_pkey DESC
LIMIT $1 OFFSET $2;

-- name: GetAttachableItemsByMimeTypeIDWithMimeType :many
SELECT sqlc.embed(t_attachable_items), sqlc.embed(m_mime_types), sqlc.embed(t_images), sqlc.embed(t_files) FROM t_attachable_items
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.where_mime_type_id
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE
	CASE WHEN @where_mime_type_id::boolean = true THEN t_attachable_items.mime_type_id = @mime_type_id ELSE TRUE END
ORDER BY
	t_attachable_items_pkey DESC
LIMIT $1 OFFSET $2;

-- name: CountAttachableItems :one
SELECT COUNT(*) FROM t_attachable_items
WHERE
	CASE WHEN @where_mime_type_id::boolean = true THEN mime_type_id = @mime_type_id ELSE TRUE END;
