-- name: CreateAttachedMessages :copyfrom
INSERT INTO t_attached_messages (message_id, attachable_item_id) VALUES ($1, $2);

-- name: CreateAttachedMessage :one
INSERT INTO t_attached_messages (message_id, attachable_item_id) VALUES ($1, $2) RETURNING *;

-- name: DeleteAttachedMessage :exec
DELETE FROM t_attached_messages WHERE message_id = $1 AND attachable_item_id = $2;

-- name: GetAttachableItemsOnMessageID :many
SELECT sqlc.embed(t_attached_messages), sqlc.embed(t_attachable_items), sqlc.embed(m_mime_types), sqlc.embed(t_images), sqlc.embed(t_files) FROM t_attached_messages
INNER JOIN t_attachable_items ON t_attached_messages.attachable_item_id = t_attachable_items.attachable_item_id
INNER JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE message_id = $1
ORDER BY
	t_attached_messages_pkey DESC
LIMIT $2 OFFSET $3;

