-- name: CreateAttachedMessages :copyfrom
INSERT INTO t_attached_messages (message_id, attachable_item_id) VALUES ($1, $2);

-- name: CreateAttachedMessage :one
INSERT INTO t_attached_messages (message_id, attachable_item_id) VALUES ($1, $2) RETURNING *;

-- name: DeleteAttachedMessage :exec
DELETE FROM t_attached_messages WHERE message_id = $1 AND attachable_item_id = $2;

-- name: GetAttachableItemsOnMessage :many
SELECT t_attachable_items.*, t_images.image_id, t_images.height as image_height, t_images.width as image_width, t_files.file_id FROM t_attached_messages
LEFT JOIN t_attachable_items ON t_attached_messages.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE message_id = $1
ORDER BY
	t_attached_messages_pkey ASC;

-- name: GetAttachableItemsOnMessageUseNumberedPaginate :many
SELECT t_attachable_items.*, t_images.image_id, t_images.height as image_height, t_images.width as image_width, t_files.file_id FROM t_attached_messages
LEFT JOIN t_attachable_items ON t_attached_messages.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE message_id = $1
ORDER BY
	t_attached_messages_pkey ASC
LIMIT $2 OFFSET $3;

-- name: GetAttachableItemsOnMessageUseKeysetPaginate :many
SELECT t_attachable_items.*, t_images.image_id, t_images.height as image_height, t_images.width as image_width, t_files.file_id FROM t_attached_messages
LEFT JOIN t_attachable_items ON t_attached_messages.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE message_id = $1
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			t_attached_messages_pkey > @cursor::int
		WHEN 'prev' THEN
			t_attached_messages_pkey < @cursor::int
	END
ORDER BY
	CASE WHEN @cursor_direction::text = 'next' THEN t_attached_messages_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN t_attached_messages_pkey END DESC
LIMIT $2;

-- name: GetPluralAttachableItemsOnMessage :many
SELECT t_attachable_items.*, t_images.image_id, t_images.height as image_height, t_images.width as image_width, t_files.file_id FROM t_attached_messages
LEFT JOIN t_attachable_items ON t_attached_messages.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE message_id = ANY(@message_ids::uuid[])
ORDER BY
	t_attached_messages_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountAttachableItemsOnMessage :one
SELECT COUNT(*) FROM t_attached_messages WHERE message_id = $1;

-- name: GetAttachedMessagesOnChatRoom :many
SELECT t_attachable_items.*, t_images.image_id, t_images.height as image_height, t_images.width as image_width, t_files.file_id FROM t_attached_messages
LEFT JOIN t_attachable_items ON t_attached_messages.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE message_id IN (
	SELECT message_id FROM t_messages WHERE chat_room_id = $1
)
ORDER BY
	t_attached_messages_pkey ASC;

-- name: GetAttachedMessagesOnChatRoomUseNumberedPaginate :many
SELECT t_attachable_items.*, t_images.image_id, t_images.height as image_height, t_images.width as image_width, t_files.file_id FROM t_attached_messages
LEFT JOIN t_attachable_items ON t_attached_messages.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE message_id IN (
	SELECT message_id FROM t_messages WHERE chat_room_id = $1
)
ORDER BY
	t_attached_messages_pkey ASC
LIMIT $2 OFFSET $3;

-- name: GetAttachedMessagesOnChatRoomUseKeysetPaginate :many
SELECT t_attachable_items.*, t_images.image_id, t_images.height as image_height, t_images.width as image_width, t_files.file_id FROM t_attached_messages
LEFT JOIN t_attachable_items ON t_attached_messages.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE message_id IN (
	SELECT message_id FROM t_messages WHERE chat_room_id = $1
)
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			t_attached_messages_pkey > @cursor::int
		WHEN 'prev' THEN
			t_attached_messages_pkey < @cursor::int
	END
ORDER BY
	CASE WHEN @cursor_direction::text = 'next' THEN t_attached_messages_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN t_attached_messages_pkey END DESC
LIMIT $2;

-- name: GetPluralAttachedMessagesOnChatRoom :many
SELECT t_attachable_items.*, t_images.image_id, t_images.height as image_height, t_images.width as image_width, t_files.file_id FROM t_attached_messages
LEFT JOIN t_attachable_items ON t_attached_messages.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE message_id IN (
	SELECT message_id FROM t_messages WHERE chat_room_id = ANY(@chat_room_ids::uuid[])
)
ORDER BY
	t_attached_messages_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountAttachedMessagesOnChatRoom :one
SELECT COUNT(*) FROM t_attached_messages
WHERE message_id IN (
	SELECT message_id FROM t_messages WHERE chat_room_id = $1
);

