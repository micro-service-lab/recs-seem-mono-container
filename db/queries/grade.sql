-- name: CreateGrades :copyfrom
INSERT INTO m_grades (key, organization_id) VALUES ($1, $2);

-- name: CreateGrade :one
INSERT INTO m_grades (key, organization_id) VALUES ($1, $2) RETURNING *;

-- name: DeleteGrade :execrows
DELETE FROM m_grades WHERE grade_id = $1;

-- name: DeleteGradeByKey :execrows
DELETE FROM m_grades WHERE key = $1;

-- name: PluralDeleteGrades :execrows
DELETE FROM m_grades WHERE grade_id = ANY(@grade_ids::uuid[]);

-- name: FindGradeByID :one
SELECT * FROM m_grades WHERE grade_id = $1;

-- name: FindGradeByIDWithOrganization :one
SELECT m_grades.*, m_organizations.name organization_name, m_organizations.description organization_description,
m_organizations.color organization_color, m_organizations.is_personal organization_is_personal,
m_organizations.is_whole organization_is_whole, m_organizations.chat_room_id organization_chat_room_id
FROM m_grades
LEFT JOIN m_organizations ON m_grades.organization_id = m_organizations.organization_id
WHERE grade_id = $1;

-- name: FindGradeByKey :one
SELECT * FROM m_grades WHERE key = $1;

-- name: FindGradeByKeyWithOrganization :one
SELECT m_grades.*, m_organizations.name organization_name, m_organizations.description organization_description,
m_organizations.color organization_color, m_organizations.is_personal organization_is_personal,
m_organizations.is_whole organization_is_whole, m_organizations.chat_room_id organization_chat_room_id
FROM m_grades
LEFT JOIN m_organizations ON m_grades.organization_id = m_organizations.organization_id
WHERE key = $1;

-- name: GetGrades :many
SELECT * FROM m_grades
ORDER BY
	m_grades_pkey ASC;

-- name: GetGradesUseNumberedPaginate :many
SELECT * FROM m_grades
ORDER BY
	m_grades_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetGradesUseKeysetPaginate :many
SELECT * FROM m_grades
WHERE
	CASE @cursor_direction::text
		WHEN 'next' THEN
			m_grades_pkey > @cursor::int
		WHEN 'prev' THEN
			m_grades_pkey < @cursor::int
	END
ORDER BY
	CASE WHEN @cursor_direction::text = 'next' THEN m_grades_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN m_grades_pkey END DESC
LIMIT $1;

-- name: GetPluralGrades :many
SELECT * FROM m_grades
WHERE grade_id = ANY(@grade_ids::uuid[])
ORDER BY
	m_grades_pkey ASC;

-- name: GetPluralGradesUseNumberedPaginate :many
SELECT * FROM m_grades
WHERE grade_id = ANY(@grade_ids::uuid[])
ORDER BY
	m_grades_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetGradesWithOrganization :many
SELECT m_grades.*, m_organizations.name organization_name, m_organizations.description organization_description,
m_organizations.color organization_color, m_organizations.is_personal organization_is_personal,
m_organizations.is_whole organization_is_whole, m_organizations.chat_room_id organization_chat_room_id
FROM m_grades
LEFT JOIN m_organizations ON m_grades.organization_id = m_organizations.organization_id
ORDER BY
	m_grades_pkey ASC;

-- name: GetGradesWithOrganizationUseNumberedPaginate :many
SELECT m_grades.*, m_organizations.name organization_name, m_organizations.description organization_description,
m_organizations.color organization_color, m_organizations.is_personal organization_is_personal,
m_organizations.is_whole organization_is_whole, m_organizations.chat_room_id organization_chat_room_id
FROM m_grades
LEFT JOIN m_organizations ON m_grades.organization_id = m_organizations.organization_id
ORDER BY
	m_grades_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetGradesWithOrganizationUseKeysetPaginate :many
SELECT m_grades.*, m_organizations.name organization_name, m_organizations.description organization_description,
m_organizations.color organization_color, m_organizations.is_personal organization_is_personal,
m_organizations.is_whole organization_is_whole, m_organizations.chat_room_id organization_chat_room_id
FROM m_grades
LEFT JOIN m_organizations ON m_grades.organization_id = m_organizations.organization_id
WHERE
	CASE @cursor_direction::text
		WHEN 'next' THEN
			m_grades_pkey > @cursor::int
		WHEN 'prev' THEN
			m_grades_pkey < @cursor::int
	END
ORDER BY
	CASE WHEN @cursor_direction::text = 'next' THEN m_grades_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN m_grades_pkey END DESC
LIMIT $1;

-- name: GetPluralGradesWithOrganization :many
SELECT m_grades.*, m_organizations.name organization_name, m_organizations.description organization_description,
m_organizations.color organization_color, m_organizations.is_personal organization_is_personal,
m_organizations.is_whole organization_is_whole, m_organizations.chat_room_id organization_chat_room_id
FROM m_grades
LEFT JOIN m_organizations ON m_grades.organization_id = m_organizations.organization_id
WHERE grade_id = ANY(@grade_ids::uuid[])
ORDER BY
	m_grades_pkey ASC;

-- name: GetPluralGradesWithOrganizationUseNumberedPaginate :many
SELECT m_grades.*, m_organizations.name organization_name, m_organizations.description organization_description,
m_organizations.color organization_color, m_organizations.is_personal organization_is_personal,
m_organizations.is_whole organization_is_whole, m_organizations.chat_room_id organization_chat_room_id
FROM m_grades
LEFT JOIN m_organizations ON m_grades.organization_id = m_organizations.organization_id
WHERE grade_id = ANY(@grade_ids::uuid[])
ORDER BY
	m_grades_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountGrades :one
SELECT COUNT(*) FROM m_grades;

