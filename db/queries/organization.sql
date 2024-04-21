-- name: CreateOrganizations :copyfrom
INSERT INTO m_organizations (name, description, is_personal, is_whole, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6);

-- name: CreateOrganization :one
INSERT INTO m_organizations (name, description, is_personal, is_whole, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: UpdateOrganization :one
UPDATE m_organizations SET name = $2, description = $3, updated_at = $4 WHERE organization_id = $1 RETURNING *;

-- name: DeleteOrganization :exec
DELETE FROM m_organizations WHERE organization_id = $1;

-- name: FindOrganizationByID :one
SELECT * FROM m_organizations WHERE organization_id = $1;

-- name: FindWholeOrganization :one
SELECT * FROM m_organizations WHERE is_whole = true;

-- name: GetOrganizations :many
SELECT * FROM m_organizations
WHERE CASE
	WHEN @where_like_name::boolean = true THEN m_organizations.name LIKE '%' || @search_name::text || '%'
END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_organizations.name END ASC,
	m_organizations_pkey DESC
LIMIT $1 OFFSET $2;

-- name: CountOrganizations :one
SELECT COUNT(*) FROM m_organizations;
