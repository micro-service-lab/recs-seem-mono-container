-- name: CreateMessages :copyfrom
INSERT INTO t_messages (chat_room_id, sender_id, body, posted_at, last_edited_at) VALUES ($1, $2, $3, $4, $5);

-- name: CreateMessage :one
INSERT INTO t_messages (chat_room_id, sender_id, body, posted_at, last_edited_at) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: UpdateMessage :one
UPDATE t_messages SET chat_room_id = $2, sender_id = $3, body = $4, last_edited_at = $5 WHERE message_id = $1 RETURNING *;

-- name: DeleteMessage :exec
DELETE FROM t_messages WHERE message_id = $1;

-- name: DeleteMessagesOnChatRoom :exec
DELETE FROM t_messages WHERE chat_room_id = $1;

-- name: PluralDeleteMessages :exec
DELETE FROM t_messages WHERE message_id = ANY($1::uuid[]);

-- name: FindMessageByID :one
SELECT * FROM t_messages WHERE message_id = $1;

-- name: FindMessageByIDWithChatRoom :one
SELECT sqlc.embed(t_messages), sqlc.embed(m_chat_rooms) FROM t_messages
LEFT JOIN m_chat_rooms ON t_messages.chat_room_id = m_chat_rooms.chat_room_id
WHERE message_id = $1;

-- name: FindMessageByIDWithSender :one
SELECT sqlc.embed(t_messages), sqlc.embed(m_members) FROM t_messages
LEFT JOIN m_members ON t_messages.sender_id = m_members.member_id
WHERE message_id = $1;

-- name: FindMessageByIDWithAll :one
SELECT sqlc.embed(t_messages), sqlc.embed(m_chat_rooms), sqlc.embed(m_members) FROM t_messages
LEFT JOIN m_chat_rooms ON t_messages.chat_room_id = m_chat_rooms.chat_room_id
LEFT JOIN m_members ON t_messages.sender_id = m_members.member_id
WHERE message_id = $1;

-- name: GetMessages :many
SELECT * FROM t_messages
WHERE
	CASE WHEN @where_in_chat_room::boolean = true THEN chat_room_id = ANY(@in_chat_room) ELSE TRUE END
AND
	CASE WHEN @where_in_sender::boolean = true THEN sender_id = ANY(@in_sender) ELSE TRUE END
AND
	CASE WHEN @where_like_body::boolean = true THEN body LIKE '%' || @search_body::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_earlier_posted_at::boolean = true THEN posted_at >= @earlier_posted_at ELSE TRUE END
AND
	CASE WHEN @where_later_posted_at::boolean = true THEN posted_at <= @later_posted_at ELSE TRUE END
AND
	CASE WHEN @where_earlier_last_edited_at::boolean = true THEN last_edited_at >= @earlier_last_edited_at ELSE TRUE END
AND
	CASE WHEN @where_later_last_edited_at::boolean = true THEN last_edited_at <= @later_last_edited_at ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'posted_at' THEN posted_at END ASC,
	CASE WHEN @order_method::text = 'r_posted_at' THEN posted_at END DESC,
	CASE WHEN @order_method::text = 'last_edited_at' THEN last_edited_at END ASC,
	CASE WHEN @order_method::text = 'r_last_edited_at' THEN last_edited_at END DESC,
	t_messages_pkey ASC;

-- name: GetMessagesUseNumberedPaginate :many
SELECT * FROM t_messages
WHERE
	CASE WHEN @where_in_chat_room::boolean = true THEN chat_room_id = ANY(@in_chat_room) ELSE TRUE END
AND
	CASE WHEN @where_in_sender::boolean = true THEN sender_id = ANY(@in_sender) ELSE TRUE END
AND
	CASE WHEN @where_like_body::boolean = true THEN body LIKE '%' || @search_body::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_earlier_posted_at::boolean = true THEN posted_at >= @earlier_posted_at ELSE TRUE END
AND
	CASE WHEN @where_later_posted_at::boolean = true THEN posted_at <= @later_posted_at ELSE TRUE END
