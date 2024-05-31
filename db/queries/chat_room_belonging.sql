-- name: CreateChatRoomBelongings :copyfrom
INSERT INTO m_chat_room_belongings (member_id, chat_room_id, added_at) VALUES ($1, $2, $3);

-- name: CreateChatRoomBelonging :one
INSERT INTO m_chat_room_belongings (member_id, chat_room_id, added_at) VALUES ($1, $2, $3) RETURNING *;

-- name: DeleteChatRoomBelonging :execrows
DELETE FROM m_chat_room_belongings WHERE member_id = $1 AND chat_room_id = $2;

-- name: DeleteChatRoomBelongingsOnMember :execrows
DELETE FROM m_chat_room_belongings WHERE member_id = $1;

-- name: DeleteChatRoomBelongingsOnMembers :execrows
DELETE FROM m_chat_room_belongings WHERE member_id = ANY(@member_ids::uuid[]);

-- name: DeleteChatRoomBelongingsOnChatRoom :execrows
DELETE FROM m_chat_room_belongings WHERE chat_room_id = $1;

-- name: DeleteChatRoomBelongingsOnChatRooms :execrows
DELETE FROM m_chat_room_belongings WHERE chat_room_id = ANY(@chat_room_ids::uuid[]);

-- name: GetMembersOnChatRoom :many
SELECT m_chat_room_belongings.*, m_members.name member_name, m_members.first_name member_first_name, m_members.last_name member_last_name, m_members.email member_email,
m_members.profile_image_id member_profile_image_id, t_images.height member_profile_image_height,
t_images.width member_profile_image_width, t_images.attachable_item_id member_profile_image_attachable_item_id,
t_attachable_items.owner_id member_profile_image_owner_id, t_attachable_items.from_outer member_profile_image_from_outer, t_attachable_items.alias member_profile_image_alias,
t_attachable_items.url member_profile_image_url, t_attachable_items.size member_profile_image_size, t_attachable_items.mime_type_id member_profile_image_mime_type_id FROM m_chat_room_belongings
LEFT JOIN m_members ON m_chat_room_belongings.member_id = m_members.member_id
LEFT JOIN t_images ON m_members.profile_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE chat_room_id = $1
AND CASE
	WHEN @where_like_name::boolean = true THEN m_members.name LIKE '%' || @search_name::text || '%'
END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_members.name END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' THEN m_members.name END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'old_add' THEN m_chat_room_belongings.added_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'late_add' THEN m_chat_room_belongings.added_at END DESC NULLS LAST,
	m_chat_room_belongings_pkey ASC;

-- name: GetMembersOnChatRoomUseNumberedPaginate :many
SELECT m_chat_room_belongings.*, m_members.name member_name, m_members.first_name member_first_name, m_members.last_name member_last_name, m_members.email member_email,
m_members.profile_image_id member_profile_image_id, t_images.height member_profile_image_height,
t_images.width member_profile_image_width, t_images.attachable_item_id member_profile_image_attachable_item_id,
t_attachable_items.owner_id member_profile_image_owner_id, t_attachable_items.from_outer member_profile_image_from_outer, t_attachable_items.alias member_profile_image_alias,
t_attachable_items.url member_profile_image_url, t_attachable_items.size member_profile_image_size, t_attachable_items.mime_type_id member_profile_image_mime_type_id FROM m_chat_room_belongings
LEFT JOIN m_members ON m_chat_room_belongings.member_id = m_members.member_id
LEFT JOIN t_images ON m_members.profile_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE chat_room_id = $1
AND CASE
	WHEN @where_like_name::boolean = true THEN m_members.name LIKE '%' || @search_name::text || '%'
END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_members.name END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' THEN m_members.name END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'old_add' THEN m_chat_room_belongings.added_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'late_add' THEN m_chat_room_belongings.added_at END DESC NULLS LAST,
	m_chat_room_belongings_pkey ASC
LIMIT $2 OFFSET $3;

