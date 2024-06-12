-- name: CreateAttachedMessages :copyfrom
INSERT INTO t_attached_messages (message_id, attachable_item_id) VALUES ($1, $2);

-- name: CreateAttachedMessage :one
INSERT INTO t_attached_messages (message_id, attachable_item_id) VALUES ($1, $2) RETURNING *;

-- name: DeleteAttachedMessage :execrows
DELETE FROM t_attached_messages WHERE attached_message_id = $1;

-- name: DeleteAttachedMessagesOnMessage :execrows
DELETE FROM t_attached_messages WHERE message_id = $1;

-- name: DeleteAttachedMessagesOnMessages :execrows
DELETE FROM t_attached_messages WHERE message_id = ANY(@message_ids::uuid[]);

-- name: PluralDeleteAttachedMessagesOnMessage :execrows
DELETE FROM t_attached_messages WHERE message_id = $1 AND attachable_item_id = ANY(@attachable_item_ids::uuid[]);

-- name: GetAttachedItemsOnMessage :many
SELECT t_attached_messages.*, t_attachable_items.url attached_item_url, t_attachable_items.alias attached_item_alias,
t_attachable_items.size attached_item_size, t_attachable_items.mime_type_id attached_item_mime_type_id,
t_attachable_items.owner_id attached_item_owner_id, t_attachable_items.from_outer attached_item_from_outer,
t_images.image_id attached_image_id, t_images.height attached_image_height,
t_images.width attached_image_width, t_files.file_id attached_file_id
FROM t_attached_messages
LEFT JOIN t_attachable_items ON t_attached_messages.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE message_id = $1
AND
	CASE WHEN @where_in_mime_type::boolean = true THEN t_attachable_items.mime_type_id = ANY(@in_mime_types::uuid[]) ELSE TRUE END
AND
	CASE WHEN @where_is_image::boolean = true THEN EXISTS (SELECT 1 FROM t_images WHERE t_images.attachable_item_id = t_attachable_items.attachable_item_id) ELSE TRUE END
AND
	CASE WHEN @where_is_file::boolean = true THEN EXISTS (SELECT 1 FROM t_files WHERE t_files.attachable_item_id = t_attachable_items.attachable_item_id) ELSE TRUE END
ORDER BY
	t_attached_messages_pkey ASC;

-- name: GetAttachedItemsOnMessageUseNumberedPaginate :many
SELECT t_attached_messages.*, t_attachable_items.url attached_item_url, t_attachable_items.alias attached_item_alias,
t_attachable_items.size attached_item_size, t_attachable_items.mime_type_id attached_item_mime_type_id,
t_attachable_items.owner_id attached_item_owner_id, t_attachable_items.from_outer attached_item_from_outer,
t_images.image_id attached_image_id, t_images.height attached_image_height,
t_images.width attached_image_width, t_files.file_id attached_file_id
FROM t_attached_messages
LEFT JOIN t_attachable_items ON t_attached_messages.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE message_id = $1
AND
	CASE WHEN @where_in_mime_type::boolean = true THEN t_attachable_items.mime_type_id = ANY(@in_mime_types::uuid[]) ELSE TRUE END
	AND
	CASE WHEN @where_is_image::boolean = true THEN EXISTS (SELECT 1 FROM t_images WHERE t_images.attachable_item_id = t_attachable_items.attachable_item_id) ELSE TRUE END
AND
	CASE WHEN @where_is_file::boolean = true THEN EXISTS (SELECT 1 FROM t_files WHERE t_files.attachable_item_id = t_attachable_items.attachable_item_id) ELSE TRUE END
ORDER BY
	t_attached_messages_pkey ASC
LIMIT $2 OFFSET $3;