AND
	CASE WHEN @where_earlier_last_edited_at::boolean = true THEN last_edited_at >= @earlier_last_edited_at ELSE TRUE END
AND
	CASE WHEN @where_later_last_edited_at::boolean = true THEN last_edited_at <= @later_last_edited_at ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'posted_at' THEN posted_at END ASC,
	CASE WHEN @order_method::text = 'r_posted_at' THEN posted_at END DESC,
	CASE WHEN @order_method::text = 'last_edited_at' THEN last_edited_at END ASC,
	CASE WHEN @order_method::text = 'r_last_edited_at' THEN last_edited_at END DESC,
	t_messages_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetMessagesUseKeysetPaginate :many
SELECT * FROM t_messages
WHERE
	CASE WHEN @where_in_chat_room::boolean = true THEN chat_room_id = ANY(@in_chat_room) ELSE TRUE END
AND
	CASE WHEN @where_in_sender::boolean = true THEN sender_id = ANY(@in_sender) ELSE TRUE END
AND
	CASE WHEN @where_like_body::boolean = true THEN body LIKE '%' || @search_body::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_earlier_posted_at::boolean = true THEN posted_at >= @earlier_posted_at ELSE TRUE END
AND
	CASE WHEN @where_later_posted_at::boolean = true THEN posted_at <= @later_posted_at ELSE TRUE END
AND
	CASE WHEN @where_earlier_last_edited_at::boolean = true THEN last_edited_at >= @earlier_last_edited_at ELSE TRUE END
AND
	CASE WHEN @where_later_last_edited_at::boolean = true THEN last_edited_at <= @later_last_edited_at ELSE TRUE END
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'posted_at' THEN posted_at > @posted_at_cursor OR (posted_at = @posted_at_cursor AND t_messages_pkey > @cursor::int)
				WHEN 'r_posted_at' THEN posted_at < @posted_at_cursor OR (posted_at = @posted_at_cursor AND t_messages_pkey > @cursor::int)
				WHEN 'last_edited_at' THEN last_edited_at > @last_edited_at_cursor OR (last_edited_at = @last_edited_at_cursor AND t_messages_pkey > @cursor::int)
				WHEN 'r_last_edited_at' THEN last_edited_at < @last_edited_at_cursor OR (last_edited_at = @last_edited_at_cursor AND t_messages_pkey > @cursor::int)
				ELSE t_messages_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'posted_at' THEN posted_at < @posted_at_cursor OR (posted_at = @posted_at_cursor AND t_messages_pkey < @cursor::int)
				WHEN 'r_posted_at' THEN posted_at > @posted_at_cursor OR (posted_at = @posted_at_cursor AND t_messages_pkey < @cursor::int)
				WHEN 'last_edited_at' THEN last_edited_at < @last_edited_at_cursor OR (last_edited_at = @last_edited_at_cursor AND t_messages_pkey < @cursor::int)
				WHEN 'r_last_edited_at' THEN last_edited_at > @last_edited_at_cursor OR (last_edited_at = @last_edited_at_cursor AND t_messages_pkey < @cursor::int)
				ELSE t_messages_pkey < @cursor::int
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'posted_at' AND @cursor_direction::text = 'next' THEN posted_at END ASC,
	CASE WHEN @order_method::text = 'posted_at' AND @cursor_direction::text = 'prev' THEN posted_at END DESC,
	CASE WHEN @order_method::text = 'r_posted_at' AND @cursor_direction::text = 'next' THEN posted_at END DESC,
	CASE WHEN @order_method::text = 'r_posted_at' AND @cursor_direction::text = 'prev' THEN posted_at END ASC,
	CASE WHEN @order_method::text = 'last_edited_at' AND @cursor_direction::text = 'next' THEN last_edited_at END ASC,
	CASE WHEN @order_method::text = 'last_edited_at' AND @cursor_direction::text = 'prev' THEN last_edited_at END DESC,
	CASE WHEN @order_method::text = 'r_last_edited_at' AND @cursor_direction::text = 'next' THEN last_edited_at END DESC,
	CASE WHEN @order_method::text = 'r_last_edited_at' AND @cursor_direction::text = 'prev' THEN last_edited_at END ASC,
	CASE WHEN @cursor_direction::text = 'next' THEN t_messages_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN t_messages_pkey END DESC
