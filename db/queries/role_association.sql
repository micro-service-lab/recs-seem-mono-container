-- name: CreateRoleAssociations :copyfrom
INSERT INTO m_role_associations (role_id, policy_id) VALUES ($1, $2);

-- name: CreateRoleAssociation :one
INSERT INTO m_role_associations (role_id, policy_id) VALUES ($1, $2) RETURNING *;

-- name: DeleteRoleAssociation :exec
DELETE FROM m_role_associations WHERE role_id = $1 AND policy_id = $2;

-- name: GetPoliciesOnRole :many
SELECT sqlc.embed(m_role_associations), sqlc.embed(m_policies) FROM m_role_associations
LEFT JOIN m_policies ON m_role_associations.policy_id = m_policies.policy_id
WHERE role_id = $1
AND
	CASE WHEN @where_like_name::boolean = true THEN m_policies.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_policies.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_policies.name END DESC,
	m_role_associations_pkey DESC;

-- name: GetPoliciesOnRoleUseNumberedPaginate :many
SELECT sqlc.embed(m_role_associations), sqlc.embed(m_policies) FROM m_role_associations
LEFT JOIN m_policies ON m_role_associations.policy_id = m_policies.policy_id
WHERE role_id = $1
AND
	CASE WHEN @where_like_name::boolean = true THEN m_policies.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_policies.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_policies.name END DESC,
	m_role_associations_pkey DESC
LIMIT $2 OFFSET $3;

-- name: GetPoliciesOnRoleUseKeysetPaginate :many
SELECT sqlc.embed(m_role_associations), sqlc.embed(m_policies) FROM m_role_associations
LEFT JOIN m_policies ON m_role_associations.policy_id = m_policies.policy_id
WHERE role_id = $1
AND
	CASE WHEN @where_like_name::boolean = true THEN m_policies.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE @cursor_direction
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'name' THEN m_policies.name > @cursor_column OR (m_policies.name = @cursor_column AND m_role_associations_pkey < @cursor)
				WHEN 'r_name' THEN m_policies.name < @cursor_column OR (m_policies.name = @cursor_column AND m_role_associations_pkey < @cursor)
				ELSE m_role_associations_pkey < @cursor
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'name' THEN m_policies.name < @cursor_column OR (m_policies.name = @cursor_column AND m_role_associations_pkey > @cursor)
				WHEN 'r_name' THEN m_policies.name > @cursor_column OR (m_policies.name = @cursor_column AND m_role_associations_pkey > @cursor)
				ELSE m_role_associations_pkey > @cursor
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_policies.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_policies.name END DESC,
	m_role_associations_pkey DESC
LIMIT $2;

-- name: CountPoliciesByRoleID :one
SELECT COUNT(*) FROM m_role_associations
LEFT JOIN m_policies ON m_role_associations.policy_id = m_policies.policy_id
WHERE role_id = $1
AND
	CASE WHEN @where_like_name::boolean = true THEN m_policies.name LIKE '%' || @search_name::text || '%' ELSE TRUE END;

-- name: GetRolesOnPolicy :many
SELECT sqlc.embed(m_role_associations), sqlc.embed(m_roles) FROM m_role_associations
LEFT JOIN m_roles ON m_role_associations.role_id = m_roles.role_id
WHERE policy_id = $1
AND
	CASE WHEN @where_like_name::boolean = true THEN m_roles.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_roles.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_roles.name END DESC,
	m_role_associations_pkey DESC;

-- name: GetRolesOnPolicyUseNumberedPaginate :many
SELECT sqlc.embed(m_role_associations), sqlc.embed(m_roles) FROM m_role_associations
LEFT JOIN m_roles ON m_role_associations.role_id = m_roles.role_id
WHERE policy_id = $1
AND
	CASE WHEN @where_like_name::boolean = true THEN m_roles.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_roles.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_roles.name END DESC,
	m_role_associations_pkey DESC
LIMIT $2 OFFSET $3;

-- name: GetRolesOnPolicyUseKeysetPaginate :many
SELECT sqlc.embed(m_role_associations), sqlc.embed(m_roles) FROM m_role_associations
LEFT JOIN m_roles ON m_role_associations.role_id = m_roles.role_id
WHERE policy_id = $1
AND
	CASE WHEN @where_like_name::boolean = true THEN m_roles.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE @cursor_direction
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'name' THEN m_roles.name > @cursor_column OR (m_roles.name = @cursor_column AND m_role_associations_pkey < @cursor)
				WHEN 'r_name' THEN m_roles.name < @cursor_column OR (m_roles.name = @cursor_column AND m_role_associations_pkey < @cursor)
				ELSE m_role_associations_pkey < @cursor
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'name' THEN m_roles.name < @cursor_column OR (m_roles.name = @cursor_column AND m_role_associations_pkey > @cursor)
				WHEN 'r_name' THEN m_roles.name > @cursor_column OR (m_roles.name = @cursor_column AND m_role_associations_pkey > @cursor)
				ELSE m_role_associations_pkey > @cursor
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_roles.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_roles.name END DESC,
	m_role_associations_pkey DESC
LIMIT $2;

-- name: CountRolesByPolicyID :one
SELECT COUNT(*) FROM m_role_associations
LEFT JOIN m_roles ON m_role_associations.role_id = m_roles.role_id
WHERE policy_id = $1
AND
	CASE WHEN @where_like_name::boolean = true THEN m_roles.name LIKE '%' || @search_name::text || '%' ELSE TRUE END;
