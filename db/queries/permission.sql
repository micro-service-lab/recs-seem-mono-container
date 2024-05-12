-- name: CreatePermissions :copyfrom
INSERT INTO m_permissions (name, description, key, permission_category_id) VALUES ($1, $2, $3, $4);

-- name: CreatePermission :one
INSERT INTO m_permissions (name, description, key, permission_category_id) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: UpdatePermission :one
UPDATE m_permissions SET name = $2, description = $3, key = $4, permission_category_id = $5 WHERE permission_id = $1 RETURNING *;

-- name: UpdatePermissionByKey :one
UPDATE m_permissions SET name = $2, description = $3, permission_category_id = $4 WHERE key = $1 RETURNING *;

-- name: DeletePermission :exec
DELETE FROM m_permissions WHERE permission_id = $1;

-- name: DeletePermissionByKey :exec
DELETE FROM m_permissions WHERE key = $1;

-- name: PluralDeletePermissions :exec
DELETE FROM m_permissions WHERE permission_id = ANY($1::uuid[]);

-- name: FindPermissionByID :one
SELECT * FROM m_permissions WHERE permission_id = $1;

-- name: FindPermissionByIDWithCategory :one
SELECT m_permissions.*,  m_permission_categories.name permission_category_name, m_permission_categories.key permission_category_key, m_permission_categories.description permission_category_description FROM m_permissions
JOIN m_permission_categories ON m_permissions.permission_category_id = m_permission_categories.permission_category_id
WHERE m_permissions.permission_id = $1;

-- name: FindPermissionByKey :one
SELECT * FROM m_permissions WHERE key = $1;

-- name: FindPermissionByKeyWithCategory :one
SELECT m_permissions.*,  m_permission_categories.name permission_category_name, m_permission_categories.key permission_category_key, m_permission_categories.description permission_category_description FROM m_permissions
JOIN m_permission_categories ON m_permissions.permission_category_id = m_permission_categories.permission_category_id
WHERE m_permissions.key = $1;

-- name: GetPermissions :many
SELECT * FROM m_permissions
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_permissions.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_in_category::boolean = true THEN permission_category_id = ANY(@in_categories::uuid[]) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_permissions.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_permissions.name END DESC,
	m_permissions_pkey ASC;

-- name: GetPermissionsUseNumberedPaginate :many
SELECT * FROM m_permissions
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_permissions.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_in_category::boolean = true THEN permission_category_id = ANY(@in_categories::uuid[]) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_permissions.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_permissions.name END DESC,
	m_permissions_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetPermissionsUseKeysetPaginate :many
SELECT * FROM m_permissions
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_permissions.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_in_category::boolean = true THEN permission_category_id = ANY(@in_categories::uuid[]) ELSE TRUE END
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'name' THEN m_permissions.name > @name_cursor OR (m_permissions.name = @name_cursor AND m_permissions_pkey > @cursor::int)
				WHEN 'r_name' THEN m_permissions.name < @name_cursor OR (m_permissions.name = @name_cursor AND m_permissions_pkey > @cursor::int)
				ELSE m_permissions_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'name' THEN m_permissions.name < @name_cursor OR (m_permissions.name = @name_cursor AND m_permissions_pkey < @cursor::int)
				WHEN 'r_name' THEN m_permissions.name > @name_cursor OR (m_permissions.name = @name_cursor AND m_permissions_pkey < @cursor::int)
				ELSE m_permissions_pkey < @cursor::int
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'next' THEN m_permissions.name END ASC,
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'prev' THEN m_permissions.name END DESC,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'next' THEN m_permissions.name END DESC,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'prev' THEN m_permissions.name END ASC,
	CASE WHEN @cursor_direction::text = 'next' THEN m_permissions_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN m_permissions_pkey END DESC
LIMIT $1;

-- name: GetPluralPermissions :many
SELECT * FROM m_permissions WHERE permission_id = ANY(@permission_ids::uuid[])
ORDER BY
	m_permissions_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetPermissionsWithCategory :many
SELECT m_permissions.*, m_permission_categories.name permission_category_name, m_permission_categories.key permission_category_key, m_permission_categories.description permission_category_description FROM m_permissions
JOIN m_permission_categories ON m_permissions.permission_category_id = m_permission_categories.permission_category_id
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_permissions.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_in_category::boolean = true THEN m_permissions.permission_category_id = ANY(@in_categories::uuid[]) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_permissions.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_permissions.name END DESC,
	m_permissions_pkey ASC;

-- name: GetPermissionsWithCategoryUseNumberedPaginate :many
SELECT m_permissions.*, m_permission_categories.name permission_category_name, m_permission_categories.key permission_category_key, m_permission_categories.description permission_category_description FROM m_permissions
JOIN m_permission_categories ON m_permissions.permission_category_id = m_permission_categories.permission_category_id
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_permissions.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_in_category::boolean = true THEN m_permissions.permission_category_id = ANY(@in_categories::uuid[]) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_permissions.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_permissions.name END DESC,
	m_permissions_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetPermissionsWithCategoryUseKeysetPaginate :many
SELECT m_permissions.*, m_permission_categories.name permission_category_name, m_permission_categories.key permission_category_key, m_permission_categories.description permission_category_description FROM m_permissions
JOIN m_permission_categories ON m_permissions.permission_category_id = m_permission_categories.permission_category_id
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_permissions.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_in_category::boolean = true THEN m_permissions.permission_category_id = ANY(@in_categories::uuid[]) ELSE TRUE END
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'name' THEN m_permissions.name > @name_cursor OR (m_permissions.name = @name_cursor AND m_permissions_pkey > @cursor::int)
				WHEN 'r_name' THEN m_permissions.name < @name_cursor OR (m_permissions.name = @name_cursor AND m_permissions_pkey > @cursor::int)
				ELSE m_permissions_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'name' THEN m_permissions.name < @name_cursor OR (m_permissions.name = @name_cursor AND m_permissions_pkey < @cursor::int)
				WHEN 'r_name' THEN m_permissions.name > @name_cursor OR (m_permissions.name = @name_cursor AND m_permissions_pkey < @cursor::int)
				ELSE m_permissions_pkey < @cursor::int
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'next' THEN m_permissions.name END ASC,
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'prev' THEN m_permissions.name END DESC,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'next' THEN m_permissions.name END DESC,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'prev' THEN m_permissions.name END ASC,
	CASE WHEN @cursor_direction::text = 'next' THEN m_permissions_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN m_permissions_pkey END DESC
LIMIT $1;

-- name: GetPluralPermissionsWithCategory :many
SELECT m_permissions.*,  m_permission_categories.name permission_category_name, m_permission_categories.key permission_category_key, m_permission_categories.description permission_category_description FROM m_permissions
JOIN m_permission_categories ON m_permissions.permission_category_id = m_permission_categories.permission_category_id
WHERE permission_id = ANY(@permission_ids::uuid[])
ORDER BY
	m_permissions_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountPermissions :one
SELECT COUNT(*) FROM m_permissions
WHERE
	CASE WHEN @where_like_name::boolean = true THEN name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_in_category::boolean = true THEN permission_category_id = ANY(@in_categories::uuid[]) ELSE TRUE END;

