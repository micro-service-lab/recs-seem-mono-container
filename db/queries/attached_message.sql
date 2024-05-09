-- name: CreateAttachedMessages :copyfrom
INSERT INTO t_attached_messages (message_id, file_url) VALUES ($1, $2);

-- name: CreateAttachedMessage :one
INSERT INTO t_attached_messages (message_id, file_url) VALUES ($1, $2) RETURNING *;

-- name: DeleteAttachedMessage :exec
DELETE FROM t_attached_messages WHERE attached_message_id = $1;

-- name: DeleteAttachedMessagesOnMessage :exec
DELETE FROM t_attached_messages WHERE message_id = $1;

-- name: DeleteAttachedMessagesOnMessages :exec
DELETE FROM t_attached_messages WHERE message_id = ANY($1::uuid[]);

-- name: GetFilesOnMessage :many
SELECT t_attachable_items.* FROM t_attached_messages
WHERE message_id = $1
ORDER BY
	t_attached_messages_pkey ASC;

-- name: GetFilesOnMessageUseNumberedPaginate :many
SELECT t_attachable_items.* FROM t_attached_messages
WHERE message_id = $1
ORDER BY
	t_attached_messages_pkey ASC
LIMIT $2 OFFSET $3;

-- name: GetFilesOnMessageUseKeysetPaginate :many
SELECT t_attachable_items.* FROM t_attached_messages
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

-- name: GetPluralFilesOnMessage :many
SELECT t_attachable_items.* FROM t_attached_messages
WHERE message_id = ANY(@message_ids::uuid[])
ORDER BY
	t_attached_messages_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountFilesOnMessage :one
SELECT COUNT(*) FROM t_attached_messages WHERE message_id = $1;

-- name: GetFilesOnChatRoom :many
SELECT t_attachable_items.* FROM t_attached_messages
WHERE message_id IN (
	SELECT message_id FROM t_messages WHERE chat_room_id = $1
)
ORDER BY
	t_attached_messages_pkey ASC;

-- name: GetFilesOnChatRoomUseNumberedPaginate :many
SELECT t_attachable_items.* FROM t_attached_messages
WHERE message_id IN (
	SELECT message_id FROM t_messages WHERE chat_room_id = $1
)
ORDER BY
	t_attached_messages_pkey ASC
LIMIT $2 OFFSET $3;

-- name: GetFilesOnChatRoomUseKeysetPaginate :many
SELECT t_attachable_items.* FROM t_attached_messages
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

-- name: GetPluralFilesOnChatRoom :many
SELECT t_attachable_items.* FROM t_attached_messages
WHERE message_id IN (
	SELECT message_id FROM t_messages WHERE chat_room_id = ANY(@chat_room_ids::uuid[])
)
ORDER BY
	t_attached_messages_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountFilesOnChatRoom :one
SELECT COUNT(*) FROM t_attached_messages
WHERE message_id IN (
	SELECT message_id FROM t_messages WHERE chat_room_id = $1
);