-- name: GetAttachedItemsOnMessageUseKeysetPaginate :many
SELECT t_attached_messages.*, t_attachable_items.url attached_item_url, t_attachable_items.alias attached_item_alias,
t_attachable_items.size attached_item_size, t_attachable_items.mime_type_id attached_item_mime_type_id,
t_attachable_items.owner_id attached_item_owner_id, t_attachable_items.from_outer attached_item_from_outer,
t_images.image_id attached_image_id, t_images.height attached_image_height,
t_images.width attached_image_width, t_files.file_id attached_file_id
FROM t_attached_messages
LEFT JOIN t_attachable_items ON t_attached_messages.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE message_id = $1
AND
	CASE WHEN @where_in_mime_type::boolean = true THEN t_attachable_items.mime_type_id = ANY(@in_mime_types::uuid[]) ELSE TRUE END
AND
	CASE WHEN @where_is_image::boolean = true THEN EXISTS (SELECT 1 FROM t_images WHERE t_images.attachable_item_id = t_attachable_items.attachable_item_id) ELSE TRUE END
AND
	CASE WHEN @where_is_file::boolean = true THEN EXISTS (SELECT 1 FROM t_files WHERE t_files.attachable_item_id = t_attachable_items.attachable_item_id) ELSE TRUE END
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
SELECT t_attached_messages.*, t_attachable_items.url attached_item_url, t_attachable_items.alias attached_item_alias,
t_attachable_items.size attached_item_size, t_attachable_items.mime_type_id attached_item_mime_type_id,
t_attachable_items.owner_id attached_item_owner_id, t_attachable_items.from_outer attached_item_from_outer,
t_images.image_id attached_image_id, t_images.height attached_image_height,
t_images.width attached_image_width, t_files.file_id attached_file_id
FROM t_attached_messages
LEFT JOIN t_attachable_items ON t_attached_messages.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE message_id = ANY(@message_ids::uuid[])
ORDER BY
	t_attached_messages_pkey ASC;

-- name: GetPluralAttachedItemsOnMessageUseNumberedPaginate :many
SELECT t_attached_messages.*, t_attachable_items.url attached_item_url, t_attachable_items.alias attached_item_alias,
t_attachable_items.size attached_item_size, t_attachable_items.mime_type_id attached_item_mime_type_id,
t_attachable_items.owner_id attached_item_owner_id, t_attachable_items.from_outer attached_item_from_outer,
t_images.image_id attached_image_id, t_images.height attached_image_height,
t_images.width attached_image_width, t_files.file_id attached_file_id
FROM t_attached_messages
LEFT JOIN t_attachable_items ON t_attached_messages.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE message_id = ANY(@message_ids::uuid[])
ORDER BY
	t_attached_messages_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetAttachedItemsOnMessageWithMimeType :many
SELECT t_attached_messages.*, t_attachable_items.url attached_item_url, t_attachable_items.alias attached_item_alias,
t_attachable_items.size attached_item_size, t_attachable_items.mime_type_id attached_item_mime_type_id,
t_attachable_items.owner_id attached_item_owner_id, t_attachable_items.from_outer attached_item_from_outer,
t_images.image_id attached_image_id, t_images.height attached_image_height,
t_images.width attached_image_width, t_files.file_id attached_file_id,
m_mime_types.name mime_type_name, m_mime_types.key mime_type_key, m_mime_types.kind mime_type_kind
FROM t_attached_messages
LEFT JOIN t_attachable_items ON t_attached_messages.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
WHERE message_id = $1
AND
	CASE WHEN @where_in_mime_type::boolean = true THEN t_attachable_items.mime_type_id = ANY(@in_mime_types::uuid[]) ELSE TRUE END
AND
	CASE WHEN @where_is_image::boolean = true THEN EXISTS (SELECT 1 FROM t_images WHERE t_images.attachable_item_id = t_attachable_items.attachable_item_id) ELSE TRUE END
AND
	CASE WHEN @where_is_file::boolean = true THEN EXISTS (SELECT 1 FROM t_files WHERE t_files.attachable_item_id = t_attachable_items.attachable_item_id) ELSE TRUE END
ORDER BY
	t_attached_messages_pkey ASC;

