-- name: CreatePolicyCategories :copyfrom
INSERT INTO m_policy_categories (name, description, key) VALUES ($1, $2, $3);

-- name: CreatePolicyCategory :one
INSERT INTO m_policy_categories (name, description, key) VALUES ($1, $2, $3) RETURNING *;

-- name: UpdatePolicyCategory :one
UPDATE m_policy_categories SET name = $2, description = $3, key = $4 WHERE policy_category_id = $1 RETURNING *;

-- name: UpdatePolicyCategoryByKey :one
UPDATE m_policy_categories SET name = $2, description = $3 WHERE key = $1 RETURNING *;

-- name: DeletePolicyCategory :execrows
DELETE FROM m_policy_categories WHERE policy_category_id = $1;

-- name: DeletePolicyCategoryByKey :execrows
DELETE FROM m_policy_categories WHERE key = $1;

-- name: PluralDeletePolicyCategories :execrows
DELETE FROM m_policy_categories WHERE policy_category_id = ANY(@policy_category_ids::uuid[]);

-- name: FindPolicyCategoryByID :one
SELECT * FROM m_policy_categories WHERE policy_category_id = $1;

-- name: FindPolicyCategoryByKey :one
SELECT * FROM m_policy_categories WHERE key = $1;

-- name: GetPolicyCategories :many
SELECT * FROM m_policy_categories
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_policy_categories.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_policy_categories.name END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' THEN m_policy_categories.name END DESC NULLS LAST,
	m_policy_categories_pkey ASC;

-- name: GetPolicyCategoriesUseNumberedPaginate :many
SELECT * FROM m_policy_categories
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_policy_categories.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_policy_categories.name END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' THEN m_policy_categories.name END DESC NULLS LAST,
	m_policy_categories_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetPolicyCategoriesUseKeysetPaginate :many
SELECT * FROM m_policy_categories
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_policy_categories.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'name' THEN name > @name_cursor OR (name = @name_cursor AND m_policy_categories_pkey > @cursor::int)
				WHEN 'r_name' THEN name < @name_cursor OR (name = @name_cursor AND m_policy_categories_pkey > @cursor::int)
				ELSE m_policy_categories_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'name' THEN name < @name_cursor OR (name = @name_cursor AND m_policy_categories_pkey < @cursor::int)
				WHEN 'r_name' THEN name > @name_cursor OR (name = @name_cursor AND m_policy_categories_pkey < @cursor::int)
				ELSE m_policy_categories_pkey < @cursor::int
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'next' THEN m_policy_categories.name END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'prev' THEN m_policy_categories.name END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'next' THEN m_policy_categories.name END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'prev' THEN m_policy_categories.name END ASC NULLS LAST,
	CASE WHEN @cursor_direction::text = 'next' THEN m_policy_categories_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN m_policy_categories_pkey END DESC
LIMIT $1;

-- name: GetPluralPolicyCategories :many
SELECT * FROM m_policy_categories
WHERE policy_category_id = ANY(@policy_category_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_policy_categories.name END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' THEN m_policy_categories.name END DESC NULLS LAST,
	m_policy_categories_pkey ASC;

-- name: GetPluralPolicyCategoriesUseNumberedPaginate :many
SELECT * FROM m_policy_categories
WHERE policy_category_id = ANY(@policy_category_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_policy_categories.name END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' THEN m_policy_categories.name END DESC NULLS LAST,
	m_policy_categories_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountPolicyCategories :one
SELECT COUNT(*) FROM m_policy_categories
WHERE
	CASE WHEN @where_like_name::boolean = true THEN name LIKE '%' || @search_name::text || '%' ELSE TRUE END;
