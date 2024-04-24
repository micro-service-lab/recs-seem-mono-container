-- name: CreateGroups :copyfrom
INSERT INTO m_groups (key, organization_id) VALUES ($1, $2);

-- name: CreateGroup :one
INSERT INTO m_groups (key, organization_id) VALUES ($1, $2) RETURNING *;

-- name: DeleteGroup :exec
DELETE FROM m_groups WHERE group_id = $1;

-- name: FindGroupByID :one
SELECT * FROM m_groups WHERE group_id = $1;

-- name: FindGroupByIDWithOrganization :one
SELECT sqlc.embed(m_groups), sqlc.embed(m_organizations) FROM m_groups
INNER JOIN m_organizations ON m_groups.organization_id = m_organizations.organization_id
WHERE group_id = $1;

-- name: FindGroupByKey :one
SELECT * FROM m_groups WHERE key = $1;

-- name: FindGroupByKeyWithOrganization :one
SELECT sqlc.embed(m_groups), sqlc.embed(m_organizations) FROM m_groups
INNER JOIN m_organizations ON m_groups.organization_id = m_organizations.organization_id
WHERE key = $1;

-- name: GetGroupsByOrganizationID :many
SELECT * FROM m_groups WHERE organization_id = $1
ORDER BY
	m_groups_pkey DESC
LIMIT $2 OFFSET $3;

-- name: GetGroupsByOrganizationIDWithOrganization :many
SELECT sqlc.embed(m_groups), sqlc.embed(m_organizations) FROM m_groups
INNER JOIN m_organizations ON m_groups.organization_id = m_organizations.organization_id
WHERE m_groups.organization_id = $1
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_groups.name END ASC,
	m_groups_pkey DESC
LIMIT $2 OFFSET $3;

-- name: CountGroupsByOrganizationID :one
SELECT COUNT(*) FROM m_groups WHERE organization_id = $1;
