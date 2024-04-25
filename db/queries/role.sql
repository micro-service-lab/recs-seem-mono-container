-- name: CreateRoles :copyfrom
INSERT INTO m_roles (name, description, created_at, updated_at) VALUES ($1, $2, $3, $4);

-- name: CreateRole :one
INSERT INTO m_roles (name, description, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: UpdateRole :one
UPDATE m_roles SET name = $2, description = $3, updated_at = $4 WHERE role_id = $1 RETURNING *;

-- name: DeleteRole :exec
DELETE FROM m_roles WHERE role_id = $1;

-- name: FindRoleByID :one
SELECT * FROM m_roles WHERE role_id = $1;

-- name: GetRoles :many
SELECT * FROM m_roles
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_roles.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_roles.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_roles.name END DESC,
	m_roles_pkey DESC

-- name: GetRolesUseNumberedPaginate :many
SELECT * FROM m_roles
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_roles.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_roles.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_roles.name END DESC,
	m_roles_pkey DESC
LIMIT $1 OFFSET $2;

-- name: GetRolesUseKeysetPaginate :many
SELECT * FROM m_roles
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_roles.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE @cursor_direction
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'name' THEN name > @cursor_column OR (name = @cursor_column AND m_roles_pkey < @cursor)
				WHEN 'r_name' THEN name < @cursor_column OR (name = @cursor_column AND m_roles_pkey < @cursor)
				ELSE m_roles_pkey < @cursor
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'name' THEN name < @cursor_column OR (name = @cursor_column AND m_roles_pkey > @cursor)
				WHEN 'r_name' THEN name > @cursor_column OR (name = @cursor_column AND m_roles_pkey > @cursor)
				ELSE m_roles_pkey > @cursor
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_roles.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_roles.name END DESC,
	m_roles_pkey DESC
LIMIT $1;

-- name: CountRoles :one
SELECT COUNT(*) FROM m_roles
WHERE
	CASE WHEN @where_like_name::boolean = true THEN name LIKE '%' || @search_name::text || '%' ELSE TRUE END;
