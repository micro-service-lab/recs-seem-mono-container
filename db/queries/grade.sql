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
LEFT JOIN m_organizations ON m_grades.organization_id = m_organizations.organization_id
WHERE grade_id = $1;

-- name: FindGradeByKey :one
SELECT * FROM m_grades WHERE key = $1;

-- name: FindGradeByKeyWithOrganization :one
SELECT sqlc.embed(m_grades), sqlc.embed(m_organizations) FROM m_grades
LEFT JOIN m_organizations ON m_grades.organization_id = m_organizations.organization_id
WHERE key = $1;

-- name: GetGrades :many
SELECT * FROM m_grades
ORDER BY
	m_grades_pkey DESC;

-- name: GetGradesUseNumberedPaginate :many
SELECT * FROM m_grades
ORDER BY
	m_grades_pkey DESC
LIMIT $1 OFFSET $2;

-- name: GetGradesUseKeysetPaginate :many
SELECT * FROM m_grades
WHERE
	CASE @cursor_direction::text
		WHEN 'next' THEN
			m_grades_pkey < @cursor::int
		WHEN 'prev' THEN
			m_grades_pkey > @cursor::int
	END
ORDER BY
	m_grades_pkey DESC
LIMIT $1;

-- name: GetPluralGrades :many
SELECT * FROM m_grades
WHERE organization_id = ANY(@organization_ids::uuid[])
ORDER BY
	m_grades_pkey DESC
LIMIT $1 OFFSET $2;

-- name: GetGradesWithOrganization :many
SELECT sqlc.embed(m_grades), sqlc.embed(m_organizations) FROM m_grades
LEFT JOIN m_organizations ON m_grades.organization_id = m_organizations.organization_id
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_organizations.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_organizations.name END DESC,
	m_grades_pkey DESC;

-- name: GetGradesWithOrganizationUseNumberedPaginate :many
SELECT sqlc.embed(m_grades), sqlc.embed(m_organizations) FROM m_grades
LEFT JOIN m_organizations ON m_grades.organization_id = m_organizations.organization_id
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_organizations.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_organizations.name END DESC,
	m_grades_pkey DESC
LIMIT $1 OFFSET $2;

-- name: GetGradesWithOrganizationUseKeysetPaginate :many
SELECT sqlc.embed(m_grades), sqlc.embed(m_organizations) FROM m_grades
LEFT JOIN m_organizations ON m_grades.organization_id = m_organizations.organization_id
WHERE
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'name' THEN name > @name_cursor OR (name = @name_cursor AND m_grades_pkey < @cursor::int)
				WHEN 'r_name' THEN name < @name_cursor OR (name = @name_cursor AND m_grades_pkey < @cursor::int)
				ELSE m_grades_pkey < @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'name' THEN name < @name_cursor OR (name = @name_cursor AND m_grades_pkey > @cursor::int)
				WHEN 'r_name' THEN name > @name_cursor OR (name = @name_cursor AND m_grades_pkey > @cursor::int)
				ELSE m_grades_pkey > @cursor::int
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_organizations.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_organizations.name END DESC,
	m_grades_pkey DESC
LIMIT $1;

-- name: GetPluralGradesWithOrganization :many
SELECT sqlc.embed(m_grades), sqlc.embed(m_organizations) FROM m_grades
LEFT JOIN m_organizations ON m_grades.organization_id = m_organizations.organization_id
WHERE organization_id = ANY(@organization_ids::uuid[])
ORDER BY
	m_grades_pkey DESC
LIMIT $1 OFFSET $2;

-- name: CountGrades :one
SELECT COUNT(*) FROM m_grades;

