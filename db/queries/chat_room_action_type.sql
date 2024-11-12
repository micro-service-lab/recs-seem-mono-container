-- name: CreateChatRoomActionTypes :copyfrom
INSERT INTO m_chat_room_action_types (name, key) VALUES ($1, $2);

-- name: CreateChatRoomActionType :one
INSERT INTO m_chat_room_action_types (name, key) VALUES ($1, $2) RETURNING *;

-- name: UpdateChatRoomActionType :one
UPDATE m_chat_room_action_types SET name = $2, key = $3 WHERE chat_room_action_type_id = $1 RETURNING *;

-- name: UpdateChatRoomActionTypeByKey :one
UPDATE m_chat_room_action_types SET name = $2 WHERE key = $1 RETURNING *;

-- name: DeleteChatRoomActionType :execrows
DELETE FROM m_chat_room_action_types WHERE chat_room_action_type_id = $1;

-- name: DeleteChatRoomActionTypeByKey :execrows
DELETE FROM m_chat_room_action_types WHERE key = $1;

-- name: PluralDeleteChatRoomActionTypes :execrows
DELETE FROM m_chat_room_action_types WHERE chat_room_action_type_id = ANY(@chat_room_action_type_ids::uuid[]);

-- name: FindChatRoomActionTypeByID :one
SELECT * FROM m_chat_room_action_types WHERE chat_room_action_type_id = $1;

-- name: FindChatRoomActionTypeByKey :one
SELECT * FROM m_chat_room_action_types WHERE key = $1;

-- name: GetChatRoomActionTypes :many
SELECT * FROM m_chat_room_action_types
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_chat_room_action_types.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN name END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' THEN name END DESC NULLS LAST,
	m_chat_room_action_types_pkey ASC;

-- name: GetChatRoomActionTypesUseNumberedPaginate :many
SELECT * FROM m_chat_room_action_types
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_chat_room_action_types.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN name END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' THEN name END DESC NULLS LAST,
	m_chat_room_action_types_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetChatRoomActionTypesUseKeysetPaginate :many
SELECT * FROM m_chat_room_action_types
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_chat_room_action_types.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'name' THEN name > @name_cursor OR (name = @name_cursor AND m_chat_room_action_types_pkey > @cursor::int)
				WHEN 'r_name' THEN name < @name_cursor OR (name = @name_cursor AND m_chat_room_action_types_pkey > @cursor::int)
				ELSE m_chat_room_action_types_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'name' THEN name < @name_cursor OR (name = @name_cursor AND m_chat_room_action_types_pkey < @cursor::int)
				WHEN 'r_name' THEN name > @name_cursor OR (name = @name_cursor AND m_chat_room_action_types_pkey < @cursor::int)
				ELSE m_chat_room_action_types_pkey < @cursor::int
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'next' THEN name END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'prev' THEN name END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'next' THEN name END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'prev' THEN name END ASC NULLS LAST,
	CASE WHEN @cursor_direction::text = 'next' THEN m_chat_room_action_types_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN m_chat_room_action_types_pkey END DESC
LIMIT $1;

-- name: GetPluralChatRoomActionTypes :many
SELECT * FROM m_chat_room_action_types
WHERE
	chat_room_action_type_id = ANY(@chat_room_action_type_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN name END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' THEN name END DESC NULLS LAST,
	m_chat_room_action_types_pkey ASC;

-- name: GetPluralChatRoomActionTypesUseNumberedPaginate :many
SELECT * FROM m_chat_room_action_types
WHERE
	chat_room_action_type_id = ANY(@chat_room_action_type_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN name END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' THEN name END DESC NULLS LAST,
	m_chat_room_action_types_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountChatRoomActionTypes :one
SELECT COUNT(*) FROM m_chat_room_action_types
WHERE
	CASE WHEN @where_like_name::boolean = true THEN name LIKE '%' || @search_name::text || '%' ELSE TRUE END;
