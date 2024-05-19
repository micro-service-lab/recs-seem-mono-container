-- name: CreateGroups :copyfrom
INSERT INTO m_groups (key, organization_id) VALUES ($1, $2);

-- name: CreateGroup :one
INSERT INTO m_groups (key, organization_id) VALUES ($1, $2) RETURNING *;

-- name: DeleteGroup :execrows
DELETE FROM m_groups WHERE group_id = $1;

-- name: DeleteGroupByKey :execrows
DELETE FROM m_groups WHERE key = $1;

-- name: PluralDeleteGroups :execrows
DELETE FROM m_groups WHERE group_id = ANY(@group_ids::uuid[]);

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
	m_groups_pkey ASC;

-- name: GetGroupsUseNumberedPaginate :many
SELECT * FROM m_groups
ORDER BY
	m_groups_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetGroupsUseKeysetPaginate :many
SELECT * FROM m_groups
WHERE
	CASE @cursor_direction::text
		WHEN 'next' THEN
			m_groups_pkey > @cursor::int
		WHEN 'prev' THEN
			m_groups_pkey < @cursor::int
	END
ORDER BY
	CASE WHEN @cursor_direction::text = 'next' THEN m_groups_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN m_groups_pkey END DESC
LIMIT $1;

-- name: GetPluralGroups :many
SELECT * FROM m_groups
WHERE organization_id = ANY(@organization_ids::uuid[])
ORDER BY
	m_groups_pkey ASC;

-- name: GetPluralGroupsUseNumberedPaginate :many
SELECT * FROM m_groups
WHERE organization_id = ANY(@organization_ids::uuid[])
ORDER BY
	m_groups_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetGroupsWithOrganization :many
SELECT sqlc.embed(m_groups), sqlc.embed(m_organizations) FROM m_groups
LEFT JOIN m_organizations ON m_groups.organization_id = m_organizations.organization_id
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_organizations.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_organizations.name END DESC,
	m_groups_pkey ASC;

-- name: GetGroupsWithOrganizationUseNumberedPaginate :many
SELECT sqlc.embed(m_groups), sqlc.embed(m_organizations) FROM m_groups
LEFT JOIN m_organizations ON m_groups.organization_id = m_organizations.organization_id
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_organizations.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_organizations.name END DESC,
	m_groups_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetGroupsWithOrganizationUseKeysetPaginate :many
SELECT sqlc.embed(m_groups), sqlc.embed(m_organizations) FROM m_groups
LEFT JOIN m_organizations ON m_groups.organization_id = m_organizations.organization_id
WHERE
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'name' THEN name > @name_cursor OR (name = @name_cursor AND m_groups_pkey > @cursor::int)
				WHEN 'r_name' THEN name < @name_cursor OR (name = @name_cursor AND m_groups_pkey > @cursor::int)
				ELSE m_groups_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'name' THEN name > @name_cursor OR (name = @name_cursor AND m_groups_pkey > @cursor::int)
				WHEN 'r_name' THEN name < @name_cursor OR (name = @name_cursor AND m_groups_pkey > @cursor::int)
				ELSE m_groups_pkey > @cursor::int
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'next' THEN m_organizations.name END ASC,
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'prev' THEN m_organizations.name END DESC,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'next' THEN m_organizations.name END ASC,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'prev' THEN m_organizations.name END DESC,
	CASE WHEN @cursor_direction::text = 'next' THEN m_groups_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN m_groups_pkey END DESC
LIMIT $1;

-- name: GetPluralGroupsWithOrganization :many
SELECT sqlc.embed(m_groups), sqlc.embed(m_organizations) FROM m_groups
LEFT JOIN m_organizations ON m_groups.organization_id = m_organizations.organization_id
WHERE group_id = ANY(@group_ids::uuid[])
ORDER BY
	m_groups_pkey ASC;

-- name: GetPluralGroupsWithOrganizationUseNumberedPaginate :many
SELECT sqlc.embed(m_groups), sqlc.embed(m_organizations) FROM m_groups
LEFT JOIN m_organizations ON m_groups.organization_id = m_organizations.organization_id
WHERE group_id = ANY(@group_ids::uuid[])
ORDER BY
	m_groups_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountGroups :one
SELECT COUNT(*) FROM m_groups;
