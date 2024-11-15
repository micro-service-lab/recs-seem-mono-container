-- name: CreateAttendStatuses :copyfrom
INSERT INTO m_attend_statuses (name, key) VALUES ($1, $2);

-- name: CreateAttendStatus :one
INSERT INTO m_attend_statuses (name, key) VALUES ($1, $2) RETURNING *;

-- name: UpdateAttendStatus :one
UPDATE m_attend_statuses SET name = $2, key = $3 WHERE attend_status_id = $1 RETURNING *;

-- name: UpdateAttendStatusByKey :one
UPDATE m_attend_statuses SET name = $2 WHERE key = $1 RETURNING *;

-- name: DeleteAttendStatus :execrows
DELETE FROM m_attend_statuses WHERE attend_status_id = $1;

-- name: DeleteAttendStatusByKey :execrows
DELETE FROM m_attend_statuses WHERE key = $1;

-- name: PluralDeleteAttendStatuses :execrows
DELETE FROM m_attend_statuses WHERE attend_status_id = ANY(@attend_status_ids::uuid[]);

-- name: FindAttendStatusByID :one
SELECT * FROM m_attend_statuses WHERE attend_status_id = $1;

-- name: FindAttendStatusByKey :one
SELECT * FROM m_attend_statuses WHERE key = $1;

-- name: GetAttendStatuses :many
SELECT * FROM m_attend_statuses
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_attend_statuses.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_attend_statuses.name END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' THEN m_attend_statuses.name END DESC NULLS LAST,
	m_attend_statuses_pkey ASC;

-- name: GetAttendStatusesUseNumberedPaginate :many
SELECT * FROM m_attend_statuses
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_attend_statuses.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_attend_statuses.name END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' THEN m_attend_statuses.name END DESC NULLS LAST,
	m_attend_statuses_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetAttendStatusesUseKeysetPaginate :many
SELECT * FROM m_attend_statuses
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_attend_statuses.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'name' THEN m_attend_statuses.name > @name_cursor OR (m_attend_statuses.name = @name_cursor AND m_attend_statuses_pkey > @cursor::int)
				WHEN 'r_name' THEN m_attend_statuses.name < @name_cursor OR (m_attend_statuses.name = @name_cursor AND m_attend_statuses_pkey > @cursor::int)
				ELSE m_attend_statuses_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'name' THEN m_attend_statuses.name < @name_cursor OR (m_attend_statuses.name = @name_cursor AND m_attend_statuses_pkey < @cursor::int)
				WHEN 'r_name' THEN m_attend_statuses.name > @name_cursor OR (m_attend_statuses.name = @name_cursor AND m_attend_statuses_pkey < @cursor::int)
				ELSE m_attend_statuses_pkey < @cursor::int
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'next' THEN m_attend_statuses.name END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'prev' THEN m_attend_statuses.name END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'next' THEN m_attend_statuses.name END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'prev' THEN m_attend_statuses.name END ASC NULLS LAST,
	CASE WHEN @cursor_direction::text = 'next' THEN m_attend_statuses_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN m_attend_statuses_pkey END DESC
LIMIT $1;

-- name: GetPluralAttendStatuses :many
SELECT * FROM m_attend_statuses
WHERE attend_status_id = ANY(@attend_status_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_attend_statuses.name END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' THEN m_attend_statuses.name END DESC NULLS LAST,
	m_attend_statuses_pkey ASC;

-- name: GetPluralAttendStatusesUseNumberedPaginate :many
SELECT * FROM m_attend_statuses
WHERE attend_status_id = ANY(@attend_status_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_attend_statuses.name END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' THEN m_attend_statuses.name END DESC NULLS LAST,
	m_attend_statuses_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountAttendStatuses :one
SELECT COUNT(*) FROM m_attend_statuses
WHERE
	CASE WHEN @where_like_name::boolean = true THEN name LIKE '%' || @search_name::text || '%' ELSE TRUE END;
