-- name: CreateAttachedMessages :copyfrom
INSERT INTO t_attached_messages (message_id, attachable_item_id) VALUES ($1, $2);

-- name: CreateAttachedMessage :one
INSERT INTO t_attached_messages (message_id, attachable_item_id) VALUES ($1, $2) RETURNING *;

-- name: DeleteAttachedMessage :execrows
DELETE FROM t_attached_messages WHERE attached_message_id = $1;

-- name: DeleteAttachedMessagesOnMessage :execrows
DELETE FROM t_attached_messages WHERE message_id = $1;

-- name: DeleteAttachedMessagesOnMessages :execrows
DELETE FROM t_attached_messages WHERE message_id = ANY($1::uuid[]);

-- name: GetAttachedItemsOnMessage :many
SELECT t_attached_messages.*, t_attachable_items.url attached_item_url,
t_attachable_items.size attached_item_size, t_attachable_items.mime_type_id attached_item_mime_type_id,
t_attachable_items.owner_id attached_item_owner_id, t_attachable_items.from_outer attached_item_from_outer
FROM t_attached_messages
LEFT JOIN t_attachable_items ON t_attached_messages.attachable_item_id = t_attachable_items.attachable_item_id
WHERE message_id = $1
ORDER BY
	t_attached_messages_pkey ASC;

-- name: GetAttachedItemsOnMessageUseNumberedPaginate :many
SELECT t_attached_messages.*, t_attachable_items.url attached_item_url,
t_attachable_items.size attached_item_size, t_attachable_items.mime_type_id attached_item_mime_type_id,
t_attachable_items.owner_id attached_item_owner_id, t_attachable_items.from_outer attached_item_from_outer
FROM t_attached_messages
LEFT JOIN t_attachable_items ON t_attached_messages.attachable_item_id = t_attachable_items.attachable_item_id
WHERE message_id = $1
ORDER BY
	t_attached_messages_pkey ASC
LIMIT $2 OFFSET $3;

-- name: GetAttachedItemsOnMessageUseKeysetPaginate :many
SELECT t_attached_messages.*, t_attachable_items.url attached_item_url,
t_attachable_items.size attached_item_size, t_attachable_items.mime_type_id attached_item_mime_type_id,
t_attachable_items.owner_id attached_item_owner_id, t_attachable_items.from_outer attached_item_from_outer
FROM t_attached_messages
LEFT JOIN t_attachable_items ON t_attached_messages.attachable_item_id = t_attachable_items.attachable_item_id
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

-- name: GetPluralAttachedItemsOnMessage :many
SELECT t_attached_messages.*, t_attachable_items.url attached_item_url,
t_attachable_items.size attached_item_size, t_attachable_items.mime_type_id attached_item_mime_type_id,
t_attachable_items.owner_id attached_item_owner_id, t_attachable_items.from_outer attached_item_from_outer
FROM t_attached_messages
LEFT JOIN t_attachable_items ON t_attached_messages.attachable_item_id = t_attachable_items.attachable_item_id
WHERE message_id = ANY(@message_ids::uuid[])
ORDER BY
	t_attached_messages_pkey ASC;

-- name: GetPluralAttachedItemsOnMessageUseNumberedPaginate :many
SELECT t_attached_messages.*, t_attachable_items.url attached_item_url,
t_attachable_items.size attached_item_size, t_attachable_items.mime_type_id attached_item_mime_type_id,
t_attachable_items.owner_id attached_item_owner_id, t_attachable_items.from_outer attached_item_from_outer
FROM t_attached_messages
LEFT JOIN t_attachable_items ON t_attached_messages.attachable_item_id = t_attachable_items.attachable_item_id
WHERE message_id = ANY(@message_ids::uuid[])
ORDER BY
	t_attached_messages_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetAttachedItemsOnMessageWithMimeType :many
SELECT t_attached_messages.*, t_attachable_items.url attached_item_url,
t_attachable_items.size attached_item_size, t_attachable_items.mime_type_id attached_item_mime_type_id,
t_attachable_items.owner_id attached_item_owner_id, t_attachable_items.from_outer attached_item_from_outer,
m_mime_types.name mime_type_name, m_mime_types.key mime_type_key, m_mime_types.kind mime_type_kind
FROM t_attached_messages
LEFT JOIN t_attachable_items ON t_attached_messages.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
WHERE message_id = $1
ORDER BY
	t_attached_messages_pkey ASC;

-- name: GetAttachedItemsOnMessageWithMimeTypeUseNumberedPaginate :many
SELECT t_attached_messages.*, t_attachable_items.url attached_item_url,
t_attachable_items.size attached_item_size, t_attachable_items.mime_type_id attached_item_mime_type_id,
t_attachable_items.owner_id attached_item_owner_id, t_attachable_items.from_outer attached_item_from_outer,
m_mime_types.name mime_type_name, m_mime_types.key mime_type_key, m_mime_types.kind mime_type_kind
FROM t_attached_messages
LEFT JOIN t_attachable_items ON t_attached_messages.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
WHERE message_id = $1
ORDER BY
	t_attached_messages_pkey ASC
LIMIT $2 OFFSET $3;

-- name: GetAttachedItemsOnMessageWithMimeTypeUseKeysetPaginate :many
SELECT t_attached_messages.*, t_attachable_items.url attached_item_url,
t_attachable_items.size attached_item_size, t_attachable_items.mime_type_id attached_item_mime_type_id,
t_attachable_items.owner_id attached_item_owner_id, t_attachable_items.from_outer attached_item_from_outer,
m_mime_types.name mime_type_name, m_mime_types.key mime_type_key, m_mime_types.kind mime_type_kind
FROM t_attached_messages
LEFT JOIN t_attachable_items ON t_attached_messages.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
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