-- name: GetAttachedItemsOnMessageWithMimeTypeUseNumberedPaginate :many
SELECT t_attached_messages.*, t_attachable_items.url attached_item_url, t_attachable_items.alias attached_item_alias,
t_attachable_items.size attached_item_size, t_attachable_items.mime_type_id attached_item_mime_type_id,
t_attachable_items.owner_id attached_item_owner_id, t_attachable_items.from_outer attached_item_from_outer,
t_images.image_id attached_image_id, t_images.height attached_image_height,
t_images.width attached_image_width, t_files.file_id attached_file_id,
m_mime_types.name mime_type_name, m_mime_types.key mime_type_key, m_mime_types.kind mime_type_kind
FROM t_attached_messages
LEFT JOIN t_attachable_items ON t_attached_messages.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
WHERE message_id = $1
AND
	CASE WHEN @where_in_mime_type::boolean = true THEN t_attachable_items.mime_type_id = ANY(@in_mime_types::uuid[]) ELSE TRUE END
AND
	CASE WHEN @where_is_image::boolean = true THEN EXISTS (SELECT 1 FROM t_images WHERE t_images.attachable_item_id = t_attachable_items.attachable_item_id) ELSE TRUE END
AND
	CASE WHEN @where_is_file::boolean = true THEN EXISTS (SELECT 1 FROM t_files WHERE t_files.attachable_item_id = t_attachable_items.attachable_item_id) ELSE TRUE END
ORDER BY
	t_attached_messages_pkey ASC
LIMIT $2 OFFSET $3;

-- name: GetAttachedItemsOnMessageWithMimeTypeUseKeysetPaginate :many
SELECT t_attached_messages.*, t_attachable_items.url attached_item_url, t_attachable_items.alias attached_item_alias,
t_attachable_items.size attached_item_size, t_attachable_items.mime_type_id attached_item_mime_type_id,
t_attachable_items.owner_id attached_item_owner_id, t_attachable_items.from_outer attached_item_from_outer,
t_images.image_id attached_image_id, t_images.height attached_image_height,
t_images.width attached_image_width, t_files.file_id attached_file_id,
m_mime_types.name mime_type_name, m_mime_types.key mime_type_key, m_mime_types.kind mime_type_kind
FROM t_attached_messages
LEFT JOIN t_attachable_items ON t_attached_messages.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
WHERE message_id = $1
AND
	CASE WHEN @where_in_mime_type::boolean = true THEN t_attachable_items.mime_type_id = ANY(@in_mime_types::uuid[]) ELSE TRUE END
AND
	CASE WHEN @where_is_image::boolean = true THEN EXISTS (SELECT 1 FROM t_images WHERE t_images.attachable_item_id = t_attachable_items.attachable_item_id) ELSE TRUE END
AND
	CASE WHEN @where_is_file::boolean = true THEN EXISTS (SELECT 1 FROM t_files WHERE t_files.attachable_item_id = t_attachable_items.attachable_item_id) ELSE TRUE END
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
SELECT t_attached_messages.*, t_attachable_items.url attached_item_url, t_attachable_items.alias attached_item_alias,
t_attachable_items.size attached_item_size, t_attachable_items.mime_type_id attached_item_mime_type_id,
t_attachable_items.owner_id attached_item_owner_id, t_attachable_items.from_outer attached_item_from_outer,
t_images.image_id attached_image_id, t_images.height attached_image_height,
t_images.width attached_image_width, t_files.file_id attached_file_id,
m_mime_types.name mime_type_name, m_mime_types.key mime_type_key, m_mime_types.kind mime_type_kind
FROM t_attached_messages
LEFT JOIN t_attachable_items ON t_attached_messages.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
WHERE message_id = ANY(@message_ids::uuid[])
ORDER BY
	t_attached_messages_pkey ASC;

