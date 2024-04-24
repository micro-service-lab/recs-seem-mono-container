-- name: CreateMembers :copyfrom
INSERT INTO m_members (login_id, password, email, name, attend_status_id, grade_id, group_id, personal_organization_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);

-- name: CreateMember :one
INSERT INTO m_members (login_id, password, email, name, attend_status_id, grade_id, group_id, personal_organization_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING *;

-- name: DeleteMember :exec
DELETE FROM m_members WHERE member_id = $1;

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

-- name: FindMemberByLoginID :one
SELECT * FROM m_members WHERE login_id = $1;

-- name: GetMembers :many
SELECT * FROM m_members
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
