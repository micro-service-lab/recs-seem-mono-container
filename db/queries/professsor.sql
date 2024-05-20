-- name: CreateProfessors :copyfrom
INSERT INTO m_professors (member_id) VALUES ($1);

-- name: CreateProfessor :one
INSERT INTO m_professors (member_id) VALUES ($1) RETURNING *;

-- name: DeleteProfessor :execrows
DELETE FROM m_professors WHERE professor_id = $1;

-- name: PluralDeleteProfessors :execrows
DELETE FROM m_professors WHERE professor_id = ANY(@professor_ids::uuid[]);

-- name: FindProfessorByID :one
SELECT * FROM m_professors WHERE professor_id = $1;

-- name: FindProfessorByIDWithMember :one
SELECT sqlc.embed(m_professors), sqlc.embed(m_members) FROM m_professors
LEFT JOIN m_members ON m_professors.member_id = m_members.member_id
WHERE professor_id = $1;

-- name: GetProfessors :many
SELECT * FROM m_professors
ORDER BY
	m_professors_pkey ASC;

-- name: GetProfessorsUseNumberedPaginate :many
SELECT * FROM m_professors
ORDER BY
	m_professors_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetProfessorsUseKeysetPaginate :many
SELECT * FROM m_professors
WHERE
	CASE @cursor_direction::text
		WHEN 'next' THEN
			m_professors_pkey > @cursor::int
		WHEN 'prev' THEN
			m_professors_pkey < @cursor::int
	END
ORDER BY
	CASE WHEN @cursor_direction::text = 'next' THEN m_professors_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN m_professors_pkey END DESC
LIMIT $1;

-- name: GetPluralProfessors :many
SELECT * FROM m_professors
WHERE member_id = ANY(@member_ids::uuid[])
ORDER BY
	m_professors_pkey ASC;

-- name: GetPluralProfessorsUseNumberedPaginate :many
SELECT * FROM m_professors
WHERE member_id = ANY(@member_ids::uuid[])
ORDER BY
	m_professors_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetProfessorsWithMember :many
SELECT m_professors.*, sqlc.embed(m_members) FROM m_professors
LEFT JOIN m_members ON m_professors.member_id = m_members.member_id
ORDER BY
	m_professors_pkey ASC;

-- name: GetProfessorsWithMemberUseNumberedPaginate :many
SELECT m_professors.*, sqlc.embed(m_members) FROM m_professors
LEFT JOIN m_members ON m_professors.member_id = m_members.member_id
ORDER BY
	m_professors_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetProfessorsWithMemberUseKeysetPaginate :many
SELECT m_professors.*, sqlc.embed(m_members) FROM m_professors
LEFT JOIN m_members ON m_professors.member_id = m_members.member_id
WHERE
	CASE @cursor_direction::text
		WHEN 'next' THEN
			m_professors_pkey > @cursor::int
		WHEN 'prev' THEN
			m_professors_pkey < @cursor::int
	END
ORDER BY
	CASE WHEN @cursor_direction::text = 'next' THEN m_professors_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN m_professors_pkey END DESC
LIMIT $1;

-- name: CountProfessors :one
SELECT COUNT(*) FROM m_professors;