-- name: GetPluralAttachedItemsOnMessageWithMimeTypeUseNumberedPaginate :many
SELECT t_attached_messages.*, t_attachable_items.url attached_item_url, t_attachable_items.alias attached_item_alias,
t_attachable_items.size attached_item_size, t_attachable_items.mime_type_id attached_item_mime_type_id,
t_attachable_items.owner_id attached_item_owner_id, t_attachable_items.from_outer attached_item_from_outer,
t_images.image_id attached_image_id, t_images.height attached_image_height,
t_images.width attached_image_width, t_files.file_id attached_file_id,
m_mime_types.name mime_type_name, m_mime_types.key mime_type_key, m_mime_types.kind mime_type_kind
FROM t_attached_messages
LEFT JOIN t_attachable_items ON t_attached_messages.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
WHERE message_id = ANY(@message_ids::uuid[])
ORDER BY
	t_attached_messages_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountAttachedItemsOnMessage :one
SELECT COUNT(*) FROM t_attached_messages WHERE message_id = $1
AND
	CASE WHEN @where_in_mime_type::boolean = true THEN t_attachable_items.mime_type_id = ANY(@in_mime_types::uuid[]) ELSE TRUE END
AND
	CASE WHEN @where_is_image::boolean = true THEN EXISTS (SELECT 1 FROM t_images WHERE t_images.attachable_item_id = t_attachable_items.attachable_item_id) ELSE TRUE END
AND
	CASE WHEN @where_is_file::boolean = true THEN EXISTS (SELECT 1 FROM t_files WHERE t_files.attachable_item_id = t_attachable_items.attachable_item_id) ELSE TRUE END;

-- name: GetAttachedItemsOnChatRoom :many
SELECT t_attached_messages.*, t_attachable_items.url attached_item_url, t_attachable_items.alias attached_item_alias,
t_attachable_items.size attached_item_size, t_attachable_items.mime_type_id attached_item_mime_type_id,
t_attachable_items.owner_id attached_item_owner_id, t_attachable_items.from_outer attached_item_from_outer,
t_images.image_id attached_image_id, t_images.height attached_image_height,
t_images.width attached_image_width, t_files.file_id attached_file_id,
t_messages.sender_id message_sender_id, t_messages.chat_room_action_id message_chat_room_action_id, t_messages.body message_body, t_messages.posted_at message_posted_at, t_messages.last_edited_at message_last_edited_at
FROM t_attached_messages
LEFT JOIN t_messages ON t_attached_messages.message_id = t_messages.message_id
LEFT JOIN t_attachable_items ON t_attached_messages.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE t_attached_messages.message_id IN (
	SELECT m.message_id FROM t_messages m WHERE m.chat_room_action_id IN (
		SELECT chat_room_action_id FROM t_chat_room_actions WHERE chat_room_id = $1
	)
)
AND
	CASE WHEN @where_in_mime_type::boolean = true THEN t_attachable_items.mime_type_id = ANY(@in_mime_types::uuid[]) ELSE TRUE END
AND
	CASE WHEN @where_is_image::boolean = true THEN EXISTS (SELECT 1 FROM t_images WHERE t_images.attachable_item_id = t_attachable_items.attachable_item_id) ELSE TRUE END
AND
	CASE WHEN @where_is_file::boolean = true THEN EXISTS (SELECT 1 FROM t_files WHERE t_files.attachable_item_id = t_attachable_items.attachable_item_id) ELSE TRUE END
ORDER BY
	t_attached_messages_pkey ASC;

-- name: GetAttachedItemsOnChatRoomUseNumberedPaginate :many
SELECT t_attached_messages.*, t_attachable_items.url attached_item_url, t_attachable_items.alias attached_item_alias,
t_attachable_items.size attached_item_size, t_attachable_items.mime_type_id attached_item_mime_type_id,
t_attachable_items.owner_id attached_item_owner_id, t_attachable_items.from_outer attached_item_from_outer,
t_images.image_id attached_image_id, t_images.height attached_image_height,
t_images.width attached_image_width, t_files.file_id attached_file_id,
t_messages.sender_id message_sender_id, t_messages.chat_room_action_id message_chat_room_action_id, t_messages.body message_body, t_messages.posted_at message_posted_at, t_messages.last_edited_at message_last_edited_at
FROM t_attached_messages
LEFT JOIN t_messages ON t_attached_messages.message_id = t_messages.message_id
LEFT JOIN t_attachable_items ON t_attached_messages.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE t_attached_messages.message_id IN (
	SELECT m.message_id FROM t_messages m WHERE m.chat_room_action_id IN (
		SELECT chat_room_action_id FROM t_chat_room_actions WHERE chat_room_id = $1
	)
)
AND
	CASE WHEN @where_in_mime_type::boolean = true THEN t_attachable_items.mime_type_id = ANY(@in_mime_types::uuid[]) ELSE TRUE END
