-- name: CreateChatRoomBelongings :copyfrom
INSERT INTO m_chat_room_belongings (member_id, chat_room_id, added_at) VALUES ($1, $2, $3);

-- name: CreateChatRoomBelonging :one
INSERT INTO m_chat_room_belongings (member_id, chat_room_id, added_at) VALUES ($1, $2, $3) RETURNING *;

-- name: DeleteChatRoomBelonging :exec
DELETE FROM m_chat_room_belongings WHERE member_id = $1 AND chat_room_id = $2;

-- name: GetMembersOnChatRoom :many
SELECT sqlc.embed(m_chat_room_belongings), sqlc.embed(m_members) FROM m_chat_room_belongings
LEFT JOIN m_members ON m_chat_room_belongings.member_id = m_members.member_id
WHERE chat_room_id = $1
AND CASE
	WHEN @where_like_name::boolean = true THEN m_members.name LIKE '%' || @search_name::text || '%'
END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_members.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_members.name END DESC,
	CASE WHEN @order_method::text = 'old_add' THEN m_chat_room_belongings.added_at END ASC,
	CASE WHEN @order_method::text = 'late_add' THEN m_chat_room_belongings.added_at END DESC,
	m_chat_room_belongings_pkey DESC;

-- name: GetMembersOnChatRoomUseNumberedPaginate :many
SELECT sqlc.embed(m_chat_room_belongings), sqlc.embed(m_members) FROM m_chat_room_belongings
LEFT JOIN m_members ON m_chat_room_belongings.member_id = m_members.member_id
WHERE chat_room_id = $1
AND CASE
	WHEN @where_like_name::boolean = true THEN m_members.name LIKE '%' || @search_name::text || '%'
END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_members.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_members.name END DESC,
	CASE WHEN @order_method::text = 'old_add' THEN m_chat_room_belongings.added_at END ASC,
	CASE WHEN @order_method::text = 'late_add' THEN m_chat_room_belongings.added_at END DESC,
	m_chat_room_belongings_pkey DESC
LIMIT $2 OFFSET $3;

-- name: GetMembersOnChatRoomUseKeysetPaginate :many
SELECT sqlc.embed(m_chat_room_belongings), sqlc.embed(m_members) FROM m_chat_room_belongings
LEFT JOIN m_members ON m_chat_room_belongings.member_id = m_members.member_id
WHERE chat_room_id = $1
AND CASE @cursor_direction
	WHEN 'next' THEN
		CASE @order_method::text
			WHEN 'name' THEN m_members.name > @cursor_column OR (m_members.name = @cursor_column AND m_chat_room_belongings_pkey < @cursor)
			WHEN 'r_name' THEN m_members.name < @cursor_column OR (m_members.name = @cursor_column AND m_chat_room_belongings_pkey < @cursor)
			WHEN 'old_add' THEN m_chat_room_belongings.added_at > @cursor_column OR (m_chat_room_belongings.added_at = @cursor_column AND m_chat_room_belongings_pkey < @cursor)
			WHEN 'late_add' THEN m_chat_room_belongings.added_at < @cursor_column OR (m_chat_room_belongings.added_at = @cursor_column AND m_chat_room_belongings_pkey < @cursor)
		END
	WHEN 'prev' THEN
		CASE @order_method::text
			WHEN 'name' THEN m_members.name < @cursor_column OR (m_members.name = @cursor_column AND m_chat_room_belongings_pkey > @cursor)
			WHEN 'r_name' THEN m_members.name > @cursor_column OR (m_members.name = @cursor_column AND m_chat_room_belongings_pkey > @cursor)
			WHEN 'old_add' THEN m_chat_room_belongings.added_at < @cursor_column OR (m_chat_room_belongings.added_at = @cursor_column AND m_chat_room_belongings_pkey > @cursor)
			WHEN 'late_add' THEN m_chat_room_belongings.added_at > @cursor_column OR (m_chat_room_belongings.added_at = @cursor_column AND m_chat_room_belongings_pkey > @cursor)
		END
