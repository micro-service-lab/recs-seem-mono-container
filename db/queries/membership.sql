-- name: CreateMemberships :copyfrom
INSERT INTO m_memberships (member_id, organization_id, work_position_id, added_at) VALUES ($1, $2, $3, $4);

-- name: CreateMembership :one
INSERT INTO m_memberships (member_id, organization_id, work_position_id, added_at) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: DeleteMembership :execrows
DELETE FROM m_memberships WHERE member_id = $1 AND organization_id = $2;

-- name: DeleteMembershipsOnMember :execrows
DELETE FROM m_memberships WHERE member_id = $1;

-- name: DeleteMembershipsOnMembers :execrows
DELETE FROM m_memberships WHERE member_id = ANY(@member_ids::uuid[]);

-- name: GetMembersOnOrganization :many
SELECT m_memberships.*, m_members.name member_name, m_members.first_name member_first_name, m_members.last_name member_last_name, m_members.email member_email,
m_members.profile_image_id member_profile_image_id, t_images.height member_profile_image_height,
t_images.width member_profile_image_width, t_images.attachable_item_id member_profile_image_attachable_item_id,
t_attachable_items.owner_id member_profile_image_owner_id, t_attachable_items.from_outer member_profile_image_from_outer, t_attachable_items.alias member_profile_image_alias,
t_attachable_items.url member_profile_image_url, t_attachable_items.size member_profile_image_size, t_attachable_items.mime_type_id member_profile_image_mime_type_id FROM m_memberships
LEFT JOIN m_members ON m_memberships.member_id = m_members.member_id
LEFT JOIN t_images ON m_members.profile_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE organization_id = $1
AND CASE
	WHEN @where_like_name::boolean = true THEN m_members.name LIKE '%' || @search_name::text || '%'
END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_members.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_members.name END DESC,
	CASE WHEN @order_method::text = 'old_add' THEN m_memberships.added_at END ASC,
	CASE WHEN @order_method::text = 'late_add' THEN m_memberships.added_at END DESC,
	m_memberships_pkey ASC;

-- name: GetMembersOnOrganizationUseNumberedPaginate :many
SELECT m_memberships.*, m_members.name member_name, m_members.first_name member_first_name, m_members.last_name member_last_name, m_members.email member_email,
m_members.profile_image_id member_profile_image_id, t_images.height member_profile_image_height,
t_images.width member_profile_image_width, t_images.attachable_item_id member_profile_image_attachable_item_id,
t_attachable_items.owner_id member_profile_image_owner_id, t_attachable_items.from_outer member_profile_image_from_outer, t_attachable_items.alias member_profile_image_alias,
t_attachable_items.url member_profile_image_url, t_attachable_items.size member_profile_image_size, t_attachable_items.mime_type_id member_profile_image_mime_type_id FROM m_memberships
LEFT JOIN m_members ON m_memberships.member_id = m_members.member_id
LEFT JOIN t_images ON m_members.profile_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE organization_id = $1
AND CASE
	WHEN @where_like_name::boolean = true THEN m_members.name LIKE '%' || @search_name::text || '%'
END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_members.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_members.name END DESC,
	CASE WHEN @order_method::text = 'old_add' THEN m_memberships.added_at END ASC,
	CASE WHEN @order_method::text = 'late_add' THEN m_memberships.added_at END DESC,
	m_memberships_pkey ASC
LIMIT $2 OFFSET $3;

-- name: GetMembersOnOrganizationUseKeysetPaginate :many
SELECT m_memberships.*, m_members.name member_name, m_members.first_name member_first_name, m_members.last_name member_last_name, m_members.email member_email,
m_members.profile_image_id member_profile_image_id, t_images.height member_profile_image_height,
t_images.width member_profile_image_width, t_images.attachable_item_id member_profile_image_attachable_item_id,
t_attachable_items.owner_id member_profile_image_owner_id, t_attachable_items.from_outer member_profile_image_from_outer, t_attachable_items.alias member_profile_image_alias,
t_attachable_items.url member_profile_image_url, t_attachable_items.size member_profile_image_size, t_attachable_items.mime_type_id member_profile_image_mime_type_id FROM m_memberships
LEFT JOIN m_members ON m_memberships.member_id = m_members.member_id
LEFT JOIN t_images ON m_members.profile_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE organization_id = $1
AND CASE
	WHEN @where_like_name::boolean = true THEN m_members.name LIKE '%' || @search_name::text || '%'
