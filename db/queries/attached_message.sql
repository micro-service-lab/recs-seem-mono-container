-- name: CreateAttachedMessages :copyfrom
INSERT INTO t_attached_messages (message_id, attachable_item_id) VALUES ($1, $2);

-- name: CreateAttachedMessage :one
INSERT INTO t_attached_messages (message_id, attachable_item_id) VALUES ($1, $2) RETURNING *;

-- name: DeleteAttachedMessage :exec
DELETE FROM t_attached_messages WHERE message_id = $1 AND attachable_item_id = $2;

-- name: GetAttachableItemsOnMessage :many
SELECT sqlc.embed(t_attached_messages), sqlc.embed(t_attachable_items), sqlc.embed(m_mime_types), sqlc.embed(t_images), sqlc.embed(t_files) FROM t_attached_messages
LEFT JOIN t_attachable_items ON t_attached_messages.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE message_id = $1
ORDER BY
	t_attached_messages_pkey DESC
LIMIT $2 OFFSET $3;

-- name: CountAttachableItemsOnMessage :one
SELECT COUNT(*) FROM t_attached_messages WHERE message_id = $1;

-- name: GetAttachedMessagesOnChatRoom :many
SELECT sqlc.embed(t_attached_messages), sqlc.embed(t_attachable_items), sqlc.embed(m_mime_types), sqlc.embed(t_images), sqlc.embed(t_files) FROM t_attached_messages
LEFT JOIN t_attachable_items ON t_attached_messages.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE message_id IN (
	SELECT message_id FROM t_messages WHERE chat_room_id = $1
)
ORDER BY
	t_attached_messages_pkey DESC
LIMIT $2 OFFSET $3;

-- name: CountAttachedMessagesOnChatRoom :one
SELECT COUNT(*) FROM t_attached_messages
WHERE message_id IN (
	SELECT message_id FROM t_messages WHERE chat_room_id = $1
);

