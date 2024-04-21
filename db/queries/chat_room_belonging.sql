-- name: CreateChatRoomBelongings :copyfrom
INSERT INTO m_chat_room_belongings (member_id, chat_room_id, added_at) VALUES ($1, $2, $3);

-- name: CreateChatRoomBelonging :one
INSERT INTO m_chat_room_belongings (member_id, chat_room_id, added_at) VALUES ($1, $2, $3) RETURNING *;

-- name: DeleteChatRoomBelonging :exec
DELETE FROM m_chat_room_belongings WHERE member_id = $1 AND chat_room_id = $2;

-- name: GetMembersOnChatRoomID :many
SELECT sqlc.embed(m_chat_room_belongings), sqlc.embed(m_members) FROM m_chat_room_belongings
INNER JOIN m_members ON m_chat_room_belongings.member_id = m_members.member_id
WHERE chat_room_id = $1
AND CASE
	WHEN @where_like_name::boolean = true THEN m_members.name LIKE '%' || @search_name::text || '%'
END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_members.name END ASC,
	CASE WHEN @order_method::text = 'add' THEN m_chat_room_belongings.added_at END DESC,
	m_chat_room_belongings_pkey DESC
LIMIT $2 OFFSET $3;

-- name: CountMembersOnChatRoomID :one
SELECT COUNT(*) FROM m_chat_room_belongings WHERE chat_room_id = $1;

-- name: GetChatRoomsByMemberID :many
SELECT sqlc.embed(m_chat_room_belongings), sqlc.embed(m_chat_rooms) FROM m_chat_room_belongings
INNER JOIN m_chat_rooms ON m_chat_room_belongings.chat_room_id = m_chat_rooms.chat_room_id
WHERE member_id = $1
AND CASE
	WHEN @where_like_name::boolean = true THEN m_members.name LIKE '%' || @search_name::text || '%'
END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_chat_rooms.name END ASC,
	CASE WHEN @order_method::text = 'add' THEN m_chat_room_belongings.added_at END DESC,
	CASE WHEN @order_method::text = 'late_chat' THEN
		(SELECT MAX(created_at) FROM t_messages m WHERE m.chat_room_id = m_chat_room_belongings.chat_room_id)
	END DESC,
	CASE WHEN @order_method::text = 'old_chat' THEN
		(SELECT MAX(created_at) FROM t_messages m WHERE m.chat_room_id = m_chat_room_belongings.chat_room_id)
	END ASC,
	m_chat_room_belongings_pkey DESC
LIMIT $2 OFFSET $3;

-- name: CountChatRoomsByMemberID :one
SELECT COUNT(*) FROM m_chat_room_belongings WHERE member_id = $1;