-- name: GetMembersOnChatRoomUseKeysetPaginate :many
SELECT m_chat_room_belongings.*, m_members.name member_name, m_members.first_name member_first_name, m_members.last_name member_last_name, m_members.email member_email,
m_members.profile_image_id member_profile_image_id, t_images.height member_profile_image_height,
t_images.width member_profile_image_width, t_images.attachable_item_id member_profile_image_attachable_item_id,
t_attachable_items.owner_id member_profile_image_owner_id, t_attachable_items.from_outer member_profile_image_from_outer, t_attachable_items.alias member_profile_image_alias,
t_attachable_items.url member_profile_image_url, t_attachable_items.size member_profile_image_size, t_attachable_items.mime_type_id member_profile_image_mime_type_id FROM m_chat_room_belongings
LEFT JOIN m_members ON m_chat_room_belongings.member_id = m_members.member_id
LEFT JOIN t_images ON m_members.profile_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE chat_room_id = $1
AND CASE
	WHEN @where_like_name::boolean = true THEN m_members.name LIKE '%' || @search_name::text || '%'
END
AND CASE @cursor_direction::text
	WHEN 'next' THEN
		CASE @order_method::text
			WHEN 'name' THEN m_members.name > @name_cursor OR (m_members.name = @name_cursor AND m_chat_room_belongings_pkey > @cursor::int)
			WHEN 'r_name' THEN m_members.name < @name_cursor OR (m_members.name = @name_cursor AND m_chat_room_belongings_pkey > @cursor::int)
			WHEN 'old_add' THEN m_chat_room_belongings.added_at > @added_at_cursor OR (m_chat_room_belongings.added_at = @added_at_cursor AND m_chat_room_belongings_pkey > @cursor::int)
			WHEN 'late_add' THEN m_chat_room_belongings.added_at < @added_at_cursor OR (m_chat_room_belongings.added_at = @added_at_cursor AND m_chat_room_belongings_pkey > @cursor::int)
			ELSE m_chat_room_belongings_pkey > @cursor::int
		END
	WHEN 'prev' THEN
		CASE @order_method::text
			WHEN 'name' THEN m_members.name < @name_cursor OR (m_members.name = @name_cursor AND m_chat_room_belongings_pkey < @cursor::int)
			WHEN 'r_name' THEN m_members.name > @name_cursor OR (m_members.name = @name_cursor AND m_chat_room_belongings_pkey < @cursor::int)
			WHEN 'old_add' THEN m_chat_room_belongings.added_at < @added_at_cursor OR (m_chat_room_belongings.added_at = @added_at_cursor AND m_chat_room_belongings_pkey < @cursor::int)
			WHEN 'late_add' THEN m_chat_room_belongings.added_at > @added_at_cursor OR (m_chat_room_belongings.added_at = @added_at_cursor AND m_chat_room_belongings_pkey < @cursor::int)
			ELSE m_chat_room_belongings_pkey < @cursor::int
		END
END
ORDER BY
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'next' THEN m_members.name END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'prev' THEN m_members.name END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'next' THEN m_members.name END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'prev' THEN m_members.name END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'old_add' AND @cursor_direction::text = 'next' THEN m_chat_room_belongings.added_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'old_add' AND @cursor_direction::text = 'prev' THEN m_chat_room_belongings.added_at END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'late_add' AND @cursor_direction::text = 'next' THEN m_chat_room_belongings.added_at END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'late_add' AND @cursor_direction::text = 'prev' THEN m_chat_room_belongings.added_at END ASC NULLS LAST,
	CASE WHEN @cursor_direction::text = 'next' THEN m_chat_room_belongings_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN m_chat_room_belongings_pkey END DESC
LIMIT $2;

