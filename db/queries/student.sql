-- name: CreateStudents :copyfrom
INSERT INTO m_students (member_id) VALUES ($1);

-- name: CreateStudent :one
INSERT INTO m_students (member_id) VALUES ($1) RETURNING *;

-- name: DeleteStudent :exec
DELETE FROM m_students WHERE student_id = $1;

-- name: FindStudentByID :one
SELECT * FROM m_students WHERE student_id = $1;

-- name: FindStudentByMemberID :one
SELECT * FROM m_students WHERE member_id = $1;

-- name: FindStudentWithMember :one
SELECT sqlc.embed(m_students), sqlc.embed(m_members) FROM m_students
INNER JOIN m_members ON m_students.member_id = m_members.member_id
WHERE student_id = $1;

-- name: GetStudents :many
SELECT * FROM m_students
ORDER BY
	m_students_pkey DESC
LIMIT $1 OFFSET $2;

-- name: GetStudentsWithMember :many
SELECT sqlc.embed(m_students), sqlc.embed(m_members) FROM m_students
INNER JOIN m_members ON m_students.member_id = m_members.member_id
ORDER BY
	m_students_pkey DESC
LIMIT $1 OFFSET $2;

-- name: CountStudents :one
SELECT COUNT(*) FROM m_students;