END
AND CASE @cursor_direction::text
	WHEN 'next' THEN
		CASE @order_method::text
			WHEN 'name' THEN m_members.name > @name_cursor OR (m_members.name = @name_cursor AND m_memberships_pkey > @cursor::int)
			WHEN 'r_name' THEN m_members.name < @name_cursor OR (m_members.name = @name_cursor AND m_memberships_pkey > @cursor::int)
			WHEN 'old_add' THEN m_memberships.added_at > @added_at_cursor OR (m_memberships.added_at = @added_at_cursor AND m_memberships_pkey > @cursor::int)
			WHEN 'late_add' THEN m_memberships.added_at < @added_at_cursor OR (m_memberships.added_at = @added_at_cursor AND m_memberships_pkey > @cursor::int)
			ELSE m_memberships_pkey > @cursor::int
		END
	WHEN 'prev' THEN
		CASE @order_method::text
			WHEN 'name' THEN m_members.name < @name_cursor OR (m_members.name = @name_cursor AND m_memberships_pkey < @cursor::int)
			WHEN 'r_name' THEN m_members.name > @name_cursor OR (m_members.name = @name_cursor AND m_memberships_pkey < @cursor::int)
			WHEN 'old_add' THEN m_memberships.added_at < @added_at_cursor OR (m_memberships.added_at = @added_at_cursor AND m_memberships_pkey < @cursor::int)
			WHEN 'late_add' THEN m_memberships.added_at > @added_at_cursor OR (m_memberships.added_at = @added_at_cursor AND m_memberships_pkey < @cursor::int)
			ELSE m_memberships_pkey < @cursor::int
		END
END
ORDER BY
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'next' THEN m_members.name END ASC,
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'prev' THEN m_members.name END DESC,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'next' THEN m_members.name END DESC,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'prev' THEN m_members.name END ASC,
	CASE WHEN @order_method::text = 'old_add' AND @cursor_direction::text = 'next' THEN m_memberships.added_at END ASC,
	CASE WHEN @order_method::text = 'old_add' AND @cursor_direction::text = 'prev' THEN m_memberships.added_at END DESC,
	CASE WHEN @order_method::text = 'late_add' AND @cursor_direction::text = 'next' THEN m_memberships.added_at END DESC,
	CASE WHEN @order_method::text = 'late_add' AND @cursor_direction::text = 'prev' THEN m_memberships.added_at END ASC,
	CASE WHEN @cursor_direction::text = 'next' THEN m_memberships_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN m_memberships_pkey END DESC
LIMIT $2;