-- name: GetPluralMembersOnChatRoom :many
SELECT m_chat_room_belongings.*, m_members.name member_name, m_members.first_name member_first_name, m_members.last_name member_last_name, m_members.email member_email,
m_members.profile_image_id member_profile_image_id, t_images.height member_profile_image_height,
t_images.width member_profile_image_width, t_images.attachable_item_id member_profile_image_attachable_item_id,
t_attachable_items.owner_id member_profile_image_owner_id, t_attachable_items.from_outer member_profile_image_from_outer, t_attachable_items.alias member_profile_image_alias,
t_attachable_items.url member_profile_image_url, t_attachable_items.size member_profile_image_size, t_attachable_items.mime_type_id member_profile_image_mime_type_id FROM m_chat_room_belongings
LEFT JOIN m_members ON m_chat_room_belongings.member_id = m_members.member_id
LEFT JOIN t_images ON m_members.profile_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE chat_room_id = ANY(@chat_room_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_members.name END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' THEN m_members.name END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'old_add' THEN m_chat_room_belongings.added_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'late_add' THEN m_chat_room_belongings.added_at END DESC NULLS LAST,
	m_chat_room_belongings_pkey ASC;

-- name: GetPluralMembersOnChatRoomUseNumberedPaginate :many
SELECT m_chat_room_belongings.*, m_members.name member_name, m_members.first_name member_first_name, m_members.last_name member_last_name, m_members.email member_email,
m_members.profile_image_id member_profile_image_id, t_images.height member_profile_image_height,
t_images.width member_profile_image_width, t_images.attachable_item_id member_profile_image_attachable_item_id,
t_attachable_items.owner_id member_profile_image_owner_id, t_attachable_items.from_outer member_profile_image_from_outer, t_attachable_items.alias member_profile_image_alias,
t_attachable_items.url member_profile_image_url, t_attachable_items.size member_profile_image_size, t_attachable_items.mime_type_id member_profile_image_mime_type_id FROM m_chat_room_belongings
LEFT JOIN m_members ON m_chat_room_belongings.member_id = m_members.member_id
LEFT JOIN t_images ON m_members.profile_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE chat_room_id = ANY(@chat_room_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_members.name END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' THEN m_members.name END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'old_add' THEN m_chat_room_belongings.added_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'late_add' THEN m_chat_room_belongings.added_at END DESC NULLS LAST,
	m_chat_room_belongings_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountMembersOnChatRoom :one
SELECT COUNT(*) FROM m_chat_room_belongings WHERE chat_room_id = $1
AND CASE WHEN @where_like_name::boolean = true THEN
		EXISTS (SELECT 1 FROM m_members WHERE m_chat_room_belongings.member_id = m_members.member_id AND m_members.name LIKE '%' || @search_name::text || '%')
	ELSE TRUE END;

-- name: GetChatRoomsOnMember :many
SELECT m_chat_room_belongings.*, m_chat_rooms.name chat_room_name, m_chat_rooms.is_private chat_room_is_private,
m_chat_rooms.from_organization chat_room_from_organization, m_chat_rooms.owner_id chat_room_owner_id,
m_chat_rooms.cover_image_id chat_room_cover_image_id, t_images.height chat_room_cover_image_height,
t_images.width chat_room_cover_image_width, t_images.attachable_item_id chat_room_cover_image_attachable_item_id,
t_attachable_items.owner_id chat_room_cover_image_owner_id, t_attachable_items.from_outer chat_room_cover_image_from_outer, t_attachable_items.alias chat_room_cover_image_alias,
t_attachable_items.url chat_room_cover_image_url, t_attachable_items.size chat_room_cover_image_size, t_attachable_items.mime_type_id chat_room_cover_image_mime_type_id,
latest_message.message_id chat_room_latest_message_id, latest_message.posted_at chat_room_latest_message_posted_at, latest_message.body chat_room_latest_message_body
FROM m_chat_room_belongings
LEFT JOIN m_chat_rooms ON m_chat_room_belongings.chat_room_id = m_chat_rooms.chat_room_id
LEFT JOIN t_images ON m_chat_rooms.cover_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN (
    SELECT sub.* FROM (
        SELECT
            cra.chat_room_id AS chat_room_id, msg.*, ROW_NUMBER() OVER (PARTITION BY cra.chat_room_id ORDER BY msg.posted_at DESC, msg.t_messages_pkey DESC) AS rn
        FROM t_messages msg
        INNER JOIN t_chat_room_actions cra ON msg.chat_room_action_id = cra.chat_room_action_id
    ) sub
    WHERE sub.rn = 1
) latest_message ON m_chat_rooms.chat_room_id = latest_message.chat_room_id
WHERE member_id = $1
AND CASE
	WHEN @where_like_name::boolean = true THEN m_chat_rooms.name LIKE '%' || @search_name::text || '%'
END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_chat_rooms.name END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' THEN m_chat_rooms.name END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'old_add' THEN m_chat_room_belongings.added_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'late_add' THEN m_chat_room_belongings.added_at END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'old_chat' THEN latest_message.posted_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'late_chat' THEN latest_message.posted_at END DESC NULLS LAST,
	latest_message.posted_at DESC NULLS LAST,
	m_chat_room_belongings_pkey ASC;

-- name: GetChatRoomsOnMemberUseNumberedPaginate :many
SELECT m_chat_room_belongings.*, m_chat_rooms.name chat_room_name, m_chat_rooms.is_private chat_room_is_private,
m_chat_rooms.from_organization chat_room_from_organization, m_chat_rooms.owner_id chat_room_owner_id,
m_chat_rooms.cover_image_id chat_room_cover_image_id, t_images.height chat_room_cover_image_height,
t_images.width chat_room_cover_image_width, t_images.attachable_item_id chat_room_cover_image_attachable_item_id,
t_attachable_items.owner_id chat_room_cover_image_owner_id, t_attachable_items.from_outer chat_room_cover_image_from_outer, t_attachable_items.alias chat_room_cover_image_alias,
t_attachable_items.url chat_room_cover_image_url, t_attachable_items.size chat_room_cover_image_size, t_attachable_items.mime_type_id chat_room_cover_image_mime_type_id,
latest_message.message_id chat_room_latest_message_id, latest_message.posted_at chat_room_latest_message_posted_at, latest_message.body chat_room_latest_message_body
FROM m_chat_room_belongings
LEFT JOIN m_chat_rooms ON m_chat_room_belongings.chat_room_id = m_chat_rooms.chat_room_id
LEFT JOIN t_images ON m_chat_rooms.cover_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN (
    SELECT sub.* FROM (
        SELECT
            cra.chat_room_id AS chat_room_id, msg.*, ROW_NUMBER() OVER (PARTITION BY cra.chat_room_id ORDER BY msg.posted_at DESC, msg.t_messages_pkey DESC) AS rn
        FROM t_messages msg
        INNER JOIN t_chat_room_actions cra ON msg.chat_room_action_id = cra.chat_room_action_id
    ) sub
    WHERE sub.rn = 1
) latest_message ON m_chat_rooms.chat_room_id = latest_message.chat_room_id
WHERE member_id = $1
AND CASE
	WHEN @where_like_name::boolean = true THEN m_chat_rooms.name LIKE '%' || @search_name::text || '%'
END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_chat_rooms.name END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' THEN m_chat_rooms.name END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'old_add' THEN m_chat_room_belongings.added_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'late_add' THEN m_chat_room_belongings.added_at END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'old_chat' THEN latest_message.posted_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'late_chat' THEN latest_message.posted_at END DESC NULLS LAST,
	latest_message.posted_at DESC NULLS LAST,
	m_chat_room_belongings_pkey ASC
LIMIT $2 OFFSET $3;

-- name: GetChatRoomsOnMemberUseKeysetPaginate :many
SELECT m_chat_room_belongings.*, m_chat_rooms.name chat_room_name, m_chat_rooms.is_private chat_room_is_private,
m_chat_rooms.from_organization chat_room_from_organization, m_chat_rooms.owner_id chat_room_owner_id,
m_chat_rooms.cover_image_id chat_room_cover_image_id, t_images.height chat_room_cover_image_height,
t_images.width chat_room_cover_image_width, t_images.attachable_item_id chat_room_cover_image_attachable_item_id,
t_attachable_items.owner_id chat_room_cover_image_owner_id, t_attachable_items.from_outer chat_room_cover_image_from_outer, t_attachable_items.alias chat_room_cover_image_alias,
t_attachable_items.url chat_room_cover_image_url, t_attachable_items.size chat_room_cover_image_size, t_attachable_items.mime_type_id chat_room_cover_image_mime_type_id,
latest_message.message_id chat_room_latest_message_id, latest_message.posted_at chat_room_latest_message_posted_at, latest_message.body chat_room_latest_message_body
FROM m_chat_room_belongings
LEFT JOIN m_chat_rooms ON m_chat_room_belongings.chat_room_id = m_chat_rooms.chat_room_id
LEFT JOIN t_images ON m_chat_rooms.cover_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN (
    SELECT sub.* FROM (
        SELECT
            cra.chat_room_id AS chat_room_id, msg.*, ROW_NUMBER() OVER (PARTITION BY cra.chat_room_id ORDER BY msg.posted_at DESC, msg.t_messages_pkey DESC) AS rn
        FROM t_messages msg
        INNER JOIN t_chat_room_actions cra ON msg.chat_room_action_id = cra.chat_room_action_id
    ) sub
    WHERE sub.rn = 1
) latest_message ON m_chat_rooms.chat_room_id = latest_message.chat_room_id
WHERE member_id = $1
AND CASE
	WHEN @where_like_name::boolean = true THEN m_chat_rooms.name LIKE '%' || @search_name::text || '%'
END
AND CASE @cursor_direction::text
	WHEN 'next' THEN
		CASE @order_method::text
			WHEN 'name' THEN m_chat_rooms.name > @name_cursor OR (m_chat_rooms.name = @name_cursor AND m_chat_room_belongings_pkey > @cursor::int)
			WHEN 'r_name' THEN m_chat_rooms.name < @name_cursor OR (m_chat_rooms.name = @name_cursor AND m_chat_room_belongings_pkey > @cursor::int)
			WHEN 'old_add' THEN m_chat_room_belongings.added_at > @add_cursor OR (m_chat_room_belongings.added_at = @add_cursor AND m_chat_room_belongings_pkey > @cursor::int)
			WHEN 'late_add' THEN m_chat_room_belongings.added_at < @add_cursor OR (m_chat_room_belongings.added_at = @add_cursor AND m_chat_room_belongings_pkey > @cursor::int)
			WHEN 'old_chat' THEN latest_message.posted_at > @chat_cursor::timestamptz OR (latest_message.posted_at = @chat_cursor::timestamptz AND m_chat_room_belongings_pkey > @cursor::int)
			WHEN 'late_chat' THEN latest_message.posted_at < @chat_cursor::timestamptz OR (latest_message.posted_at = @chat_cursor::timestamptz AND m_chat_room_belongings_pkey > @cursor::int)
			ELSE m_chat_room_belongings_pkey > @cursor::int
		END
	WHEN 'prev' THEN
		CASE @order_method::text
			WHEN 'name' THEN m_chat_rooms.name < @name_cursor OR (m_chat_rooms.name = @name_cursor AND m_chat_room_belongings_pkey < @cursor::int)
			WHEN 'r_name' THEN m_chat_rooms.name > @name_cursor OR (m_chat_rooms.name = @name_cursor AND m_chat_room_belongings_pkey < @cursor::int)
			WHEN 'old_add' THEN m_chat_room_belongings.added_at < @add_cursor OR (m_chat_room_belongings.added_at = @add_cursor AND m_chat_room_belongings_pkey < @cursor::int)
			WHEN 'late_add' THEN m_chat_room_belongings.added_at > @add_cursor OR (m_chat_room_belongings.added_at = @add_cursor AND m_chat_room_belongings_pkey < @cursor::int)
			WHEN 'old_chat' THEN latest_message.posted_at < @chat_cursor::timestamptz OR (latest_message.posted_at = @chat_cursor::timestamptz AND m_chat_room_belongings_pkey < @cursor::int)
			WHEN 'late_chat' THEN latest_message.posted_at > @chat_cursor::timestamptz OR (latest_message.posted_at = @chat_cursor::timestamptz AND m_chat_room_belongings_pkey < @cursor::int)
			ELSE m_chat_room_belongings_pkey < @cursor::int
		END
END
ORDER BY
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'next' THEN m_chat_rooms.name END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'prev' THEN m_chat_rooms.name END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'next' THEN m_chat_rooms.name END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'prev' THEN m_chat_rooms.name END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'old_add' AND @cursor_direction::text = 'next' THEN m_chat_room_belongings.added_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'old_add' AND @cursor_direction::text = 'prev' THEN m_chat_room_belongings.added_at END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'late_add' AND @cursor_direction::text = 'next' THEN m_chat_room_belongings.added_at END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'late_add' AND @cursor_direction::text = 'prev' THEN m_chat_room_belongings.added_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'old_chat' AND @cursor_direction::text = 'next' THEN latest_message.posted_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'old_chat' AND @cursor_direction::text = 'prev' THEN latest_message.posted_at END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'late_chat' AND @cursor_direction::text = 'next' THEN latest_message.posted_at END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'late_chat' AND @cursor_direction::text = 'prev' THEN latest_message.posted_at END ASC NULLS LAST,
	CASE WHEN @cursor_direction::text = 'next' THEN m_chat_room_belongings_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN m_chat_room_belongings_pkey END DESC
LIMIT $2;

-- name: GetPluralChatRoomsOnMember :many
SELECT m_chat_room_belongings.*, m_chat_rooms.name chat_room_name, m_chat_rooms.is_private chat_room_is_private,
m_chat_rooms.from_organization chat_room_from_organization, m_chat_rooms.owner_id chat_room_owner_id,
m_chat_rooms.cover_image_id chat_room_cover_image_id, t_images.height chat_room_cover_image_height,
t_images.width chat_room_cover_image_width, t_images.attachable_item_id chat_room_cover_image_attachable_item_id,
t_attachable_items.owner_id chat_room_cover_image_owner_id, t_attachable_items.from_outer chat_room_cover_image_from_outer, t_attachable_items.alias chat_room_cover_image_alias,
t_attachable_items.url chat_room_cover_image_url, t_attachable_items.size chat_room_cover_image_size, t_attachable_items.mime_type_id chat_room_cover_image_mime_type_id,
latest_message.message_id chat_room_latest_message_id, latest_message.posted_at chat_room_latest_message_posted_at, latest_message.body chat_room_latest_message_body
FROM m_chat_room_belongings
LEFT JOIN m_chat_rooms ON m_chat_room_belongings.chat_room_id = m_chat_rooms.chat_room_id
LEFT JOIN t_images ON m_chat_rooms.cover_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN (
    SELECT sub.* FROM (
        SELECT
            cra.chat_room_id AS chat_room_id, msg.*, ROW_NUMBER() OVER (PARTITION BY cra.chat_room_id ORDER BY msg.posted_at DESC, msg.t_messages_pkey DESC) AS rn
        FROM t_messages msg
        INNER JOIN t_chat_room_actions cra ON msg.chat_room_action_id = cra.chat_room_action_id
    ) sub
    WHERE sub.rn = 1
) latest_message ON m_chat_rooms.chat_room_id = latest_message.chat_room_id
WHERE member_id = ANY(@member_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_chat_rooms.name END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' THEN m_chat_rooms.name END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'old_add' THEN m_chat_room_belongings.added_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'late_add' THEN m_chat_room_belongings.added_at END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'old_chat' THEN latest_message.posted_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'late_chat' THEN latest_message.posted_at END DESC NULLS LAST,
	m_chat_room_belongings_pkey ASC;

-- name: GetPluralChatRoomsOnMemberUseNumberedPaginate :many
SELECT m_chat_room_belongings.*, m_chat_rooms.name chat_room_name, m_chat_rooms.is_private chat_room_is_private,
m_chat_rooms.from_organization chat_room_from_organization, m_chat_rooms.owner_id chat_room_owner_id,
m_chat_rooms.cover_image_id chat_room_cover_image_id, t_images.height chat_room_cover_image_height,
t_images.width chat_room_cover_image_width, t_images.attachable_item_id chat_room_cover_image_attachable_item_id,
t_attachable_items.owner_id chat_room_cover_image_owner_id, t_attachable_items.from_outer chat_room_cover_image_from_outer, t_attachable_items.alias chat_room_cover_image_alias,
t_attachable_items.url chat_room_cover_image_url, t_attachable_items.size chat_room_cover_image_size, t_attachable_items.mime_type_id chat_room_cover_image_mime_type_id,
latest_message.message_id chat_room_latest_message_id, latest_message.posted_at chat_room_latest_message_posted_at, latest_message.body chat_room_latest_message_body
FROM m_chat_room_belongings
LEFT JOIN m_chat_rooms ON m_chat_room_belongings.chat_room_id = m_chat_rooms.chat_room_id
LEFT JOIN t_images ON m_chat_rooms.cover_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN (
    SELECT sub.* FROM (
        SELECT
            cra.chat_room_id AS chat_room_id, msg.*, ROW_NUMBER() OVER (PARTITION BY cra.chat_room_id ORDER BY msg.posted_at DESC, msg.t_messages_pkey DESC) AS rn
        FROM t_messages msg
        INNER JOIN t_chat_room_actions cra ON msg.chat_room_action_id = cra.chat_room_action_id
    ) sub
    WHERE sub.rn = 1
) latest_message ON m_chat_rooms.chat_room_id = latest_message.chat_room_id
WHERE member_id = ANY(@member_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_chat_rooms.name END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' THEN m_chat_rooms.name END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'old_add' THEN m_chat_room_belongings.added_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'late_add' THEN m_chat_room_belongings.added_at END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'old_chat' THEN latest_message.posted_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'late_chat' THEN latest_message.posted_at END DESC NULLS LAST,
	m_chat_room_belongings_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountChatRoomsOnMember :one
SELECT COUNT(*) FROM m_chat_room_belongings WHERE member_id = $1
AND CASE WHEN @where_like_name::boolean = true THEN
		EXISTS (SELECT 1 FROM m_chat_rooms WHERE m_chat_room_belongings.chat_room_id = m_chat_rooms.chat_room_id AND m_chat_rooms.name LIKE '%' || @search_name::text || '%')
	ELSE TRUE END;
