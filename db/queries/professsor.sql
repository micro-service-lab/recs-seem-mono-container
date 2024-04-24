-- name: CreateProfessors :copyfrom
INSERT INTO m_professors (member_id) VALUES ($1);

-- name: CreateProfessor :one
INSERT INTO m_professors (member_id) VALUES ($1) RETURNING *;

-- name: DeleteProfessor :exec
DELETE FROM m_professors WHERE professor_id = $1;

-- name: FindProfessorByID :one
SELECT * FROM m_professors WHERE professor_id = $1;

-- name: FindProfessorByMemberID :one
SELECT * FROM m_professors WHERE member_id = $1;

-- name: FindProfessorWithMember :one
SELECT sqlc.embed(m_professors), sqlc.embed(m_members) FROM m_professors
INNER JOIN m_members ON m_professors.member_id = m_members.member_id
WHERE professor_id = $1;

-- name: GetProfessors :many
SELECT * FROM m_professors
ORDER BY
	m_professors_pkey DESC
LIMIT $1 OFFSET $2;

-- name: GetProfessorsWithMember :many
SELECT sqlc.embed(m_professors), sqlc.embed(m_members) FROM m_professors
INNER JOIN m_members ON m_professors.member_id = m_members.member_id
ORDER BY
	m_professors_pkey DESC
LIMIT $1 OFFSET $2;

-- name: CountProfessors :one
SELECT COUNT(*) FROM m_professors;