-- name: GetPluralMembersOnOrganization :many
SELECT m_memberships.*, m_members.name member_name, m_members.first_name member_first_name, m_members.last_name member_last_name, m_members.email member_email,
m_members.profile_image_id member_profile_image_id, t_images.height member_profile_image_height,
t_images.width member_profile_image_width, t_images.attachable_item_id member_profile_image_attachable_item_id,
t_attachable_items.owner_id member_profile_image_owner_id, t_attachable_items.from_outer member_profile_image_from_outer, t_attachable_items.alias member_profile_image_alias,
t_attachable_items.url member_profile_image_url, t_attachable_items.size member_profile_image_size, t_attachable_items.mime_type_id member_profile_image_mime_type_id FROM m_memberships
LEFT JOIN m_members ON m_memberships.member_id = m_members.member_id
LEFT JOIN t_images ON m_members.profile_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE organization_id = ANY(@organization_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_members.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_members.name END DESC,
	CASE WHEN @order_method::text = 'old_add' THEN m_memberships.added_at END ASC,
	CASE WHEN @order_method::text = 'late_add' THEN m_memberships.added_at END DESC,
	m_memberships_pkey ASC;

-- name: GetPluralMembersOnOrganizationUseNumberedPaginate :many
SELECT m_memberships.*, m_members.name member_name, m_members.first_name member_first_name, m_members.last_name member_last_name, m_members.email member_email,
m_members.profile_image_id member_profile_image_id, t_images.height member_profile_image_height,
t_images.width member_profile_image_width, t_images.attachable_item_id member_profile_image_attachable_item_id,
t_attachable_items.owner_id member_profile_image_owner_id, t_attachable_items.from_outer member_profile_image_from_outer, t_attachable_items.alias member_profile_image_alias,
t_attachable_items.url member_profile_image_url, t_attachable_items.size member_profile_image_size, t_attachable_items.mime_type_id member_profile_image_mime_type_id FROM m_memberships
LEFT JOIN m_members ON m_memberships.member_id = m_members.member_id
LEFT JOIN t_images ON m_members.profile_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE organization_id = ANY(@organization_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_members.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_members.name END DESC,
	CASE WHEN @order_method::text = 'old_add' THEN m_memberships.added_at END ASC,
	CASE WHEN @order_method::text = 'late_add' THEN m_memberships.added_at END DESC,
	m_memberships_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountMembersOnOrganization :one
SELECT COUNT(*) FROM m_memberships WHERE organization_id = $1
AND CASE WHEN @where_like_name::boolean = true THEN
		EXISTS (SELECT 1 FROM m_members WHERE m_memberships.member_id = m_members.member_id AND m_members.name LIKE '%' || @search_name::text || '%')
	ELSE TRUE END;

-- name: GetOrganizationsOnMember :many
SELECT m_memberships.*, m_organizations.name organization_name, m_organizations.description organization_description,
m_organizations.color organization_color, m_organizations.is_personal organization_is_personal, m_organizations.is_whole organization_is_whole,
m_organizations.chat_room_id organization_chat_room_id
FROM m_memberships
LEFT JOIN m_organizations ON m_memberships.organization_id = m_organizations.organization_id
WHERE member_id = $1
AND CASE
	WHEN @where_like_name::boolean = true THEN m_organization.name LIKE '%' || @search_name::text || '%'
END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_organizations.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_organizations.name END DESC,
	CASE WHEN @order_method::text = 'old_add' THEN m_memberships.added_at END ASC,
	CASE WHEN @order_method::text = 'late_add' THEN m_memberships.added_at END DESC,
	m_memberships_pkey ASC;

-- name: GetOrganizationsOnMemberUseNumberedPaginate :many
SELECT m_memberships.*, m_organizations.name organization_name, m_organizations.description organization_description,
m_organizations.color organization_color, m_organizations.is_personal organization_is_personal, m_organizations.is_whole organization_is_whole,
m_organizations.chat_room_id organization_chat_room_id
FROM m_memberships
LEFT JOIN m_organizations ON m_memberships.organization_id = m_organizations.organization_id
WHERE member_id = $1
AND CASE
	WHEN @where_like_name::boolean = true THEN m_organization.name LIKE '%' || @search_name::text || '%'
END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_organizations.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_organizations.name END DESC,
	CASE WHEN @order_method::text = 'old_add' THEN m_memberships.added_at END ASC,
	CASE WHEN @order_method::text = 'late_add' THEN m_memberships.added_at END DESC,
	m_memberships_pkey ASC
LIMIT $2 OFFSET $3;

-- name: GetOrganizationsOnMemberUseKeysetPaginate :many
SELECT m_memberships.*, m_organizations.name organization_name, m_organizations.description organization_description,
m_organizations.color organization_color, m_organizations.is_personal organization_is_personal, m_organizations.is_whole organization_is_whole,
m_organizations.chat_room_id organization_chat_room_id
FROM m_memberships
LEFT JOIN m_organizations ON m_memberships.organization_id = m_organizations.organization_id
WHERE member_id = $1
AND CASE
	WHEN @where_like_name::boolean = true THEN m_organization.name LIKE '%' || @search_name::text || '%'
END
AND CASE @cursor_direction::text
	WHEN 'next' THEN
		CASE @order_method::text
			WHEN 'name' THEN m_organizations.name > @name_cursor OR (m_organizations.name = @name_cursor AND m_memberships_pkey > @cursor::int)
			WHEN 'r_name' THEN m_organizations.name < @name_cursor OR (m_organizations.name = @name_cursor AND m_memberships_pkey > @cursor::int)
			WHEN 'old_add' THEN m_memberships.added_at > @add_cursor OR (m_memberships.added_at = @add_cursor AND m_memberships_pkey > @cursor::int)
			WHEN 'late_add' THEN m_memberships.added_at < @add_cursor OR (m_memberships.added_at = @add_cursor AND m_memberships_pkey > @cursor::int)
			ELSE m_memberships_pkey > @cursor::int
		END
	WHEN 'prev' THEN
		CASE @order_method::text
			WHEN 'name' THEN m_organizations.name < @name_cursor OR (m_organizations.name = @name_cursor AND m_memberships_pkey < @cursor::int)
			WHEN 'r_name' THEN m_organizations.name > @name_cursor OR (m_organizations.name = @name_cursor AND m_memberships_pkey < @cursor::int)
			WHEN 'old_add' THEN m_memberships.added_at < @add_cursor OR (m_memberships.added_at = @add_cursor AND m_memberships_pkey < @cursor::int)
			WHEN 'late_add' THEN m_memberships.added_at > @add_cursor OR (m_memberships.added_at = @add_cursor AND m_memberships_pkey < @cursor::int)
			ELSE m_memberships_pkey < @cursor::int
		END
END
ORDER BY
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'next' THEN m_organizations.name END ASC,
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'prev' THEN m_organizations.name END DESC,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'next' THEN m_organizations.name END DESC,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'prev' THEN m_organizations.name END ASC,
	CASE WHEN @order_method::text = 'old_add' AND @cursor_direction::text = 'next' THEN m_memberships.added_at END ASC,
	CASE WHEN @order_method::text = 'old_add' AND @cursor_direction::text = 'prev' THEN m_memberships.added_at END DESC,
	CASE WHEN @order_method::text = 'late_add' AND @cursor_direction::text = 'next' THEN m_memberships.added_at END DESC,
	CASE WHEN @order_method::text = 'late_add' AND @cursor_direction::text = 'prev' THEN m_memberships.added_at END ASC,
	CASE WHEN @cursor_direction::text = 'next' THEN m_memberships_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN m_memberships_pkey END DESC
LIMIT $2;

-- name: GetPluralOrganizationsOnMember :many
SELECT m_memberships.*, m_organizations.name organization_name, m_organizations.description organization_description,
m_organizations.color organization_color, m_organizations.is_personal organization_is_personal, m_organizations.is_whole organization_is_whole,
m_organizations.chat_room_id organization_chat_room_id
FROM m_memberships
LEFT JOIN m_organizations ON m_memberships.organization_id = m_organizations.organization_id
WHERE member_id = ANY(@member_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_organizations.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_organizations.name END DESC,
	CASE WHEN @order_method::text = 'old_add' THEN m_memberships.added_at END ASC,
	CASE WHEN @order_method::text = 'late_add' THEN m_memberships.added_at END DESC,
	m_memberships_pkey ASC;

-- name: GetPluralOrganizationsOnMemberUseNumberedPaginate :many
SELECT m_memberships.*, m_organizations.name organization_name, m_organizations.description organization_description,
m_organizations.color organization_color, m_organizations.is_personal organization_is_personal, m_organizations.is_whole organization_is_whole,
m_organizations.chat_room_id organization_chat_room_id
FROM m_memberships
LEFT JOIN m_organizations ON m_memberships.organization_id = m_organizations.organization_id
WHERE member_id = ANY(@member_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_organizations.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_organizations.name END DESC,
	CASE WHEN @order_method::text = 'old_add' THEN m_memberships.added_at END ASC,
	CASE WHEN @order_method::text = 'late_add' THEN m_memberships.added_at END DESC,
	m_memberships_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountOrganizationsOnMember :one
SELECT COUNT(*) FROM m_memberships WHERE member_id = $1
AND CASE WHEN @where_like_name::boolean = true THEN
		EXISTS (SELECT 1 FROM m_organizations WHERE m_memberships.organization_id = m_organizations.organization_id AND m_organizations.name LIKE '%' || @search_name::text || '%')
	ELSE TRUE END;
