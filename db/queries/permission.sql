-- name: CreatePermissions :copyfrom
INSERT INTO m_permissions (name, description, key, permission_category_id) VALUES ($1, $2, $3, $4);

-- name: CreatePermission :one
INSERT INTO m_permissions (name, description, key, permission_category_id) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: UpdatePermission :one
UPDATE m_permissions SET name = $2, description = $3 WHERE permission_id = $1 RETURNING *;

-- name: UpdatePermissionKey :one
UPDATE m_permissions SET key = $2 WHERE permission_id = $1 RETURNING *;

-- name: UpdatePermissionCategoryID :one
UPDATE m_permissions SET permission_category_id = $2 WHERE permission_id = $1 RETURNING *;

-- name: UpdatePermissionByKey :one
UPDATE m_permissions SET name = $2, description = $3 WHERE key = $1 RETURNING *;

-- name: DeletePermission :exec
DELETE FROM m_permissions WHERE permission_id = $1;

-- name: DeletePermissionByKey :exec
DELETE FROM m_permissions WHERE key = $1;

-- name: FindPermissionById :one
SELECT * FROM m_permissions WHERE permission_id = $1;

-- name: FindPermissionByIDWithCategory :one
SELECT m_permissions.*, m_permission_categories.* FROM m_permissions
JOIN m_permission_categories ON m_permissions.permission_category_id = m_permission_categories.permission_category_id
WHERE m_permissions.permission_id = $1;

-- name: FindPermissionByKey :one
SELECT * FROM m_permissions WHERE key = $1;

-- name: GetPermissionByKeyWithCategory :one
SELECT m_permissions.*, m_permission_categories.* FROM m_permissions
JOIN m_permission_categories ON m_permissions.permission_category_id = m_permission_categories.permission_category_id
WHERE m_permissions.key = $1;

-- name: GetPermissions :many
SELECT * FROM m_permissions
WHERE CASE
	WHEN @where_like_name::boolean = true THEN m_permissions.name LIKE '%' || @search_name::text || '%'
END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_permissions.name END ASC,
	m_permissions_pkey DESC
LIMIT $1 OFFSET $2;

-- name: GetPermissionsWithCategory :many
SELECT m_permissions.*, m_permission_categories.* FROM m_permissions
JOIN m_permission_categories ON m_permissions.permission_category_id = m_permission_categories.permission_category_id
WHERE CASE
	WHEN @where_like_name::boolean = true THEN m_permissions.name LIKE '%' || @search_name::text || '%'
END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_permissions.name END ASC,
	m_permissions_pkey DESC
LIMIT $1 OFFSET $2;

-- name: GetPermissionsByCategory :many
SELECT * FROM m_permissions
WHERE permission_category_id = $1
AND CASE
	WHEN @where_like_name::boolean = true THEN m_permissions.name LIKE '%' || @search_name::text || '%'
END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_permissions.name END ASC,
	m_permissions_pkey DESC
LIMIT $2 OFFSET $3;

-- name: GetPermissionsByCategories :many
SELECT * FROM m_permissions
WHERE permission_category_id = ANY($1::uuid[])
AND CASE
	WHEN @where_like_name::boolean = true THEN m_permissions.name LIKE '%' || @search_name::text || '%'
END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_permissions.name END ASC,
	m_permissions_pkey DESC
LIMIT $2 OFFSET $3;