LIMIT $1;

-- name: GetPluralMessages :many
SELECT * FROM t_messages WHERE message_id = ANY(@message_ids::uuid[])
ORDER BY
	t_messages_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetMessagesWithChatRoom :many
SELECT sqlc.embed(t_messages), sqlc.embed(m_chat_rooms) FROM t_messages
LEFT JOIN m_chat_rooms ON t_messages.chat_room_id = m_chat_rooms.chat_room_id
WHERE
	CASE WHEN @where_in_chat_room::boolean = true THEN t_messages.chat_room_id = ANY(@in_chat_room) ELSE TRUE END
AND
	CASE WHEN @where_in_sender::boolean = true THEN sender_id = ANY(@in_sender) ELSE TRUE END
AND
	CASE WHEN @where_like_body::boolean = true THEN body LIKE '%' || @search_body::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_earlier_posted_at::boolean = true THEN posted_at >= @earlier_posted_at ELSE TRUE END
AND
	CASE WHEN @where_later_posted_at::boolean = true THEN posted_at <= @later_posted_at ELSE TRUE END
AND
	CASE WHEN @where_earlier_last_edited_at::boolean = true THEN last_edited_at >= @earlier_last_edited_at ELSE TRUE END
AND
	CASE WHEN @where_later_last_edited_at::boolean = true THEN last_edited_at <= @later_last_edited_at ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'posted_at' THEN posted_at END ASC,
	CASE WHEN @order_method::text = 'r_posted_at' THEN posted_at END DESC,
	CASE WHEN @order_method::text = 'last_edited_at' THEN last_edited_at END ASC,
	CASE WHEN @order_method::text = 'r_last_edited_at' THEN last_edited_at END DESC,
	t_messages_pkey ASC;

-- name: GetMessagesWithChatRoomUseNumberedPaginate :many
SELECT sqlc.embed(t_messages), sqlc.embed(m_chat_rooms) FROM t_messages
LEFT JOIN m_chat_rooms ON t_messages.chat_room_id = m_chat_rooms.chat_room_id
WHERE
	CASE WHEN @where_in_chat_room::boolean = true THEN t_messages.chat_room_id = ANY(@in_chat_room) ELSE TRUE END
AND
	CASE WHEN @where_in_sender::boolean = true THEN sender_id = ANY(@in_sender) ELSE TRUE END
AND
	CASE WHEN @where_like_body::boolean = true THEN body LIKE '%' || @search_body::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_earlier_posted_at::boolean = true THEN posted_at >= @earlier_posted_at ELSE TRUE END
AND
	CASE WHEN @where_later_posted_at::boolean = true THEN posted_at <= @later_posted_at ELSE TRUE END
AND
	CASE WHEN @where_earlier_last_edited_at::boolean = true THEN last_edited_at >= @earlier_last_edited_at ELSE TRUE END
AND
	CASE WHEN @where_later_last_edited_at::boolean = true THEN last_edited_at <= @later_last_edited_at ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'posted_at' THEN posted_at END ASC,
	CASE WHEN @order_method::text = 'r_posted_at' THEN posted_at END DESC,
	CASE WHEN @order_method::text = 'last_edited_at' THEN last_edited_at END ASC,
	CASE WHEN @order_method::text = 'r_last_edited_at' THEN last_edited_at END DESC,
	t_messages_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetMessagesWithChatRoomUseKeysetPaginate :many
SELECT sqlc.embed(t_messages), sqlc.embed(m_chat_rooms) FROM t_messages
LEFT JOIN m_chat_rooms ON t_messages.chat_room_id = m_chat_rooms.chat_room_id
WHERE
	CASE WHEN @where_in_chat_room::boolean = true THEN t_messages.chat_room_id = ANY(@in_chat_room) ELSE TRUE END
