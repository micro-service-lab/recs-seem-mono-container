-- name: CreateGrades :copyfrom
INSERT INTO m_grades (key, organization_id) VALUES ($1, $2);

-- name: CreateGrade :one
INSERT INTO m_grades (key, organization_id) VALUES ($1, $2) RETURNING *;

-- name: DeleteGrade :exec
DELETE FROM m_grades WHERE grade_id = $1;

-- name: DeleteGradeByKey :exec
DELETE FROM m_grades WHERE key = $1;

-- name: FindGradeByID :one
SELECT * FROM m_grades WHERE grade_id = $1;

-- name: FindGradeByIDWithOrganization :one
SELECT sqlc.embed(m_grades), sqlc.embed(m_organizations) FROM m_grades
INNER JOIN m_organizations ON m_grades.organization_id = m_organizations.organization_id
WHERE grade_id = $1;

-- name: FindGradeByKey :one
SELECT * FROM m_grades WHERE key = $1;

-- name: FindGradeByKeyWithOrganization :one
SELECT sqlc.embed(m_grades), sqlc.embed(m_organizations) FROM m_grades
INNER JOIN m_organizations ON m_grades.organization_id = m_organizations.organization_id
WHERE key = $1;

-- name: GetGrades :many
SELECT * FROM m_grades
ORDER BY
	m_grades_pkey DESC
LIMIT $1 OFFSET $2;

-- name: GetGradesWithOrganization :many
SELECT sqlc.embed(m_grades), sqlc.embed(m_organizations) FROM m_grades
INNER JOIN m_organizations ON m_grades.organization_id = m_organizations.organization_id
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_organizations.name END ASC,
	m_grades_pkey DESC
LIMIT $1 OFFSET $2;

-- name: CountGrades :one
SELECT COUNT(*) FROM m_grades;

