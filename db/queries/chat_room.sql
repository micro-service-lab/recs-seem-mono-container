-- name: CreateChatRooms :copyfrom
INSERT INTO m_chat_rooms (name, is_private, cover_image_id, owner_id, from_organization, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: CreateChatRoom :one
INSERT INTO m_chat_rooms (name, is_private, cover_image_id, owner_id, from_organization, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: UpdateChatRoom :one
UPDATE m_chat_rooms SET name = $2, cover_image_id = $3, updated_at = $4 WHERE chat_room_id = $1 RETURNING *;

-- name: UpdateChatRoomOwner :one
UPDATE m_chat_rooms SET owner_id = $2, updated_at = $3 WHERE chat_room_id = $1 RETURNING *;

-- name: DeleteChatRoom :execrows
DELETE FROM m_chat_rooms WHERE chat_room_id = $1;

-- name: PluralDeleteChatRooms :execrows
DELETE FROM m_chat_rooms WHERE chat_room_id = ANY(@chat_room_ids::uuid[]);

-- name: FindChatRoomByID :one
SELECT * FROM m_chat_rooms WHERE chat_room_id = $1;

-- name: FindChatRoomByIDWithCoverImage :one
SELECT m_chat_rooms.*, t_images.height cover_image_height, t_images.width cover_image_width, t_images.attachable_item_id cover_image_attachable_item_id,
t_attachable_items.owner_id cover_image_owner_id, t_attachable_items.from_outer cover_image_from_outer, t_attachable_items.alias cover_image_alias,
t_attachable_items.url cover_image_url, t_attachable_items.size cover_image_size, t_attachable_items.mime_type_id cover_image_mime_type_id
FROM m_chat_rooms
LEFT JOIN t_images ON m_chat_rooms.cover_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE chat_room_id = $1;

-- name: FindChatRoomOnPrivate :one
SELECT * FROM m_chat_rooms
WHERE (SELECT COUNT(*) FROM m_chat_room_belongings WHERE m_chat_room_belongings.chat_room_id = m_chat_rooms.chat_room_id AND
(m_chat_room_belongings.member_id = sqlc.arg(owner_id) OR m_chat_room_belongings.member_id = sqlc.arg(member_id))) = 2
AND is_private = true;

-- name: FindChatRoomOnPrivateWithMember :one
SELECT m_chat_rooms.*, m_chat_room_belongings.added_at member_added_at, m_members.member_id, m_members.name member_name, m_members.first_name member_first_name, m_members.last_name member_last_name, m_members.email member_email,
m_members.profile_image_id member_profile_image_id, m_members.grade_id member_grade_id, m_members.group_id member_group_id
FROM m_chat_rooms
LEFT JOIN m_chat_room_belongings ON m_chat_rooms.chat_room_id = m_chat_room_belongings.chat_room_id
LEFT JOIN m_members ON m_chat_room_belongings.member_id = m_members.member_id
WHERE (SELECT COUNT(*) FROM m_chat_room_belongings WHERE m_chat_room_belongings.chat_room_id = m_chat_rooms.chat_room_id AND
(m_chat_room_belongings.member_id = sqlc.arg(owner_id) OR m_chat_room_belongings.member_id = sqlc.arg(member_id))) = 2
AND is_private = true
AND m_chat_rooms_belongs.member_id <> sqlc.arg(owner_id);

-- name: FindChatRoomByIDWithOwner :one
SELECT m_chat_rooms.*, m_members.name owner_name, m_members.first_name owner_first_name, m_members.last_name owner_last_name, m_members.email owner_email,
m_members.profile_image_id owner_profile_image_id, m_members.grade_id owner_grade_id, m_members.group_id owner_group_id
FROM m_chat_rooms
LEFT JOIN m_members ON m_chat_rooms.owner_id = m_members.member_id
WHERE chat_room_id = $1;

-- name: GetChatRooms :many
SELECT * FROM m_chat_rooms
WHERE
	CASE WHEN @where_in_owner::boolean = true THEN owner_id = ANY(@in_owner::uuid[]) ELSE TRUE END
AND
	CASE WHEN @where_is_private::boolean = true THEN is_private = @is_private ELSE TRUE END
AND
	CASE WHEN @where_is_from_organization::boolean = true THEN from_organization = @is_from_organization ELSE TRUE END
AND
	CASE WHEN @where_from_organizations::boolean = true THEN (SELECT chat_room_id FROM m_organizations WHERE organization_id = ANY(@from_organizations::uuid[])) = chat_room_id ELSE TRUE END
ORDER BY
	m_chat_rooms_pkey ASC;

-- name: GetChatRoomsUseNumberedPaginate :many
SELECT * FROM m_chat_rooms
WHERE
	CASE WHEN @where_in_owner::boolean = true THEN owner_id = ANY(@in_owner::uuid[]) ELSE TRUE END
AND
	CASE WHEN @where_is_private::boolean = true THEN is_private = @is_private ELSE TRUE END
AND
	CASE WHEN @where_is_from_organization::boolean = true THEN from_organization = @is_from_organization ELSE TRUE END
AND
	CASE WHEN @where_from_organizations::boolean = true THEN (SELECT chat_room_id FROM m_organizations WHERE organization_id = ANY(@from_organizations::uuid[])) = chat_room_id ELSE TRUE END
ORDER BY
	m_chat_rooms_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetChatRoomsUseKeysetPaginate :many
SELECT * FROM m_chat_rooms
WHERE
	CASE WHEN @where_in_owner::boolean = true THEN owner_id = ANY(@in_owner::uuid[]) ELSE TRUE END
AND
	CASE WHEN @where_is_private::boolean = true THEN is_private = @is_private ELSE TRUE END
AND
	CASE WHEN @where_is_from_organization::boolean = true THEN from_organization = @is_from_organization ELSE TRUE END
AND
	CASE WHEN @where_from_organizations::boolean = true THEN (SELECT chat_room_id FROM m_organizations WHERE organization_id = ANY(@from_organizations::uuid[])) = chat_room_id ELSE TRUE END
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			m_chat_rooms_pkey > @cursor::int
		WHEN 'prev' THEN
			m_chat_rooms_pkey < @cursor::int
	END
ORDER BY
	CASE WHEN @cursor_direction::text = 'next' THEN m_chat_rooms_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN m_chat_rooms_pkey END DESC
LIMIT $1;

-- name: GetPluralChatRooms :many
SELECT * FROM m_chat_rooms
WHERE chat_room_id = ANY(@chat_room_ids::uuid[])
ORDER BY
	m_chat_rooms_pkey ASC;

-- name: GetPluralChatRoomsUseNumberedPaginate :many
SELECT * FROM m_chat_rooms
WHERE chat_room_id = ANY(@chat_room_ids::uuid[])
ORDER BY
	m_chat_rooms_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetChatRoomsWithCoverImage :many
SELECT m_chat_rooms.*, t_images.height cover_image_height, t_images.width cover_image_width, t_images.attachable_item_id cover_image_attachable_item_id,
t_attachable_items.owner_id cover_image_owner_id, t_attachable_items.from_outer cover_image_from_outer, t_attachable_items.alias cover_image_alias,
t_attachable_items.url cover_image_url, t_attachable_items.size cover_image_size, t_attachable_items.mime_type_id cover_image_mime_type_id
FROM m_chat_rooms
LEFT JOIN t_images ON m_chat_rooms.cover_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE
	CASE WHEN @where_in_owner::boolean = true THEN owner_id = ANY(@in_owner::uuid[]) ELSE TRUE END
AND
	CASE WHEN @where_is_private::boolean = true THEN is_private = @is_private ELSE TRUE END
AND
	CASE WHEN @where_is_from_organization::boolean = true THEN from_organization = @is_from_organization ELSE TRUE END
AND
	CASE WHEN @where_from_organizations::boolean = true THEN (SELECT chat_room_id FROM m_organizations WHERE organization_id = ANY(@from_organizations::uuid[])) = chat_room_id ELSE TRUE END
ORDER BY
	m_chat_rooms_pkey ASC;

-- name: GetChatRoomsWithCoverImageUseNumberedPaginate :many
SELECT m_chat_rooms.*, t_images.height cover_image_height, t_images.width cover_image_width, t_images.attachable_item_id cover_image_attachable_item_id,
t_attachable_items.owner_id cover_image_owner_id, t_attachable_items.from_outer cover_image_from_outer, t_attachable_items.alias cover_image_alias,
t_attachable_items.url cover_image_url, t_attachable_items.size cover_image_size, t_attachable_items.mime_type_id cover_image_mime_type_id
FROM m_chat_rooms
LEFT JOIN t_images ON m_chat_rooms.cover_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE
	CASE WHEN @where_in_owner::boolean = true THEN owner_id = ANY(@in_owner::uuid[]) ELSE TRUE END
AND
	CASE WHEN @where_is_private::boolean = true THEN is_private = @is_private ELSE TRUE END
AND
	CASE WHEN @where_is_from_organization::boolean = true THEN from_organization = @is_from_organization ELSE TRUE END
AND
	CASE WHEN @where_from_organizations::boolean = true THEN (SELECT chat_room_id FROM m_organizations WHERE organization_id = ANY(@from_organizations::uuid[])) = chat_room_id ELSE TRUE END
ORDER BY
	m_chat_rooms_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetChatRoomsWithCoverImageUseKeysetPaginate :many
SELECT m_chat_rooms.*, t_images.height cover_image_height, t_images.width cover_image_width, t_images.attachable_item_id cover_image_attachable_item_id,
t_attachable_items.owner_id cover_image_owner_id, t_attachable_items.from_outer cover_image_from_outer, t_attachable_items.alias cover_image_alias,
t_attachable_items.url cover_image_url, t_attachable_items.size cover_image_size, t_attachable_items.mime_type_id cover_image_mime_type_id
FROM m_chat_rooms
LEFT JOIN t_images ON m_chat_rooms.cover_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE
	CASE WHEN @where_in_owner::boolean = true THEN owner_id = ANY(@in_owner::uuid[]) ELSE TRUE END
AND
	CASE WHEN @where_is_private::boolean = true THEN is_private = @is_private ELSE TRUE END
AND
	CASE WHEN @where_is_from_organization::boolean = true THEN from_organization = @is_from_organization ELSE TRUE END
AND
	CASE WHEN @where_from_organizations::boolean = true THEN (SELECT chat_room_id FROM m_organizations WHERE organization_id = ANY(@from_organizations::uuid[])) = chat_room_id ELSE TRUE END
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			m_chat_rooms_pkey > @cursor::int
		WHEN 'prev' THEN
			m_chat_rooms_pkey < @cursor::int
	END
ORDER BY
	CASE WHEN @cursor_direction::text = 'next' THEN m_chat_rooms_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN m_chat_rooms_pkey END DESC
LIMIT $1;

-- name: GetPluralChatRoomsWithCoverImage :many
SELECT m_chat_rooms.*, t_images.height cover_image_height, t_images.width cover_image_width, t_images.attachable_item_id cover_image_attachable_item_id,
t_attachable_items.owner_id cover_image_owner_id, t_attachable_items.from_outer cover_image_from_outer, t_attachable_items.alias cover_image_alias,
t_attachable_items.url cover_image_url, t_attachable_items.size cover_image_size, t_attachable_items.mime_type_id cover_image_mime_type_id
FROM m_chat_rooms
LEFT JOIN t_images ON m_chat_rooms.cover_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE chat_room_id = ANY(@chat_room_ids::uuid[])
ORDER BY
	m_chat_rooms_pkey ASC;

-- name: GetPluralChatRoomsWithCoverImageUseNumberedPaginate :many
SELECT m_chat_rooms.*, t_images.height cover_image_height, t_images.width cover_image_width, t_images.attachable_item_id cover_image_attachable_item_id,
t_attachable_items.owner_id cover_image_owner_id, t_attachable_items.from_outer cover_image_from_outer, t_attachable_items.alias cover_image_alias,
t_attachable_items.url cover_image_url, t_attachable_items.size cover_image_size, t_attachable_items.mime_type_id cover_image_mime_type_id
FROM m_chat_rooms
LEFT JOIN t_images ON m_chat_rooms.cover_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE chat_room_id = ANY(@chat_room_ids::uuid[])
ORDER BY
	m_chat_rooms_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetChatRoomsWithOwner :many
SELECT m_chat_rooms.*, m_members.name owner_name, m_members.first_name owner_first_name, m_members.last_name owner_last_name, m_members.email owner_email,
m_members.profile_image_id owner_profile_image_id, m_members.grade_id owner_grade_id, m_members.group_id owner_group_id
FROM m_chat_rooms
LEFT JOIN m_members ON m_chat_rooms.owner_id = m_members.member_id
WHERE
	CASE WHEN @where_in_owner::boolean THEN owner_id = ANY(@in_owner::uuid[]) ELSE TRUE END
AND
	CASE WHEN @where_is_private::boolean THEN is_private = @is_private ELSE TRUE END
AND
	CASE WHEN @where_is_from_organization::boolean = true THEN from_organization = @is_from_organization ELSE TRUE END
AND
	CASE WHEN @where_from_organizations::boolean = true THEN (SELECT chat_room_id FROM m_organizations WHERE organization_id = ANY(@from_organizations::uuid[])) = chat_room_id ELSE TRUE END
ORDER BY
	m_chat_rooms_pkey ASC;

-- name: GetChatRoomsWithOwnerUseNumberedPaginate :many
SELECT m_chat_rooms.*, m_members.name owner_name, m_members.first_name owner_first_name, m_members.last_name owner_last_name, m_members.email owner_email,
m_members.profile_image_id owner_profile_image_id, m_members.grade_id owner_grade_id, m_members.group_id owner_group_id
FROM m_chat_rooms
LEFT JOIN m_members ON m_chat_rooms.owner_id = m_members.member_id
WHERE
	CASE WHEN @where_in_owner::boolean THEN owner_id = ANY(@in_owner::uuid[]) ELSE TRUE END
AND
	CASE WHEN @where_is_private::boolean THEN is_private = @is_private ELSE TRUE END
AND
	CASE WHEN @where_is_from_organization::boolean = true THEN from_organization = @is_from_organization ELSE TRUE END
AND
	CASE WHEN @where_from_organizations::boolean = true THEN (SELECT chat_room_id FROM m_organizations WHERE organization_id = ANY(@from_organizations::uuid[])) = chat_room_id ELSE TRUE END
ORDER BY
	m_chat_rooms_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetChatRoomsWithOwnerUseKeysetPaginate :many
SELECT m_chat_rooms.*, m_members.name owner_name, m_members.first_name owner_first_name, m_members.last_name owner_last_name, m_members.email owner_email,
m_members.profile_image_id owner_profile_image_id, m_members.grade_id owner_grade_id, m_members.group_id owner_group_id
FROM m_chat_rooms
LEFT JOIN m_members ON m_chat_rooms.owner_id = m_members.member_id
WHERE
	CASE WHEN @where_in_owner::boolean = true THEN owner_id = ANY(@in_owner::uuid[]) ELSE TRUE END
AND
	CASE WHEN @where_is_private::boolean = true THEN is_private = @is_private ELSE TRUE END
AND
	CASE WHEN @where_is_from_organization::boolean = true THEN from_organization = @is_from_organization ELSE TRUE END
AND
	CASE WHEN @where_from_organizations::boolean = true THEN (SELECT chat_room_id FROM m_organizations WHERE organization_id = ANY(@from_organizations::uuid[])) = chat_room_id ELSE TRUE END
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			m_chat_rooms_pkey > @cursor::int
		WHEN 'prev' THEN
			m_chat_rooms_pkey < @cursor::int
	END
ORDER BY
	CASE WHEN @cursor_direction::text = 'next' THEN m_chat_rooms_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN m_chat_rooms_pkey END DESC
LIMIT $1;

-- name: GetPluralChatRoomsWithOwner :many
SELECT m_chat_rooms.*, m_members.name owner_name, m_members.first_name owner_first_name, m_members.last_name owner_last_name, m_members.email owner_email,
m_members.profile_image_id owner_profile_image_id, m_members.grade_id owner_grade_id, m_members.group_id owner_group_id
FROM m_chat_rooms
LEFT JOIN m_members ON m_chat_rooms.owner_id = m_members.member_id
WHERE chat_room_id = ANY(@chat_room_ids::uuid[])
ORDER BY
	m_chat_rooms_pkey ASC;

-- name: GetPluralChatRoomsWithOwnerUseNumberedPaginate :many
SELECT m_chat_rooms.*, m_members.name owner_name, m_members.first_name owner_first_name, m_members.last_name owner_last_name, m_members.email owner_email,
m_members.profile_image_id owner_profile_image_id, m_members.grade_id owner_grade_id, m_members.group_id owner_group_id
FROM m_chat_rooms
LEFT JOIN m_members ON m_chat_rooms.owner_id = m_members.member_id
WHERE chat_room_id = ANY(@chat_room_ids::uuid[])
ORDER BY
	m_chat_rooms_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountChatRooms :one
SELECT COUNT(*) FROM m_chat_rooms
WHERE
	CASE WHEN @where_in_owner::boolean = true THEN owner_id = ANY(@in_owner::uuid[]) ELSE TRUE END
AND
	CASE WHEN @where_is_private::boolean = true THEN is_private = @is_private ELSE TRUE END
AND
	CASE WHEN @where_is_from_organization::boolean = true THEN from_organization = @is_from_organization ELSE TRUE END
AND
	CASE WHEN @where_from_organizations::boolean = true THEN (SELECT chat_room_id FROM m_organizations WHERE organization_id = ANY(@from_organizations::uuid[])) = chat_room_id ELSE TRUE END;