AND
	CASE WHEN @where_in_sender::boolean = true THEN sender_id = ANY(@in_sender) ELSE TRUE END
AND
	CASE WHEN @where_like_body::boolean = true THEN body LIKE '%' || @search_body::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_earlier_posted_at::boolean = true THEN posted_at >= @earlier_posted_at ELSE TRUE END
AND
	CASE WHEN @where_later_posted_at::boolean = true THEN posted_at <= @later_posted_at ELSE TRUE END
AND
	CASE WHEN @where_earlier_last_edited_at::boolean = true THEN last_edited_at >= @earlier_last_edited_at ELSE TRUE END
AND
	CASE WHEN @where_later_last_edited_at::boolean = true THEN last_edited_at <= @later_last_edited_at ELSE TRUE END
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'posted_at' THEN posted_at > @posted_at_cursor OR (posted_at = @posted_at_cursor AND t_messages_pkey > @cursor::int)
				WHEN 'r_posted_at' THEN posted_at < @posted_at_cursor OR (posted_at = @posted_at_cursor AND t_messages_pkey > @cursor::int)
				WHEN 'last_edited_at' THEN last_edited_at > @last_edited_at_cursor OR (last_edited_at = @last_edited_at_cursor AND t_messages_pkey > @cursor::int)
				WHEN 'r_last_edited_at' THEN last_edited_at < @last_edited_at_cursor OR (last_edited_at = @last_edited_at_cursor AND t_messages_pkey > @cursor::int)
				ELSE t_messages_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'posted_at' THEN posted_at < @posted_at_cursor OR (posted_at = @posted_at_cursor AND t_messages_pkey < @cursor::int)
				WHEN 'r_posted_at' THEN posted_at > @posted_at_cursor OR (posted_at = @posted_at_cursor AND t_messages_pkey < @cursor::int)
				WHEN 'last_edited_at' THEN last_edited_at < @last_edited_at_cursor OR (last_edited_at = @last_edited_at_cursor AND t_messages_pkey < @cursor::int)
				WHEN 'r_last_edited_at' THEN last_edited_at > @last_edited_at_cursor OR (last_edited_at = @last_edited_at_cursor AND t_messages_pkey < @cursor::int)
				ELSE t_messages_pkey < @cursor::int
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'posted_at' AND @cursor_direction::text = 'next' THEN posted_at END ASC,
	CASE WHEN @order_method::text = 'posted_at' AND @cursor_direction::text = 'prev' THEN posted_at END DESC,
	CASE WHEN @order_method::text = 'r_posted_at' AND @cursor_direction::text = 'next' THEN posted_at END DESC,
	CASE WHEN @order_method::text = 'r_posted_at' AND @cursor_direction::text = 'prev' THEN posted_at END ASC,
	CASE WHEN @order_method::text = 'last_edited_at' AND @cursor_direction::text = 'next' THEN last_edited_at END ASC,
	CASE WHEN @order_method::text = 'last_edited_at' AND @cursor_direction::text = 'prev' THEN last_edited_at END DESC,
	CASE WHEN @order_method::text = 'r_last_edited_at' AND @cursor_direction::text = 'next' THEN last_edited_at END DESC,
	CASE WHEN @order_method::text = 'r_last_edited_at' AND @cursor_direction::text = 'prev' THEN last_edited_at END ASC,
	CASE WHEN @cursor_direction::text = 'next' THEN t_messages_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN t_messages_pkey END DESC
LIMIT $1;

-- name: GetPluralMessagesWithChatRoom :many
SELECT sqlc.embed(t_messages), sqlc.embed(m_chat_rooms) FROM t_messages
LEFT JOIN m_chat_rooms ON t_messages.chat_room_id = m_chat_rooms.chat_room_id
WHERE message_id = ANY(@message_ids::uuid[])
ORDER BY
	t_messages_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetMessagesWithSender :many
