-- name: CreateMembers :copyfrom
INSERT INTO m_members (login_id, password, email, name, first_name, last_name, attend_status_id, grade_id, group_id, profile_image_id, role_id, personal_organization_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14);

-- name: CreateMember :one
INSERT INTO m_members (login_id, password, email, name, first_name, last_name, attend_status_id, grade_id, group_id, profile_image_id, role_id, personal_organization_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING *;

-- name: DeleteMember :execrows
DELETE FROM m_members WHERE member_id = $1;

-- name: PluralDeleteMembers :execrows
DELETE FROM m_members WHERE member_id = ANY(@member_ids::uuid[]);

-- name: UpdateMember :one
UPDATE m_members SET email = $2, name = $3, first_name = $4, last_name = $5, profile_image_id = $6, updated_at = $7 WHERE member_id = $1 RETURNING *;

-- name: UpdateMemberRole :one
UPDATE m_members SET role_id = $2, updated_at = $3 WHERE member_id = $1 RETURNING *;

-- name: UpdateMemberPassword :one
UPDATE m_members SET password = $2, updated_at = $3 WHERE member_id = $1 RETURNING *;

-- name: UpdateMemberLoginID :one
UPDATE m_members SET login_id = $2, updated_at = $3 WHERE member_id = $1 RETURNING *;

-- name: UpdateMemberAttendStatus :one
UPDATE m_members SET attend_status_id = $2, updated_at = $3 WHERE member_id = $1 RETURNING *;

-- name: UpdateMemberGrade :one
UPDATE m_members SET grade_id = $2, updated_at = $3 WHERE member_id = $1 RETURNING *;

-- name: UpdateMemberGroup :one
UPDATE m_members SET group_id = $2, updated_at = $3 WHERE member_id = $1 RETURNING *;

-- name: FindMemberByID :one
SELECT * FROM m_members WHERE member_id = $1;

-- name: FindMemberByIDWithAttendStatus :one
SELECT m_members.*, m_attend_statuses.name attend_status_name, m_attend_statuses.key attend_status_key FROM m_members
LEFT JOIN m_attend_statuses ON m_members.attend_status_id = m_attend_statuses.attend_status_id
WHERE member_id = $1;

-- name: FindMemberWithProfileImage :one
SELECT m_members.*, t_images.height profile_image_height,
t_images.width profile_image_width, t_images.attachable_item_id profile_image_attachable_item_id,
t_attachable_items.owner_id profile_image_owner_id, t_attachable_items.from_outer profile_image_from_outer,
t_attachable_items.url profile_image_url, t_attachable_items.size profile_image_size, t_attachable_items.mime_type_id profile_image_mime_type_id
FROM m_members
LEFT JOIN t_images ON m_members.profile_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE member_id = $1;

-- name: FindMemberByIDWithCrew :one
SELECT m_members.*, m_grades.key grade_key, m_grades.organization_id grade_organization_id, grag.name grade_organization_name, grag.description grade_organization_description,
grag.color grade_organization_color, grag.is_personal grade_organization_is_personal,
grag.is_whole grade_organization_is_whole, grag.chat_room_id grade_organization_chat_room_id,
m_groups.key group_key, m_groups.organization_id group_organization_id, grog.name group_organization_name, grog.description group_organization_description,
grog.color group_organization_color, grog.is_personal group_organization_is_personal,
grog.is_whole group_organization_is_whole, grog.chat_room_id group_organization_chat_room_id FROM m_members
LEFT JOIN m_grades ON m_members.grade_id = m_grades.grade_id
LEFT JOIN m_organizations grag ON m_grades.organization_id = grag.organization_id
LEFT JOIN m_groups ON m_members.group_id = m_groups.group_id
LEFT JOIN m_organizations grog ON m_groups.organization_id = grog.organization_id
WHERE member_id = $1;

-- name: FindMemberByIDWithPersonalOrganization :one
SELECT m_members.*, m_organizations.name organization_name, m_organizations.description organization_description,
m_organizations.color organization_color, m_organizations.is_personal organization_is_personal,
m_organizations.is_whole organization_is_whole, m_organizations.chat_room_id organization_chat_room_id FROM m_members
LEFT JOIN m_organizations ON m_members.personal_organization_id = m_organizations.organization_id
WHERE member_id = $1;

-- name: FindMemberByIDWithRole :one
SELECT m_members.*, m_roles.name role_name, m_roles.description role_description FROM m_members
LEFT JOIN m_roles ON m_members.role_id = m_roles.role_id
WHERE member_id = $1;

-- name: FindMemberByIDWithProfileImage :one
SELECT m_members.*, t_images.height profile_image_height,
t_images.width profile_image_width, t_images.attachable_item_id profile_image_attachable_item_id,
t_attachable_items.owner_id profile_image_owner_id, t_attachable_items.from_outer profile_image_from_outer,
t_attachable_items.url profile_image_url, t_attachable_items.size profile_image_size, t_attachable_items.mime_type_id profile_image_mime_type_id
FROM m_members
LEFT JOIN t_images ON m_members.profile_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE member_id = $1;

-- name: FindMemberByIDWithDetail :one
SELECT m_members.*, m_students.student_id, m_professors.professor_id FROM m_members
LEFT JOIN m_students ON m_members.member_id = m_students.member_id
LEFT JOIN m_professors ON m_members.member_id = m_professor.member_id
WHERE m_members.member_id = $1;

-- name: FindMemberByLoginID :one
SELECT * FROM m_members WHERE login_id = $1;

-- name: GetMembers :many
SELECT * FROM m_members
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_members.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_has_policy::boolean = true THEN (SELECT COUNT(*) FROM m_role_associations WHERE role_id = m_members.role_id AND m_role_associations.policy_id = ANY(@has_policy_ids::uuid[])) > 0 ELSE TRUE END
AND
	CASE WHEN @when_in_attend_status::boolean = true THEN m_members.attend_status_id = ANY(@in_attend_status_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @when_in_grade::boolean = true THEN m_members.grade_id = ANY(@in_grade_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @when_in_group::boolean = true THEN m_members.group_id = ANY(@in_group_ids::uuid[]) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_members.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_members.name END DESC,
	m_members_pkey ASC;

-- name: GetMembersUseNumberedPaginate :many
SELECT * FROM m_members
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_members.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_has_policy::boolean = true THEN (SELECT COUNT(*) FROM m_role_associations WHERE role_id = m_members.role_id AND m_role_associations.policy_id = ANY(@has_policy_ids::uuid[])) > 0 ELSE TRUE END
AND
	CASE WHEN @when_in_attend_status::boolean = true THEN m_members.attend_status_id = ANY(@in_attend_status_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @when_in_grade::boolean = true THEN m_members.grade_id = ANY(@in_grade_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @when_in_group::boolean = true THEN m_members.group_id = ANY(@in_group_ids::uuid[]) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_members.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_members.name END DESC,
	m_members_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetMembersUseKeysetPaginate :many
SELECT * FROM m_members
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_members.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_has_policy::boolean = true THEN (SELECT COUNT(*) FROM m_role_associations WHERE role_id = m_members.role_id AND m_role_associations.policy_id = ANY(@has_policy_ids::uuid[])) > 0 ELSE TRUE END
AND
	CASE WHEN @when_in_attend_status::boolean = true THEN m_members.attend_status_id = ANY(@in_attend_status_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @when_in_grade::boolean = true THEN m_members.grade_id = ANY(@in_grade_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @when_in_group::boolean = true THEN m_members.group_id = ANY(@in_group_ids::uuid[]) ELSE TRUE END
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'name' THEN m_members.name > @name_cursor OR (m_members.name = @name_cursor AND m_members_pkey > @cursor::int)
				WHEN 'r_name' THEN m_members.name < @name_cursor OR (m_members.name = @name_cursor AND m_members_pkey > @cursor::int)
				ELSE m_members_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'name' THEN m_members.name < @name_cursor OR (m_members.name = @name_cursor AND m_members_pkey < @cursor::int)
				WHEN 'r_name' THEN m_members.name > @name_cursor OR (m_members.name = @name_cursor AND m_members_pkey < @cursor::int)
				ELSE m_members_pkey < @cursor::int
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'next' THEN m_members.name END ASC,
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'prev' THEN m_members.name END DESC,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'next' THEN m_members.name END DESC,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'prev' THEN m_members.name END ASC,
	CASE WHEN @cursor_direction::text = 'next' THEN m_members_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN m_members_pkey END DESC
LIMIT $1;

-- name: GetPluralMembers :many
SELECT * FROM m_members WHERE member_id = ANY(@member_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_members.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_members.name END DESC,
	m_members_pkey ASC;

-- name: GetPluralMembersUseNumberedPaginate :many
SELECT * FROM m_members WHERE member_id = ANY(@member_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_members.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_members.name END DESC,
	m_members_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetMembersWithAttendStatus :many
SELECT m_members.*, m_attend_statuses.name attend_status_name, m_attend_statuses.key attend_status_key FROM m_members
LEFT JOIN m_attend_statuses ON m_members.attend_status_id = m_attend_statuses.attend_status_id
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_members.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_has_policy::boolean = true THEN (SELECT COUNT(*) FROM m_role_associations WHERE role_id = m_members.role_id AND m_role_associations.policy_id = ANY(@has_policy_ids::uuid[])) > 0 ELSE TRUE END
AND
	CASE WHEN @when_in_attend_status::boolean = true THEN m_members.attend_status_id = ANY(@in_attend_status_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @when_in_grade::boolean = true THEN m_members.grade_id = ANY(@in_grade_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @when_in_group::boolean = true THEN m_members.group_id = ANY(@in_group_ids::uuid[]) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_members.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_members.name END DESC,
	m_members_pkey ASC;

-- name: GetMembersWithAttendStatusUseNumberedPaginate :many
SELECT m_members.*, m_attend_statuses.name attend_status_name, m_attend_statuses.key attend_status_key FROM m_members
LEFT JOIN m_attend_statuses ON m_members.attend_status_id = m_attend_statuses.attend_status_id
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_members.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_has_policy::boolean = true THEN (SELECT COUNT(*) FROM m_role_associations WHERE role_id = m_members.role_id AND m_role_associations.policy_id = ANY(@has_policy_ids::uuid[])) > 0 ELSE TRUE END
AND
	CASE WHEN @when_in_attend_status::boolean = true THEN m_members.attend_status_id = ANY(@in_attend_status_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @when_in_grade::boolean = true THEN m_members.grade_id = ANY(@in_grade_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @when_in_group::boolean = true THEN m_members.group_id = ANY(@in_group_ids::uuid[]) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_members.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_members.name END DESC,
	m_members_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetMembersWithAttendStatusUseKeysetPaginate :many
SELECT m_members.*, m_attend_statuses.name attend_status_name, m_attend_statuses.key attend_status_key FROM m_members
LEFT JOIN m_attend_statuses ON m_members.attend_status_id = m_attend_statuses.attend_status_id
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_members.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_has_policy::boolean = true THEN (SELECT COUNT(*) FROM m_role_associations WHERE role_id = m_members.role_id AND m_role_associations.policy_id = ANY(@has_policy_ids::uuid[])) > 0 ELSE TRUE END
AND
	CASE WHEN @when_in_attend_status::boolean = true THEN m_members.attend_status_id = ANY(@in_attend_status_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @when_in_grade::boolean = true THEN m_members.grade_id = ANY(@in_grade_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @when_in_group::boolean = true THEN m_members.group_id = ANY(@in_group_ids::uuid[]) ELSE TRUE END
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'name' THEN m_members.name > @name_cursor OR (m_members.name = @name_cursor AND m_members_pkey > @cursor::int)
				WHEN 'r_name' THEN m_members.name < @name_cursor OR (m_members.name = @name_cursor AND m_members_pkey > @cursor::int)
				ELSE m_members_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'name' THEN m_members.name < @name_cursor OR (m_members.name = @name_cursor AND m_members_pkey < @cursor::int)
				WHEN 'r_name' THEN m_members.name > @name_cursor OR (m_members.name = @name_cursor AND m_members_pkey < @cursor::int)
				ELSE m_members_pkey < @cursor::int
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'next' THEN m_members.name END ASC,
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'prev' THEN m_members.name END DESC,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'next' THEN m_members.name END DESC,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'prev' THEN m_members.name END ASC,
	CASE WHEN @cursor_direction::text = 'next' THEN m_members_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN m_members_pkey END DESC
LIMIT $1;

-- name: GetPluralMembersWithAttendStatus :many
SELECT m_members.*, m_attend_statuses.name attend_status_name, m_attend_statuses.key attend_status_key FROM m_members
LEFT JOIN m_attend_statuses ON m_members.attend_status_id = m_attend_statuses.attend_status_id
WHERE member_id = ANY(@member_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_members.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_members.name END DESC,
	m_members_pkey ASC;

-- name: GetPluralMembersWithAttendStatusUseNumberedPaginate :many
SELECT m_members.*, m_attend_statuses.name attend_status_name, m_attend_statuses.key attend_status_key FROM m_members
LEFT JOIN m_attend_statuses ON m_members.attend_status_id = m_attend_statuses.attend_status_id
WHERE member_id = ANY(@member_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_members.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_members.name END DESC,
	m_members_pkey ASC
LIMIT $1 OFFSET $2;


-- name: GetMembersWithProfileImage :many
SELECT m_members.*, t_images.height profile_image_height,
t_images.width profile_image_width, t_images.attachable_item_id profile_image_attachable_item_id,
t_attachable_items.owner_id profile_image_owner_id, t_attachable_items.from_outer profile_image_from_outer,
t_attachable_items.url profile_image_url, t_attachable_items.size profile_image_size, t_attachable_items.mime_type_id profile_image_mime_type_id
FROM m_members
LEFT JOIN t_images ON m_members.profile_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_members.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_has_policy::boolean = true THEN (SELECT COUNT(*) FROM m_role_associations WHERE role_id = m_members.role_id AND m_role_associations.policy_id = ANY(@has_policy_ids::uuid[])) > 0 ELSE TRUE END
AND
	CASE WHEN @when_in_attend_status::boolean = true THEN m_members.attend_status_id = ANY(@in_attend_status_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @when_in_grade::boolean = true THEN m_members.grade_id = ANY(@in_grade_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @when_in_group::boolean = true THEN m_members.group_id = ANY(@in_group_ids::uuid[]) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_members.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_members.name END DESC,
	m_members_pkey ASC;

-- name: GetMembersWithProfileImageUseNumberedPaginate :many
SELECT m_members.*, t_images.height profile_image_height,
t_images.width profile_image_width, t_images.attachable_item_id profile_image_attachable_item_id,
t_attachable_items.owner_id profile_image_owner_id, t_attachable_items.from_outer profile_image_from_outer,
t_attachable_items.url profile_image_url, t_attachable_items.size profile_image_size, t_attachable_items.mime_type_id profile_image_mime_type_id
FROM m_members
LEFT JOIN t_images ON m_members.profile_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_members.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_has_policy::boolean = true THEN (SELECT COUNT(*) FROM m_role_associations WHERE role_id = m_members.role_id AND m_role_associations.policy_id = ANY(@has_policy_ids::uuid[])) > 0 ELSE TRUE END
AND
	CASE WHEN @when_in_attend_status::boolean = true THEN m_members.attend_status_id = ANY(@in_attend_status_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @when_in_grade::boolean = true THEN m_members.grade_id = ANY(@in_grade_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @when_in_group::boolean = true THEN m_members.group_id = ANY(@in_group_ids::uuid[]) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_members.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_members.name END DESC,
	m_members_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetMembersWithProfileImageUseKeysetPaginate :many
SELECT m_members.*, t_images.height profile_image_height,
t_images.width profile_image_width, t_images.attachable_item_id profile_image_attachable_item_id,
t_attachable_items.owner_id profile_image_owner_id, t_attachable_items.from_outer profile_image_from_outer,
t_attachable_items.url profile_image_url, t_attachable_items.size profile_image_size, t_attachable_items.mime_type_id profile_image_mime_type_id
FROM m_members
LEFT JOIN t_images ON m_members.profile_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_members.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_has_policy::boolean = true THEN (SELECT COUNT(*) FROM m_role_associations WHERE role_id = m_members.role_id AND m_role_associations.policy_id = ANY(@has_policy_ids::uuid[])) > 0 ELSE TRUE END
AND
	CASE WHEN @when_in_attend_status::boolean = true THEN m_members.attend_status_id = ANY(@in_attend_status_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @when_in_grade::boolean = true THEN m_members.grade_id = ANY(@in_grade_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @when_in_group::boolean = true THEN m_members.group_id = ANY(@in_group_ids::uuid[]) ELSE TRUE END
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'name' THEN m_members.name > @name_cursor OR (m_members.name = @name_cursor AND m_members_pkey > @cursor::int)
				WHEN 'r_name' THEN m_members.name < @name_cursor OR (m_members.name = @name_cursor AND m_members_pkey > @cursor::int)
				ELSE m_members_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'name' THEN m_members.name < @name_cursor OR (m_members.name = @name_cursor AND m_members_pkey < @cursor::int)
				WHEN 'r_name' THEN m_members.name > @name_cursor OR (m_members.name = @name_cursor AND m_members_pkey < @cursor::int)
				ELSE m_members_pkey < @cursor::int
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'next' THEN m_members.name END ASC,
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'prev' THEN m_members.name END DESC,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'next' THEN m_members.name END DESC,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'prev' THEN m_members.name END ASC,
	CASE WHEN @cursor_direction::text = 'next' THEN m_members_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN m_members_pkey END DESC
LIMIT $1;

-- name: GetPluralMembersWithProfileImage :many
SELECT m_members.*, t_images.height profile_image_height,
t_images.width profile_image_width, t_images.attachable_item_id profile_image_attachable_item_id,
t_attachable_items.owner_id profile_image_owner_id, t_attachable_items.from_outer profile_image_from_outer,
t_attachable_items.url profile_image_url, t_attachable_items.size profile_image_size, t_attachable_items.mime_type_id profile_image_mime_type_id
FROM m_members
LEFT JOIN t_images ON m_members.profile_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE member_id = ANY(@member_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_members.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_members.name END DESC,
	m_members_pkey ASC;

-- name: GetPluralMembersWithProfileImageUseNumberedPaginate :many
SELECT m_members.*, t_images.height profile_image_height,
t_images.width profile_image_width, t_images.attachable_item_id profile_image_attachable_item_id,
t_attachable_items.owner_id profile_image_owner_id, t_attachable_items.from_outer profile_image_from_outer,
t_attachable_items.url profile_image_url, t_attachable_items.size profile_image_size, t_attachable_items.mime_type_id profile_image_mime_type_id
FROM m_members
LEFT JOIN t_images ON m_members.profile_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE member_id = ANY(@member_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_members.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_members.name END DESC,
	m_members_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetMembersWithCrew :many
SELECT m_members.*, m_grades.key grade_key, m_grades.organization_id grade_organization_id, grag.name grade_organization_name, grag.description grade_organization_description,
grag.color grade_organization_color, grag.is_personal grade_organization_is_personal,
grag.is_whole grade_organization_is_whole, grag.chat_room_id grade_organization_chat_room_id,
m_groups.key group_key, m_groups.organization_id group_organization_id, grog.name group_organization_name, grog.description group_organization_description,
grog.color group_organization_color, grog.is_personal group_organization_is_personal,
grog.is_whole group_organization_is_whole, grog.chat_room_id group_organization_chat_room_id FROM m_members
LEFT JOIN m_grades ON m_members.grade_id = m_grades.grade_id
LEFT JOIN m_organizations grag ON m_grades.organization_id = grag.organization_id
LEFT JOIN m_groups ON m_members.group_id = m_groups.group_id
LEFT JOIN m_organizations grog ON m_groups.organization_id = grog.organization_id
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_members.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_has_policy::boolean = true THEN (SELECT COUNT(*) FROM m_role_associations WHERE role_id = m_members.role_id AND m_role_associations.policy_id = ANY(@has_policy_ids::uuid[])) > 0 ELSE TRUE END
AND
	CASE WHEN @when_in_attend_status::boolean = true THEN m_members.attend_status_id = ANY(@in_attend_status_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @when_in_grade::boolean = true THEN m_members.grade_id = ANY(@in_grade_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @when_in_group::boolean = true THEN m_members.group_id = ANY(@in_group_ids::uuid[]) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_members.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_members.name END DESC,
	m_members_pkey ASC;

-- name: GetMembersWithCrewUseNumberedPaginate :many
SELECT m_members.*, m_grades.key grade_key, m_grades.organization_id grade_organization_id, grag.name grade_organization_name, grag.description grade_organization_description,
grag.color grade_organization_color, grag.is_personal grade_organization_is_personal,
grag.is_whole grade_organization_is_whole, grag.chat_room_id grade_organization_chat_room_id,
m_groups.key group_key, m_groups.organization_id group_organization_id, grog.name group_organization_name, grog.description group_organization_description,
grog.color group_organization_color, grog.is_personal group_organization_is_personal,
grog.is_whole group_organization_is_whole, grog.chat_room_id group_organization_chat_room_id FROM m_members
LEFT JOIN m_grades ON m_members.grade_id = m_grades.grade_id
LEFT JOIN m_organizations grag ON m_grades.organization_id = grag.organization_id
LEFT JOIN m_groups ON m_members.group_id = m_groups.group_id
LEFT JOIN m_organizations grog ON m_groups.organization_id = grog.organization_id
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_members.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_has_policy::boolean = true THEN (SELECT COUNT(*) FROM m_role_associations WHERE role_id = m_members.role_id AND m_role_associations.policy_id = ANY(@has_policy_ids::uuid[])) > 0 ELSE TRUE END
AND
	CASE WHEN @when_in_attend_status::boolean = true THEN m_members.attend_status_id = ANY(@in_attend_status_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @when_in_grade::boolean = true THEN m_members.grade_id = ANY(@in_grade_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @when_in_group::boolean = true THEN m_members.group_id = ANY(@in_group_ids::uuid[]) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_members.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_members.name END DESC,
	m_members_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetMembersWithCrewUseKeysetPaginate :many
SELECT m_members.*, m_grades.key grade_key, m_grades.organization_id grade_organization_id, grag.name grade_organization_name, grag.description grade_organization_description,
grag.color grade_organization_color, grag.is_personal grade_organization_is_personal,
grag.is_whole grade_organization_is_whole, grag.chat_room_id grade_organization_chat_room_id,
m_groups.key group_key, m_groups.organization_id group_organization_id, grog.name group_organization_name, grog.description group_organization_description,
grog.color group_organization_color, grog.is_personal group_organization_is_personal,
grog.is_whole group_organization_is_whole, grog.chat_room_id group_organization_chat_room_id FROM m_members
LEFT JOIN m_grades ON m_members.grade_id = m_grades.grade_id
LEFT JOIN m_organizations grag ON m_grades.organization_id = grag.organization_id
LEFT JOIN m_groups ON m_members.group_id = m_groups.group_id
LEFT JOIN m_organizations grog ON m_groups.organization_id = grog.organization_id
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_members.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_has_policy::boolean = true THEN (SELECT COUNT(*) FROM m_role_associations WHERE role_id = m_members.role_id AND m_role_associations.policy_id = ANY(@has_policy_ids::uuid[])) > 0 ELSE TRUE END
AND
	CASE WHEN @when_in_attend_status::boolean = true THEN m_members.attend_status_id = ANY(@in_attend_status_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @when_in_grade::boolean = true THEN m_members.grade_id = ANY(@in_grade_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @when_in_group::boolean = true THEN m_members.group_id = ANY(@in_group_ids::uuid[]) ELSE TRUE END
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'name' THEN m_members.name > @name_cursor OR (m_members.name = @name_cursor AND m_members_pkey > @cursor::int)
				WHEN 'r_name' THEN m_members.name < @name_cursor OR (m_members.name = @name_cursor AND m_members_pkey > @cursor::int)
				ELSE m_members_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'name' THEN m_members.name < @name_cursor OR (m_members.name = @name_cursor AND m_members_pkey < @cursor::int)
				WHEN 'r_name' THEN m_members.name > @name_cursor OR (m_members.name = @name_cursor AND m_members_pkey < @cursor::int)
				ELSE m_members_pkey < @cursor::int
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'next' THEN m_members.name END ASC,
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'prev' THEN m_members.name END DESC,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'next' THEN m_members.name END DESC,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'prev' THEN m_members.name END ASC,
	CASE WHEN @cursor_direction::text = 'next' THEN m_members_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN m_members_pkey END DESC
LIMIT $1;

-- name: GetPluralMembersWithCrew :many
SELECT m_members.*, m_grades.key grade_key, m_grades.organization_id grade_organization_id, grag.name grade_organization_name, grag.description grade_organization_description,
grag.color grade_organization_color, grag.is_personal grade_organization_is_personal,
grag.is_whole grade_organization_is_whole, grag.chat_room_id grade_organization_chat_room_id,
m_groups.key group_key, m_groups.organization_id group_organization_id, grog.name group_organization_name, grog.description group_organization_description,
grog.color group_organization_color, grog.is_personal group_organization_is_personal,
grog.is_whole group_organization_is_whole, grog.chat_room_id group_organization_chat_room_id FROM m_members
LEFT JOIN m_grades ON m_members.grade_id = m_grades.grade_id
LEFT JOIN m_organizations grag ON m_grades.organization_id = grag.organization_id
LEFT JOIN m_groups ON m_members.group_id = m_groups.group_id
LEFT JOIN m_organizations grog ON m_groups.organization_id = grog.organization_id
WHERE member_id = ANY(@member_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_members.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_members.name END DESC,
	m_members_pkey ASC;

-- name: GetPluralMembersWithCrewUseNumberedPaginate :many
SELECT m_members.*, m_grades.key grade_key, m_grades.organization_id grade_organization_id, grag.name grade_organization_name, grag.description grade_organization_description,
grag.color grade_organization_color, grag.is_personal grade_organization_is_personal,
grag.is_whole grade_organization_is_whole, grag.chat_room_id grade_organization_chat_room_id,
m_groups.key group_key, m_groups.organization_id group_organization_id, grog.name group_organization_name, grog.description group_organization_description,
grog.color group_organization_color, grog.is_personal group_organization_is_personal,
grog.is_whole group_organization_is_whole, grog.chat_room_id group_organization_chat_room_id FROM m_members
LEFT JOIN m_grades ON m_members.grade_id = m_grades.grade_id
LEFT JOIN m_organizations grag ON m_grades.organization_id = grag.organization_id
LEFT JOIN m_groups ON m_members.group_id = m_groups.group_id
LEFT JOIN m_organizations grog ON m_groups.organization_id = grog.organization_id
WHERE member_id = ANY(@member_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_members.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_members.name END DESC,
	m_members_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetMembersWithPersonalOrganization :many
SELECT m_members.*, m_organizations.name organization_name, m_organizations.description organization_description,
m_organizations.color organization_color, m_organizations.is_personal organization_is_personal,
m_organizations.is_whole organization_is_whole, m_organizations.chat_room_id organization_chat_room_id FROM m_members
LEFT JOIN m_organizations ON m_members.personal_organization_id = m_organizations.organization_id
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_members.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_has_policy::boolean = true THEN (SELECT COUNT(*) FROM m_role_associations WHERE role_id = m_members.role_id AND m_role_associations.policy_id = ANY(@has_policy_ids::uuid[])) > 0 ELSE TRUE END
AND
	CASE WHEN @when_in_attend_status::boolean = true THEN m_members.attend_status_id = ANY(@in_attend_status_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @when_in_grade::boolean = true THEN m_members.grade_id = ANY(@in_grade_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @when_in_group::boolean = true THEN m_members.group_id = ANY(@in_group_ids::uuid[]) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_members.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_members.name END DESC,
	m_members_pkey ASC;

-- name: GetMembersWithPersonalOrganizationUseNumberedPaginate :many
SELECT m_members.*, m_organizations.name organization_name, m_organizations.description organization_description,
m_organizations.color organization_color, m_organizations.is_personal organization_is_personal,
m_organizations.is_whole organization_is_whole, m_organizations.chat_room_id organization_chat_room_id FROM m_members
LEFT JOIN m_organizations ON m_members.personal_organization_id = m_organizations.organization_id
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_members.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_has_policy::boolean = true THEN (SELECT COUNT(*) FROM m_role_associations WHERE role_id = m_members.role_id AND m_role_associations.policy_id = ANY(@has_policy_ids::uuid[])) > 0 ELSE TRUE END
AND
	CASE WHEN @when_in_attend_status::boolean = true THEN m_members.attend_status_id = ANY(@in_attend_status_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @when_in_grade::boolean = true THEN m_members.grade_id = ANY(@in_grade_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @when_in_group::boolean = true THEN m_members.group_id = ANY(@in_group_ids::uuid[]) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_members.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_members.name END DESC,
	m_members_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetMembersWithPersonalOrganizationUseKeysetPaginate :many
SELECT m_members.*, m_organizations.name organization_name, m_organizations.description organization_description,
m_organizations.color organization_color, m_organizations.is_personal organization_is_personal,
m_organizations.is_whole organization_is_whole, m_organizations.chat_room_id organization_chat_room_id FROM m_members
LEFT JOIN m_organizations ON m_members.personal_organization_id = m_organizations.organization_id
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_members.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_has_policy::boolean = true THEN (SELECT COUNT(*) FROM m_role_associations WHERE role_id = m_members.role_id AND m_role_associations.policy_id = ANY(@has_policy_ids::uuid[])) > 0 ELSE TRUE END
AND
	CASE WHEN @when_in_attend_status::boolean = true THEN m_members.attend_status_id = ANY(@in_attend_status_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @when_in_grade::boolean = true THEN m_members.grade_id = ANY(@in_grade_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @when_in_group::boolean = true THEN m_members.group_id = ANY(@in_group_ids::uuid[]) ELSE TRUE END
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'name' THEN m_members.name > @name_cursor OR (m_members.name = @name_cursor AND m_members_pkey > @cursor::int)
				WHEN 'r_name' THEN m_members.name < @name_cursor OR (m_members.name = @name_cursor AND m_members_pkey > @cursor::int)
				ELSE m_members_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'name' THEN m_members.name < @name_cursor OR (m_members.name = @name_cursor AND m_members_pkey < @cursor::int)
				WHEN 'r_name' THEN m_members.name > @name_cursor OR (m_members.name = @name_cursor AND m_members_pkey < @cursor::int)
				ELSE m_members_pkey < @cursor::int
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'next' THEN m_members.name END ASC,
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'prev' THEN m_members.name END DESC,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'next' THEN m_members.name END DESC,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'prev' THEN m_members.name END ASC,
	CASE WHEN @cursor_direction::text = 'next' THEN m_members_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN m_members_pkey END DESC
LIMIT $1;

-- name: GetPluralMembersWithPersonalOrganization :many
SELECT m_members.*, m_organizations.name organization_name, m_organizations.description organization_description,
m_organizations.color organization_color, m_organizations.is_personal organization_is_personal,
m_organizations.is_whole organization_is_whole, m_organizations.chat_room_id organization_chat_room_id FROM m_members
LEFT JOIN m_organizations ON m_members.personal_organization_id = m_organizations.organization_id
WHERE member_id = ANY(@member_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_members.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_members.name END DESC,
	m_members_pkey ASC;

-- name: GetPluralMembersWithPersonalOrganizationUseNumberedPaginate :many
SELECT m_members.*, m_organizations.name organization_name, m_organizations.description organization_description,
m_organizations.color organization_color, m_organizations.is_personal organization_is_personal,
m_organizations.is_whole organization_is_whole, m_organizations.chat_room_id organization_chat_room_id FROM m_members
LEFT JOIN m_organizations ON m_members.personal_organization_id = m_organizations.organization_id
WHERE member_id = ANY(@member_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_members.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_members.name END DESC,
	m_members_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetMembersWithRole :many
SELECT m_members.*, m_roles.name role_name, m_roles.description role_description FROM m_members
LEFT JOIN m_roles ON m_members.role_id = m_roles.role_id
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_members.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_has_policy::boolean = true THEN (SELECT COUNT(*) FROM m_role_associations WHERE role_id = m_members.role_id AND m_role_associations.policy_id = ANY(@has_policy_ids::uuid[])) > 0 ELSE TRUE END
AND
	CASE WHEN @when_in_attend_status::boolean = true THEN m_members.attend_status_id = ANY(@in_attend_status_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @when_in_grade::boolean = true THEN m_members.grade_id = ANY(@in_grade_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @when_in_group::boolean = true THEN m_members.group_id = ANY(@in_group_ids::uuid[]) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_members.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_members.name END DESC,
	m_members_pkey ASC;

-- name: GetMembersWithRoleUseNumberedPaginate :many
SELECT m_members.*, m_roles.name role_name, m_roles.description role_description FROM m_members
LEFT JOIN m_roles ON m_members.role_id = m_roles.role_id
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_members.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_has_policy::boolean = true THEN (SELECT COUNT(*) FROM m_role_associations WHERE role_id = m_members.role_id AND m_role_associations.policy_id = ANY(@has_policy_ids::uuid[])) > 0 ELSE TRUE END
AND
	CASE WHEN @when_in_attend_status::boolean = true THEN m_members.attend_status_id = ANY(@in_attend_status_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @when_in_grade::boolean = true THEN m_members.grade_id = ANY(@in_grade_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @when_in_group::boolean = true THEN m_members.group_id = ANY(@in_group_ids::uuid[]) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_members.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_members.name END DESC,
	m_members_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetMembersWithRoleUseKeysetPaginate :many
SELECT m_members.*, m_roles.name role_name, m_roles.description role_description FROM m_members
LEFT JOIN m_roles ON m_members.role_id = m_roles.role_id
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_members.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_has_policy::boolean = true THEN (SELECT COUNT(*) FROM m_role_associations WHERE role_id = m_members.role_id AND m_role_associations.policy_id = ANY(@has_policy_ids::uuid[])) > 0 ELSE TRUE END
AND
	CASE WHEN @when_in_attend_status::boolean = true THEN m_members.attend_status_id = ANY(@in_attend_status_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @when_in_grade::boolean = true THEN m_members.grade_id = ANY(@in_grade_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @when_in_group::boolean = true THEN m_members.group_id = ANY(@in_group_ids::uuid[]) ELSE TRUE END
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'name' THEN m_members.name > @name_cursor OR (m_members.name = @name_cursor AND m_members_pkey > @cursor::int)
				WHEN 'r_name' THEN m_members.name < @name_cursor OR (m_members.name = @name_cursor AND m_members_pkey > @cursor::int)
				ELSE m_members_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'name' THEN m_members.name < @name_cursor OR (m_members.name = @name_cursor AND m_members_pkey < @cursor::int)
				WHEN 'r_name' THEN m_members.name > @name_cursor OR (m_members.name = @name_cursor AND m_members_pkey < @cursor::int)
				ELSE m_members_pkey < @cursor::int
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'next' THEN m_members.name END ASC,
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'prev' THEN m_members.name END DESC,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'next' THEN m_members.name END DESC,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'prev' THEN m_members.name END ASC,
	CASE WHEN @cursor_direction::text = 'next' THEN m_members_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN m_members_pkey END DESC
LIMIT $1;

-- name: GetPluralMembersWithRole :many
SELECT m_members.*, m_roles.name role_name, m_roles.description role_description FROM m_members
LEFT JOIN m_roles ON m_members.role_id = m_roles.role_id
WHERE member_id = ANY(@member_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_members.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_members.name END DESC,
	m_members_pkey ASC;

-- name: GetPluralMembersWithRoleUseNumberedPaginate :many
SELECT m_members.*, m_roles.name role_name, m_roles.description role_description FROM m_members
LEFT JOIN m_roles ON m_members.role_id = m_roles.role_id
WHERE member_id = ANY(@member_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_members.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_members.name END DESC,
	m_members_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetMembersWithDetail :many
SELECT m_members.*, m_students.student_id, m_professors.professor_id FROM m_members
LEFT JOIN m_students ON m_members.member_id = m_students.member_id
LEFT JOIN m_professors ON m_members.member_id = m_professors.member_id
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_members.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_has_policy::boolean = true THEN (SELECT COUNT(*) FROM m_role_associations WHERE role_id = m_members.role_id AND m_role_associations.policy_id = ANY(@has_policy_ids::uuid[])) > 0 ELSE TRUE END
AND
	CASE WHEN @when_in_attend_status::boolean = true THEN m_members.attend_status_id = ANY(@in_attend_status_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @when_in_grade::boolean = true THEN m_members.grade_id = ANY(@in_grade_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @when_in_group::boolean = true THEN m_members.group_id = ANY(@in_group_ids::uuid[]) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_members.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_members.name END DESC,
	m_members_pkey ASC;

-- name: GetMembersWithDetailUseNumberedPaginate :many
SELECT m_members.*, m_students.student_id, m_professors.professor_id FROM m_members
LEFT JOIN m_students ON m_members.member_id = m_students.member_id
LEFT JOIN m_professors ON m_members.member_id = m_professors.member_id
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_members.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_has_policy::boolean = true THEN (SELECT COUNT(*) FROM m_role_associations WHERE role_id = m_members.role_id AND m_role_associations.policy_id = ANY(@has_policy_ids::uuid[])) > 0 ELSE TRUE END
AND
	CASE WHEN @when_in_attend_status::boolean = true THEN m_members.attend_status_id = ANY(@in_attend_status_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @when_in_grade::boolean = true THEN m_members.grade_id = ANY(@in_grade_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @when_in_group::boolean = true THEN m_members.group_id = ANY(@in_group_ids::uuid[]) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_members.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_members.name END DESC,
	m_members_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetMembersWithDetailUseKeysetPaginate :many
SELECT m_members.*, m_students.student_id, m_professors.professor_id FROM m_members
LEFT JOIN m_students ON m_members.member_id = m_students.member_id
LEFT JOIN m_professors ON m_members.member_id = m_professors.member_id
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_members.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_has_policy::boolean = true THEN (SELECT COUNT(*) FROM m_role_associations WHERE role_id = m_members.role_id AND m_role_associations.policy_id = ANY(@has_policy_ids::uuid[])) > 0 ELSE TRUE END
AND
	CASE WHEN @when_in_attend_status::boolean = true THEN m_members.attend_status_id = ANY(@in_attend_status_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @when_in_grade::boolean = true THEN m_members.grade_id = ANY(@in_grade_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @when_in_group::boolean = true THEN m_members.group_id = ANY(@in_group_ids::uuid[]) ELSE TRUE END
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'name' THEN m_members.name > @name_cursor OR (m_members.name = @name_cursor AND m_members_pkey > @cursor::int)
				WHEN 'r_name' THEN m_members.name < @name_cursor OR (m_members.name = @name_cursor AND m_members_pkey > @cursor::int)
				ELSE m_members_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'name' THEN m_members.name < @name_cursor OR (m_members.name = @name_cursor AND m_members_pkey < @cursor::int)
				WHEN 'r_name' THEN m_members.name > @name_cursor OR (m_members.name = @name_cursor AND m_members_pkey < @cursor::int)
				ELSE m_members_pkey < @cursor::int
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'next' THEN m_members.name END ASC,
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'prev' THEN m_members.name END DESC,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'next' THEN m_members.name END DESC,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'prev' THEN m_members.name END ASC,
	CASE WHEN @cursor_direction::text = 'next' THEN m_members_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN m_members_pkey END DESC
LIMIT $1;

-- name: GetPluralMembersWithDetail :many
SELECT m_members.*, m_students.student_id, m_professors.professor_id FROM m_members
LEFT JOIN m_students ON m_members.member_id = m_students.member_id
LEFT JOIN m_professors ON m_members.member_id = m_professors.member_id
WHERE member_id = ANY(@member_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_members.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_members.name END DESC,
	m_members_pkey ASC;


-- name: GetPluralMembersWithDetailUseNumberedPaginate :many
SELECT m_members.*, m_students.student_id, m_professors.professor_id FROM m_members
LEFT JOIN m_students ON m_members.member_id = m_students.member_id
LEFT JOIN m_professors ON m_members.member_id = m_professors.member_id
WHERE member_id = ANY(@member_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_members.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_members.name END DESC,
	m_members_pkey ASC
LIMIT $1 OFFSET $2;


-- name: CountMembers :one
SELECT count(*) FROM m_members
WHERE
	CASE WHEN @where_like_name::boolean = true THEN name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_has_policy::boolean = true THEN (SELECT COUNT(*) FROM m_role_associations WHERE role_id = role_id AND m_role_associations.policy_id = ANY(@has_policy_ids::uuid[])) > 0 ELSE TRUE END
AND
	CASE WHEN @when_in_attend_status::boolean = true THEN attend_status_id = ANY(@in_attend_status_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @when_in_grade::boolean = true THEN grade_id = ANY(@in_grade_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @when_in_group::boolean = true THEN group_id = ANY(@in_group_ids::uuid[]) ELSE TRUE END;
