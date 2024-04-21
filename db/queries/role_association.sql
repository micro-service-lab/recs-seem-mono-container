-- name: CreateRoleAssociations :copyfrom
INSERT INTO m_role_associations (role_id, policy_id) VALUES ($1, $2);

-- name: CreateRoleAssociation :one
INSERT INTO m_role_associations (role_id, policy_id) VALUES ($1, $2) RETURNING *;

-- name: DeleteRoleAssociation :exec
DELETE FROM m_role_associations WHERE role_id = $1 AND policy_id = $2;

-- name: GetPoliciesByRoleID :many
SELECT sqlc.embed(m_role_associations), sqlc.embed(m_policies) FROM m_role_associations
INNER JOIN m_policies ON m_role_associations.policy_id = m_policies.policy_id
WHERE role_id = $1
AND CASE
	WHEN @where_like_name::boolean = true THEN m_policies.name LIKE '%' || @search_name::text || '%'
END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_policies.name END ASC,
	m_role_associations_pkey DESC
LIMIT $2 OFFSET $3;

-- name: CountPoliciesByRoleID :one
SELECT COUNT(*) FROM m_role_associations WHERE role_id = $1;

-- name: GetRolesByPolicyID :many
SELECT sqlc.embed(m_role_associations), sqlc.embed(m_roles) FROM m_role_associations
INNER JOIN m_roles ON m_role_associations.role_id = m_roles.role_id
WHERE policy_id = $1
AND CASE
	WHEN @where_like_name::boolean = true THEN m_roles.name LIKE '%' || @search_name::text || '%'
END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_roles.name END ASC,
	m_role_associations_pkey DESC
LIMIT $2 OFFSET $3;

-- name: CountRolesByPolicyID :one
SELECT COUNT(*) FROM m_role_associations WHERE policy_id = $1;
