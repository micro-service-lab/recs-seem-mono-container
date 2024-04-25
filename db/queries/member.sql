-- name: CreateMembers :copyfrom
INSERT INTO m_members (login_id, password, email, name, attend_status_id, grade_id, group_id, role_id, personal_organization_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);

-- name: CreateMember :one
INSERT INTO m_members (login_id, password, email, name, attend_status_id, grade_id, group_id, role_id, personal_organization_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING *;

-- name: DeleteMember :exec
DELETE FROM m_members WHERE member_id = $1;

-- name: UpdateMember :one
UPDATE m_members SET email = $2, name = $3, updated_at = $4 WHERE member_id = $1 RETURNING *;

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
INNER JOIN m_attend_statuses ON m_members.attend_status_id = m_attend_statuses.attend_status_id
WHERE member_id = $1;

-- name: FindMemberByIDWithGrade :one
SELECT sqlc.embed(m_members), sqlc.embed(m_grades) FROM m_members
INNER JOIN m_grades ON m_members.grade_id = m_grades.grade_id
INNER JOIN m_organizations ON m_grades.organization_id = m_organizations.organization_id
WHERE member_id = $1;

-- name: FindMemberByIDWithGroup :one
SELECT sqlc.embed(m_members), sqlc.embed(m_groups) FROM m_members
INNER JOIN m_groups ON m_members.group_id = m_groups.group_id
INNER JOIN m_organizations ON m_groups.organization_id = m_organizations.organization_id
WHERE member_id = $1;

-- name: FindMemberByIDWithPersonalOrganization :one
SELECT sqlc.embed(m_members), sqlc.embed(m_organizations) FROM m_members
INNER JOIN m_organizations ON m_members.personal_organization_id = m_organizations.organization_id
WHERE member_id = $1;

-- name: FindMemberByIDWithRole :one
SELECT sqlc.embed(m_members), sqlc.embed(m_roles) FROM m_members
LEFT JOIN m_roles ON m_members.role_id = m_roles.role_id
WHERE member_id = $1;

-- name: FindMemberByIDWithDetailRole :one
SELECT sqlc.embed(m_members), sqlc.embed(m_roles), sqlc.embed(m_role_associations), sqlc.embed(m_policies) FROM m_members
LEFT JOIN m_roles ON m_members.role_id = m_roles.role_id
LEFT JOIN m_role_associations ON m_members.role_id = m_role_associations.role_id
LEFT JOIN m_policies ON m_role_associations.policy_id = m_policies.policy_id
WHERE member_id = $1;

-- name: FindMemberByLoginID :one
SELECT * FROM m_members WHERE login_id = $1;

-- name: GetMembers :many
SELECT * FROM m_members
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_members.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_has_policy::boolean = true THEN (SELECT COUNT(*) FROM m_role_associations WHERE role_id = m_members.role_id AND m_role_associations.policy_id = ANY(@has_policy_ids::uuid[])) > 0 ELSE TRUE END
ORDER BY
	m_members_pkey DESC
LIMIT $1 OFFSET $2;

-- name: GetMembersWithAttendStatus :many
SELECT sqlc.embed(m_members), sqlc.embed(m_attend_statuses) FROM m_members
INNER JOIN m_attend_statuses ON m_members.attend_status_id = m_attend_statuses.attend_status_id
ORDER BY
	m_members_pkey DESC
LIMIT $1 OFFSET $2;

-- name: GetMembersWithGrade :many
SELECT sqlc.embed(m_members), sqlc.embed(m_grades) FROM m_members
INNER JOIN m_grades ON m_members.grade_id = m_grades.grade_id
INNER JOIN m_organizations ON m_grades.organization_id = m_organizations.organization_id
ORDER BY
	m_members_pkey DESC
LIMIT $1 OFFSET $2;

-- name: GetMembersWithGroup :many
SELECT sqlc.embed(m_members), sqlc.embed(m_groups) FROM m_members
INNER JOIN m_groups ON m_members.group_id = m_groups.group_id
INNER JOIN m_organizations ON m_groups.organization_id = m_organizations.organization_id
ORDER BY
	m_members_pkey DESC
LIMIT $1 OFFSET $2;

-- name: GetMembersWithPersonalOrganization :many
SELECT sqlc.embed(m_members), sqlc.embed(m_organizations) FROM m_members
INNER JOIN m_organizations ON m_members.personal_organization_id = m_organizations.organization_id
ORDER BY
	m_members_pkey DESC
LIMIT $1 OFFSET $2;

-- name: GetMembersWithRole :many
SELECT sqlc.embed(m_members), sqlc.embed(m_roles) FROM m_members
LEFT JOIN m_roles ON m_members.role_id = m_roles.role_id
ORDER BY
	m_members_pkey DESC
LIMIT $1 OFFSET $2;

-- name: CountMembers :one
SELECT COUNT(*) FROM m_members;