END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_members.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_members.name END DESC,
	CASE WHEN @order_method::text = 'old_add' THEN m_chat_room_belongings.added_at END ASC,
	CASE WHEN @order_method::text = 'late_add' THEN m_chat_room_belongings.added_at END DESC,
	m_chat_room_belongings_pkey DESC
LIMIT $2;

-- name: CountMembersOnChatRoom :one
SELECT COUNT(*) FROM m_chat_room_belongings WHERE chat_room_id = $1
AND CASE WHEN @where_like_name::boolean = true THEN
		EXISTS (SELECT 1 FROM m_members WHERE m_chat_room_belongings.member_id = m_members.member_id AND m_members.name LIKE '%' || @search_name::text || '%')
	ELSE TRUE END;

-- name: GetChatRoomsOnMember :many
SELECT sqlc.embed(m_chat_room_belongings), sqlc.embed(m_chat_rooms) FROM m_chat_room_belongings
LEFT JOIN m_chat_rooms ON m_chat_room_belongings.chat_room_id = m_chat_rooms.chat_room_id
WHERE member_id = $1
AND CASE
	WHEN @where_like_name::boolean = true THEN m_members.name LIKE '%' || @search_name::text || '%'
END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_chat_rooms.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_chat_rooms.name END DESC,
	CASE WHEN @order_method::text = 'old_add' THEN m_chat_room_belongings.added_at END ASC,
	CASE WHEN @order_method::text = 'late_add' THEN m_chat_room_belongings.added_at END DESC,
	CASE WHEN @order_method::text = 'old_chat' THEN
		(SELECT MAX(created_at) FROM t_messages m WHERE m.chat_room_id = m_chat_room_belongings.chat_room_id)
	END ASC,
	CASE WHEN @order_method::text = 'late_chat' THEN
		(SELECT MAX(created_at) FROM t_messages m WHERE m.chat_room_id = m_chat_room_belongings.chat_room_id)
	END DESC,
	m_chat_room_belongings_pkey DESC;

-- name: GetChatRoomsOnMemberUseNumberedPaginate :many
SELECT sqlc.embed(m_chat_room_belongings), sqlc.embed(m_chat_rooms) FROM m_chat_room_belongings
LEFT JOIN m_chat_rooms ON m_chat_room_belongings.chat_room_id = m_chat_rooms.chat_room_id
WHERE member_id = $1
AND CASE
	WHEN @where_like_name::boolean = true THEN m_members.name LIKE '%' || @search_name::text || '%'
END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_chat_rooms.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_chat_rooms.name END DESC,
	CASE WHEN @order_method::text = 'old_add' THEN m_chat_room_belongings.added_at END ASC,
	CASE WHEN @order_method::text = 'late_add' THEN m_chat_room_belongings.added_at END DESC,
	CASE WHEN @order_method::text = 'old_chat' THEN
		(SELECT MAX(created_at) FROM t_messages m WHERE m.chat_room_id = m_chat_room_belongings.chat_room_id)
	END ASC,
	CASE WHEN @order_method::text = 'late_chat' THEN
		(SELECT MAX(created_at) FROM t_messages m WHERE m.chat_room_id = m_chat_room_belongings.chat_room_id)
	END DESC,
	m_chat_room_belongings_pkey DESC
LIMIT $2 OFFSET $3;