SELECT sqlc.embed(t_messages), sqlc.embed(m_members) FROM t_messages
LEFT JOIN m_members ON t_messages.sender_id = m_members.member_id
WHERE
	CASE WHEN @where_in_chat_room::boolean = true THEN t_messages.chat_room_id = ANY(@in_chat_room) ELSE TRUE END
AND
	CASE WHEN @where_in_sender::boolean = true THEN sender_id = ANY(@in_sender) ELSE TRUE END
AND
	CASE WHEN @where_like_body::boolean = true THEN body LIKE '%' || @search_body::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_earlier_posted_at::boolean = true THEN posted_at >= @earlier_posted_at ELSE TRUE END
AND
	CASE WHEN @where_later_posted_at::boolean = true THEN posted_at <= @later_posted_at ELSE TRUE END
AND
	CASE WHEN @where_earlier_last_edited_at::boolean = true THEN last_edited_at >= @earlier_last_edited_at ELSE TRUE END
AND
	CASE WHEN @where_later_last_edited_at::boolean = true THEN last_edited_at <= @later_last_edited_at ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'posted_at' THEN posted_at END ASC,
	CASE WHEN @order_method::text = 'r_posted_at' THEN posted_at END DESC,
	CASE WHEN @order_method::text = 'last_edited_at' THEN last_edited_at END ASC,
	CASE WHEN @order_method::text = 'r_last_edited_at' THEN last_edited_at END DESC,
	t_messages_pkey ASC;

-- name: GetMessagesWithSenderUseNumberedPaginate :many
SELECT sqlc.embed(t_messages), sqlc.embed(m_members) FROM t_messages
LEFT JOIN m_members ON t_messages.sender_id = m_members.member_id
WHERE
	CASE WHEN @where_in_chat_room::boolean = true THEN t_messages.chat_room_id = ANY(@in_chat_room) ELSE TRUE END
AND
	CASE WHEN @where_in_sender::boolean = true THEN sender_id = ANY(@in_sender) ELSE TRUE END
AND
	CASE WHEN @where_like_body::boolean = true THEN body LIKE '%' || @search_body::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_earlier_posted_at::boolean = true THEN posted_at >= @earlier_posted_at ELSE TRUE END
AND
	CASE WHEN @where_later_posted_at::boolean = true THEN posted_at <= @later_posted_at ELSE TRUE END
AND
	CASE WHEN @where_earlier_last_edited_at::boolean = true THEN last_edited_at >= @earlier_last_edited_at ELSE TRUE END
AND
	CASE WHEN @where_later_last_edited_at::boolean = true THEN last_edited_at <= @later_last_edited_at ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'posted_at' THEN posted_at END ASC,
	CASE WHEN @order_method::text = 'r_posted_at' THEN posted_at END DESC,
	CASE WHEN @order_method::text = 'last_edited_at' THEN last_edited_at END ASC,
	CASE WHEN @order_method::text = 'r_last_edited_at' THEN last_edited_at END DESC,
	t_messages_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetMessagesWithSenderUseKeysetPaginate :many
SELECT sqlc.embed(t_messages), sqlc.embed(m_members) FROM t_messages
LEFT JOIN m_members ON t_messages.sender_id = m_members.member_id
WHERE
	CASE WHEN @where_in_chat_room::boolean = true THEN t_messages.chat_room_id = ANY(@in_chat_room) ELSE TRUE END
AND
	CASE WHEN @where_in_sender::boolean = true THEN sender_id = ANY(@in_sender) ELSE TRUE END
AND
	CASE WHEN @where_like_body::boolean = true THEN body LIKE '%' || @search_body::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_earlier_posted_at::boolean = true THEN posted_at >= @earlier_posted_at ELSE TRUE END
AND
	CASE WHEN @where_later_posted_at::boolean = true THEN posted_at <= @later_posted_at ELSE TRUE END
AND
	CASE WHEN @where_earlier_last_edited_at::boolean = true THEN last_edited_at >= @earlier_last_edited_at ELSE TRUE END
