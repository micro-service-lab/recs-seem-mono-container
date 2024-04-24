-- name: CreateMessages :copyfrom
INSERT INTO t_messages (chat_room_id, sender_id, body, posted_at, last_edited_at) VALUES ($1, $2, $3, $4, $5);

-- name: CreateMessage :one
INSERT INTO t_messages (chat_room_id, sender_id, body, posted_at, last_edited_at) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: UpdateMessage :one
UPDATE t_messages SET chat_room_id = $2, sender_id = $3, body = $4, posted_at = $5, last_edited_at = $6 WHERE message_id = $1 RETURNING *;

-- name: DeleteMessage :exec
DELETE FROM t_messages WHERE message_id = $1;

-- name: FindMessageByID :one
SELECT * FROM t_messages WHERE message_id = $1;

-- name: FindMessageByIDWithSender :one
SELECT sqlc.embed(t_messages), sqlc.embed(m_members) FROM t_messages
INNER JOIN m_members ON t_messages.sender_id = m_members.member_id
WHERE message_id = $1;

-- name: GetMessages :many
SELECT * FROM t_messages
WHERE
	CASE WHEN @where_in_chat_room::boolean = true THEN chat_room_id = ANY(@in_chat_room) ELSE TRUE END
AND
	CASE WHEN @where_in_sender::boolean = true THEN sender_id = ANY(@in_sender) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'posted_at' THEN posted_at END ASC,
	CASE WHEN @order_method::text = 'r_posted_at' THEN posted_at END DESC,
	CASE WHEN @order_method::text = 'last_edited_at' THEN last_edited_at END ASC,
	CASE WHEN @order_method::text = 'r_last_edited_at' THEN last_edited_at END DESC,
	t_messages_pkey DESC
LIMIT $1 OFFSET $2;

-- name: CountMessages :one
SELECT COUNT(*) FROM t_messages
WHERE
	CASE WHEN @where_in_chat_room::boolean = true THEN chat_room_id = ANY(@in_chat_room) ELSE TRUE END
AND
	CASE WHEN @where_in_sender::boolean = true THEN sender_id = ANY(@in_sender) ELSE TRUE END;
