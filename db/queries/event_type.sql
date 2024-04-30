-- name: CreateEventTypes :copyfrom
INSERT INTO m_event_types (name, key, color) VALUES ($1, $2, $3);

-- name: CreateEventType :one
INSERT INTO m_event_types (name, key, color) VALUES ($1, $2, $3) RETURNING *;

-- name: UpdateEventType :one
UPDATE m_event_types SET name = $2, key = $3, color = $4 WHERE event_type_id = $1 RETURNING *;

-- name: UpdateEventTypeByKey :one
UPDATE m_event_types SET name = $2, color = $3 WHERE key = $1 RETURNING *;

-- name: DeleteEventType :exec
DELETE FROM m_event_types WHERE event_type_id = $1;

-- name: DeleteEventTypeByKey :exec
DELETE FROM m_event_types WHERE key = $1;

-- name: FindEventTypeByID :one
SELECT * FROM m_event_types WHERE event_type_id = $1;

-- name: FindEventTypeByKey :one
SELECT * FROM m_event_types WHERE key = $1;

-- name: GetEventTypes :many
SELECT * FROM m_event_types
WHERE
	CASE WHEN @where_like_name::boolean = true THEN name LIKE '%' || @search_name::text || '%' ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN name END DESC,
	m_event_types_pkey ASC;

-- name: GetEventTypesUseNumberedPaginate :many
SELECT * FROM m_event_types
WHERE
	CASE WHEN @where_like_name::boolean = true THEN name LIKE '%' || @search_name::text || '%' ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN name END DESC,
	m_event_types_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetEventTypesUseKeysetPaginate :many
SELECT * FROM m_event_types
WHERE
	CASE WHEN @where_like_name::boolean = true THEN name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'name' THEN name > @name_cursor OR (name = @name_cursor AND m_event_types_pkey > @cursor::int)
				WHEN 'r_name' THEN name < @name_cursor OR (name = @name_cursor AND m_event_types_pkey > @cursor::int)
				ELSE m_event_types_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'name' THEN name < @name_cursor OR (name = @name_cursor AND m_event_types_pkey < @cursor::int)
				WHEN 'r_name' THEN name > @name_cursor OR (name = @name_cursor AND m_event_types_pkey < @cursor::int)
				ELSE m_event_types_pkey < @cursor::int
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'next' THEN name END ASC,
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'prev' THEN name END DESC,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'next' THEN name END ASC,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'prev' THEN name END DESC,
	CASE WHEN @cursor_direction::text = 'next' THEN m_event_types_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN m_event_types_pkey END DESC
LIMIT $1;

-- name: GetPluralEventTypes :many
SELECT * FROM m_event_types WHERE event_type_id = ANY(@event_type_ids::uuid[])
ORDER BY
	m_event_types_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountEventTypes :one
SELECT COUNT(*) FROM m_event_types
WHERE
	CASE WHEN @where_like_name::boolean = true THEN name LIKE '%' || @search_name::text || '%' ELSE TRUE END;