AND
	CASE WHEN @where_is_image::boolean = true THEN EXISTS (SELECT 1 FROM t_images WHERE t_images.attachable_item_id = t_attachable_items.attachable_item_id) ELSE TRUE END
AND
	CASE WHEN @where_is_file::boolean = true THEN EXISTS (SELECT 1 FROM t_files WHERE t_files.attachable_item_id = t_attachable_items.attachable_item_id) ELSE TRUE END
ORDER BY
	t_attached_messages_pkey ASC
LIMIT $2 OFFSET $3;

-- name: GetAttachedItemsOnChatRoomUseKeysetPaginate :many
SELECT t_attached_messages.*, t_attachable_items.url attached_item_url, t_attachable_items.alias attached_item_alias,
t_attachable_items.size attached_item_size, t_attachable_items.mime_type_id attached_item_mime_type_id,
t_attachable_items.owner_id attached_item_owner_id, t_attachable_items.from_outer attached_item_from_outer,
t_images.image_id attached_image_id, t_images.height attached_image_height,
t_images.width attached_image_width, t_files.file_id attached_file_id,
t_messages.sender_id message_sender_id, t_messages.chat_room_action_id message_chat_room_action_id, t_messages.body message_body, t_messages.posted_at message_posted_at, t_messages.last_edited_at message_last_edited_at
FROM t_attached_messages
LEFT JOIN t_messages ON t_attached_messages.message_id = t_messages.message_id
LEFT JOIN t_attachable_items ON t_attached_messages.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE t_attached_messages.message_id IN (
	SELECT m.message_id FROM t_messages m WHERE m.chat_room_action_id IN (
		SELECT chat_room_action_id FROM t_chat_room_actions WHERE chat_room_id = $1
	)
)
AND
	CASE WHEN @where_in_mime_type::boolean = true THEN t_attachable_items.mime_type_id = ANY(@in_mime_types::uuid[]) ELSE TRUE END
AND
	CASE WHEN @where_is_image::boolean = true THEN EXISTS (SELECT 1 FROM t_images WHERE t_images.attachable_item_id = t_attachable_items.attachable_item_id) ELSE TRUE END
AND
	CASE WHEN @where_is_file::boolean = true THEN EXISTS (SELECT 1 FROM t_files WHERE t_files.attachable_item_id = t_attachable_items.attachable_item_id) ELSE TRUE END
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

-- name: CountAttachedItemsOnChatRoom :one
SELECT COUNT(*) FROM t_attached_messages
WHERE t_attached_messages.message_id IN (
	SELECT m.message_id FROM t_messages m WHERE m.chat_room_action_id IN (
		SELECT chat_room_action_id FROM t_chat_room_actions WHERE chat_room_id = $1
	)
)
AND
	CASE WHEN @where_in_mime_type::boolean = true THEN t_attachable_items.mime_type_id = ANY(@in_mime_types::uuid[]) ELSE TRUE END
AND
	CASE WHEN @where_is_image::boolean = true THEN EXISTS (SELECT 1 FROM t_images WHERE t_images.attachable_item_id = t_attachable_items.attachable_item_id) ELSE TRUE END
AND
	CASE WHEN @where_is_file::boolean = true THEN EXISTS (SELECT 1 FROM t_files WHERE t_files.attachable_item_id = t_attachable_items.attachable_item_id) ELSE TRUE END;

