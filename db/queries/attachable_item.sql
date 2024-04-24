-- name: CreateAttachableItems :copyfrom
INSERT INTO t_attachable_items (size, mime_type_id) VALUES ($1, $2);

-- name: CreateAttachableItem :one
INSERT INTO t_attachable_items (size, mime_type_id) VALUES ($1, $2) RETURNING *;

-- name: DeleteAttachableItem :exec
DELETE FROM t_attachable_items WHERE attachable_item_id = $1;

-- name: FindAttachableItemByID :one
SELECT * FROM t_attachable_items WHERE attachable_item_id = $1;

-- name: FindAttachableItemByIDWithMimeType :one
SELECT sqlc.embed(t_attachable_items), sqlc.embed(m_mime_types) FROM t_attachable_items
INNER JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
WHERE attachable_item_id = $1;

-- name: GetAttachableItemsByMimeTypeID :many
SELECT * FROM t_attachable_items WHERE mime_type_id = $1
ORDER BY
	t_attachable_items_pkey DESC
LIMIT $2 OFFSET $3;

-- name: CountAttachableItems :one
SELECT COUNT(*) FROM t_attachable_items;

-- name: GetAttachableItemsByMimeTypeIDWithMimeType :many
SELECT sqlc.embed(t_attachable_items), sqlc.embed(m_mime_types) FROM t_attachable_items
INNER JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
WHERE t_attachable_items.mime_type_id = $1
ORDER BY
	t_attachable_items_pkey DESC
LIMIT $2 OFFSET $3;

-- name: CountAttachableItemsByMimeTypeID :one
SELECT COUNT(*) FROM t_attachable_items WHERE mime_type_id = $1;
