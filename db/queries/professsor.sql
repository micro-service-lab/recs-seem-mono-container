-- name: CreateProfessors :copyfrom
INSERT INTO m_professors (member_id) VALUES ($1);

-- name: CreateProfessor :one
INSERT INTO m_professors (member_id) VALUES ($1) RETURNING *;

-- name: DeleteProfessor :exec
DELETE FROM m_professors WHERE professor_id = $1;

-- name: FindProfessorByID :one
SELECT * FROM m_professors WHERE professor_id = $1;

-- name: FindProfessorByIDWithMember :one
SELECT sqlc.embed(m_professors), sqlc.embed(m_members) FROM m_professors
LEFT JOIN m_members ON m_professors.member_id = m_members.member_id
WHERE professor_id = $1;

-- name: GetProfessors :many
SELECT * FROM m_professors
ORDER BY
	m_professors_pkey DESC;

-- name: GetProfessorsUseNumberedPaginate :many
SELECT * FROM m_professors
ORDER BY
	m_professors_pkey DESC
LIMIT $1 OFFSET $2;

-- name: GetProfessorsUseKeysetPaginate :many
SELECT * FROM m_professors
WHERE
	CASE @cursor_direction
		WHEN 'next' THEN
			m_professors_pkey < @cursor
		WHEN 'prev' THEN
			m_professors_pkey > @cursor
	END
ORDER BY
	m_professors_pkey DESC
LIMIT $1;

-- name: CountProfessors :one
SELECT COUNT(*) FROM m_professors;
