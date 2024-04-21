-- name: CreatePolicyCategories :copyfrom
INSERT INTO m_policy_categories (name, description, key) VALUES ($1, $2, $3);

-- name: CreatePolicyCategory :one
INSERT INTO m_policy_categories (name, description, key) VALUES ($1, $2, $3) RETURNING *;

-- name: UpdatePolicyCategory :one
UPDATE m_policy_categories SET name = $2, description = $3 WHERE policy_category_id = $1 RETURNING *;

-- name: UpdatePolicyCategoryKey :one
UPDATE m_policy_categories SET key = $2 WHERE policy_category_id = $1 RETURNING *;

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
WHERE CASE
	WHEN @where_like_name::boolean = true THEN m_policy_categories.name LIKE '%' || @search_name::text || '%'
END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_policy_categories.name END ASC,
	m_policy_categories_pkey DESC
LIMIT $1 OFFSET $2;

-- name: GetPolicyCategoriesByKeys :many
SELECT * FROM m_policy_categories WHERE key = ANY(@keys::varchar[])
AND CASE
	WHEN @where_like_name::boolean = true THEN m_policy_categories.name LIKE '%' || @search_name::text || '%'
END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_policy_categories.name END ASC,
	m_policy_categories_pkey DESC
LIMIT $1 OFFSET $2;
