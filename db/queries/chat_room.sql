-- name: CreateChatRooms :copyfrom
INSERT INTO m_chat_rooms (name, is_private, cover_image_id, owner_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6);

-- name: CreateChatRoom :one
INSERT INTO m_chat_rooms (name, is_private, cover_image_id, owner_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: UpdateChatRoom :one
UPDATE m_chat_rooms SET name = $2, is_private = $3, cover_image_id = $4, owner_id = $5, updated_at = $6 WHERE chat_room_id = $1 RETURNING *;

-- name: DeleteChatRoom :exec
DELETE FROM m_chat_rooms WHERE chat_room_id = $1;

-- name: FindChatRoomByID :one
SELECT * FROM m_chat_rooms WHERE chat_room_id = $1;

-- name: FindChatRoomByIDWithOwner :one
SELECT sqlc.embed(m_chat_rooms), sqlc.embed(m_members) FROM m_chat_rooms
INNER JOIN m_members ON m_chat_rooms.owner_id = m_members.member_id
WHERE chat_room_id = $1;

-- name: GetChatRooms :many
SELECT * FROM m_chat_rooms
WHERE
	CASE WHEN @where_in_owner::boolean = true THEN owner_id = ANY(@in_owner) ELSE TRUE END
AND
	CASE WHEN @where_is_private::boolean = true THEN is_private = @is_private ELSE TRUE END
ORDER BY
	m_chat_rooms_pkey DESC
LIMIT $1 OFFSET $2;

-- name: GetChatRoomsWithOwner :many
SELECT sqlc.embed(m_chat_rooms), sqlc.embed(m_members) FROM m_chat_rooms
INNER JOIN m_members ON m_chat_rooms.owner_id = m_members.member_id
WHERE
	CASE WHEN @where_in_owner::boolean THEN owner_id = ANY(@in_owner) ELSE TRUE END
AND
	CASE WHEN @where_is_private::boolean THEN is_private = @is_private ELSE TRUE END
ORDER BY
	m_chat_rooms_pkey DESC
LIMIT $1 OFFSET $2;

-- name: CountChatRooms :one
SELECT COUNT(*) FROM m_chat_rooms
WHERE
	CASE WHEN @where_in_owner::boolean = true THEN owner_id = ANY(@in_owner) ELSE TRUE END
AND
	CASE WHEN @where_is_private::boolean = true THEN is_private = @is_private ELSE TRUE END;
