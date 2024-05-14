-- name: CreateMembers :copyfrom
INSERT INTO m_members (login_id, password, email, name, attend_status_id, grade_id, group_id, profile_image_id, role_id, personal_organization_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);

-- name: CreateMember :one
INSERT INTO m_members (login_id, password, email, name, attend_status_id, grade_id, group_id, profile_image_id, role_id, personal_organization_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING *;

-- name: DeleteMember :execrows
DELETE FROM m_members WHERE member_id = $1;

-- name: PluralDeleteMembers :execrows
DELETE FROM m_members WHERE member_id = ANY($1::uuid[]);

-- name: UpdateMember :one
UPDATE m_members SET email = $2, name = $3, profile_image_id = $4, updated_at = $5 WHERE member_id = $1 RETURNING *;

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
SELECT sqlc.embed(m_members), sqlc.embed(m_attend_statuses) FROM m_members
LEFT JOIN m_attend_statuses ON m_members.attend_status_id = m_attend_statuses.attend_status_id
WHERE member_id = $1;

-- name: FindMemberByIDWithGrade :one
SELECT sqlc.embed(m_members), sqlc.embed(m_grades) FROM m_members
LEFT JOIN m_grades ON m_members.grade_id = m_grades.grade_id
LEFT JOIN m_organizations ON m_grades.organization_id = m_organizations.organization_id
WHERE member_id = $1;

-- name: FindMemberByIDWithGroup :one
SELECT sqlc.embed(m_members), sqlc.embed(m_groups) FROM m_members
LEFT JOIN m_groups ON m_members.group_id = m_groups.group_id
LEFT JOIN m_organizations ON m_groups.organization_id = m_organizations.organization_id
WHERE member_id = $1;

-- name: FindMemberByIDWithPersonalOrganization :one
SELECT sqlc.embed(m_members), sqlc.embed(m_organizations) FROM m_members
LEFT JOIN m_organizations ON m_members.personal_organization_id = m_organizations.organization_id
WHERE member_id = $1;

-- name: FindMemberByIDWithRole :one
SELECT sqlc.embed(m_members), sqlc.embed(m_roles) FROM m_members
LEFT JOIN m_roles ON m_members.role_id = m_roles.role_id
WHERE member_id = $1;

-- name: FindMemberWithDetail :one
SELECT sqlc.embed(m_members), sqlc.embed(m_students), sqlc.embed(m_professors) FROM m_members
LEFT JOIN m_students ON m_members.member_id = m_students.member_id
LEFT JOIN m_professors ON m_members.member_id = m_professor.member_id
WHERE m_members.member_id = $1;

-- name: FindMemberWithAll :one
SELECT sqlc.embed(m_members), sqlc.embed(m_attend_statuses), sqlc.embed(m_grades), sqlc.embed(m_groups), sqlc.embed(m_organizations), sqlc.embed(m_roles) FROM m_members
LEFT JOIN m_attend_statuses ON m_members.attend_status_id = m_attend_statuses.attend_status_id
LEFT JOIN m_grades ON m_members.grade_id = m_grades.grade_id
LEFT JOIN m_groups ON m_members.group_id = m_groups.group_id
LEFT JOIN m_organizations ON m_members.personal_organization_id = m_organizations.organization_id
LEFT JOIN m_roles ON m_members.role_id = m_roles.role_id
WHERE member_id = $1;

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
	m_members_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetMembersWithAttendStatus :many
SELECT sqlc.embed(m_members), sqlc.embed(m_attend_statuses) FROM m_members
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
SELECT sqlc.embed(m_members), sqlc.embed(m_attend_statuses) FROM m_members
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
SELECT sqlc.embed(m_members), sqlc.embed(m_attend_statuses) FROM m_members
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
SELECT sqlc.embed(m_members), sqlc.embed(m_attend_statuses) FROM m_members
LEFT JOIN m_attend_statuses ON m_members.attend_status_id = m_attend_statuses.attend_status_id
WHERE member_id = ANY(@member_ids::uuid[])
ORDER BY
	m_members_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetMembersWithGrade :many
SELECT sqlc.embed(m_members), sqlc.embed(m_grades) FROM m_members
LEFT JOIN m_grades ON m_members.grade_id = m_grades.grade_id
LEFT JOIN m_organizations ON m_grades.organization_id = m_organizations.organization_id
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

-- name: GetMembersWithGradeUseNumberedPaginate :many
SELECT sqlc.embed(m_members), sqlc.embed(m_grades) FROM m_members
LEFT JOIN m_grades ON m_members.grade_id = m_grades.grade_id
LEFT JOIN m_organizations ON m_grades.organization_id = m_organizations.organization_id
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

-- name: GetMembersWithGradeUseKeysetPaginate :many
SELECT sqlc.embed(m_members), sqlc.embed(m_grades) FROM m_members
LEFT JOIN m_grades ON m_members.grade_id = m_grades.grade_id
LEFT JOIN m_organizations ON m_grades.organization_id = m_organizations.organization_id
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

-- name: GetPluralMembersWithGrade :many
SELECT sqlc.embed(m_members), sqlc.embed(m_grades) FROM m_members
LEFT JOIN m_grades ON m_members.grade_id = m_grades.grade_id
LEFT JOIN m_organizations ON m_grades.organization_id = m_organizations.organization_id
WHERE member_id = ANY(@member_ids::uuid[])
ORDER BY
	m_members_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetMembersWithGroup :many
SELECT sqlc.embed(m_members), sqlc.embed(m_groups) FROM m_members
LEFT JOIN m_groups ON m_members.group_id = m_groups.group_id
LEFT JOIN m_organizations ON m_groups.organization_id = m_organizations.organization_id
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

-- name: GetMembersWithGroupUseNumberedPaginate :many
SELECT sqlc.embed(m_members), sqlc.embed(m_groups) FROM m_members
LEFT JOIN m_groups ON m_members.group_id = m_groups.group_id
LEFT JOIN m_organizations ON m_groups.organization_id = m_organizations.organization_id
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

-- name: GetMembersWithGroupUseKeysetPaginate :many
SELECT sqlc.embed(m_members), sqlc.embed(m_groups) FROM m_members
LEFT JOIN m_groups ON m_members.group_id = m_groups.group_id
LEFT JOIN m_organizations ON m_groups.organization_id = m_organizations.organization_id
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

-- name: GetPluralMembersWithGroup :many
SELECT sqlc.embed(m_members), sqlc.embed(m_groups) FROM m_members
LEFT JOIN m_groups ON m_members.group_id = m_groups.group_id
LEFT JOIN m_organizations ON m_groups.organization_id = m_organizations.organization_id
WHERE member_id = ANY(@member_ids::uuid[])
ORDER BY
	m_members_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetMembersWithPersonalOrganization :many
SELECT sqlc.embed(m_members), sqlc.embed(m_organizations) FROM m_members
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
SELECT sqlc.embed(m_members), sqlc.embed(m_organizations) FROM m_members
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
SELECT sqlc.embed(m_members), sqlc.embed(m_organizations) FROM m_members
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
SELECT sqlc.embed(m_members), sqlc.embed(m_organizations) FROM m_members
LEFT JOIN m_organizations ON m_members.personal_organization_id = m_organizations.organization_id
WHERE member_id = ANY(@member_ids::uuid[])
ORDER BY
	m_members_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetMembersWithRole :many
SELECT sqlc.embed(m_members), sqlc.embed(m_roles) FROM m_members
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
SELECT sqlc.embed(m_members), sqlc.embed(m_roles) FROM m_members
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
SELECT sqlc.embed(m_members), sqlc.embed(m_roles) FROM m_members
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
SELECT sqlc.embed(m_members), sqlc.embed(m_roles) FROM m_members
LEFT JOIN m_roles ON m_members.role_id = m_roles.role_id
WHERE member_id = ANY(@member_ids::uuid[])
ORDER BY
	m_members_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetMembersWithAll :many
SELECT sqlc.embed(m_members), sqlc.embed(m_attend_statuses), sqlc.embed(m_grades), sqlc.embed(m_groups), sqlc.embed(m_organizations), sqlc.embed(m_roles) FROM m_members
LEFT JOIN m_attend_statuses ON m_members.attend_status_id = m_attend_statuses.attend_status_id
LEFT JOIN m_grades ON m_members.grade_id = m_grades.grade_id
LEFT JOIN m_groups ON m_members.group_id = m_groups.group_id
LEFT JOIN m_organizations ON m_members.personal_organization_id = m_organizations.organization_id
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

-- name: GetMembersWithAllUseNumberedPaginate :many
SELECT sqlc.embed(m_members), sqlc.embed(m_attend_statuses), sqlc.embed(m_grades), sqlc.embed(m_groups), sqlc.embed(m_organizations), sqlc.embed(m_roles) FROM m_members
LEFT JOIN m_attend_statuses ON m_members.attend_status_id = m_attend_statuses.attend_status_id
LEFT JOIN m_grades ON m_members.grade_id = m_grades.grade_id
LEFT JOIN m_groups ON m_members.group_id = m_groups.group_id
LEFT JOIN m_organizations ON m_members.personal_organization_id = m_organizations.organization_id
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

-- name: GetMembersWithAllUseKeysetPaginate :many
SELECT sqlc.embed(m_members), sqlc.embed(m_attend_statuses), sqlc.embed(m_grades), sqlc.embed(m_groups), sqlc.embed(m_organizations), sqlc.embed(m_roles) FROM m_members
LEFT JOIN m_attend_statuses ON m_members.attend_status_id = m_attend_statuses.attend_status_id
LEFT JOIN m_grades ON m_members.grade_id = m_grades.grade_id
LEFT JOIN m_groups ON m_members.group_id = m_groups.group_id
LEFT JOIN m_organizations ON m_members.personal_organization_id = m_organizations.organization_id
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

-- name: GetPluralMembersWithAll :many
SELECT sqlc.embed(m_members), sqlc.embed(m_attend_statuses), sqlc.embed(m_grades), sqlc.embed(m_groups), sqlc.embed(m_organizations), sqlc.embed(m_roles) FROM m_members
LEFT JOIN m_attend_statuses ON m_members.attend_status_id = m_attend_statuses.attend_status_id
LEFT JOIN m_grades ON m_members.grade_id = m_grades.grade_id
LEFT JOIN m_groups ON m_members.group_id = m_groups.group_id
LEFT JOIN m_organizations ON m_members.personal_organization_id = m_organizations.organization_id
LEFT JOIN m_roles ON m_members.role_id = m_roles.role_id
WHERE member_id = ANY(@member_ids::uuid[])
ORDER BY
	m_members_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountMembers :one
SELECT COUNT(*) FROM m_members
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_members.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_has_policy::boolean = true THEN (SELECT COUNT(*) FROM m_role_associations WHERE role_id = m_members.role_id AND m_role_associations.policy_id = ANY(@has_policy_ids::uuid[])) > 0 ELSE TRUE END
AND
	CASE WHEN @when_in_attend_status::boolean = true THEN m_members.attend_status_id = ANY(@in_attend_status_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @when_in_grade::boolean = true THEN m_members.grade_id = ANY(@in_grade_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @when_in_group::boolean = true THEN m_members.group_id = ANY(@in_group_ids::uuid[]) ELSE TRUE END;
