-- name: CreateOrganizations :copyfrom
INSERT INTO m_organizations (name, description, color, is_personal, is_whole, chat_room_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: CreateOrganization :one
INSERT INTO m_organizations (name, description, color, is_personal, is_whole, chat_room_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

-- name: UpdateOrganization :one
UPDATE m_organizations SET name = $2, color = $3, description = $4, updated_at = $5 WHERE organization_id = $1 RETURNING *;

-- name: DeleteOrganization :execrows
DELETE FROM m_organizations WHERE organization_id = $1;

-- name: PluralDeleteOrganizations :execrows
DELETE FROM m_organizations WHERE organization_id = ANY(@organization_ids::uuid[]);

-- name: FindOrganizationByID :one
SELECT * FROM m_organizations WHERE organization_id = $1;

-- name: FindOrganizationByIDWithDetail :one
SELECT m_organizations.*, m_groups.group_id, m_groups.key group_key, m_grades.key grade_key, m_grades.grade_id FROM m_organizations
LEFT JOIN m_groups ON m_organizations.organization_id = m_groups.organization_id
LEFT JOIN m_grades ON m_organizations.organization_id = m_grades.organization_id
WHERE m_organizations.organization_id = $1;

-- name: FindOrganizationByIDWithChatRoom :one
SELECT m_organizations.*, m_chat_rooms.name chat_room_name, m_chat_rooms.is_private chat_room_is_private,
m_chat_rooms.from_organization chat_room_from_organization, m_chat_rooms.owner_id chat_room_owner_id,
m_chat_rooms.cover_image_id chat_room_cover_image_id, t_images.height chat_room_cover_image_height,
t_images.width chat_room_cover_image_width, t_images.attachable_item_id chat_room_cover_image_attachable_item_id,
t_attachable_items.owner_id chat_room_cover_image_owner_id, t_attachable_items.from_outer chat_room_cover_image_from_outer,
t_attachable_items.url chat_room_cover_image_url, t_attachable_items.size chat_room_cover_image_size, t_attachable_items.mime_type_id chat_room_cover_image_mime_type_id
FROM m_organizations
LEFT JOIN m_chat_rooms ON m_organizations.chat_room_id = m_chat_rooms.chat_room_id
LEFT JOIN t_images ON m_chat_rooms.cover_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE m_organizations.organization_id = $1;

-- name: FindOrganizationByIDWithChatRoomAndDetail :one
SELECT m_organizations.*, m_groups.group_id, m_groups.key group_key, m_grades.key grade_key, m_grades.grade_id, m_chat_rooms.name chat_room_name, m_chat_rooms.is_private chat_room_is_private,
m_chat_rooms.from_organization chat_room_from_organization, m_chat_rooms.owner_id chat_room_owner_id,
m_chat_rooms.cover_image_id chat_room_cover_image_id, t_images.height chat_room_cover_image_height,
t_images.width chat_room_cover_image_width, t_images.attachable_item_id chat_room_cover_image_attachable_item_id,
t_attachable_items.owner_id chat_room_cover_image_owner_id, t_attachable_items.from_outer chat_room_cover_image_from_outer,
t_attachable_items.url chat_room_cover_image_url, t_attachable_items.size chat_room_cover_image_size, t_attachable_items.mime_type_id chat_room_cover_image_mime_type_id
FROM m_organizations
LEFT JOIN m_groups ON m_organizations.organization_id = m_groups.organization_id
LEFT JOIN m_grades ON m_organizations.organization_id = m_grades.organization_id
LEFT JOIN m_chat_rooms ON m_organizations.chat_room_id = m_chat_rooms.chat_room_id
LEFT JOIN t_images ON m_chat_rooms.cover_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE m_organizations.organization_id = $1;

-- name: FindWholeOrganization :one
SELECT * FROM m_organizations WHERE is_whole = true;

-- name: FindWholeOrganizationWithChatRoom :one
SELECT m_organizations.*, m_chat_rooms.name chat_room_name, m_chat_rooms.is_private chat_room_is_private,
m_chat_rooms.from_organization chat_room_from_organization, m_chat_rooms.owner_id chat_room_owner_id,
m_chat_rooms.cover_image_id chat_room_cover_image_id, t_images.height chat_room_cover_image_height,
t_images.width chat_room_cover_image_width, t_images.attachable_item_id chat_room_cover_image_attachable_item_id,
t_attachable_items.owner_id chat_room_cover_image_owner_id, t_attachable_items.from_outer chat_room_cover_image_from_outer,
t_attachable_items.url chat_room_cover_image_url, t_attachable_items.size chat_room_cover_image_size, t_attachable_items.mime_type_id chat_room_cover_image_mime_type_id
FROM m_organizations
LEFT JOIN m_chat_rooms ON m_organizations.chat_room_id = m_chat_rooms.chat_room_id
LEFT JOIN t_images ON m_chat_rooms.cover_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE m_organizations.is_whole = true;

-- name: FindPersonalOrganization :one
SELECT m_organizations.* FROM m_organizations
WHERE m_organizations.is_personal = true AND EXISTS (SELECT * FROM m_members WHERE m_members.personal_organization_id = m_organizations.organization_id AND m_members.member_id = $1);

-- name: GetOrganizations :many
SELECT * FROM m_organizations
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_organizations.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_is_whole::boolean = true THEN m_organizations.is_whole = @is_whole ELSE TRUE END
AND
	CASE WHEN @where_is_personal::boolean = true THEN m_organizations.is_personal = @is_personal AND EXISTS (SELECT * FROM m_members WHERE m_members.personal_organization_id = m_organizations.organization_id AND m_members.member_id = @personal_member_id::uuid) ELSE TRUE END
AND
	CASE WHEN @where_is_group::boolean = true THEN EXISTS (SELECT * FROM m_groups WHERE m_groups.organization_id = m_organizations.organization_id) ELSE TRUE END
AND
	CASE WHEN @where_is_grade::boolean = true THEN EXISTS (SELECT * FROM m_grades WHERE m_grades.organization_id = m_organizations.organization_id) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_organizations.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_organizations.name END DESC,
	m_organizations_pkey ASC;

-- name: GetOrganizationsUseNumberedPaginate :many
SELECT * FROM m_organizations
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_organizations.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_is_whole::boolean = true THEN m_organizations.is_whole = @is_whole ELSE TRUE END
AND
	CASE WHEN @where_is_personal::boolean = true THEN m_organizations.is_personal = @is_personal AND EXISTS (SELECT * FROM m_members WHERE m_members.personal_organization_id = m_organizations.organization_id AND m_members.member_id = @personal_member_id::uuid) ELSE TRUE END
AND
	CASE WHEN @where_is_group::boolean = true THEN EXISTS (SELECT * FROM m_groups WHERE m_groups.organization_id = m_organizations.organization_id) ELSE TRUE END
AND
	CASE WHEN @where_is_grade::boolean = true THEN EXISTS (SELECT * FROM m_grades WHERE m_grades.organization_id = m_organizations.organization_id) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_organizations.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_organizations.name END DESC,
	m_organizations_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetOrganizationsUseKeysetPaginate :many
SELECT * FROM m_organizations
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_organizations.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_is_whole::boolean = true THEN m_organizations.is_whole = @is_whole ELSE TRUE END
AND
	CASE WHEN @where_is_personal::boolean = true THEN m_organizations.is_personal = @is_personal AND EXISTS (SELECT * FROM m_members WHERE m_members.personal_organization_id = m_organizations.organization_id AND m_members.member_id = @personal_member_id::uuid) ELSE TRUE END
AND
	CASE WHEN @where_is_group::boolean = true THEN EXISTS (SELECT * FROM m_groups WHERE m_groups.organization_id = m_organizations.organization_id) ELSE TRUE END
AND
	CASE WHEN @where_is_grade::boolean = true THEN EXISTS (SELECT * FROM m_grades WHERE m_grades.organization_id = m_organizations.organization_id) ELSE TRUE END
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

-- name: GetPluralOrganizations :many
SELECT * FROM m_organizations WHERE organization_id = ANY(@organization_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_organizations.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_organizations.name END DESC,
	m_organizations_pkey ASC;

-- name: GetPluralOrganizationsUseNumberedPaginate :many
SELECT * FROM m_organizations WHERE organization_id = ANY(@organization_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_organizations.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_organizations.name END DESC,
	m_organizations_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetOrganizationsWithDetail :many
SELECT m_organizations.*, m_groups.group_id, m_groups.key group_key, m_grades.key grade_key, m_grades.grade_id FROM m_organizations
LEFT JOIN m_groups ON m_organizations.organization_id = m_groups.organization_id
LEFT JOIN m_grades ON m_organizations.organization_id = m_grades.organization_id
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_organizations.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_is_whole::boolean = true THEN m_organizations.is_whole = @is_whole ELSE TRUE END
AND
	CASE WHEN @where_is_personal::boolean = true THEN m_organizations.is_personal = @is_personal AND EXISTS (SELECT * FROM m_members WHERE m_members.personal_organization_id = m_organizations.organization_id AND m_members.member_id = @personal_member_id::uuid) ELSE TRUE END
AND
	CASE WHEN @where_is_group::boolean = true THEN EXISTS (SELECT * FROM m_groups WHERE m_groups.organization_id = m_organizations.organization_id) ELSE TRUE END
AND
	CASE WHEN @where_is_grade::boolean = true THEN EXISTS (SELECT * FROM m_grades WHERE m_grades.organization_id = m_organizations.organization_id) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_organizations.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_organizations.name END DESC,
	m_organizations_pkey ASC;

-- name: GetOrganizationsWithDetailUseNumberedPaginate :many
SELECT m_organizations.*, m_groups.group_id, m_groups.key group_key, m_grades.key grade_key, m_grades.grade_id FROM m_organizations
LEFT JOIN m_groups ON m_organizations.organization_id = m_groups.organization_id
LEFT JOIN m_grades ON m_organizations.organization_id = m_grades.organization_id
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_organizations.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_is_whole::boolean = true THEN m_organizations.is_whole = @is_whole ELSE TRUE END
AND
	CASE WHEN @where_is_personal::boolean = true THEN m_organizations.is_personal = @is_personal AND EXISTS (SELECT * FROM m_members WHERE m_members.personal_organization_id = m_organizations.organization_id AND m_members.member_id = @personal_member_id::uuid) ELSE TRUE END
AND
	CASE WHEN @where_is_group::boolean = true THEN EXISTS (SELECT * FROM m_groups WHERE m_groups.organization_id = m_organizations.organization_id) ELSE TRUE END
AND
	CASE WHEN @where_is_grade::boolean = true THEN EXISTS (SELECT * FROM m_grades WHERE m_grades.organization_id = m_organizations.organization_id) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_organizations.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_organizations.name END DESC,
	m_organizations_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetOrganizationsWithDetailUseKeysetPaginate :many
SELECT m_organizations.*, m_groups.group_id, m_groups.key group_key, m_grades.key grade_key, m_grades.grade_id FROM m_organizations
LEFT JOIN m_groups ON m_organizations.organization_id = m_groups.organization_id
LEFT JOIN m_grades ON m_organizations.organization_id = m_grades.organization_id
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_organizations.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_is_whole::boolean = true THEN m_organizations.is_whole = @is_whole ELSE TRUE END
AND
	CASE WHEN @where_is_personal::boolean = true THEN m_organizations.is_personal = @is_personal AND EXISTS (SELECT * FROM m_members WHERE m_members.personal_organization_id = m_organizations.organization_id AND m_members.member_id = @personal_member_id::uuid) ELSE TRUE END
AND
	CASE WHEN @where_is_group::boolean = true THEN EXISTS (SELECT * FROM m_groups WHERE m_groups.organization_id = m_organizations.organization_id) ELSE TRUE END
AND
	CASE WHEN @where_is_grade::boolean = true THEN EXISTS (SELECT * FROM m_grades WHERE m_grades.organization_id = m_organizations.organization_id) ELSE TRUE END
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
SELECT m_organizations.*, m_groups.group_id, m_groups.key group_key, m_grades.key grade_key, m_grades.grade_id FROM m_organizations
LEFT JOIN m_groups ON m_organizations.organization_id = m_groups.organization_id
LEFT JOIN m_grades ON m_organizations.organization_id = m_grades.organization_id
WHERE organization_id = ANY(@organization_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_organizations.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_organizations.name END DESC,
	m_organizations_pkey ASC;

-- name: GetPluralOrganizationsWithDetailUseNumberedPaginate :many
SELECT m_organizations.*, m_groups.group_id, m_groups.key group_key, m_grades.key grade_key, m_grades.grade_id FROM m_organizations
LEFT JOIN m_groups ON m_organizations.organization_id = m_groups.organization_id
LEFT JOIN m_grades ON m_organizations.organization_id = m_grades.organization_id
WHERE organization_id = ANY(@organization_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_organizations.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_organizations.name END DESC,
	m_organizations_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetOrganizationsWithChatRoom :many
SELECT m_organizations.*, m_chat_rooms.name chat_room_name, m_chat_rooms.is_private chat_room_is_private,
m_chat_rooms.from_organization chat_room_from_organization, m_chat_rooms.owner_id chat_room_owner_id,
m_chat_rooms.cover_image_id chat_room_cover_image_id, t_images.height chat_room_cover_image_height,
t_images.width chat_room_cover_image_width, t_images.attachable_item_id chat_room_cover_image_attachable_item_id,
t_attachable_items.owner_id chat_room_cover_image_owner_id, t_attachable_items.from_outer chat_room_cover_image_from_outer,
t_attachable_items.url chat_room_cover_image_url, t_attachable_items.size chat_room_cover_image_size, t_attachable_items.mime_type_id chat_room_cover_image_mime_type_id
FROM m_organizations
LEFT JOIN m_chat_rooms ON m_organizations.chat_room_id = m_chat_rooms.chat_room_id
LEFT JOIN t_images ON m_chat_rooms.cover_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_organizations.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_is_whole::boolean = true THEN m_organizations.is_whole = @is_whole ELSE TRUE END
AND
	CASE WHEN @where_is_personal::boolean = true THEN m_organizations.is_personal = @is_personal AND EXISTS (SELECT * FROM m_members WHERE m_members.personal_organization_id = m_organizations.organization_id AND m_members.member_id = @personal_member_id::uuid) ELSE TRUE END
AND
	CASE WHEN @where_is_group::boolean = true THEN EXISTS (SELECT * FROM m_groups WHERE m_groups.organization_id = m_organizations.organization_id) ELSE TRUE END
AND
	CASE WHEN @where_is_grade::boolean = true THEN EXISTS (SELECT * FROM m_grades WHERE m_grades.organization_id = m_organizations.organization_id) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_organizations.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_organizations.name END DESC,
	m_organizations_pkey ASC;

-- name: GetOrganizationsWithChatRoomUseNumberedPaginate :many
SELECT m_organizations.*, m_chat_rooms.name chat_room_name, m_chat_rooms.is_private chat_room_is_private,
m_chat_rooms.from_organization chat_room_from_organization, m_chat_rooms.owner_id chat_room_owner_id,
m_chat_rooms.cover_image_id chat_room_cover_image_id, t_images.height chat_room_cover_image_height,
t_images.width chat_room_cover_image_width, t_images.attachable_item_id chat_room_cover_image_attachable_item_id,
t_attachable_items.owner_id chat_room_cover_image_owner_id, t_attachable_items.from_outer chat_room_cover_image_from_outer,
t_attachable_items.url chat_room_cover_image_url, t_attachable_items.size chat_room_cover_image_size, t_attachable_items.mime_type_id chat_room_cover_image_mime_type_id
FROM m_organizations
LEFT JOIN m_chat_rooms ON m_organizations.chat_room_id = m_chat_rooms.chat_room_id
LEFT JOIN t_images ON m_chat_rooms.cover_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_organizations.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_is_whole::boolean = true THEN m_organizations.is_whole = @is_whole ELSE TRUE END
AND
	CASE WHEN @where_is_personal::boolean = true THEN m_organizations.is_personal = @is_personal AND EXISTS (SELECT * FROM m_members WHERE m_members.personal_organization_id = m_organizations.organization_id AND m_members.member_id = @personal_member_id::uuid) ELSE TRUE END
AND
	CASE WHEN @where_is_group::boolean = true THEN EXISTS (SELECT * FROM m_groups WHERE m_groups.organization_id = m_organizations.organization_id) ELSE TRUE END
AND
	CASE WHEN @where_is_grade::boolean = true THEN EXISTS (SELECT * FROM m_grades WHERE m_grades.organization_id = m_organizations.organization_id) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_organizations.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_organizations.name END DESC,
	m_organizations_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetOrganizationsWithChatRoomUseKeysetPaginate :many
SELECT m_organizations.*, m_chat_rooms.name chat_room_name, m_chat_rooms.is_private chat_room_is_private,
m_chat_rooms.from_organization chat_room_from_organization, m_chat_rooms.owner_id chat_room_owner_id,
m_chat_rooms.cover_image_id chat_room_cover_image_id, t_images.height chat_room_cover_image_height,
t_images.width chat_room_cover_image_width, t_images.attachable_item_id chat_room_cover_image_attachable_item_id,
t_attachable_items.owner_id chat_room_cover_image_owner_id, t_attachable_items.from_outer chat_room_cover_image_from_outer,
t_attachable_items.url chat_room_cover_image_url, t_attachable_items.size chat_room_cover_image_size, t_attachable_items.mime_type_id chat_room_cover_image_mime_type_id
FROM m_organizations
LEFT JOIN m_chat_rooms ON m_organizations.chat_room_id = m_chat_rooms.chat_room_id
LEFT JOIN t_images ON m_chat_rooms.cover_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_organizations.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_is_whole::boolean = true THEN m_organizations.is_whole = @is_whole ELSE TRUE END
AND
	CASE WHEN @where_is_personal::boolean = true THEN m_organizations.is_personal = @is_personal AND EXISTS (SELECT * FROM m_members WHERE m_members.personal_organization_id = m_organizations.organization_id AND m_members.member_id = @personal_member_id::uuid) ELSE TRUE END
AND
	CASE WHEN @where_is_group::boolean = true THEN EXISTS (SELECT * FROM m_groups WHERE m_groups.organization_id = m_organizations.organization_id) ELSE TRUE END
AND
	CASE WHEN @where_is_grade::boolean = true THEN EXISTS (SELECT * FROM m_grades WHERE m_grades.organization_id = m_organizations.organization_id) ELSE TRUE END
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
SELECT m_organizations.*, m_chat_rooms.name chat_room_name, m_chat_rooms.is_private chat_room_is_private,
m_chat_rooms.from_organization chat_room_from_organization, m_chat_rooms.owner_id chat_room_owner_id,
m_chat_rooms.cover_image_id chat_room_cover_image_id, t_images.height chat_room_cover_image_height,
t_images.width chat_room_cover_image_width, t_images.attachable_item_id chat_room_cover_image_attachable_item_id,
t_attachable_items.owner_id chat_room_cover_image_owner_id, t_attachable_items.from_outer chat_room_cover_image_from_outer,
t_attachable_items.url chat_room_cover_image_url, t_attachable_items.size chat_room_cover_image_size, t_attachable_items.mime_type_id chat_room_cover_image_mime_type_id
FROM m_organizations
LEFT JOIN m_chat_rooms ON m_organizations.chat_room_id = m_chat_rooms.chat_room_id
LEFT JOIN t_images ON m_chat_rooms.cover_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE organization_id = ANY(@organization_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_organizations.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_organizations.name END DESC,
	m_organizations_pkey ASC;

-- name: GetPluralOrganizationsWithChatRoomUseNumberedPaginate :many
SELECT m_organizations.*, m_chat_rooms.name chat_room_name, m_chat_rooms.is_private chat_room_is_private,
m_chat_rooms.from_organization chat_room_from_organization, m_chat_rooms.owner_id chat_room_owner_id,
m_chat_rooms.cover_image_id chat_room_cover_image_id, t_images.height chat_room_cover_image_height,
t_images.width chat_room_cover_image_width, t_images.attachable_item_id chat_room_cover_image_attachable_item_id,
t_attachable_items.owner_id chat_room_cover_image_owner_id, t_attachable_items.from_outer chat_room_cover_image_from_outer,
t_attachable_items.url chat_room_cover_image_url, t_attachable_items.size chat_room_cover_image_size, t_attachable_items.mime_type_id chat_room_cover_image_mime_type_id
FROM m_organizations
LEFT JOIN m_chat_rooms ON m_organizations.chat_room_id = m_chat_rooms.chat_room_id
LEFT JOIN t_images ON m_chat_rooms.cover_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE organization_id = ANY(@organization_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_organizations.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_organizations.name END DESC,
	m_organizations_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetOrganizationsWithChatRoomAndDetail :many
SELECT m_organizations.*, m_groups.group_id, m_groups.key group_key, m_grades.key grade_key, m_grades.grade_id, m_chat_rooms.name chat_room_name, m_chat_rooms.is_private chat_room_is_private,
m_chat_rooms.from_organization chat_room_from_organization, m_chat_rooms.owner_id chat_room_owner_id,
m_chat_rooms.cover_image_id chat_room_cover_image_id, t_images.height chat_room_cover_image_height,
t_images.width chat_room_cover_image_width, t_images.attachable_item_id chat_room_cover_image_attachable_item_id,
t_attachable_items.owner_id chat_room_cover_image_owner_id, t_attachable_items.from_outer chat_room_cover_image_from_outer,
t_attachable_items.url chat_room_cover_image_url, t_attachable_items.size chat_room_cover_image_size, t_attachable_items.mime_type_id chat_room_cover_image_mime_type_id
FROM m_organizations
LEFT JOIN m_groups ON m_organizations.organization_id = m_groups.organization_id
LEFT JOIN m_grades ON m_organizations.organization_id = m_grades.organization_id
LEFT JOIN m_chat_rooms ON m_organizations.chat_room_id = m_chat_rooms.chat_room_id
LEFT JOIN t_images ON m_chat_rooms.cover_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_organizations.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_is_whole::boolean = true THEN m_organizations.is_whole = @is_whole ELSE TRUE END
AND
	CASE WHEN @where_is_personal::boolean = true THEN m_organizations.is_personal = @is_personal AND EXISTS (SELECT * FROM m_members WHERE m_members.personal_organization_id = m_organizations.organization_id AND m_members.member_id = @personal_member_id::uuid) ELSE TRUE END
AND
	CASE WHEN @where_is_group::boolean = true THEN EXISTS (SELECT * FROM m_groups WHERE m_groups.organization_id = m_organizations.organization_id) ELSE TRUE END
AND
	CASE WHEN @where_is_grade::boolean = true THEN EXISTS (SELECT * FROM m_grades WHERE m_grades.organization_id = m_organizations.organization_id) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_organizations.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_organizations.name END DESC,
	m_organizations_pkey ASC;

-- name: GetOrganizationsWithChatRoomAndDetailUseNumberedPaginate :many
SELECT m_organizations.*, m_groups.group_id, m_groups.key group_key, m_grades.key grade_key, m_grades.grade_id, m_chat_rooms.name chat_room_name, m_chat_rooms.is_private chat_room_is_private,
m_chat_rooms.from_organization chat_room_from_organization, m_chat_rooms.owner_id chat_room_owner_id,
m_chat_rooms.cover_image_id chat_room_cover_image_id, t_images.height chat_room_cover_image_height,
t_images.width chat_room_cover_image_width, t_images.attachable_item_id chat_room_cover_image_attachable_item_id,
t_attachable_items.owner_id chat_room_cover_image_owner_id, t_attachable_items.from_outer chat_room_cover_image_from_outer,
t_attachable_items.url chat_room_cover_image_url, t_attachable_items.size chat_room_cover_image_size, t_attachable_items.mime_type_id chat_room_cover_image_mime_type_id
FROM m_organizations
LEFT JOIN m_groups ON m_organizations.organization_id = m_groups.organization_id
LEFT JOIN m_grades ON m_organizations.organization_id = m_grades.organization_id
LEFT JOIN m_chat_rooms ON m_organizations.chat_room_id = m_chat_rooms.chat_room_id
LEFT JOIN t_images ON m_chat_rooms.cover_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_organizations.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_is_whole::boolean = true THEN m_organizations.is_whole = @is_whole ELSE TRUE END
AND
	CASE WHEN @where_is_personal::boolean = true THEN m_organizations.is_personal = @is_personal AND EXISTS (SELECT * FROM m_members WHERE m_members.personal_organization_id = m_organizations.organization_id AND m_members.member_id = @personal_member_id::uuid) ELSE TRUE END
AND
	CASE WHEN @where_is_group::boolean = true THEN EXISTS (SELECT * FROM m_groups WHERE m_groups.organization_id = m_organizations.organization_id) ELSE TRUE END
AND
	CASE WHEN @where_is_grade::boolean = true THEN EXISTS (SELECT * FROM m_grades WHERE m_grades.organization_id = m_organizations.organization_id) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_organizations.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_organizations.name END DESC,
	m_organizations_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetOrganizationsWithChatRoomAndDetailUseKeysetPaginate :many
SELECT m_organizations.*, m_groups.group_id, m_groups.key group_key, m_grades.key grade_key, m_grades.grade_id, m_chat_rooms.name chat_room_name, m_chat_rooms.is_private chat_room_is_private,
m_chat_rooms.from_organization chat_room_from_organization, m_chat_rooms.owner_id chat_room_owner_id,
m_chat_rooms.cover_image_id chat_room_cover_image_id, t_images.height chat_room_cover_image_height,
t_images.width chat_room_cover_image_width, t_images.attachable_item_id chat_room_cover_image_attachable_item_id,
t_attachable_items.owner_id chat_room_cover_image_owner_id, t_attachable_items.from_outer chat_room_cover_image_from_outer,
t_attachable_items.url chat_room_cover_image_url, t_attachable_items.size chat_room_cover_image_size, t_attachable_items.mime_type_id chat_room_cover_image_mime_type_id
FROM m_organizations
LEFT JOIN m_groups ON m_organizations.organization_id = m_groups.organization_id
LEFT JOIN m_grades ON m_organizations.organization_id = m_grades.organization_id
LEFT JOIN m_chat_rooms ON m_organizations.chat_room_id = m_chat_rooms.chat_room_id
LEFT JOIN t_images ON m_chat_rooms.cover_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_organizations.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_is_whole::boolean = true THEN m_organizations.is_whole = @is_whole ELSE TRUE END
AND
	CASE WHEN @where_is_personal::boolean = true THEN m_organizations.is_personal = @is_personal AND EXISTS (SELECT * FROM m_members WHERE m_members.personal_organization_id = m_organizations.organization_id AND m_members.member_id = @personal_member_id::uuid) ELSE TRUE END
AND
	CASE WHEN @where_is_group::boolean = true THEN EXISTS (SELECT * FROM m_groups WHERE m_groups.organization_id = m_organizations.organization_id) ELSE TRUE END
AND
	CASE WHEN @where_is_grade::boolean = true THEN EXISTS (SELECT * FROM m_grades WHERE m_grades.organization_id = m_organizations.organization_id) ELSE TRUE END
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

-- name: GetPluralOrganizationsWithChatRoomAndDetail :many
SELECT m_organizations.*, m_groups.group_id, m_groups.key group_key, m_grades.key grade_key, m_grades.grade_id, m_chat_rooms.name chat_room_name, m_chat_rooms.is_private chat_room_is_private,
m_chat_rooms.from_organization chat_room_from_organization, m_chat_rooms.owner_id chat_room_owner_id,
m_chat_rooms.cover_image_id chat_room_cover_image_id, t_images.height chat_room_cover_image_height,
t_images.width chat_room_cover_image_width, t_images.attachable_item_id chat_room_cover_image_attachable_item_id,
t_attachable_items.owner_id chat_room_cover_image_owner_id, t_attachable_items.from_outer chat_room_cover_image_from_outer,
t_attachable_items.url chat_room_cover_image_url, t_attachable_items.size chat_room_cover_image_size, t_attachable_items.mime_type_id chat_room_cover_image_mime_type_id
FROM m_organizations
LEFT JOIN m_groups ON m_organizations.organization_id = m_groups.organization_id
LEFT JOIN m_grades ON m_organizations.organization_id = m_grades.organization_id
LEFT JOIN m_chat_rooms ON m_organizations.chat_room_id = m_chat_rooms.chat_room_id
LEFT JOIN t_images ON m_chat_rooms.cover_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE organization_id = ANY(@organization_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_organizations.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_organizations.name END DESC,
	m_organizations_pkey ASC;

-- name: GetPluralOrganizationsWithChatRoomAndDetailUseNumberedPaginate :many
SELECT m_organizations.*, m_groups.group_id, m_groups.key group_key, m_grades.key grade_key, m_grades.grade_id, m_chat_rooms.name chat_room_name, m_chat_rooms.is_private chat_room_is_private,
m_chat_rooms.from_organization chat_room_from_organization, m_chat_rooms.owner_id chat_room_owner_id,
m_chat_rooms.cover_image_id chat_room_cover_image_id, t_images.height chat_room_cover_image_height,
t_images.width chat_room_cover_image_width, t_images.attachable_item_id chat_room_cover_image_attachable_item_id,
t_attachable_items.owner_id chat_room_cover_image_owner_id, t_attachable_items.from_outer chat_room_cover_image_from_outer,
t_attachable_items.url chat_room_cover_image_url, t_attachable_items.size chat_room_cover_image_size, t_attachable_items.mime_type_id chat_room_cover_image_mime_type_id
FROM m_organizations
LEFT JOIN m_groups ON m_organizations.organization_id = m_groups.organization_id
LEFT JOIN m_grades ON m_organizations.organization_id = m_grades.organization_id
LEFT JOIN m_chat_rooms ON m_organizations.chat_room_id = m_chat_rooms.chat_room_id
LEFT JOIN t_images ON m_chat_rooms.cover_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE organization_id = ANY(@organization_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_organizations.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_organizations.name END DESC,
	m_organizations_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountOrganizations :one
SELECT COUNT(*) FROM m_organizations
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_organizations.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_is_whole::boolean = true THEN m_organizations.is_whole = @is_whole ELSE TRUE END
AND
	CASE WHEN @where_is_personal::boolean = true THEN m_organizations.is_personal = @is_personal AND EXISTS (SELECT * FROM m_members WHERE m_members.personal_organization_id = m_organizations.organization_id AND m_members.member_id = @personal_member_id::uuid) ELSE TRUE END
AND
	CASE WHEN @where_is_group::boolean = true THEN EXISTS (SELECT * FROM m_groups WHERE m_groups.organization_id = m_organizations.organization_id) ELSE TRUE END
AND
	CASE WHEN @where_is_grade::boolean = true THEN EXISTS (SELECT * FROM m_grades WHERE m_grades.organization_id = m_organizations.organization_id) ELSE TRUE END;
