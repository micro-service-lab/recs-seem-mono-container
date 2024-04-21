-- name: CreatePolicies :copyfrom
INSERT INTO m_policies (name, description, key, policy_category_id) VALUES ($1, $2, $3, $4);

-- name: CreatePolicy :one
INSERT INTO m_policies (name, description, key, policy_category_id) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: UpdatePolicy :one
UPDATE m_policies SET name = $2, description = $3 WHERE policy_id = $1 RETURNING *;

-- name: UpdatePolicyKey :one
UPDATE m_policies SET key = $2 WHERE policy_id = $1 RETURNING *;

-- name: UpdatePolicyCategoryID :one
UPDATE m_policies SET policy_category_id = $2 WHERE policy_id = $1 RETURNING *;

-- name: UpdatePolicyByKey :one
UPDATE m_policies SET name = $2, description = $3 WHERE key = $1 RETURNING *;

-- name: DeletePolicy :exec
DELETE FROM m_policies WHERE policy_id = $1;

-- name: DeletePolicyByKey :exec
DELETE FROM m_policies WHERE key = $1;

-- name: FindPolicyById :one
SELECT * FROM m_policies WHERE policy_id = $1;

-- name: FindPolicyByIDWithCategory :one
SELECT m_policies.*, m_policy_categories.* FROM m_policies
JOIN m_policy_categories ON m_policies.policy_category_id = m_policy_categories.policy_category_id
WHERE m_policies.policy_id = $1;

-- name: FindPolicyByKey :one
SELECT * FROM m_policies WHERE key = $1;

-- name: GetPolicyByKeyWithCategory :one
SELECT m_policies.*, m_policy_categories.* FROM m_policies
JOIN m_policy_categories ON m_policies.policy_category_id = m_policy_categories.policy_category_id
WHERE m_policies.key = $1;

-- name: GetPolicies :many
SELECT * FROM m_policies
WHERE CASE
	WHEN @where_like_name::boolean = true THEN m_policies.name LIKE '%' || @search_name::text || '%'
END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_policies.name END ASC,
	m_policies_pkey DESC
LIMIT $1 OFFSET $2;

-- name: GetPoliciesWithCategory :many
SELECT m_policies.*, m_policy_categories.* FROM m_policies
JOIN m_policy_categories ON m_policies.policy_category_id = m_policy_categories.policy_category_id
WHERE CASE
	WHEN @where_like_name::boolean = true THEN m_policies.name LIKE '%' || @search_name::text || '%'
END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_policies.name END ASC,
	m_policies_pkey DESC
LIMIT $1 OFFSET $2;

-- name: GetPoliciesByKeys :many
SELECT * FROM m_policies WHERE key = ANY(@keys::varchar[])
AND CASE
	WHEN @where_like_name::boolean = true THEN m_policies.name LIKE '%' || @search_name::text || '%'
END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_policies.name END ASC,
	m_policies_pkey DESC
LIMIT $1 OFFSET $2;

-- name: GetPoliciesByKeysWithCategory :many
SELECT m_policies.*, m_policy_categories.* FROM m_policies
JOIN m_policy_categories ON m_policies.policy_category_id = m_policy_categories.policy_category_id
WHERE m_policies.key = ANY(@keys::varchar[])
AND CASE
	WHEN @where_like_name::boolean = true THEN m_policies.name LIKE '%' || @search_name::text || '%'
END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_policies.name END ASC,
	m_policies_pkey DESC
LIMIT $1 OFFSET $2;

-- name: GetPoliciesByCategory :many
SELECT * FROM m_policies WHERE policy_category_id = $1
AND CASE
	WHEN @where_like_name::boolean = true THEN m_policies.name LIKE '%' || @search_name::text || '%'
END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_policies.name END ASC,
	m_policies_pkey DESC
LIMIT $2 OFFSET $3;

-- name: GetPoliciesByCategories :many
SELECT * FROM m_policies WHERE policy_category_id = ANY(@policy_category_ids::uuid[])
AND CASE
	WHEN @where_like_name::boolean = true THEN m_policies.name LIKE '%' || @search_name::text || '%'
END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_policies.name END ASC,
	m_policies_pkey DESC
LIMIT $1 OFFSET $2;


