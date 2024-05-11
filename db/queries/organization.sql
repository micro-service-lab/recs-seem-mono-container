-- name: CreateOrganizations :copyfrom
INSERT INTO m_organizations (name, description, color, is_personal, is_whole, chat_room_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: CreateOrganization :one
INSERT INTO m_organizations (name, description, color, is_personal, is_whole, chat_room_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

-- name: UpdateOrganization :one
UPDATE m_organizations SET name = $2, color = $3, description = $4, updated_at = $5 WHERE organization_id = $1 RETURNING *;

-- name: DeleteOrganization :exec
DELETE FROM m_organizations WHERE organization_id = $1;

-- name: PluralDeleteOrganizations :exec
DELETE FROM m_organizations WHERE organization_id = ANY($1::uuid[]);

-- name: FindOrganizationByID :one
SELECT * FROM m_organizations WHERE organization_id = $1;

-- name: FindOrganizationByIDWithDetail :one
SELECT sqlc.embed(m_organizations), sqlc.embed(m_groups), sqlc.embed(m_grades) FROM m_organizations
LEFT JOIN m_groups ON m_organizations.organization_id = m_groups.organization_id
LEFT JOIN m_grades ON m_organizations.organization_id = m_grades.organization_id
WHERE m_organizations.organization_id = $1;

-- name: FindOrganizationByIDWithChatRoom :one
SELECT m_organizations.*, sqlc.embed(m_chat_rooms) FROM m_organizations
LEFT JOIN m_chat_rooms ON m_organizations.chat_room_id = m_chat_rooms.chat_room_id
WHERE m_organizations.organization_id = $1;

-- name: FindOrganizationByIDWithAll :one
SELECT sqlc.embed(m_organizations), sqlc.embed(m_groups), sqlc.embed(m_grades), sqlc.embed(m_chat_rooms) FROM m_organizations
LEFT JOIN m_groups ON m_organizations.organization_id = m_groups.organization_id
LEFT JOIN m_grades ON m_organizations.organization_id = m_grades.organization_id
LEFT JOIN m_chat_rooms ON m_organizations.chat_room_id = m_chat_rooms.chat_room_id
WHERE m_organizations.organization_id = $1;

-- name: FindWholeOrganization :one
SELECT * FROM m_organizations WHERE is_whole = true;

-- name: FindWholeOrganizationWithChatRoom :one
SELECT m_organizations.*, sqlc.embed(m_chat_rooms) FROM m_organizations
LEFT JOIN m_chat_rooms ON m_organizations.chat_room_id = m_chat_rooms.chat_room_id
WHERE m_organizations.is_whole = true;

-- name: FindPersonalOrganization :one
SELECT * FROM m_organizations
LEFT JOIN m_members ON m_organizations.organization_id = m_members.personal_organization_id
WHERE m_organizations.is_personal = true AND m_members.member_id = $1;

-- name: GetOrganizations :many
SELECT * FROM m_organizations
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_organizations.name LIKE '%' || @search_name::text || '%' END
AND
	CASE WHEN @where_is_whole::boolean = true THEN m_organizations.is_whole = @is_whole END
AND
	CASE WHEN @where_is_personal::boolean = true THEN m_organizations.is_personal = @is_personal AND EXISTS (SELECT * FROM m_members WHERE m_members.personal_organization_id = m_organizations.organization_id AND m_members.member_id = @personal_member_id::uuid) END
AND
	CASE WHEN @where_is_group::boolean = true THEN EXISTS (SELECT * FROM m_groups WHERE m_groups.organization_id = m_organizations.organization_id) END
AND
	CASE WHEN @where_is_grade::boolean = true THEN EXISTS (SELECT * FROM m_grades WHERE m_grades.organization_id = m_organizations.organization_id) END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_organizations.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_organizations.name END DESC,
	m_organizations_pkey ASC;

-- name: GetOrganizationsUseNumberedPaginate :many
SELECT * FROM m_organizations
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_organizations.name LIKE '%' || @search_name::text || '%' END
AND
	CASE WHEN @where_is_whole::boolean = true THEN m_organizations.is_whole = @is_whole END
AND
	CASE WHEN @where_is_personal::boolean = true THEN m_organizations.is_personal = @is_personal AND EXISTS (SELECT * FROM m_members WHERE m_members.personal_organization_id = m_organizations.organization_id AND m_members.member_id = @personal_member_id::uuid) END
AND
	CASE WHEN @where_is_group::boolean = true THEN EXISTS (SELECT * FROM m_groups WHERE m_groups.organization_id = m_organizations.organization_id) END
AND
	CASE WHEN @where_is_grade::boolean = true THEN EXISTS (SELECT * FROM m_grades WHERE m_grades.organization_id = m_organizations.organization_id) END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_organizations.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_organizations.name END DESC,
	m_organizations_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetOrganizationsUseKeysetPaginate :many
SELECT * FROM m_organizations
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_organizations.name LIKE '%' || @search_name::text || '%' END
AND
	CASE WHEN @where_is_whole::boolean = true THEN m_organizations.is_whole = @is_whole END
AND
	CASE WHEN @where_is_personal::boolean = true THEN m_organizations.is_personal = @is_personal AND EXISTS (SELECT * FROM m_members WHERE m_members.personal_organization_id = m_organizations.organization_id AND m_members.member_id = @personal_member_id::uuid) END
AND
	CASE WHEN @where_is_group::boolean = true THEN EXISTS (SELECT * FROM m_groups WHERE m_groups.organization_id = m_organizations.organization_id) END
AND
	CASE WHEN @where_is_grade::boolean = true THEN EXISTS (SELECT * FROM m_grades WHERE m_grades.organization_id = m_organizations.organization_id) END
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'name' THEN m_organizations.name > @name_cursor OR (m_organizations.name = @name_cursor AND m_organizations_pkey > @cursor::int)
				WHEN 'r_name' THEN m_organizations.name < @name_cursor OR (m_organizations.name = @name_cursor AND m_organizations_pkey > @cursor::int)
				ELSE m_organizations_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'name' THEN m_organizations.name < @name_cursor OR (m_organizations.name = @name_cursor AND m_organizations_pkey < @cursor::int)
				WHEN 'r_name' THEN m_organizations.name > @name_cursor OR (m_organizations.name = @name_cursor AND m_organizations_pkey < @cursor::int)
				ELSE m_organizations_pkey < @cursor::int
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'next' THEN m_organizations.name END ASC,
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'prev' THEN m_organizations.name END DESC,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'next' THEN m_organizations.name END DESC
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'prev' THEN m_organizations.name END ASC,
	CASE WHEN @cursor_direction::text = 'next' THEN m_organizations_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN m_organizations_pkey END DESC
LIMIT $1;

-- name: GetPluralOrganizations :many
SELECT * FROM m_organizations WHERE organization_id = ANY(@organization_ids::uuid[])
ORDER BY
	m_organizations_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetOrganizationsWithDetail :many
SELECT sqlc.embed(m_organizations), sqlc.embed(m_groups), sqlc.embed(m_grades) FROM m_organizations
LEFT JOIN m_groups ON m_organizations.organization_id = m_groups.organization_id
LEFT JOIN m_grades ON m_organizations.organization_id = m_grades.organization_id
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_organizations.name LIKE '%' || @search_name::text || '%' END
AND
	CASE WHEN @where_is_whole::boolean = true THEN m_organizations.is_whole = @is_whole END
AND
	CASE WHEN @where_is_personal::boolean = true THEN m_organizations.is_personal = @is_personal AND EXISTS (SELECT * FROM m_members WHERE m_members.personal_organization_id = m_organizations.organization_id AND m_members.member_id = @personal_member_id::uuid) END
AND
	CASE WHEN @where_is_group::boolean = true THEN EXISTS (SELECT * FROM m_groups WHERE m_groups.organization_id = m_organizations.organization_id) END
AND
	CASE WHEN @where_is_grade::boolean = true THEN EXISTS (SELECT * FROM m_grades WHERE m_grades.organization_id = m_organizations.organization_id) END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_organizations.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_organizations.name END DESC,
	m_organizations_pkey ASC;

-- name: GetOrganizationsWithDetailUseNumberedPaginate :many
SELECT sqlc.embed(m_organizations), sqlc.embed(m_groups), sqlc.embed(m_grades) FROM m_organizations
LEFT JOIN m_groups ON m_organizations.organization_id = m_groups.organization_id
LEFT JOIN m_grades ON m_organizations.organization_id = m_grades.organization_id
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_organizations.name LIKE '%' || @search_name::text || '%' END
AND
	CASE WHEN @where_is_whole::boolean = true THEN m_organizations.is_whole = @is_whole END
AND
	CASE WHEN @where_is_personal::boolean = true THEN m_organizations.is_personal = @is_personal AND EXISTS (SELECT * FROM m_members WHERE m_members.personal_organization_id = m_organizations.organization_id AND m_members.member_id = @personal_member_id::uuid) END
AND
	CASE WHEN @where_is_group::boolean = true THEN EXISTS (SELECT * FROM m_groups WHERE m_groups.organization_id = m_organizations.organization_id) END
AND
	CASE WHEN @where_is_grade::boolean = true THEN EXISTS (SELECT * FROM m_grades WHERE m_grades.organization_id = m_organizations.organization_id) END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_organizations.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_organizations.name END DESC,
	m_organizations_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetOrganizationsWithDetailUseKeysetPaginate :many
SELECT sqlc.embed(m_organizations), sqlc.embed(m_groups), sqlc.embed(m_grades) FROM m_organizations
LEFT JOIN m_groups ON m_organizations.organization_id = m_groups.organization_id
LEFT JOIN m_grades ON m_organizations.organization_id = m_grades.organization_id
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_organizations.name LIKE '%' || @search_name::text || '%' END
AND
	CASE WHEN @where_is_whole::boolean = true THEN m_organizations.is_whole = @is_whole END
AND
	CASE WHEN @where_is_personal::boolean = true THEN m_organizations.is_personal = @is_personal AND EXISTS (SELECT * FROM m_members WHERE m_members.personal_organization_id = m_organizations.organization_id AND m_members.member_id = @personal_member_id::uuid) END
AND
	CASE WHEN @where_is_group::boolean = true THEN EXISTS (SELECT * FROM m_groups WHERE m_groups.organization_id = m_organizations.organization_id) END
AND
	CASE WHEN @where_is_grade::boolean = true THEN EXISTS (SELECT * FROM m_grades WHERE m_grades.organization_id = m_organizations.organization_id) END
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'name' THEN m_organizations.name > @name_cursor OR (m_organizations.name = @name_cursor AND m_organizations_pkey > @cursor::int)
				WHEN 'r_name' THEN m_organizations.name < @name_cursor OR (m_organizations.name = @name_cursor AND m_organizations_pkey > @cursor::int)
				ELSE m_organizations_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'name' THEN m_organizations.name < @name_cursor OR (m_organizations.name = @name_cursor AND m_organizations_pkey < @cursor::int)
				WHEN 'r_name' THEN m_organizations.name > @name_cursor OR (m_organizations.name = @name_cursor AND m_organizations_pkey < @cursor::int)
				ELSE m_organizations_pkey < @cursor::int
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'next' THEN m_organizations.name END ASC,
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'prev' THEN m_organizations.name END DESC,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'next' THEN m_organizations.name END DESC,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'prev' THEN m_organizations.name END ASC,
	CASE WHEN @cursor_direction::text = 'next' THEN m_organizations_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN m_organizations_pkey END DESC
LIMIT $1;

-- name: GetPluralOrganizationsWithDetail :many
SELECT sqlc.embed(m_organizations), sqlc.embed(m_groups), sqlc.embed(m_grades) FROM m_organizations
LEFT JOIN m_groups ON m_organizations.organization_id = m_groups.organization_id
LEFT JOIN m_grades ON m_organizations.organization_id = m_grades.organization_id
WHERE organization_id = ANY(@organization_ids::uuid[])
ORDER BY
	m_organizations_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetOrganizationsWithChatRoom :many
SELECT m_organizations.*, sqlc.embed(m_chat_rooms) FROM m_organizations
LEFT JOIN m_chat_rooms ON m_organizations.chat_room_id = m_chat_rooms.chat_room_id
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_organizations.name LIKE '%' || @search_name::text || '%' END
AND
	CASE WHEN @where_is_whole::boolean = true THEN m_organizations.is_whole = @is_whole END
AND
	CASE WHEN @where_is_personal::boolean = true THEN m_organizations.is_personal = @is_personal AND EXISTS (SELECT * FROM m_members WHERE m_members.personal_organization_id = m_organizations.organization_id AND m_members.member_id = @personal_member_id::uuid) END
AND
	CASE WHEN @where_is_group::boolean = true THEN EXISTS (SELECT * FROM m_groups WHERE m_groups.organization_id = m_organizations.organization_id) END
AND
	CASE WHEN @where_is_grade::boolean = true THEN EXISTS (SELECT * FROM m_grades WHERE m_grades.organization_id = m_organizations.organization_id) END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_organizations.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_organizations.name END DESC,
	m_organizations_pkey ASC;

-- name: GetOrganizationsWithChatRoomUseNumberedPaginate :many
SELECT m_organizations.*, sqlc.embed(m_chat_rooms) FROM m_organizations
LEFT JOIN m_chat_rooms ON m_organizations.chat_room_id = m_chat_rooms.chat_room_id
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_organizations.name LIKE '%' || @search_name::text || '%' END
AND
	CASE WHEN @where_is_whole::boolean = true THEN m_organizations.is_whole = @is_whole END
AND
	CASE WHEN @where_is_personal::boolean = true THEN m_organizations.is_personal = @is_personal AND EXISTS (SELECT * FROM m_members WHERE m_members.personal_organization_id = m_organizations.organization_id AND m_members.member_id = @personal_member_id::uuid) END
AND
	CASE WHEN @where_is_group::boolean = true THEN EXISTS (SELECT * FROM m_groups WHERE m_groups.organization_id = m_organizations.organization_id) END
AND
	CASE WHEN @where_is_grade::boolean = true THEN EXISTS (SELECT * FROM m_grades WHERE m_grades.organization_id = m_organizations.organization_id) END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_organizations.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_organizations.name END DESC,
	m_organizations_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetOrganizationsWithChatRoomUseKeysetPaginate :many
SELECT m_organizations.*, sqlc.embed(m_chat_rooms) FROM m_organizations
LEFT JOIN m_chat_rooms ON m_organizations.chat_room_id = m_chat_rooms.chat_room_id
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_organizations.name LIKE '%' || @search_name::text || '%' END
AND
	CASE WHEN @where_is_whole::boolean = true THEN m_organizations.is_whole = @is_whole END
AND
	CASE WHEN @where_is_personal::boolean = true THEN m_organizations.is_personal = @is_personal AND EXISTS (SELECT * FROM m_members WHERE m_members.personal_organization_id = m_organizations.organization_id AND m_members.member_id = @personal_member_id::uuid) END
AND
	CASE WHEN @where_is_group::boolean = true THEN EXISTS (SELECT * FROM m_groups WHERE m_groups.organization_id = m_organizations.organization_id) END
AND
	CASE WHEN @where_is_grade::boolean = true THEN EXISTS (SELECT * FROM m_grades WHERE m_grades.organization_id = m_organizations.organization_id) END
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'name' THEN m_organizations.name > @name_cursor OR (m_organizations.name = @name_cursor AND m_organizations_pkey > @cursor::int)
				WHEN 'r_name' THEN m_organizations.name < @name_cursor OR (m_organizations.name = @name_cursor AND m_organizations_pkey > @cursor::int)
				ELSE m_organizations_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'name' THEN m_organizations.name < @name_cursor OR (m_organizations.name = @name_cursor AND m_organizations_pkey < @cursor::int)
				WHEN 'r_name' THEN m_organizations.name > @name_cursor OR (m_organizations.name = @name_cursor AND m_organizations_pkey < @cursor::int)
				ELSE m_organizations_pkey < @cursor::int
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'next' THEN m_organizations.name END ASC,
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'prev' THEN m_organizations.name END DESC,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'next' THEN m_organizations.name END DESC,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'prev' THEN m_organizations.name END ASC,
	CASE WHEN @cursor_direction::text = 'next' THEN m_organizations_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN m_organizations_pkey END DESC
LIMIT $1;

-- name: GetPluralOrganizationsWithChatRoom :many
SELECT m_organizations.*, sqlc.embed(m_chat_rooms) FROM m_organizations
LEFT JOIN m_chat_rooms ON m_organizations.chat_room_id = m_chat_rooms.chat_room_id
WHERE organization_id = ANY(@organization_ids::uuid[])
ORDER BY
	m_organizations_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetOrganizationsWithAll :many
SELECT sqlc.embed(m_organizations), sqlc.embed(m_groups), sqlc.embed(m_grades), sqlc.embed(m_chat_rooms) FROM m_organizations
LEFT JOIN m_groups ON m_organizations.organization_id = m_groups.organization_id
LEFT JOIN m_grades ON m_organizations.organization_id = m_grades.organization_id
LEFT JOIN m_chat_rooms ON m_organizations.chat_room_id = m_chat_rooms.chat_room_id
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_organizations.name LIKE '%' || @search_name::text || '%' END
AND
	CASE WHEN @where_is_whole::boolean = true THEN m_organizations.is_whole = @is_whole END
AND
	CASE WHEN @where_is_personal::boolean = true THEN m_organizations.is_personal = @is_personal AND EXISTS (SELECT * FROM m_members WHERE m_members.personal_organization_id = m_organizations.organization_id AND m_members.member_id = @personal_member_id::uuid) END
AND
	CASE WHEN @where_is_group::boolean = true THEN EXISTS (SELECT * FROM m_groups WHERE m_groups.organization_id = m_organizations.organization_id) END
AND
	CASE WHEN @where_is_grade::boolean = true THEN EXISTS (SELECT * FROM m_grades WHERE m_grades.organization_id = m_organizations.organization_id) END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_organizations.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_organizations.name END DESC,
	m_organizations_pkey ASC;

-- name: GetOrganizationsWithAllUseNumberedPaginate :many
SELECT sqlc.embed(m_organizations), sqlc.embed(m_groups), sqlc.embed(m_grades), sqlc.embed(m_chat_rooms) FROM m_organizations
LEFT JOIN m_groups ON m_organizations.organization_id = m_groups.organization_id
LEFT JOIN m_grades ON m_organizations.organization_id = m_grades.organization_id
LEFT JOIN m_chat_rooms ON m_organizations.chat_room_id = m_chat_rooms.chat_room_id
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_organizations.name LIKE '%' || @search_name::text || '%' END
AND
	CASE WHEN @where_is_whole::boolean = true THEN m_organizations.is_whole = @is_whole END
AND
	CASE WHEN @where_is_personal::boolean = true THEN m_organizations.is_personal = @is_personal AND EXISTS (SELECT * FROM m_members WHERE m_members.personal_organization_id = m_organizations.organization_id AND m_members.member_id = @personal_member_id::uuid) END
AND
	CASE WHEN @where_is_group::boolean = true THEN EXISTS (SELECT * FROM m_groups WHERE m_groups.organization_id = m_organizations.organization_id) END
AND
	CASE WHEN @where_is_grade::boolean = true THEN EXISTS (SELECT * FROM m_grades WHERE m_grades.organization_id = m_organizations.organization_id) END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_organizations.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_organizations.name END DESC,
	m_organizations_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetOrganizationsWithAllUseKeysetPaginate :many
SELECT sqlc.embed(m_organizations), sqlc.embed(m_groups), sqlc.embed(m_grades), sqlc.embed(m_chat_rooms) FROM m_organizations
LEFT JOIN m_groups ON m_organizations.organization_id = m_groups.organization_id
LEFT JOIN m_grades ON m_organizations.organization_id = m_grades.organization_id
LEFT JOIN m_chat_rooms ON m_organizations.chat_room_id = m_chat_rooms.chat_room_id
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_organizations.name LIKE '%' || @search_name::text || '%' END
AND
	CASE WHEN @where_is_whole::boolean = true THEN m_organizations.is_whole = @is_whole END
AND
	CASE WHEN @where_is_personal::boolean = true THEN m_organizations.is_personal = @is_personal AND EXISTS (SELECT * FROM m_members WHERE m_members.personal_organization_id = m_organizations.organization_id AND m_members.member_id = @personal_member_id::uuid) END
AND
	CASE WHEN @where_is_group::boolean = true THEN EXISTS (SELECT * FROM m_groups WHERE m_groups.organization_id = m_organizations.organization_id) END
AND
	CASE WHEN @where_is_grade::boolean = true THEN EXISTS (SELECT * FROM m_grades WHERE m_grades.organization_id = m_organizations.organization_id) END
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'name' THEN m_organizations.name > @name_cursor OR (m_organizations.name = @name_cursor AND m_organizations_pkey > @cursor::int)
				WHEN 'r_name' THEN m_organizations.name < @name_cursor OR (m_organizations.name = @name_cursor AND m_organizations_pkey > @cursor::int)
				ELSE m_organizations_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'name' THEN m_organizations.name < @name_cursor OR (m_organizations.name = @name_cursor AND m_organizations_pkey < @cursor::int)
				WHEN 'r_name' THEN m_organizations.name > @name_cursor OR (m_organizations.name = @name_cursor AND m_organizations_pkey < @cursor::int)
				ELSE m_organizations_pkey < @cursor::int
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'next' THEN m_organizations.name END ASC,
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'prev' THEN m_organizations.name END DESC,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'next' THEN m_organizations.name END DESC,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'prev' THEN m_organizations.name END ASC,
	CASE WHEN @cursor_direction::text = 'next' THEN m_organizations_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN m_organizations_pkey END DESC
LIMIT $1;

-- name: GetPluralOrganizationsWithAll :many
SELECT sqlc.embed(m_organizations), sqlc.embed(m_groups), sqlc.embed(m_grades), sqlc.embed(m_chat_rooms) FROM m_organizations
LEFT JOIN m_groups ON m_organizations.organization_id = m_groups.organization_id
LEFT JOIN m_grades ON m_organizations.organization_id = m_grades.organization_id
LEFT JOIN m_chat_rooms ON m_organizations.chat_room_id = m_chat_rooms.chat_room_id
WHERE organization_id = ANY(@organization_ids::uuid[])
ORDER BY
	m_organizations_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountOrganizations :one
SELECT COUNT(*) FROM m_organizations
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_organizations.name LIKE '%' || @search_name::text || '%' END
AND
	CASE WHEN @where_is_whole::boolean = true THEN m_organizations.is_whole = @is_whole END
AND
	CASE WHEN @where_is_personal::boolean = true THEN m_organizations.is_personal = @is_personal AND EXISTS (SELECT * FROM m_members WHERE m_members.personal_organization_id = m_organizations.organization_id AND m_members.member_id = @personal_member_id::uuid) END
AND
	CASE WHEN @where_is_group::boolean = true THEN EXISTS (SELECT * FROM m_groups WHERE m_groups.organization_id = m_organizations.organization_id) END
AND
	CASE WHEN @where_is_grade::boolean = true THEN EXISTS (SELECT * FROM m_grades WHERE m_grades.organization_id = m_organizations.organization_id) END;
