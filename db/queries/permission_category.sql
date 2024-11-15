-- name: CreatePermissionCategories :copyfrom
INSERT INTO m_permission_categories (name, description, key) VALUES ($1, $2, $3);

-- name: CreatePermissionCategory :one
INSERT INTO m_permission_categories (name, description, key) VALUES ($1, $2, $3) RETURNING *;

-- name: UpdatePermissionCategory :one
UPDATE m_permission_categories SET name = $2, description = $3, key = $4 WHERE permission_category_id = $1 RETURNING *;

-- name: UpdatePermissionCategoryByKey :one
UPDATE m_permission_categories SET name = $2, description = $3 WHERE key = $1 RETURNING *;

-- name: DeletePermissionCategory :execrows
DELETE FROM m_permission_categories WHERE permission_category_id = $1;

-- name: DeletePermissionCategoryByKey :execrows
DELETE FROM m_permission_categories WHERE key = $1;

-- name: PluralDeletePermissionCategories :execrows
DELETE FROM m_permission_categories WHERE permission_category_id = ANY(@permission_category_ids::uuid[]);

-- name: FindPermissionCategoryByID :one
SELECT * FROM m_permission_categories WHERE permission_category_id = $1;

-- name: FindPermissionCategoryByKey :one
SELECT * FROM m_permission_categories WHERE key = $1;

-- name: GetPermissionCategories :many
SELECT * FROM m_permission_categories
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_permission_categories.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_permission_categories.name END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' THEN m_permission_categories.name END DESC NULLS LAST,
	m_permission_categories_pkey ASC;

-- name: GetPermissionCategoriesUseNumberedPaginate :many
SELECT * FROM m_permission_categories
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_permission_categories.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_permission_categories.name END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' THEN m_permission_categories.name END DESC NULLS LAST,
	m_permission_categories_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetPermissionCategoriesUseKeysetPaginate :many
SELECT * FROM m_permission_categories
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_permission_categories.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'name' THEN name > @name_cursor OR (name = @name_cursor AND m_permission_categories_pkey > @cursor::int)
				WHEN 'r_name' THEN name < @name_cursor OR (name = @name_cursor AND m_permission_categories_pkey > @cursor::int)
				ELSE m_permission_categories_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'name' THEN name < @name_cursor OR (name = @name_cursor AND m_permission_categories_pkey < @cursor::int)
				WHEN 'r_name' THEN name > @name_cursor OR (name = @name_cursor AND m_permission_categories_pkey < @cursor::int)
				ELSE m_permission_categories_pkey < @cursor::int
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'next' THEN m_permission_categories.name END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'prev' THEN m_permission_categories.name END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'next' THEN m_permission_categories.name END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'prev' THEN m_permission_categories.name END ASC NULLS LAST,
	CASE WHEN @cursor_direction::text = 'next' THEN m_permission_categories_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN m_permission_categories_pkey END DESC
LIMIT $1;

-- name: GetPluralPermissionCategories :many
SELECT * FROM m_permission_categories
WHERE permission_category_id = ANY(@permission_category_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_permission_categories.name END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' THEN m_permission_categories.name END DESC NULLS LAST,
	m_permission_categories_pkey ASC;

-- name: GetPluralPermissionCategoriesUseNumberedPaginate :many
SELECT * FROM m_permission_categories
WHERE permission_category_id = ANY(@permission_category_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_permission_categories.name END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' THEN m_permission_categories.name END DESC NULLS LAST,
	m_permission_categories_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountPermissionCategories :one
SELECT COUNT(*) FROM m_permission_categories
WHERE
	CASE WHEN @where_like_name::boolean = true THEN name LIKE '%' || @search_name::text || '%' ELSE TRUE END;

