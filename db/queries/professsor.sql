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
SELECT m_professors.*, m_members.name member_name, m_members.first_name member_first_name, m_members.last_name member_last_name, m_members.email member_email,
m_members.profile_image_id member_profile_image_id, t_images.height member_profile_image_height,
t_images.width member_profile_image_width, t_images.attachable_item_id member_profile_image_attachable_item_id,
t_attachable_items.owner_id member_profile_image_owner_id, t_attachable_items.from_outer member_profile_image_from_outer, t_attachable_items.alias member_profile_image_alias,
t_attachable_items.url member_profile_image_url, t_attachable_items.size member_profile_image_size, t_attachable_items.mime_type_id member_profile_image_mime_type_id FROM m_professors
LEFT JOIN m_members ON m_professors.member_id = m_members.member_id
LEFT JOIN t_images ON m_members.profile_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
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
WHERE professor_id = ANY(@professor_ids::uuid[])
ORDER BY
	m_professors_pkey ASC;

-- name: GetPluralProfessorsUseNumberedPaginate :many
SELECT * FROM m_professors
WHERE professor_id = ANY(@professor_ids::uuid[])
ORDER BY
	m_professors_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetProfessorsWithMember :many
SELECT m_professors.*, m_members.name member_name, m_members.first_name member_first_name, m_members.last_name member_last_name, m_members.email member_email,
m_members.profile_image_id member_profile_image_id, t_images.height member_profile_image_height,
t_images.width member_profile_image_width, t_images.attachable_item_id member_profile_image_attachable_item_id,
t_attachable_items.owner_id member_profile_image_owner_id, t_attachable_items.from_outer member_profile_image_from_outer, t_attachable_items.alias member_profile_image_alias,
t_attachable_items.url member_profile_image_url, t_attachable_items.size member_profile_image_size, t_attachable_items.mime_type_id member_profile_image_mime_type_id FROM m_professors
LEFT JOIN m_members ON m_professors.member_id = m_members.member_id
LEFT JOIN t_images ON m_members.profile_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
ORDER BY
	m_professors_pkey ASC;

-- name: GetProfessorsWithMemberUseNumberedPaginate :many
SELECT m_professors.*, m_members.name member_name, m_members.first_name member_first_name, m_members.last_name member_last_name, m_members.email member_email,
m_members.profile_image_id member_profile_image_id, t_images.height member_profile_image_height,
t_images.width member_profile_image_width, t_images.attachable_item_id member_profile_image_attachable_item_id,
t_attachable_items.owner_id member_profile_image_owner_id, t_attachable_items.from_outer member_profile_image_from_outer, t_attachable_items.alias member_profile_image_alias,
t_attachable_items.url member_profile_image_url, t_attachable_items.size member_profile_image_size, t_attachable_items.mime_type_id member_profile_image_mime_type_id FROM m_professors
LEFT JOIN m_members ON m_professors.member_id = m_members.member_id
LEFT JOIN t_images ON m_members.profile_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
ORDER BY
	m_professors_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetProfessorsWithMemberUseKeysetPaginate :many
SELECT m_professors.*, m_members.name member_name, m_members.first_name member_first_name, m_members.last_name member_last_name, m_members.email member_email,
m_members.profile_image_id member_profile_image_id, t_images.height member_profile_image_height,
t_images.width member_profile_image_width, t_images.attachable_item_id member_profile_image_attachable_item_id,
t_attachable_items.owner_id member_profile_image_owner_id, t_attachable_items.from_outer member_profile_image_from_outer, t_attachable_items.alias member_profile_image_alias,
t_attachable_items.url member_profile_image_url, t_attachable_items.size member_profile_image_size, t_attachable_items.mime_type_id member_profile_image_mime_type_id FROM m_professors
LEFT JOIN m_members ON m_professors.member_id = m_members.member_id
LEFT JOIN t_images ON m_members.profile_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
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

-- name: GetPluralProfessorsWithMember :many
SELECT m_professors.*, m_members.name member_name, m_members.first_name member_first_name, m_members.last_name member_last_name, m_members.email member_email,
m_members.profile_image_id member_profile_image_id, t_images.height member_profile_image_height,
t_images.width member_profile_image_width, t_images.attachable_item_id member_profile_image_attachable_item_id,
t_attachable_items.owner_id member_profile_image_owner_id, t_attachable_items.from_outer member_profile_image_from_outer, t_attachable_items.alias member_profile_image_alias,
t_attachable_items.url member_profile_image_url, t_attachable_items.size member_profile_image_size, t_attachable_items.mime_type_id member_profile_image_mime_type_id FROM m_professors
LEFT JOIN m_members ON m_professors.member_id = m_members.member_id
LEFT JOIN t_images ON m_members.profile_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE professor_id = ANY(@professor_ids::uuid[])
ORDER BY
	m_professors_pkey ASC;

-- name: GetPluralProfessorsWithMemberUseNumberedPaginate :many
SELECT m_professors.*, m_members.name member_name, m_members.first_name member_first_name, m_members.last_name member_last_name, m_members.email member_email,
m_members.profile_image_id member_profile_image_id, t_images.height member_profile_image_height,
t_images.width member_profile_image_width, t_images.attachable_item_id member_profile_image_attachable_item_id,
t_attachable_items.owner_id member_profile_image_owner_id, t_attachable_items.from_outer member_profile_image_from_outer, t_attachable_items.alias member_profile_image_alias,
t_attachable_items.url member_profile_image_url, t_attachable_items.size member_profile_image_size, t_attachable_items.mime_type_id member_profile_image_mime_type_id FROM m_professors
LEFT JOIN m_members ON m_professors.member_id = m_members.member_id
LEFT JOIN t_images ON m_members.profile_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE professor_id = ANY(@professor_ids::uuid[])
ORDER BY
	m_professors_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountProfessors :one
SELECT COUNT(*) FROM m_professors;
