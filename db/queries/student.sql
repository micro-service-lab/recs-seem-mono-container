-- name: CreateStudents :copyfrom
INSERT INTO m_students (member_id) VALUES ($1);

-- name: CreateStudent :one
INSERT INTO m_students (member_id) VALUES ($1) RETURNING *;

-- name: DeleteStudent :execrows
DELETE FROM m_students WHERE student_id = $1;

-- name: PluralDeleteStudents :execrows
DELETE FROM m_students WHERE student_id = ANY(@student_ids::uuid[]);

-- name: FindStudentByID :one
SELECT * FROM m_students WHERE student_id = $1;

-- name: FindStudentByIDWithMember :one
SELECT m_students.*, m_members.name member_name, m_members.first_name member_first_name, m_members.last_name member_last_name, m_members.email member_email, m_members.grade_id member_grade_id, m_members.group_id member_group_id,
m_members.profile_image_id member_profile_image_id, t_images.height member_profile_image_height,
t_images.width member_profile_image_width, t_images.attachable_item_id member_profile_image_attachable_item_id,
t_attachable_items.owner_id member_profile_image_owner_id, t_attachable_items.from_outer member_profile_image_from_outer, t_attachable_items.alias member_profile_image_alias,
t_attachable_items.url member_profile_image_url, t_attachable_items.size member_profile_image_size, t_attachable_items.mime_type_id member_profile_image_mime_type_id FROM m_students
LEFT JOIN m_members ON m_students.member_id = m_members.member_id
LEFT JOIN t_images ON m_members.profile_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
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
WHERE student_id = ANY(@student_ids::uuid[])
ORDER BY
	m_students_pkey ASC;

-- name: GetPluralStudentsUseNumberedPaginate :many
SELECT * FROM m_students
WHERE student_id = ANY(@student_ids::uuid[])
ORDER BY
	m_students_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetStudentsWithMember :many
SELECT m_students.*, m_members.name member_name, m_members.first_name member_first_name, m_members.last_name member_last_name, m_members.email member_email, m_members.grade_id member_grade_id, m_members.group_id member_group_id,
m_members.profile_image_id member_profile_image_id, t_images.height member_profile_image_height,
t_images.width member_profile_image_width, t_images.attachable_item_id member_profile_image_attachable_item_id,
t_attachable_items.owner_id member_profile_image_owner_id, t_attachable_items.from_outer member_profile_image_from_outer, t_attachable_items.alias member_profile_image_alias,
t_attachable_items.url member_profile_image_url, t_attachable_items.size member_profile_image_size, t_attachable_items.mime_type_id member_profile_image_mime_type_id FROM m_students
LEFT JOIN m_members ON m_students.member_id = m_members.member_id
LEFT JOIN t_images ON m_members.profile_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
ORDER BY
	m_students_pkey ASC;

-- name: GetStudentsWithMemberUseNumberedPaginate :many
SELECT m_students.*, m_members.name member_name, m_members.first_name member_first_name, m_members.last_name member_last_name, m_members.email member_email, m_members.grade_id member_grade_id, m_members.group_id member_group_id,
m_members.profile_image_id member_profile_image_id, t_images.height member_profile_image_height,
t_images.width member_profile_image_width, t_images.attachable_item_id member_profile_image_attachable_item_id,
t_attachable_items.owner_id member_profile_image_owner_id, t_attachable_items.from_outer member_profile_image_from_outer, t_attachable_items.alias member_profile_image_alias,
t_attachable_items.url member_profile_image_url, t_attachable_items.size member_profile_image_size, t_attachable_items.mime_type_id member_profile_image_mime_type_id FROM m_students
LEFT JOIN m_members ON m_students.member_id = m_members.member_id
LEFT JOIN t_images ON m_members.profile_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
ORDER BY
	m_students_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetStudentsWithMemberUseKeysetPaginate :many
SELECT m_students.*, m_members.name member_name, m_members.first_name member_first_name, m_members.last_name member_last_name, m_members.email member_email, m_members.grade_id member_grade_id, m_members.group_id member_group_id,
m_members.profile_image_id member_profile_image_id, t_images.height member_profile_image_height,
t_images.width member_profile_image_width, t_images.attachable_item_id member_profile_image_attachable_item_id,
t_attachable_items.owner_id member_profile_image_owner_id, t_attachable_items.from_outer member_profile_image_from_outer, t_attachable_items.alias member_profile_image_alias,
t_attachable_items.url member_profile_image_url, t_attachable_items.size member_profile_image_size, t_attachable_items.mime_type_id member_profile_image_mime_type_id FROM m_students
LEFT JOIN m_members ON m_students.member_id = m_members.member_id
LEFT JOIN t_images ON m_members.profile_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
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

-- name: GetPluralStudentsWithMember :many
SELECT m_students.*, m_members.name member_name, m_members.first_name member_first_name, m_members.last_name member_last_name, m_members.email member_email, m_members.grade_id member_grade_id, m_members.group_id member_group_id,
m_members.profile_image_id member_profile_image_id, t_images.height member_profile_image_height,
t_images.width member_profile_image_width, t_images.attachable_item_id member_profile_image_attachable_item_id,
t_attachable_items.owner_id member_profile_image_owner_id, t_attachable_items.from_outer member_profile_image_from_outer, t_attachable_items.alias member_profile_image_alias,
t_attachable_items.url member_profile_image_url, t_attachable_items.size member_profile_image_size, t_attachable_items.mime_type_id member_profile_image_mime_type_id FROM m_students
LEFT JOIN m_members ON m_students.member_id = m_members.member_id
LEFT JOIN t_images ON m_members.profile_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE student_id = ANY(@student_ids::uuid[])
ORDER BY
	m_students_pkey ASC;

-- name: GetPluralStudentsWithMemberUseNumberedPaginate :many
SELECT m_students.*, m_members.name member_name, m_members.first_name member_first_name, m_members.last_name member_last_name, m_members.email member_email, m_members.grade_id member_grade_id, m_members.group_id member_group_id,
m_members.profile_image_id member_profile_image_id, t_images.height member_profile_image_height,
t_images.width member_profile_image_width, t_images.attachable_item_id member_profile_image_attachable_item_id,
t_attachable_items.owner_id member_profile_image_owner_id, t_attachable_items.from_outer member_profile_image_from_outer, t_attachable_items.alias member_profile_image_alias,
t_attachable_items.url member_profile_image_url, t_attachable_items.size member_profile_image_size, t_attachable_items.mime_type_id member_profile_image_mime_type_id FROM m_students
LEFT JOIN m_members ON m_students.member_id = m_members.member_id
LEFT JOIN t_images ON m_members.profile_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE student_id = ANY(@student_ids::uuid[])
ORDER BY
	m_students_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountStudents :one
SELECT COUNT(*) FROM m_students;