-- name: GetChatRoomsOnMemberUseKeysetPaginate :many
SELECT sqlc.embed(m_chat_room_belongings), sqlc.embed(m_chat_rooms) FROM m_chat_room_belongings
LEFT JOIN m_chat_rooms ON m_chat_room_belongings.chat_room_id = m_chat_rooms.chat_room_id
WHERE member_id = $1
AND CASE @cursor_direction
	WHEN 'next' THEN
		CASE @order_method::text
			WHEN 'name' THEN m_chat_rooms.name > @cursor_column OR (m_chat_rooms.name = @cursor_column AND m_chat_room_belongings_pkey < @cursor)
			WHEN 'r_name' THEN m_chat_rooms.name < @cursor_column OR (m_chat_rooms.name = @cursor_column AND m_chat_room_belongings_pkey < @cursor)
			WHEN 'old_add' THEN m_chat_room_belongings.added_at > @cursor_column OR (m_chat_room_belongings.added_at = @cursor_column AND m_chat_room_belongings_pkey < @cursor)
			WHEN 'late_add' THEN m_chat_room_belongings.added_at < @cursor_column OR (m_chat_room_belongings.added_at = @cursor_column AND m_chat_room_belongings_pkey < @cursor)
			WHEN 'old_chat' THEN
				(SELECT MAX(created_at) FROM t_messages m WHERE m.chat_room_id = m_chat_room_belongings.chat_room_id) > @cursor_column
				OR ((SELECT MAX(created_at) FROM t_messages m WHERE m.chat_room_id = m_chat_room_belongings.chat_room_id) = @cursor_column AND m_chat_room_belongings_pkey < @cursor)
			WHEN 'late_chat' THEN
				(SELECT MAX(created_at) FROM t_messages m WHERE m.chat_room_id = m_chat_room_belongings.chat_room_id) < @cursor_column
				OR ((SELECT MAX(created_at) FROM t_messages m WHERE m.chat_room_id = m_chat_room_belongings.chat_room_id) = @cursor_column AND m_chat_room_belongings_pkey < @cursor)
		END
	WHEN 'prev' THEN
		CASE @order_method::text
			WHEN 'name' THEN m_chat_rooms.name < @cursor_column OR (m_chat_rooms.name = @cursor_column AND m_chat_room_belongings_pkey > @cursor)
			WHEN 'r_name' THEN m_chat_rooms.name > @cursor_column OR (m_chat_rooms.name = @cursor_column AND m_chat_room_belongings_pkey > @cursor)
			WHEN 'old_add' THEN m_chat_room_belongings.added_at < @cursor_column OR (m_chat_room_belongings.added = @cursor_column AND m_chat_room_belongings_pkey > @cursor)
			WHEN 'late_add' THEN m_chat_room_belongings.added_at > @cursor_column OR (m_chat_room_belongings.added_at = @cursor_column AND m_chat_room_belongings_pkey > @cursor)
			WHEN 'old_chat' THEN
				(SELECT MAX(created_at) FROM t_messages m WHERE m.chat_room_id = m_chat_room_belongings.chat_room_id) < @cursor_column
				OR ((SELECT MAX(created_at) FROM t_messages m WHERE m.chat_room_id = m_chat_room_belongings.chat_room_id) = @cursor_column AND m_chat_room_belongings_pkey > @cursor)
			WHEN 'late_chat' THEN
				(SELECT MAX(created_at) FROM t_messages m WHERE m.chat_room_id = m_chat_room_belongings.chat_room_id) > @cursor_column
				OR ((SELECT MAX(created_at) FROM t_messages m WHERE m.chat_room_id = m_chat_room_belongings.chat_room_id) = @cursor_column AND m_chat_room_belongings_pkey > @cursor)
		END
END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_chat_rooms.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_chat_rooms.name END DESC,
	CASE WHEN @order_method::text = 'old_add' THEN m_chat_room_belongings.added_at END ASC,
	CASE WHEN @order_method::text = 'late_add' THEN m_chat_room_belongings.added_at END DESC,
	CASE WHEN @order_method::text = 'old_chat' THEN
		(SELECT MAX(created_at) FROM t_messages m WHERE m.chat_room_id = m_chat_room_belongings.chat_room_id) END ASC,
	CASE WHEN @order_method::text = 'late_chat' THEN
		(SELECT MAX(created_at) FROM t_messages m WHERE m.chat_room_id = m_chat_room_belongings.chat_room_id) END DESC,
	m_chat_room_belongings_pkey DESC
LIMIT $2;

-- name: CountChatRoomsOnMember :one
SELECT COUNT(*) FROM m_chat_room_belongings WHERE member_id = $1
AND CASE WHEN @where_like_name::boolean = true THEN
		EXISTS (SELECT 1 FROM m_chat_rooms WHERE m_chat_room_belongings.chat_room_id = m_chat_rooms.chat_room_id AND m_chat_rooms.name LIKE '%' || @search_name::text || '%')
	ELSE TRUE END;