AND
	CASE WHEN @where_later_last_edited_at::boolean = true THEN last_edited_at <= @later_last_edited_at ELSE TRUE END
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'posted_at' THEN posted_at > @posted_at_cursor OR (posted_at = @posted_at_cursor AND t_messages_pkey > @cursor::int)
				WHEN 'r_posted_at' THEN posted_at < @posted_at_cursor OR (posted_at = @posted_at_cursor AND t_messages_pkey > @cursor::int)
				WHEN 'last_edited_at' THEN last_edited_at > @last_edited_at_cursor OR (last_edited_at = @last_edited_at_cursor AND t_messages_pkey > @cursor::int)
				WHEN 'r_last_edited_at' THEN last_edited_at < @last_edited_at_cursor OR (last_edited_at = @last_edited_at_cursor AND t_messages_pkey > @cursor::int)
				ELSE t_messages_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'posted_at' THEN posted_at < @posted_at_cursor OR (posted_at = @posted_at_cursor AND t_messages_pkey < @cursor::int)
				WHEN 'r_posted_at' THEN posted_at > @posted_at_cursor OR (posted_at = @posted_at_cursor AND t_messages_pkey < @cursor::int)
				WHEN 'last_edited_at' THEN last_edited_at < @last_edited_at_cursor OR (last_edited_at = @last_edited_at_cursor AND t_messages_pkey < @cursor::int)
				WHEN 'r_last_edited_at' THEN last_edited_at > @last_edited_at_cursor OR (last_edited_at = @last_edited_at_cursor AND t_messages_pkey < @cursor::int)
				ELSE t_messages_pkey < @cursor::int
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'posted_at' AND @cursor_direction::text = 'next' THEN posted_at END ASC,
	CASE WHEN @order_method::text = 'posted_at' AND @cursor_direction::text = 'prev' THEN posted_at END DESC,
	CASE WHEN @order_method::text = 'r_posted_at' AND @cursor_direction::text = 'next' THEN posted_at END DESC,
	CASE WHEN @order_method::text = 'r_posted_at' AND @cursor_direction::text = 'prev' THEN posted_at END ASC,
	CASE WHEN @order_method::text = 'last_edited_at' AND @cursor_direction::text = 'next' THEN last_edited_at END ASC,
	CASE WHEN @order_method::text = 'last_edited_at' AND @cursor_direction::text = 'prev' THEN last_edited_at END DESC,
	CASE WHEN @order_method::text = 'r_last_edited_at' AND @cursor_direction::text = 'next' THEN last_edited_at END DESC,
	CASE WHEN @order_method::text = 'r_last_edited_at' AND @cursor_direction::text = 'prev' THEN last_edited_at END ASC,
	CASE WHEN @cursor_direction::text = 'next' THEN t_messages_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN t_messages_pkey END DESC
LIMIT $1;

-- name: GetPluralMessagesWithSender :many
SELECT sqlc.embed(t_messages), sqlc.embed(m_members) FROM t_messages
LEFT JOIN m_members ON t_messages.sender_id = m_members.member_id
WHERE message_id = ANY(@message_ids::uuid[])
ORDER BY
	t_messages_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetMessagesWithAll :many
SELECT sqlc.embed(t_messages), sqlc.embed(m_chat_rooms), sqlc.embed(m_members) FROM t_messages
LEFT JOIN m_chat_rooms ON t_messages.chat_room_id = m_chat_rooms.chat_room_id
LEFT JOIN m_members ON t_messages.sender_id = m_members.member_id
WHERE
	CASE WHEN @where_in_chat_room::boolean = true THEN t_messages.chat_room_id = ANY(@in_chat_room) ELSE TRUE END
AND
	CASE WHEN @where_in_sender::boolean = true THEN sender_id = ANY(@in_sender) ELSE TRUE END
AND
	CASE WHEN @where_like_body::boolean = true THEN body LIKE '%' || @search_body::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_earlier_posted_at::boolean = true THEN posted_at >= @earlier_posted_at ELSE TRUE END
