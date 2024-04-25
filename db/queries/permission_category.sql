-- name: CreatePermissionCategories :copyfrom
INSERT INTO m_permission_categories (name, description, key) VALUES ($1, $2, $3);

-- name: CreatePermissionCategory :one
INSERT INTO m_permission_categories (name, description, key) VALUES ($1, $2, $3) RETURNING *;

-- name: UpdatePermissionCategory :one
UPDATE m_permission_categories SET name = $2, description = $3, key = $4 WHERE permission_category_id = $1 RETURNING *;

-- name: UpdatePermissionCategoryByKey :one
UPDATE m_permission_categories SET name = $2, description = $3 WHERE key = $1 RETURNING *;

-- name: DeletePermissionCategory :exec
DELETE FROM m_permission_categories WHERE permission_category_id = $1;

-- name: DeletePermissionCategoryByKey :exec
DELETE FROM m_permission_categories WHERE key = $1;

-- name: FindPermissionCategoryByID :one
SELECT * FROM m_permission_categories WHERE permission_category_id = $1;

-- name: FindPermissionCategoryByKey :one
SELECT * FROM m_permission_categories WHERE key = $1;

-- name: GetPermissionCategories :many
SELECT * FROM m_permission_categories
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_permission_categories.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_permission_categories.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_permission_categories.name END DESC,
	m_permission_categories_pkey DESC
LIMIT $1 OFFSET $2;

-- name: CountPermissionCategories :one
SELECT COUNT(*) FROM m_permission_categories
WHERE
	CASE WHEN @where_like_name::boolean = true THEN name LIKE '%' || @search_name::text || '%' ELSE TRUE END;

