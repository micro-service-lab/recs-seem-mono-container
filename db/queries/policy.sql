-- name: CreatePolicies :copyfrom
INSERT INTO m_policies (name, description, key, policy_category_id) VALUES ($1, $2, $3, $4);

-- name: CreatePolicy :one
INSERT INTO m_policies (name, description, key, policy_category_id) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: UpdatePolicy :one
UPDATE m_policies SET name = $2, description = $3, key = $4, policy_category_id = $5 WHERE policy_id = $1 RETURNING *;

-- name: UpdatePolicyByKey :one
UPDATE m_policies SET name = $2, description = $3, policy_category_id = $4 WHERE key = $1 RETURNING *;

-- name: DeletePolicy :execrows
DELETE FROM m_policies WHERE policy_id = $1;

-- name: DeletePolicyByKey :execrows
DELETE FROM m_policies WHERE key = $1;

-- name: PluralDeletePolicies :execrows
DELETE FROM m_policies WHERE policy_id = ANY(@policy_ids::uuid[]);

-- name: FindPolicyByID :one
SELECT * FROM m_policies WHERE policy_id = $1;

-- name: FindPolicyByIDWithCategory :one
SELECT m_policies.*, m_policy_categories.name policy_category_name, m_policy_categories.key policy_category_key, m_policy_categories.description policy_category_description FROM m_policies
JOIN m_policy_categories ON m_policies.policy_category_id = m_policy_categories.policy_category_id
WHERE m_policies.policy_id = $1;

-- name: FindPolicyByKey :one
SELECT * FROM m_policies WHERE key = $1;

-- name: FindPolicyByKeyWithCategory :one
SELECT m_policies.*, m_policy_categories.name policy_category_name, m_policy_categories.key policy_category_key, m_policy_categories.description policy_category_description FROM m_policies
JOIN m_policy_categories ON m_policies.policy_category_id = m_policy_categories.policy_category_id
WHERE m_policies.key = $1;

-- name: GetPolicies :many
SELECT * FROM m_policies
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_policies.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_in_category::boolean = true THEN policy_category_id = ANY(@in_categories::uuid[]) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_policies.name END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' THEN m_policies.name END DESC NULLS LAST,
	m_policies_pkey ASC;

-- name: GetPoliciesUseNumberedPaginate :many
SELECT * FROM m_policies
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_policies.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_in_category::boolean = true THEN policy_category_id = ANY(@in_categories::uuid[]) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_policies.name END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' THEN m_policies.name END DESC NULLS LAST,
	m_policies_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetPoliciesUseKeysetPaginate :many
SELECT * FROM m_policies
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_policies.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_in_category::boolean = true THEN policy_category_id = ANY(@in_categories::uuid[]) ELSE TRUE END
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'name' THEN m_policies.name > @name_cursor OR (m_policies.name = @name_cursor AND m_policies_pkey > @cursor::int)
				WHEN 'r_name' THEN m_policies.name < @name_cursor OR (m_policies.name = @name_cursor AND m_policies_pkey > @cursor::int)
				ELSE m_policies_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'name' THEN m_policies.name < @name_cursor OR (m_policies.name = @name_cursor AND m_policies_pkey < @cursor::int)
				WHEN 'r_name' THEN m_policies.name > @name_cursor OR (m_policies.name = @name_cursor AND m_policies_pkey < @cursor::int)
				ELSE m_policies_pkey < @cursor::int
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'next' THEN m_policies.name END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'prev' THEN m_policies.name END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'next' THEN m_policies.name END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'prev' THEN m_policies.name END ASC NULLS LAST,
	CASE WHEN @cursor_direction::text = 'next' THEN m_policies_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN m_policies_pkey END DESC
LIMIT $1;

