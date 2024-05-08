-- name: CreateStudents :copyfrom
INSERT INTO m_students (member_id) VALUES ($1);

-- name: CreateStudent :one
INSERT INTO m_students (member_id) VALUES ($1) RETURNING *;

-- name: DeleteStudent :exec
DELETE FROM m_students WHERE student_id = $1;

-- name: PluralDeleteStudents :exec
DELETE FROM m_students WHERE student_id = ANY($1::uuid[]);

-- name: FindStudentByID :one
SELECT * FROM m_students WHERE student_id = $1;

-- name: FindStudentByIDWithMember :one
SELECT m_students.*, sqlc.embed(m_members) FROM m_students
LEFT JOIN m_members ON m_students.member_id = m_members.member_id
WHERE student_id = $1;

-- name: GetStudents :many
SELECT * FROM m_students
ORDER BY
	m_students_pkey ASC;

-- name: GetStudentsUseNumberedPaginate :many
SELECT * FROM m_students
ORDER BY
	m_students_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetStudentsUseKeysetPaginate :many
SELECT * FROM m_students
WHERE
	CASE @cursor_direction::text
		WHEN 'next' THEN
			m_students_pkey > @cursor::int
		WHEN 'prev' THEN
			m_students_pkey < @cursor::int
	END
ORDER BY
	CASE WHEN @cursor_direction::text = 'next' THEN m_students_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN m_students_pkey END DESC
LIMIT $1;

-- name: GetPluralStudents :many
SELECT * FROM m_students
WHERE member_id = ANY(@member_ids::uuid[])
ORDER BY
	m_students_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountStudents :one
SELECT COUNT(*) FROM m_students;