AND
	CASE WHEN @where_later_posted_at::boolean = true THEN posted_at <= @later_posted_at ELSE TRUE END
AND
	CASE WHEN @where_earlier_last_edited_at::boolean = true THEN last_edited_at >= @earlier_last_edited_at ELSE TRUE END
AND
	CASE WHEN @where_later_last_edited_at::boolean = true THEN last_edited_at <= @later_last_edited_at ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'posted_at' THEN posted_at END ASC,
	CASE WHEN @order_method::text = 'r_posted_at' THEN posted_at END DESC,
	CASE WHEN @order_method::text = 'last_edited_at' THEN last_edited_at END ASC,
	CASE WHEN @order_method::text = 'r_last_edited_at' THEN last_edited_at END DESC,
	t_messages_pkey ASC;

-- name: GetMessagesWithAllUseNumberedPaginate :many
SELECT sqlc.embed(t_messages), sqlc.embed(m_chat_rooms), sqlc.embed(m_members) FROM t_messages
LEFT JOIN m_chat_rooms ON t_messages.chat_room_id = m_chat_rooms.chat_room_id
LEFT JOIN m_members ON t_messages.sender_id = m_members.member_id
WHERE
	CASE WHEN @where_in_chat_room::boolean = true THEN t_messages.chat_room_id = ANY(@in_chat_room) ELSE TRUE END
AND
	CASE WHEN @where_in_sender::boolean = true THEN sender_id = ANY(@in_sender) ELSE TRUE END
AND
	CASE WHEN @where_like_body::boolean = true THEN body LIKE '%' || @search_body::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_earlier_posted_at::boolean = true THEN posted_at >= @earlier_posted_at ELSE TRUE END
AND
	CASE WHEN @where_later_posted_at::boolean = true THEN posted_at <= @later_posted_at ELSE TRUE END
AND
	CASE WHEN @where_earlier_last_edited_at::boolean = true THEN last_edited_at >= @earlier_last_edited_at ELSE TRUE END
AND
	CASE WHEN @where_later_last_edited_at::boolean = true THEN last_edited_at <= @later_last_edited_at ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'posted_at' THEN posted_at END ASC,
	CASE WHEN @order_method::text = 'r_posted_at' THEN posted_at END DESC,
	CASE WHEN @order_method::text = 'last_edited_at' THEN last_edited_at END ASC,
	CASE WHEN @order_method::text = 'r_last_edited_at' THEN last_edited_at END DESC,
	t_messages_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetMessagesWithAllUseKeysetPaginate :many
SELECT sqlc.embed(t_messages), sqlc.embed(m_chat_rooms), sqlc.embed(m_members) FROM t_messages
LEFT JOIN m_chat_rooms ON t_messages.chat_room_id = m_chat_rooms.chat_room_id
LEFT JOIN m_members ON t_messages.sender_id = m_members.member_id
WHERE
	CASE WHEN @where_in_chat_room::boolean = true THEN t_messages.chat_room_id = ANY(@in_chat_room) ELSE TRUE END
AND
	CASE WHEN @where_in_sender::boolean = true THEN sender_id = ANY(@in_sender) ELSE TRUE END
AND
	CASE WHEN @where_like_body::boolean = true THEN body LIKE '%' || @search_body::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_earlier_posted_at::boolean = true THEN posted_at >= @earlier_posted_at ELSE TRUE END
AND
	CASE WHEN @where_later_posted_at::boolean = true THEN posted_at <= @later_posted_at ELSE TRUE END
AND
	CASE WHEN @where_earlier_last_edited_at::boolean = true THEN last_edited_at >= @earlier_last_edited_at ELSE TRUE END