-- name: GetPluralAttachedItemsOnMessageWithMimeType :many
SELECT t_attached_messages.*, t_attachable_items.url attached_item_url,
t_attachable_items.size attached_item_size, t_attachable_items.mime_type_id attached_item_mime_type_id,
t_attachable_items.owner_id attached_item_owner_id, t_attachable_items.from_outer attached_item_from_outer,
m_mime_types.name mime_type_name, m_mime_types.key mime_type_key, m_mime_types.kind mime_type_kind
FROM t_attached_messages
LEFT JOIN t_attachable_items ON t_attached_messages.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
WHERE message_id = ANY(@message_ids::uuid[])
ORDER BY
	t_attached_messages_pkey ASC;

-- name: GetPluralAttachedItemsOnMessageWithMimeTypeUseNumberedPaginate :many
SELECT t_attached_messages.*, t_attachable_items.url attached_item_url,
t_attachable_items.size attached_item_size, t_attachable_items.mime_type_id attached_item_mime_type_id,
t_attachable_items.owner_id attached_item_owner_id, t_attachable_items.from_outer attached_item_from_outer,
m_mime_types.name mime_type_name, m_mime_types.key mime_type_key, m_mime_types.kind mime_type_kind
FROM t_attached_messages
LEFT JOIN t_attachable_items ON t_attached_messages.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
WHERE message_id = ANY(@message_ids::uuid[])
ORDER BY
	t_attached_messages_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountAttachedItemsOnMessage :one
SELECT COUNT(*) FROM t_attached_messages WHERE message_id = $1;

-- name: GetAttachedItemsOnChatRoom :many
SELECT t_attached_messages.*, t_attachable_items.url attached_item_url,
t_attachable_items.size attached_item_size, t_attachable_items.mime_type_id attached_item_mime_type_id,
t_attachable_items.owner_id attached_item_owner_id, t_attachable_items.from_outer attached_item_from_outer
FROM t_attached_messages
LEFT JOIN t_attachable_items ON t_attached_messages.attachable_item_id = t_attachable_items.attachable_item_id
WHERE message_id IN (
	SELECT message_id FROM t_messages WHERE chat_room_id = $1
)
ORDER BY
	t_attached_messages_pkey ASC;

-- name: GetAttachedItemsOnChatRoomUseNumberedPaginate :many
SELECT t_attached_messages.*, t_attachable_items.url attached_item_url,
t_attachable_items.size attached_item_size, t_attachable_items.mime_type_id attached_item_mime_type_id,
t_attachable_items.owner_id attached_item_owner_id, t_attachable_items.from_outer attached_item_from_outer
FROM t_attached_messages
LEFT JOIN t_attachable_items ON t_attached_messages.attachable_item_id = t_attachable_items.attachable_item_id
WHERE message_id IN (
	SELECT message_id FROM t_messages WHERE chat_room_id = $1
)
ORDER BY
	t_attached_messages_pkey ASC
LIMIT $2 OFFSET $3;

-- name: GetAttachedItemsOnChatRoomUseKeysetPaginate :many
SELECT t_attached_messages.*, t_attachable_items.url attached_item_url,
t_attachable_items.size attached_item_size, t_attachable_items.mime_type_id attached_item_mime_type_id,
t_attachable_items.owner_id attached_item_owner_id, t_attachable_items.from_outer attached_item_from_outer
FROM t_attached_messages
LEFT JOIN t_attachable_items ON t_attached_messages.attachable_item_id = t_attachable_items.attachable_item_id
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

-- name: GetPluralAttachedItemsOnChatRoom :many
SELECT t_attached_messages.*, t_attachable_items.url attached_item_url,
t_attachable_items.size attached_item_size, t_attachable_items.mime_type_id attached_item_mime_type_id,
t_attachable_items.owner_id attached_item_owner_id, t_attachable_items.from_outer attached_item_from_outer
FROM t_attached_messages
LEFT JOIN t_attachable_items ON t_attached_messages.attachable_item_id = t_attachable_items.attachable_item_id
WHERE message_id IN (
	SELECT message_id FROM t_messages WHERE chat_room_id = ANY(@chat_room_ids::uuid[])
)
ORDER BY
	t_attached_messages_pkey ASC;

-- name: GetPluralAttachedItemsOnChatRoomUseNumberedPaginate :many
SELECT t_attached_messages.*, t_attachable_items.url attached_item_url,
t_attachable_items.size attached_item_size, t_attachable_items.mime_type_id attached_item_mime_type_id,
t_attachable_items.owner_id attached_item_owner_id, t_attachable_items.from_outer attached_item_from_outer
FROM t_attached_messages
LEFT JOIN t_attachable_items ON t_attached_messages.attachable_item_id = t_attachable_items.attachable_item_id
WHERE message_id IN (
	SELECT message_id FROM t_messages WHERE chat_room_id = ANY(@chat_room_ids::uuid[])
)
ORDER BY
	t_attached_messages_pkey ASC
LIMIT $1 OFFSET $2;


-- name: CountAttachedItemsOnChatRoom :one
SELECT COUNT(*) FROM t_attached_messages
WHERE message_id IN (
	SELECT message_id FROM t_messages WHERE chat_room_id = $1
);

