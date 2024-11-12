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
SELECT m_groups.*, m_organizations.name organization_name, m_organizations.description organization_description,
m_organizations.color organization_color, m_organizations.is_personal organization_is_personal,
m_organizations.is_whole organization_is_whole, m_organizations.chat_room_id organization_chat_room_id
FROM m_groups
LEFT JOIN m_organizations ON m_groups.organization_id = m_organizations.organization_id
WHERE group_id = $1;

-- name: FindGroupByKey :one
SELECT * FROM m_groups WHERE key = $1;

-- name: FindGroupByKeyWithOrganization :one
SELECT m_groups.*, m_organizations.name organization_name, m_organizations.description organization_description,
m_organizations.color organization_color, m_organizations.is_personal organization_is_personal,
m_organizations.is_whole organization_is_whole, m_organizations.chat_room_id organization_chat_room_id
FROM m_groups
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
WHERE group_id = ANY(@group_ids::uuid[])
ORDER BY
	m_groups_pkey ASC;

-- name: GetPluralGroupsUseNumberedPaginate :many
SELECT * FROM m_groups
WHERE group_id = ANY(@group_ids::uuid[])
ORDER BY
	m_groups_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetGroupsWithOrganization :many
SELECT m_groups.*, m_organizations.name organization_name, m_organizations.description organization_description,
m_organizations.color organization_color, m_organizations.is_personal organization_is_personal,
m_organizations.is_whole organization_is_whole, m_organizations.chat_room_id organization_chat_room_id
FROM m_groups
LEFT JOIN m_organizations ON m_groups.organization_id = m_organizations.organization_id
ORDER BY
	m_groups_pkey ASC;

-- name: GetGroupsWithOrganizationUseNumberedPaginate :many
SELECT m_groups.*, m_organizations.name organization_name, m_organizations.description organization_description,
m_organizations.color organization_color, m_organizations.is_personal organization_is_personal,
m_organizations.is_whole organization_is_whole, m_organizations.chat_room_id organization_chat_room_id
FROM m_groups
LEFT JOIN m_organizations ON m_groups.organization_id = m_organizations.organization_id
ORDER BY
	m_groups_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetGroupsWithOrganizationUseKeysetPaginate :many
SELECT m_groups.*, m_organizations.name organization_name, m_organizations.description organization_description,
m_organizations.color organization_color, m_organizations.is_personal organization_is_personal,
m_organizations.is_whole organization_is_whole, m_organizations.chat_room_id organization_chat_room_id
FROM m_groups
LEFT JOIN m_organizations ON m_groups.organization_id = m_organizations.organization_id
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

-- name: GetPluralGroupsWithOrganization :many
SELECT m_groups.*, m_organizations.name organization_name, m_organizations.description organization_description,
m_organizations.color organization_color, m_organizations.is_personal organization_is_personal,
m_organizations.is_whole organization_is_whole, m_organizations.chat_room_id organization_chat_room_id
FROM m_groups
LEFT JOIN m_organizations ON m_groups.organization_id = m_organizations.organization_id
WHERE group_id = ANY(@group_ids::uuid[])
ORDER BY
	m_groups_pkey ASC;

-- name: GetPluralGroupsWithOrganizationUseNumberedPaginate :many
SELECT m_groups.*, m_organizations.name organization_name, m_organizations.description organization_description,
m_organizations.color organization_color, m_organizations.is_personal organization_is_personal,
m_organizations.is_whole organization_is_whole, m_organizations.chat_room_id organization_chat_room_id
FROM m_groups
LEFT JOIN m_organizations ON m_groups.organization_id = m_organizations.organization_id
WHERE group_id = ANY(@group_ids::uuid[])
ORDER BY
	m_groups_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountGroups :one
SELECT COUNT(*) FROM m_groups;