AND
	CASE WHEN @where_later_last_edited_at::boolean = true THEN last_edited_at <= @later_last_edited_at ELSE TRUE END
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'posted_at' THEN posted_at > @posted_at_cursor OR (posted_at = @posted_at_cursor AND t_messages_pkey > @cursor::int)
				WHEN 'r_posted_at' THEN posted_at < @posted_at_cursor OR (posted_at = @posted_at_cursor AND t_messages_pkey > @cursor::int)
				WHEN 'last_edited_at' THEN last_edited_at > @last_edited_at_cursor OR (last_edited_at = @last_edited_at_cursor AND t_messages_pkey > @cursor::int)
				WHEN 'r_last_edited_at' THEN last_edited_at < @last_edited_at_cursor OR (last_edited_at = @last_edited_at_cursor AND t_messages_pkey > @cursor::int)
				ELSE t_messages_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'posted_at' THEN posted_at < @posted_at_cursor OR (posted_at = @posted_at_cursor AND t_messages_pkey < @cursor::int)
				WHEN 'r_posted_at' THEN posted_at > @posted_at_cursor OR (posted_at = @posted_at_cursor AND t_messages_pkey < @cursor::int)
				WHEN 'last_edited_at' THEN last_edited_at < @last_edited_at_cursor OR (last_edited_at = @last_edited_at_cursor AND t_messages_pkey < @cursor::int)
				WHEN 'r_last_edited_at' THEN last_edited_at > @last_edited_at_cursor OR (last_edited_at = @last_edited_at_cursor AND t_messages_pkey < @cursor::int)
				ELSE t_messages_pkey < @cursor::int
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'posted_at' AND @cursor_direction::text = 'next' THEN posted_at END ASC,
	CASE WHEN @order_method::text = 'posted_at' AND @cursor_direction::text = 'prev' THEN posted_at END DESC,
	CASE WHEN @order_method::text = 'r_posted_at' AND @cursor_direction::text = 'next' THEN posted_at END DESC,
	CASE WHEN @order_method::text = 'r_posted_at' AND @cursor_direction::text = 'prev' THEN posted_at END ASC,
	CASE WHEN @order_method::text = 'last_edited_at' AND @cursor_direction::text = 'next' THEN last_edited_at END ASC,
	CASE WHEN @order_method::text = 'last_edited_at' AND @cursor_direction::text = 'prev' THEN last_edited_at END DESC,
	CASE WHEN @order_method::text = 'r_last_edited_at' AND @cursor_direction::text = 'next' THEN last_edited_at END DESC,
	CASE WHEN @order_method::text = 'r_last_edited_at' AND @cursor_direction::text = 'prev' THEN last_edited_at END ASC,
	CASE WHEN @cursor_direction::text = 'next' THEN t_messages_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN t_messages_pkey END DESC
LIMIT $1;

-- name: GetPluralMessagesWithAll :many
SELECT sqlc.embed(t_messages), sqlc.embed(m_chat_rooms), sqlc.embed(m_members) FROM t_messages
LEFT JOIN m_chat_rooms ON t_messages.chat_room_id = m_chat_rooms.chat_room_id
LEFT JOIN m_members ON t_messages.sender_id = m_members.member_id
WHERE message_id = ANY(@message_ids::uuid[])
ORDER BY
	t_messages_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountMessages :one
SELECT COUNT(*) FROM t_messages
WHERE
	CASE WHEN @where_in_chat_room::boolean = true THEN chat_room_id = ANY(@in_chat_room) ELSE TRUE END
AND
	CASE WHEN @where_in_sender::boolean = true THEN sender_id = ANY(@in_sender) ELSE TRUE END
AND
	CASE WHEN @where_like_body::boolean = true THEN body LIKE '%' || @search_body::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_earlier_posted_at::boolean = true THEN posted_at >= @earlier_posted_at ELSE TRUE END
AND
	CASE WHEN @where_later_posted_at::boolean = true THEN posted_at <= @later_posted_at ELSE TRUE END
AND
	CASE WHEN @where_earlier_last_edited_at::boolean = true THEN last_edited_at >= @earlier_last_edited_at ELSE TRUE END
AND
	CASE WHEN @where_later_last_edited_at::boolean = true THEN last_edited_at <= @later_last_edited_at ELSE TRUE END;
