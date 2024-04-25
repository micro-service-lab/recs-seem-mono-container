-- name: CreateGroups :copyfrom
INSERT INTO m_groups (key, organization_id) VALUES ($1, $2);

-- name: CreateGroup :one
INSERT INTO m_groups (key, organization_id) VALUES ($1, $2) RETURNING *;

-- name: DeleteGroup :exec
DELETE FROM m_groups WHERE group_id = $1;

-- name: DeleteGroupByKey :exec
DELETE FROM m_groups WHERE key = $1;

-- name: FindGroupByID :one
SELECT * FROM m_groups WHERE group_id = $1;

-- name: FindGroupByIDWithOrganization :one
SELECT sqlc.embed(m_groups), sqlc.embed(m_organizations) FROM m_groups
LEFT JOIN m_organizations ON m_groups.organization_id = m_organizations.organization_id
WHERE group_id = $1;

-- name: FindGroupByKey :one
SELECT * FROM m_groups WHERE key = $1;

-- name: FindGroupByKeyWithOrganization :one
SELECT sqlc.embed(m_groups), sqlc.embed(m_organizations) FROM m_groups
LEFT JOIN m_organizations ON m_groups.organization_id = m_organizations.organization_id
WHERE key = $1;

-- name: GetGroups :many
SELECT * FROM m_groups
ORDER BY
	m_groups_pkey DESC
LIMIT $1 OFFSET $2;

-- name: GetGroupsWithOrganization :many
SELECT sqlc.embed(m_groups), sqlc.embed(m_organizations) FROM m_groups
LEFT JOIN m_organizations ON m_groups.organization_id = m_organizations.organization_id
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_organizations.name END ASC,
	m_groups_pkey DESC
LIMIT $1 OFFSET $2;

-- name: CountGroups :one
SELECT COUNT(*) FROM m_groups;
