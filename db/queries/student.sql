-- name: CreateStudents :copyfrom
INSERT INTO m_students (member_id) VALUES ($1);

-- name: CreateStudent :one
INSERT INTO m_students (member_id) VALUES ($1) RETURNING *;

-- name: DeleteStudent :exec
DELETE FROM m_students WHERE student_id = $1;

-- name: FindStudentByID :one
SELECT * FROM m_students WHERE student_id = $1;

-- name: FindStudentByIDWithMember :one
SELECT sqlc.embed(m_students), sqlc.embed(m_members) FROM m_students
LEFT JOIN m_members ON m_students.member_id = m_members.member_id
WHERE student_id = $1;

-- name: GetStudents :many
SELECT * FROM m_students
ORDER BY
	m_students_pkey DESC;

-- name: GetStudentsUseNumberedPaginate :many
SELECT * FROM m_students
ORDER BY
	m_students_pkey DESC
LIMIT $1 OFFSET $2;

-- name: GetStudentsUseKeysetPaginate :many
SELECT * FROM m_students
WHERE
	CASE @cursor_direction
		WHEN 'next' THEN
			m_students_pkey < @cursor
		WHEN 'prev' THEN
			m_students_pkey > @cursor
	END
ORDER BY
	m_students_pkey DESC
LIMIT $1;

-- name: CountStudents :one
SELECT COUNT(*) FROM m_students;

