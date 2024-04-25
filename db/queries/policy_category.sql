-- name: CreatePolicyCategories :copyfrom
INSERT INTO m_policy_categories (name, description, key) VALUES ($1, $2, $3);

-- name: CreatePolicyCategory :one
INSERT INTO m_policy_categories (name, description, key) VALUES ($1, $2, $3) RETURNING *;

-- name: UpdatePolicyCategory :one
UPDATE m_policy_categories SET name = $2, description = $3, key = $4 WHERE policy_category_id = $1 RETURNING *;

-- name: UpdatePolicyCategoryByKey :one
UPDATE m_policy_categories SET name = $2, description = $3 WHERE key = $1 RETURNING *;

-- name: DeletePolicyCategory :exec
DELETE FROM m_policy_categories WHERE policy_category_id = $1;

-- name: DeletePolicyCategoryByKey :exec
DELETE FROM m_policy_categories WHERE key = $1;

-- name: FindPolicyCategoryById :one
SELECT * FROM m_policy_categories WHERE policy_category_id = $1;

-- name: FindPolicyCategoryByKey :one
SELECT * FROM m_policy_categories WHERE key = $1;

-- name: GetPolicyCategories :many
SELECT * FROM m_policy_categories
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_policy_categories.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_policy_categories.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_policy_categories.name END DESC,
	m_policy_categories_pkey DESC;

-- name: GetPolicyCategoriesUseNumberedPaginate :many
SELECT * FROM m_policy_categories
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_policy_categories.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_policy_categories.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_policy_categories.name END DESC,
	m_policy_categories_pkey DESC
LIMIT $1 OFFSET $2;

-- name: GetPolicyCategoriesUseKeysetPaginate :many
SELECT * FROM m_policy_categories
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_policy_categories.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE @cursor_direction
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'name' THEN name > @cursor_column OR (name = @cursor_column AND m_policy_categories_pkey < @cursor)
				WHEN 'r_name' THEN name < @cursor_column OR (name = @cursor_column AND m_policy_categories_pkey < @cursor)
				ELSE m_policy_categories_pkey < @cursor
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'name' THEN name < @cursor_column OR (name = @cursor_column AND m_policy_categories_pkey > @cursor)
				WHEN 'r_name' THEN name > @cursor_column OR (name = @cursor_column AND m_policy_categories_pkey > @cursor)
				ELSE m_policy_categories_pkey > @cursor
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_policy_categories.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_policy_categories.name END DESC,
	m_policy_categories_pkey DESC;

-- name: CountPolicyCategories :one
SELECT COUNT(*) FROM m_policy_categories
WHERE
	CASE WHEN @where_like_name::boolean = true THEN name LIKE '%' || @search_name::text || '%' ELSE TRUE END;