-- name: GetPluralPolicies :many
SELECT * FROM m_policies WHERE policy_id = ANY(@policy_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_policies.name END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' THEN m_policies.name END DESC NULLS LAST,
	m_policies_pkey ASC;

-- name: GetPluralPoliciesUseNumberedPaginate :many
SELECT * FROM m_policies WHERE policy_id = ANY(@policy_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_policies.name END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' THEN m_policies.name END DESC NULLS LAST,
	m_policies_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetPoliciesWithCategory :many
SELECT m_policies.*, m_policy_categories.name policy_category_name, m_policy_categories.key policy_category_key, m_policy_categories.description policy_category_description FROM m_policies
JOIN m_policy_categories ON m_policies.policy_category_id = m_policy_categories.policy_category_id
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_policies.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_in_category::boolean = true THEN m_policies.policy_category_id = ANY(@in_categories::uuid[]) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_policies.name END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' THEN m_policies.name END DESC NULLS LAST,
	m_policies_pkey ASC;

-- name: GetPoliciesWithCategoryUseNumberedPaginate :many
SELECT m_policies.*, m_policy_categories.name policy_category_name, m_policy_categories.key policy_category_key, m_policy_categories.description policy_category_description FROM m_policies
JOIN m_policy_categories ON m_policies.policy_category_id = m_policy_categories.policy_category_id
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_policies.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_in_category::boolean = true THEN  m_policies.policy_category_id = ANY(@in_categories::uuid[]) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_policies.name END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' THEN m_policies.name END DESC NULLS LAST,
	m_policies_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetPoliciesWithCategoryUseKeysetPaginate :many
SELECT m_policies.*, m_policy_categories.name policy_category_name, m_policy_categories.key policy_category_key, m_policy_categories.description policy_category_description FROM m_policies
JOIN m_policy_categories ON m_policies.policy_category_id = m_policy_categories.policy_category_id
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_policies.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_in_category::boolean = true THEN  m_policies.policy_category_id = ANY(@in_categories::uuid[]) ELSE TRUE END
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'name' THEN m_policies.name > @name_cursor OR (m_policies.name = @name_cursor AND m_policies_pkey > @cursor::int)
				WHEN 'r_name' THEN m_policies.name < @name_cursor OR (m_policies.name = @name_cursor AND m_policies_pkey > @cursor::int)
				ELSE m_policies_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'name' THEN m_policies.name < @name_cursor OR (m_policies.name = @name_cursor AND m_policies_pkey < @cursor::int)
				WHEN 'r_name' THEN m_policies.name > @name_cursor OR (m_policies.name = @name_cursor AND m_policies_pkey < @cursor::int)
				ELSE m_policies_pkey < @cursor::int
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'next' THEN m_policies.name END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'prev' THEN m_policies.name END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'next' THEN m_policies.name END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'prev' THEN m_policies.name END ASC NULLS LAST,
	CASE WHEN @cursor_direction::text = 'next' THEN m_policies_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN m_policies_pkey END DESC
LIMIT $1;

-- name: GetPluralPoliciesWithCategory :many
SELECT m_policies.*, m_policy_categories.name policy_category_name, m_policy_categories.key policy_category_key, m_policy_categories.description policy_category_description FROM m_policies
JOIN m_policy_categories ON m_policies.policy_category_id = m_policy_categories.policy_category_id
WHERE policy_id = ANY(@policy_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_policies.name END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' THEN m_policies.name END DESC NULLS LAST,
	m_policies_pkey ASC;

-- name: GetPluralPoliciesWithCategoryUseNumberedPaginate :many
SELECT m_policies.*, m_policy_categories.name policy_category_name, m_policy_categories.key policy_category_key, m_policy_categories.description policy_category_description FROM m_policies
JOIN m_policy_categories ON m_policies.policy_category_id = m_policy_categories.policy_category_id
WHERE policy_id = ANY(@policy_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_policies.name END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' THEN m_policies.name END DESC NULLS LAST,
	m_policies_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountPolicies :one
SELECT COUNT(*) FROM m_policies
WHERE
	CASE WHEN @where_like_name::boolean = true THEN name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_in_category::boolean = true THEN policy_category_id = ANY(@in_categories::uuid[]) ELSE TRUE END;


