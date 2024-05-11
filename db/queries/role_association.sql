-- name: CreateRoleAssociations :copyfrom
INSERT INTO m_role_associations (role_id, policy_id) VALUES ($1, $2);

-- name: CreateRoleAssociation :one
INSERT INTO m_role_associations (role_id, policy_id) VALUES ($1, $2) RETURNING *;

-- name: DeleteRoleAssociation :exec
DELETE FROM m_role_associations WHERE role_id = $1 AND policy_id = $2;

-- name: DeleteRoleAssociationsOnRole :exec
DELETE FROM m_role_associations WHERE role_id = $1;

-- name: DeleteRoleAssociationsOnRoles :exec
DELETE FROM m_role_associations WHERE role_id = ANY($1::uuid[]);

-- name: PluralDeleteRoleAssociationsOnRole :exec
DELETE FROM m_role_associations WHERE role_id = $1 AND policy_id = ANY($2::uuid[]);

-- name: DeleteRoleAssociationsOnPolicy :exec
DELETE FROM m_role_associations WHERE policy_id = $1;

-- name: DeleteRoleAssociationsOnPolicies :exec
DELETE FROM m_role_associations WHERE policy_id = ANY($1::uuid[]);

-- name: PluralDeleteRoleAssociationsOnPolicy :exec
DELETE FROM m_role_associations WHERE policy_id = $1 AND role_id = ANY($2::uuid[]);

-- name: GetPoliciesOnRole :many
SELECT sqlc.embed(m_role_associations), sqlc.embed(m_policies) FROM m_role_associations
LEFT JOIN m_policies ON m_role_associations.policy_id = m_policies.policy_id
WHERE role_id = $1
AND
	CASE WHEN @where_like_name::boolean = true THEN m_policies.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_policies.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_policies.name END DESC,
	m_role_associations_pkey ASC;

-- name: GetPoliciesOnRoleUseNumberedPaginate :many
SELECT sqlc.embed(m_role_associations), sqlc.embed(m_policies) FROM m_role_associations
LEFT JOIN m_policies ON m_role_associations.policy_id = m_policies.policy_id
WHERE role_id = $1
AND
	CASE WHEN @where_like_name::boolean = true THEN m_policies.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_policies.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_policies.name END DESC,
	m_role_associations_pkey ASC
LIMIT $2 OFFSET $3;

-- name: GetPoliciesOnRoleUseKeysetPaginate :many
SELECT sqlc.embed(m_role_associations), sqlc.embed(m_policies) FROM m_role_associations
LEFT JOIN m_policies ON m_role_associations.policy_id = m_policies.policy_id
WHERE role_id = $1
AND
	CASE WHEN @where_like_name::boolean = true THEN m_policies.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'name' THEN m_policies.name > @name_cursor OR (m_policies.name = @name_cursor AND m_role_associations_pkey > @cursor::int)
				WHEN 'r_name' THEN m_policies.name < @name_cursor OR (m_policies.name = @name_cursor AND m_role_associations_pkey > @cursor::int)
				ELSE m_role_associations_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'name' THEN m_policies.name < @name_cursor OR (m_policies.name = @name_cursor AND m_role_associations_pkey < @cursor::int)
				WHEN 'r_name' THEN m_policies.name > @name_cursor OR (m_policies.name = @name_cursor AND m_role_associations_pkey < @cursor::int)
				ELSE m_role_associations_pkey < @cursor::int
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'next' THEN m_policies.name END ASC,
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'prev' THEN m_policies.name END DESC,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'next' THEN m_policies.name END DESC
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'prev' THEN m_policies.name END ASC,
	CASE WHEN @cursor_direction::text = 'next' THEN m_role_associations_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN m_role_associations_pkey END DESC
LIMIT $2;

-- name: GetPluralPoliciesOnRole :many
SELECT sqlc.embed(m_role_associations), sqlc.embed(m_policies) FROM m_role_associations
LEFT JOIN m_policies ON m_role_associations.policy_id = m_policies.policy_id
WHERE
	role_id = ANY(@role_ids::uuid[])
ORDER BY
	m_role_associations_pkey ASC
LIMIT $1 OFFSET $2;

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
	m_role_associations_pkey ASC;

-- name: GetRolesOnPolicyUseNumberedPaginate :many
SELECT sqlc.embed(m_role_associations), sqlc.embed(m_roles) FROM m_role_associations
LEFT JOIN m_roles ON m_role_associations.role_id = m_roles.role_id
WHERE policy_id = $1
AND
	CASE WHEN @where_like_name::boolean = true THEN m_roles.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_roles.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_roles.name END DESC,
	m_role_associations_pkey ASC
LIMIT $2 OFFSET $3;

-- name: GetRolesOnPolicyUseKeysetPaginate :many
SELECT sqlc.embed(m_role_associations), sqlc.embed(m_roles) FROM m_role_associations
LEFT JOIN m_roles ON m_role_associations.role_id = m_roles.role_id
WHERE policy_id = $1
AND
	CASE WHEN @where_like_name::boolean = true THEN m_roles.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'name' THEN m_roles.name > @name_cursor OR (m_roles.name = @name_cursor AND m_role_associations_pkey > @cursor::int)
				WHEN 'r_name' THEN m_roles.name < @name_cursor OR (m_roles.name = @name_cursor AND m_role_associations_pkey > @cursor::int)
				ELSE m_role_associations_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'name' THEN m_roles.name < @name_cursor OR (m_roles.name = @name_cursor AND m_role_associations_pkey < @cursor::int)
				WHEN 'r_name' THEN m_roles.name > @name_cursor OR (m_roles.name = @name_cursor AND m_role_associations_pkey < @cursor::int)
				ELSE m_role_associations_pkey < @cursor::int
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'next' THEN m_roles.name END ASC,
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'prev' THEN m_roles.name END DESC,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'next' THEN m_roles.name END DESC,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'prev' THEN m_roles.name END ASC,
	CASE WHEN @cursor_direction::text = 'next' THEN m_role_associations_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN m_role_associations_pkey END DESC
LIMIT $2;

-- name: GetPluralRolesOnPolicy :many
SELECT sqlc.embed(m_role_associations), sqlc.embed(m_roles) FROM m_role_associations
LEFT JOIN m_roles ON m_role_associations.role_id = m_roles.role_id
WHERE
	policy_id = ANY(@policy_ids::uuid[])
ORDER BY
	m_role_associations_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountRolesByPolicyID :one
SELECT COUNT(*) FROM m_role_associations
LEFT JOIN m_roles ON m_role_associations.role_id = m_roles.role_id
WHERE policy_id = $1
AND
	CASE WHEN @where_like_name::boolean = true THEN m_roles.name LIKE '%' || @search_name::text || '%' ELSE TRUE END;
