-- name: CreateWorkPositions :copyfrom
INSERT INTO m_work_positions (name, description, created_at, updated_at) VALUES ($1, $2, $3, $4);

-- name: CreateWorkPosition :one
INSERT INTO m_work_positions (name, description, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: UpdateWorkPosition :one
UPDATE m_work_positions SET name = $2, description = $3, updated_at = $4 WHERE work_position_id = $1 RETURNING *;

-- name: DeleteWorkPosition :exec
DELETE FROM m_work_positions WHERE work_position_id = $1;

-- name: FindWorkPositionByID :one
SELECT * FROM m_work_positions WHERE work_position_id = $1;

-- name: GetWorkPositions :many
SELECT * FROM m_work_positions
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_work_positions.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_work_positions.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_work_positions.name END DESC,
	m_work_positions_pkey DESC;

-- name: GetWorkPositionsUseNumberedPaginate :many
SELECT * FROM m_work_positions
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_work_positions.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_work_positions.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_work_positions.name END DESC,
	m_work_positions_pkey DESC
LIMIT $1 OFFSET $2;

-- name: GetWorkPositionsUseKeysetPaginate :many
SELECT * FROM m_work_positions
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_work_positions.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE @cursor_direction
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'name' THEN name > @cursor_column OR (name = @cursor_column AND m_work_positions_pkey < @cursor)
				WHEN 'r_name' THEN name < @cursor_column OR (name = @cursor_column AND m_work_positions_pkey < @cursor)
				ELSE m_work_positions_pkey < @cursor
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'name' THEN name < @cursor_column OR (name = @cursor_column AND m_work_positions_pkey > @cursor)
				WHEN 'r_name' THEN name > @cursor_column OR (name = @cursor_column AND m_work_positions_pkey > @cursor)
				ELSE m_work_positions_pkey > @cursor
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_work_positions.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_work_positions.name END DESC,
	m_work_positions_pkey DESC
LIMIT $1;

-- name: CountWorkPositions :one
SELECT COUNT(*) FROM m_work_positions
WHERE
	CASE WHEN @where_like_name::boolean = true THEN name LIKE '%' || @search_name::text || '%' ELSE TRUE END;
