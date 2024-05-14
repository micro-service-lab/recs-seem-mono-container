-- name: CreateWorkPositions :copyfrom
INSERT INTO m_work_positions (name, organization_id, description, created_at, updated_at) VALUES ($1, $2, $3, $4, $5);

-- name: CreateWorkPosition :one
INSERT INTO m_work_positions (name, organization_id, description, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: UpdateWorkPosition :one
UPDATE m_work_positions SET name = $2, description = $3, updated_at = $4 WHERE work_position_id = $1 RETURNING *;

-- name: DeleteWorkPosition :execrows
DELETE FROM m_work_positions WHERE work_position_id = $1;

-- name: PluralDeleteWorkPositions :execrows
DELETE FROM m_work_positions WHERE work_position_id = ANY($1::uuid[]);

-- name: FindWorkPositionByID :one
SELECT * FROM m_work_positions WHERE work_position_id = $1;

-- name: GetWorkPositions :many
SELECT * FROM m_work_positions
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_work_positions.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_in_organization_id::boolean = true THEN m_work_positions.organization_id = ANY(@organization_ids::uuid[]) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_work_positions.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_work_positions.name END DESC,
	m_work_positions_pkey ASC;

-- name: GetWorkPositionsUseNumberedPaginate :many
SELECT * FROM m_work_positions
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_work_positions.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_in_organization_id::boolean = true THEN m_work_positions.organization_id = ANY(@organization_ids::uuid[]) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_work_positions.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_work_positions.name END DESC,
	m_work_positions_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetWorkPositionsUseKeysetPaginate :many
SELECT * FROM m_work_positions
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_work_positions.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_in_organization_id::boolean = true THEN m_work_positions.organization_id = ANY(@organization_ids::uuid[]) ELSE TRUE END
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'name' THEN name > @name_cursor OR (name = @name_cursor AND m_work_positions_pkey > @cursor::int)
				WHEN 'r_name' THEN name < @name_cursor OR (name = @name_cursor AND m_work_positions_pkey > @cursor::int)
				ELSE m_work_positions_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'name' THEN name < @name_cursor OR (name = @name_cursor AND m_work_positions_pkey < @cursor::int)
				WHEN 'r_name' THEN name > @name_cursor OR (name = @name_cursor AND m_work_positions_pkey < @cursor::int)
				ELSE m_work_positions_pkey < @cursor::int
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'next' THEN m_work_positions.name END ASC,
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'prev' THEN m_work_positions.name END DESC,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'next' THEN m_work_positions.name END DESC,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'prev' THEN m_work_positions.name END ASC,
	CASE WHEN @cursor_direction::text = 'next' THEN m_work_positions_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN m_work_positions_pkey END DESC
LIMIT $1;

-- name: GetPluckWorkPositions :many
SELECT work_position_id, name FROM m_work_positions
WHERE
	work_position_id = ANY(@work_position_ids::uuid[])
ORDER BY
	m_work_positions_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountWorkPositions :one
SELECT COUNT(*) FROM m_work_positions
WHERE
	CASE WHEN @where_like_name::boolean = true THEN name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_in_organization_id::boolean = true THEN m_work_positions.organization_id = ANY(@organization_ids::uuid[]) ELSE TRUE END;
