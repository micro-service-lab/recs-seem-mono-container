-- name: CreateMessages :copyfrom
INSERT INTO t_messages (chat_room_action_id, sender_id, body, posted_at, last_edited_at) VALUES ($1, $2, $3, $4, $5);

-- name: CreateMessage :one
INSERT INTO t_messages (chat_room_action_id, sender_id, body, posted_at, last_edited_at) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: UpdateMessage :one
UPDATE t_messages SET body = $2, last_edited_at = $3 WHERE message_id = $1 RETURNING *;

-- name: DeleteMessage :execrows
DELETE FROM t_messages WHERE message_id = $1;

-- name: DeleteMessagesOnChatRoom :execrows
DELETE FROM t_messages WHERE (SELECT chat_room_id FROM t_chat_room_actions WHERE t_chat_room_actions.chat_room_action_id = t_messages.chat_room_action_id) = $1;

-- name: PluralDeleteMessages :execrows
DELETE FROM t_messages WHERE message_id = ANY(@member_ids::uuid[]);

-- name: FindMessageByID :one
SELECT * FROM t_messages WHERE message_id = $1;

-- name: FindMessageByIDWithChatRoom :one
SELECT t_messages.*, m_chat_rooms.chat_room_id, m_chat_rooms.name chat_room_name, m_chat_rooms.is_private chat_room_is_private,
m_chat_rooms.from_organization chat_room_from_organization, m_chat_rooms.owner_id chat_room_owner_id,
m_chat_rooms.cover_image_id chat_room_cover_image_id, t_images.height chat_room_cover_image_height,
t_images.width chat_room_cover_image_width, t_images.attachable_item_id chat_room_cover_image_attachable_item_id,
t_attachable_items.owner_id chat_room_cover_image_owner_id, t_attachable_items.from_outer chat_room_cover_image_from_outer, t_attachable_items.alias chat_room_cover_image_alias,
t_attachable_items.url chat_room_cover_image_url, t_attachable_items.size chat_room_cover_image_size, t_attachable_items.mime_type_id chat_room_cover_image_mime_type_id FROM t_messages
LEFT JOIN t_chat_room_actions ON t_messages.chat_room_action_id = t_chat_room_actions.chat_room_action_id
LEFT JOIN m_chat_rooms ON t_messages.chat_room_id = m_chat_rooms.chat_room_id
LEFT JOIN t_images ON m_chat_rooms.cover_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE message_id = $1;

-- name: FindMessageByIDWithSender :one
SELECT t_messages.*, m_members.name member_name, m_members.first_name member_first_name, m_members.last_name member_last_name, m_members.email member_email, m_members.grade_id member_grade_id, m_members.group_id member_group_id,
m_members.profile_image_id member_profile_image_id, t_images.height member_profile_image_height,
t_images.width member_profile_image_width, t_images.attachable_item_id member_profile_image_attachable_item_id,
t_attachable_items.owner_id member_profile_image_owner_id, t_attachable_items.from_outer member_profile_image_from_outer, t_attachable_items.alias member_profile_image_alias,
t_attachable_items.url member_profile_image_url, t_attachable_items.size member_profile_image_size, t_attachable_items.mime_type_id member_profile_image_mime_type_id FROM t_messages
LEFT JOIN m_members ON t_messages.sender_id = m_members.member_id
LEFT JOIN t_images ON m_members.profile_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE message_id = $1;

-- name: FindMessageByIDWithChatRoomAction :one
SELECT t_messages.*, t_chat_room_actions.chat_room_id chat_room_action_chat_room_id, t_chat_room_actions.chat_room_action_type_id chat_room_action_action_type_id, t_chat_room_actions.acted_at chat_room_action_acted_at
FROM t_messages
LEFT JOIN t_chat_room_actions ON t_messages.chat_room_action_id = t_chat_room_actions.chat_room_action_id
WHERE message_id = $1;

-- name: GetMessages :many
SELECT * FROM t_messages
WHERE
	CASE WHEN @where_in_chat_room::boolean = true THEN chat_room_id = ANY(@in_chat_room::uuid[]) ELSE TRUE END
AND
	CASE WHEN @where_in_sender::boolean = true THEN sender_id = ANY(@in_sender::uuid[]) ELSE TRUE END
AND
	CASE WHEN @where_like_body::boolean = true THEN body LIKE '%' || @search_body::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_earlier_posted_at::boolean = true THEN posted_at >= @earlier_posted_at ELSE TRUE END
AND
	CASE WHEN @where_later_posted_at::boolean = true THEN posted_at <= @later_posted_at ELSE TRUE END
AND
	CASE WHEN @where_earlier_last_edited_at::boolean = true THEN last_edited_at >= @earlier_last_edited_at ELSE TRUE END
AND
	CASE WHEN @where_later_last_edited_at::boolean = true THEN last_edited_at <= @later_last_edited_at ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'posted_at' THEN posted_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_posted_at' THEN posted_at END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'last_edited_at' THEN last_edited_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_last_edited_at' THEN last_edited_at END DESC NULLS LAST,
	t_messages_pkey ASC;

-- name: GetMessagesUseNumberedPaginate :many
SELECT * FROM t_messages
WHERE
	CASE WHEN @where_in_chat_room::boolean = true THEN chat_room_id = ANY(@in_chat_room::uuid[]) ELSE TRUE END
AND
	CASE WHEN @where_in_sender::boolean = true THEN sender_id = ANY(@in_sender::uuid[]) ELSE TRUE END
AND
	CASE WHEN @where_like_body::boolean = true THEN body LIKE '%' || @search_body::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_earlier_posted_at::boolean = true THEN posted_at >= @earlier_posted_at ELSE TRUE END
AND
	CASE WHEN @where_later_posted_at::boolean = true THEN posted_at <= @later_posted_at ELSE TRUE END
AND
	CASE WHEN @where_earlier_last_edited_at::boolean = true THEN last_edited_at >= @earlier_last_edited_at ELSE TRUE END
AND
	CASE WHEN @where_later_last_edited_at::boolean = true THEN last_edited_at <= @later_last_edited_at ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'posted_at' THEN posted_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_posted_at' THEN posted_at END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'last_edited_at' THEN last_edited_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_last_edited_at' THEN last_edited_at END DESC NULLS LAST,
	t_messages_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetMessagesUseKeysetPaginate :many
SELECT * FROM t_messages
WHERE
	CASE WHEN @where_in_chat_room::boolean = true THEN chat_room_id = ANY(@in_chat_room::uuid[]) ELSE TRUE END
AND
	CASE WHEN @where_in_sender::boolean = true THEN sender_id = ANY(@in_sender::uuid[]) ELSE TRUE END
AND
	CASE WHEN @where_like_body::boolean = true THEN body LIKE '%' || @search_body::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_earlier_posted_at::boolean = true THEN posted_at >= @earlier_posted_at ELSE TRUE END
AND
	CASE WHEN @where_later_posted_at::boolean = true THEN posted_at <= @later_posted_at ELSE TRUE END
AND
	CASE WHEN @where_earlier_last_edited_at::boolean = true THEN last_edited_at >= @earlier_last_edited_at ELSE TRUE END
AND
	CASE WHEN @where_later_last_edited_at::boolean = true THEN last_edited_at <= @later_last_edited_at ELSE TRUE END
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'posted_at' THEN posted_at > @posted_at_cursor OR (posted_at = @posted_at_cursor AND t_messages_pkey > @cursor::int)
				WHEN 'r_posted_at' THEN posted_at < @posted_at_cursor OR (posted_at = @posted_at_cursor AND t_messages_pkey > @cursor::int)
				WHEN 'last_edited_at' THEN last_edited_at > @last_edited_at_cursor OR (last_edited_at = @last_edited_at_cursor AND t_messages_pkey > @cursor::int)
				WHEN 'r_last_edited_at' THEN last_edited_at < @last_edited_at_cursor OR (last_edited_at = @last_edited_at_cursor AND t_messages_pkey > @cursor::int)
				ELSE t_messages_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'posted_at' THEN posted_at < @posted_at_cursor OR (posted_at = @posted_at_cursor AND t_messages_pkey < @cursor::int)
				WHEN 'r_posted_at' THEN posted_at > @posted_at_cursor OR (posted_at = @posted_at_cursor AND t_messages_pkey < @cursor::int)
				WHEN 'last_edited_at' THEN last_edited_at < @last_edited_at_cursor OR (last_edited_at = @last_edited_at_cursor AND t_messages_pkey < @cursor::int)
				WHEN 'r_last_edited_at' THEN last_edited_at > @last_edited_at_cursor OR (last_edited_at = @last_edited_at_cursor AND t_messages_pkey < @cursor::int)
				ELSE t_messages_pkey < @cursor::int
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'posted_at' AND @cursor_direction::text = 'next' THEN posted_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'posted_at' AND @cursor_direction::text = 'prev' THEN posted_at END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'r_posted_at' AND @cursor_direction::text = 'next' THEN posted_at END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'r_posted_at' AND @cursor_direction::text = 'prev' THEN posted_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'last_edited_at' AND @cursor_direction::text = 'next' THEN last_edited_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'last_edited_at' AND @cursor_direction::text = 'prev' THEN last_edited_at END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'r_last_edited_at' AND @cursor_direction::text = 'next' THEN last_edited_at END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'r_last_edited_at' AND @cursor_direction::text = 'prev' THEN last_edited_at END ASC NULLS LAST,
	CASE WHEN @cursor_direction::text = 'next' THEN t_messages_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN t_messages_pkey END DESC
LIMIT $1;

-- name: GetPluralMessages :many
SELECT * FROM t_messages WHERE message_id = ANY(@message_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'posted_at' THEN posted_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_posted_at' THEN posted_at END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'last_edited_at' THEN last_edited_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_last_edited_at' THEN last_edited_at END DESC NULLS LAST,
	t_messages_pkey ASC;

-- name: GetPluralMessagesUseNumberedPaginate :many
SELECT * FROM t_messages WHERE message_id = ANY(@message_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'posted_at' THEN posted_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_posted_at' THEN posted_at END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'last_edited_at' THEN last_edited_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_last_edited_at' THEN last_edited_at END DESC NULLS LAST,
	t_messages_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetMessagesWithChatRoom :many
SELECT t_messages.*, m_chat_rooms.chat_room_id, m_chat_rooms.name chat_room_name, m_chat_rooms.is_private chat_room_is_private,
m_chat_rooms.from_organization chat_room_from_organization, m_chat_rooms.owner_id chat_room_owner_id,
m_chat_rooms.cover_image_id chat_room_cover_image_id, t_images.height chat_room_cover_image_height,
t_images.width chat_room_cover_image_width, t_images.attachable_item_id chat_room_cover_image_attachable_item_id,
t_attachable_items.owner_id chat_room_cover_image_owner_id, t_attachable_items.from_outer chat_room_cover_image_from_outer, t_attachable_items.alias chat_room_cover_image_alias,
t_attachable_items.url chat_room_cover_image_url, t_attachable_items.size chat_room_cover_image_size, t_attachable_items.mime_type_id chat_room_cover_image_mime_type_id FROM t_messages
LEFT JOIN t_chat_room_actions ON t_messages.chat_room_action_id = t_chat_room_actions.chat_room_action_id
LEFT JOIN m_chat_rooms ON t_messages.chat_room_id = m_chat_rooms.chat_room_id
LEFT JOIN t_images ON m_chat_rooms.cover_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE
	CASE WHEN @where_in_chat_room::boolean = true THEN t_chat_room_actions in (SELECT chat_room_action_id FROM t_chat_room_actions WHERE chat_room_id = ANY(@in_chat_room::uuid[])) ELSE TRUE END
AND
	CASE WHEN @where_in_sender::boolean = true THEN sender_id = ANY(@in_sender::uuid[]) ELSE TRUE END
AND
	CASE WHEN @where_like_body::boolean = true THEN body LIKE '%' || @search_body::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_earlier_posted_at::boolean = true THEN posted_at >= @earlier_posted_at ELSE TRUE END
AND
	CASE WHEN @where_later_posted_at::boolean = true THEN posted_at <= @later_posted_at ELSE TRUE END
AND
	CASE WHEN @where_earlier_last_edited_at::boolean = true THEN last_edited_at >= @earlier_last_edited_at ELSE TRUE END
AND
	CASE WHEN @where_later_last_edited_at::boolean = true THEN last_edited_at <= @later_last_edited_at ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'posted_at' THEN posted_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_posted_at' THEN posted_at END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'last_edited_at' THEN last_edited_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_last_edited_at' THEN last_edited_at END DESC NULLS LAST,
	t_messages_pkey ASC;

-- name: GetMessagesWithChatRoomUseNumberedPaginate :many
SELECT t_messages.*, m_chat_rooms.chat_room_id, m_chat_rooms.name chat_room_name, m_chat_rooms.is_private chat_room_is_private,
m_chat_rooms.from_organization chat_room_from_organization, m_chat_rooms.owner_id chat_room_owner_id,
m_chat_rooms.cover_image_id chat_room_cover_image_id, t_images.height chat_room_cover_image_height,
t_images.width chat_room_cover_image_width, t_images.attachable_item_id chat_room_cover_image_attachable_item_id,
t_attachable_items.owner_id chat_room_cover_image_owner_id, t_attachable_items.from_outer chat_room_cover_image_from_outer, t_attachable_items.alias chat_room_cover_image_alias,
t_attachable_items.url chat_room_cover_image_url, t_attachable_items.size chat_room_cover_image_size, t_attachable_items.mime_type_id chat_room_cover_image_mime_type_id FROM t_messages
LEFT JOIN t_chat_room_actions ON t_messages.chat_room_action_id = t_chat_room_actions.chat_room_action_id
LEFT JOIN m_chat_rooms ON t_messages.chat_room_id = m_chat_rooms.chat_room_id
LEFT JOIN t_images ON m_chat_rooms.cover_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE
	CASE WHEN @where_in_chat_room::boolean = true THEN t_chat_room_actions in (SELECT chat_room_action_id FROM t_chat_room_actions WHERE chat_room_id = ANY(@in_chat_room::uuid[])) ELSE TRUE END
AND
	CASE WHEN @where_in_sender::boolean = true THEN sender_id = ANY(@in_sender::uuid[]) ELSE TRUE END
AND
	CASE WHEN @where_like_body::boolean = true THEN body LIKE '%' || @search_body::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_earlier_posted_at::boolean = true THEN posted_at >= @earlier_posted_at ELSE TRUE END
AND
	CASE WHEN @where_later_posted_at::boolean = true THEN posted_at <= @later_posted_at ELSE TRUE END
AND
	CASE WHEN @where_earlier_last_edited_at::boolean = true THEN last_edited_at >= @earlier_last_edited_at ELSE TRUE END
AND
	CASE WHEN @where_later_last_edited_at::boolean = true THEN last_edited_at <= @later_last_edited_at ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'posted_at' THEN posted_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_posted_at' THEN posted_at END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'last_edited_at' THEN last_edited_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_last_edited_at' THEN last_edited_at END DESC NULLS LAST,
	t_messages_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetMessagesWithChatRoomUseKeysetPaginate :many
SELECT t_messages.*, m_chat_rooms.chat_room_id, m_chat_rooms.name chat_room_name, m_chat_rooms.is_private chat_room_is_private,
m_chat_rooms.from_organization chat_room_from_organization, m_chat_rooms.owner_id chat_room_owner_id,
m_chat_rooms.cover_image_id chat_room_cover_image_id, t_images.height chat_room_cover_image_height,
t_images.width chat_room_cover_image_width, t_images.attachable_item_id chat_room_cover_image_attachable_item_id,
t_attachable_items.owner_id chat_room_cover_image_owner_id, t_attachable_items.from_outer chat_room_cover_image_from_outer, t_attachable_items.alias chat_room_cover_image_alias,
t_attachable_items.url chat_room_cover_image_url, t_attachable_items.size chat_room_cover_image_size, t_attachable_items.mime_type_id chat_room_cover_image_mime_type_id FROM t_messages
LEFT JOIN t_chat_room_actions ON t_messages.chat_room_action_id = t_chat_room_actions.chat_room_action_id
LEFT JOIN m_chat_rooms ON t_messages.chat_room_id = m_chat_rooms.chat_room_id
LEFT JOIN t_images ON m_chat_rooms.cover_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE
	CASE WHEN @where_in_chat_room::boolean = true THEN t_chat_room_actions in (SELECT chat_room_action_id FROM t_chat_room_actions WHERE chat_room_id = ANY(@in_chat_room::uuid[])) ELSE TRUE END
AND
	CASE WHEN @where_in_sender::boolean = true THEN sender_id = ANY(@in_sender::uuid[]) ELSE TRUE END
AND
	CASE WHEN @where_like_body::boolean = true THEN body LIKE '%' || @search_body::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_earlier_posted_at::boolean = true THEN posted_at >= @earlier_posted_at ELSE TRUE END
AND
	CASE WHEN @where_later_posted_at::boolean = true THEN posted_at <= @later_posted_at ELSE TRUE END
AND
	CASE WHEN @where_earlier_last_edited_at::boolean = true THEN last_edited_at >= @earlier_last_edited_at ELSE TRUE END
AND
	CASE WHEN @where_later_last_edited_at::boolean = true THEN last_edited_at <= @later_last_edited_at ELSE TRUE END
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'posted_at' THEN posted_at > @posted_at_cursor OR (posted_at = @posted_at_cursor AND t_messages_pkey > @cursor::int)
				WHEN 'r_posted_at' THEN posted_at < @posted_at_cursor OR (posted_at = @posted_at_cursor AND t_messages_pkey > @cursor::int)
				WHEN 'last_edited_at' THEN last_edited_at > @last_edited_at_cursor OR (last_edited_at = @last_edited_at_cursor AND t_messages_pkey > @cursor::int)
				WHEN 'r_last_edited_at' THEN last_edited_at < @last_edited_at_cursor OR (last_edited_at = @last_edited_at_cursor AND t_messages_pkey > @cursor::int)
				ELSE t_messages_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'posted_at' THEN posted_at < @posted_at_cursor OR (posted_at = @posted_at_cursor AND t_messages_pkey < @cursor::int)
				WHEN 'r_posted_at' THEN posted_at > @posted_at_cursor OR (posted_at = @posted_at_cursor AND t_messages_pkey < @cursor::int)
				WHEN 'last_edited_at' THEN last_edited_at < @last_edited_at_cursor OR (last_edited_at = @last_edited_at_cursor AND t_messages_pkey < @cursor::int)
				WHEN 'r_last_edited_at' THEN last_edited_at > @last_edited_at_cursor OR (last_edited_at = @last_edited_at_cursor AND t_messages_pkey < @cursor::int)
				ELSE t_messages_pkey < @cursor::int
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'posted_at' AND @cursor_direction::text = 'next' THEN posted_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'posted_at' AND @cursor_direction::text = 'prev' THEN posted_at END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'r_posted_at' AND @cursor_direction::text = 'next' THEN posted_at END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'r_posted_at' AND @cursor_direction::text = 'prev' THEN posted_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'last_edited_at' AND @cursor_direction::text = 'next' THEN last_edited_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'last_edited_at' AND @cursor_direction::text = 'prev' THEN last_edited_at END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'r_last_edited_at' AND @cursor_direction::text = 'next' THEN last_edited_at END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'r_last_edited_at' AND @cursor_direction::text = 'prev' THEN last_edited_at END ASC NULLS LAST,
	CASE WHEN @cursor_direction::text = 'next' THEN t_messages_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN t_messages_pkey END DESC
LIMIT $1;

-- name: GetPluralMessagesWithChatRoom :many
SELECT t_messages.*, m_chat_rooms.chat_room_id, m_chat_rooms.name chat_room_name, m_chat_rooms.is_private chat_room_is_private,
m_chat_rooms.from_organization chat_room_from_organization, m_chat_rooms.owner_id chat_room_owner_id,
m_chat_rooms.cover_image_id chat_room_cover_image_id, t_images.height chat_room_cover_image_height,
t_images.width chat_room_cover_image_width, t_images.attachable_item_id chat_room_cover_image_attachable_item_id,
t_attachable_items.owner_id chat_room_cover_image_owner_id, t_attachable_items.from_outer chat_room_cover_image_from_outer, t_attachable_items.alias chat_room_cover_image_alias,
t_attachable_items.url chat_room_cover_image_url, t_attachable_items.size chat_room_cover_image_size, t_attachable_items.mime_type_id chat_room_cover_image_mime_type_id FROM t_messages
LEFT JOIN t_chat_room_actions ON t_messages.chat_room_action_id = t_chat_room_actions.chat_room_action_id
LEFT JOIN m_chat_rooms ON t_messages.chat_room_id = m_chat_rooms.chat_room_id
LEFT JOIN t_images ON m_chat_rooms.cover_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE message_id = ANY(@message_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'posted_at' THEN posted_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_posted_at' THEN posted_at END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'last_edited_at' THEN last_edited_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_last_edited_at' THEN last_edited_at END DESC NULLS LAST,
	t_messages_pkey ASC;

-- name: GetPluralMessagesWithChatRoomUseNumberedPaginate :many
SELECT t_messages.*, m_chat_rooms.chat_room_id, m_chat_rooms.name chat_room_name, m_chat_rooms.is_private chat_room_is_private,
m_chat_rooms.from_organization chat_room_from_organization, m_chat_rooms.owner_id chat_room_owner_id,
m_chat_rooms.cover_image_id chat_room_cover_image_id, t_images.height chat_room_cover_image_height,
t_images.width chat_room_cover_image_width, t_images.attachable_item_id chat_room_cover_image_attachable_item_id,
t_attachable_items.owner_id chat_room_cover_image_owner_id, t_attachable_items.from_outer chat_room_cover_image_from_outer, t_attachable_items.alias chat_room_cover_image_alias,
t_attachable_items.url chat_room_cover_image_url, t_attachable_items.size chat_room_cover_image_size, t_attachable_items.mime_type_id chat_room_cover_image_mime_type_id FROM t_messages
LEFT JOIN t_chat_room_actions ON t_messages.chat_room_action_id = t_chat_room_actions.chat_room_action_id
LEFT JOIN m_chat_rooms ON t_messages.chat_room_id = m_chat_rooms.chat_room_id
LEFT JOIN t_images ON m_chat_rooms.cover_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE message_id = ANY(@message_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'posted_at' THEN posted_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_posted_at' THEN posted_at END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'last_edited_at' THEN last_edited_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_last_edited_at' THEN last_edited_at END DESC NULLS LAST,
	t_messages_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetMessagesWithSender :many
SELECT t_messages.*, m_members.name member_name, m_members.first_name member_first_name, m_members.last_name member_last_name, m_members.email member_email, m_members.grade_id member_grade_id, m_members.group_id member_group_id,
m_members.profile_image_id member_profile_image_id, t_images.height member_profile_image_height,
t_images.width member_profile_image_width, t_images.attachable_item_id member_profile_image_attachable_item_id,
t_attachable_items.owner_id member_profile_image_owner_id, t_attachable_items.from_outer member_profile_image_from_outer, t_attachable_items.alias member_profile_image_alias,
t_attachable_items.url member_profile_image_url, t_attachable_items.size member_profile_image_size, t_attachable_items.mime_type_id member_profile_image_mime_type_id FROM t_messages
LEFT JOIN m_members ON t_messages.sender_id = m_members.member_id
LEFT JOIN t_images ON m_members.profile_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE
	CASE WHEN @where_in_chat_room::boolean = true THEN t_chat_room_actions in (SELECT chat_room_action_id FROM t_chat_room_actions WHERE chat_room_id = ANY(@in_chat_room::uuid[])) ELSE TRUE END
AND
	CASE WHEN @where_in_sender::boolean = true THEN sender_id = ANY(@in_sender::uuid[]) ELSE TRUE END
AND
	CASE WHEN @where_like_body::boolean = true THEN body LIKE '%' || @search_body::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_earlier_posted_at::boolean = true THEN posted_at >= @earlier_posted_at ELSE TRUE END
AND
	CASE WHEN @where_later_posted_at::boolean = true THEN posted_at <= @later_posted_at ELSE TRUE END
AND
	CASE WHEN @where_earlier_last_edited_at::boolean = true THEN last_edited_at >= @earlier_last_edited_at ELSE TRUE END
AND
	CASE WHEN @where_later_last_edited_at::boolean = true THEN last_edited_at <= @later_last_edited_at ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'posted_at' THEN posted_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_posted_at' THEN posted_at END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'last_edited_at' THEN last_edited_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_last_edited_at' THEN last_edited_at END DESC NULLS LAST,
	t_messages_pkey ASC;

-- name: GetMessagesWithSenderUseNumberedPaginate :many
SELECT t_messages.*, m_members.name member_name, m_members.first_name member_first_name, m_members.last_name member_last_name, m_members.email member_email, m_members.grade_id member_grade_id, m_members.group_id member_group_id,
m_members.profile_image_id member_profile_image_id, t_images.height member_profile_image_height,
t_images.width member_profile_image_width, t_images.attachable_item_id member_profile_image_attachable_item_id,
t_attachable_items.owner_id member_profile_image_owner_id, t_attachable_items.from_outer member_profile_image_from_outer, t_attachable_items.alias member_profile_image_alias,
t_attachable_items.url member_profile_image_url, t_attachable_items.size member_profile_image_size, t_attachable_items.mime_type_id member_profile_image_mime_type_id FROM t_messages
LEFT JOIN m_members ON t_messages.sender_id = m_members.member_id
LEFT JOIN t_images ON m_members.profile_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE
	CASE WHEN @where_in_chat_room::boolean = true THEN t_chat_room_actions in (SELECT chat_room_action_id FROM t_chat_room_actions WHERE chat_room_id = ANY(@in_chat_room::uuid[])) ELSE TRUE END
AND
	CASE WHEN @where_in_sender::boolean = true THEN sender_id = ANY(@in_sender::uuid[]) ELSE TRUE END
AND
	CASE WHEN @where_like_body::boolean = true THEN body LIKE '%' || @search_body::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_earlier_posted_at::boolean = true THEN posted_at >= @earlier_posted_at ELSE TRUE END
AND
	CASE WHEN @where_later_posted_at::boolean = true THEN posted_at <= @later_posted_at ELSE TRUE END
AND
	CASE WHEN @where_earlier_last_edited_at::boolean = true THEN last_edited_at >= @earlier_last_edited_at ELSE TRUE END
AND
	CASE WHEN @where_later_last_edited_at::boolean = true THEN last_edited_at <= @later_last_edited_at ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'posted_at' THEN posted_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_posted_at' THEN posted_at END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'last_edited_at' THEN last_edited_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_last_edited_at' THEN last_edited_at END DESC NULLS LAST,
	t_messages_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetMessagesWithSenderUseKeysetPaginate :many
SELECT t_messages.*, m_members.name member_name, m_members.first_name member_first_name, m_members.last_name member_last_name, m_members.email member_email, m_members.grade_id member_grade_id, m_members.group_id member_group_id,
m_members.profile_image_id member_profile_image_id, t_images.height member_profile_image_height,
t_images.width member_profile_image_width, t_images.attachable_item_id member_profile_image_attachable_item_id,
t_attachable_items.owner_id member_profile_image_owner_id, t_attachable_items.from_outer member_profile_image_from_outer, t_attachable_items.alias member_profile_image_alias,
t_attachable_items.url member_profile_image_url, t_attachable_items.size member_profile_image_size, t_attachable_items.mime_type_id member_profile_image_mime_type_id FROM t_messages
LEFT JOIN m_members ON t_messages.sender_id = m_members.member_id
LEFT JOIN t_images ON m_members.profile_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE
	CASE WHEN @where_in_chat_room::boolean = true THEN t_chat_room_actions in (SELECT chat_room_action_id FROM t_chat_room_actions WHERE chat_room_id = ANY(@in_chat_room::uuid[])) ELSE TRUE END
AND
	CASE WHEN @where_in_sender::boolean = true THEN sender_id = ANY(@in_sender::uuid[]) ELSE TRUE END
AND
	CASE WHEN @where_like_body::boolean = true THEN body LIKE '%' || @search_body::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_earlier_posted_at::boolean = true THEN posted_at >= @earlier_posted_at ELSE TRUE END
AND
	CASE WHEN @where_later_posted_at::boolean = true THEN posted_at <= @later_posted_at ELSE TRUE END
AND
	CASE WHEN @where_earlier_last_edited_at::boolean = true THEN last_edited_at >= @earlier_last_edited_at ELSE TRUE END
AND
	CASE WHEN @where_later_last_edited_at::boolean = true THEN last_edited_at <= @later_last_edited_at ELSE TRUE END
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'posted_at' THEN posted_at > @posted_at_cursor OR (posted_at = @posted_at_cursor AND t_messages_pkey > @cursor::int)
				WHEN 'r_posted_at' THEN posted_at < @posted_at_cursor OR (posted_at = @posted_at_cursor AND t_messages_pkey > @cursor::int)
				WHEN 'last_edited_at' THEN last_edited_at > @last_edited_at_cursor OR (last_edited_at = @last_edited_at_cursor AND t_messages_pkey > @cursor::int)
				WHEN 'r_last_edited_at' THEN last_edited_at < @last_edited_at_cursor OR (last_edited_at = @last_edited_at_cursor AND t_messages_pkey > @cursor::int)
				ELSE t_messages_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'posted_at' THEN posted_at < @posted_at_cursor OR (posted_at = @posted_at_cursor AND t_messages_pkey < @cursor::int)
				WHEN 'r_posted_at' THEN posted_at > @posted_at_cursor OR (posted_at = @posted_at_cursor AND t_messages_pkey < @cursor::int)
				WHEN 'last_edited_at' THEN last_edited_at < @last_edited_at_cursor OR (last_edited_at = @last_edited_at_cursor AND t_messages_pkey < @cursor::int)
				WHEN 'r_last_edited_at' THEN last_edited_at > @last_edited_at_cursor OR (last_edited_at = @last_edited_at_cursor AND t_messages_pkey < @cursor::int)
				ELSE t_messages_pkey < @cursor::int
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'posted_at' AND @cursor_direction::text = 'next' THEN posted_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'posted_at' AND @cursor_direction::text = 'prev' THEN posted_at END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'r_posted_at' AND @cursor_direction::text = 'next' THEN posted_at END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'r_posted_at' AND @cursor_direction::text = 'prev' THEN posted_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'last_edited_at' AND @cursor_direction::text = 'next' THEN last_edited_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'last_edited_at' AND @cursor_direction::text = 'prev' THEN last_edited_at END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'r_last_edited_at' AND @cursor_direction::text = 'next' THEN last_edited_at END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'r_last_edited_at' AND @cursor_direction::text = 'prev' THEN last_edited_at END ASC NULLS LAST,
	CASE WHEN @cursor_direction::text = 'next' THEN t_messages_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN t_messages_pkey END DESC
LIMIT $1;

-- name: GetPluralMessagesWithSender :many
SELECT t_messages.*, m_members.name member_name, m_members.first_name member_first_name, m_members.last_name member_last_name, m_members.email member_email, m_members.grade_id member_grade_id, m_members.group_id member_group_id,
m_members.profile_image_id member_profile_image_id, t_images.height member_profile_image_height,
t_images.width member_profile_image_width, t_images.attachable_item_id member_profile_image_attachable_item_id,
t_attachable_items.owner_id member_profile_image_owner_id, t_attachable_items.from_outer member_profile_image_from_outer, t_attachable_items.alias member_profile_image_alias,
t_attachable_items.url member_profile_image_url, t_attachable_items.size member_profile_image_size, t_attachable_items.mime_type_id member_profile_image_mime_type_id FROM t_messages
LEFT JOIN m_members ON t_messages.sender_id = m_members.member_id
LEFT JOIN t_images ON m_members.profile_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE message_id = ANY(@message_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'posted_at' THEN posted_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_posted_at' THEN posted_at END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'last_edited_at' THEN last_edited_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_last_edited_at' THEN last_edited_at END DESC NULLS LAST,
	t_messages_pkey ASC;

-- name: GetPluralMessagesWithSenderUseNumberedPaginate :many
SELECT t_messages.*, m_members.name member_name, m_members.first_name member_first_name, m_members.last_name member_last_name, m_members.email member_email, m_members.grade_id member_grade_id, m_members.group_id member_group_id,
m_members.profile_image_id member_profile_image_id, t_images.height member_profile_image_height,
t_images.width member_profile_image_width, t_images.attachable_item_id member_profile_image_attachable_item_id,
t_attachable_items.owner_id member_profile_image_owner_id, t_attachable_items.from_outer member_profile_image_from_outer, t_attachable_items.alias member_profile_image_alias,
t_attachable_items.url member_profile_image_url, t_attachable_items.size member_profile_image_size, t_attachable_items.mime_type_id member_profile_image_mime_type_id FROM t_messages
LEFT JOIN m_members ON t_messages.sender_id = m_members.member_id
LEFT JOIN t_images ON m_members.profile_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE message_id = ANY(@message_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'posted_at' THEN posted_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_posted_at' THEN posted_at END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'last_edited_at' THEN last_edited_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_last_edited_at' THEN last_edited_at END DESC NULLS LAST,
	t_messages_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetPluralMessagesWithSenderByChatRoomActionIDs :many
SELECT t_messages.*, m_members.name member_name, m_members.first_name member_first_name, m_members.last_name member_last_name, m_members.email member_email, m_members.grade_id member_grade_id, m_members.group_id member_group_id,
m_members.profile_image_id member_profile_image_id, t_images.height member_profile_image_height,
t_images.width member_profile_image_width, t_images.attachable_item_id member_profile_image_attachable_item_id,
t_attachable_items.owner_id member_profile_image_owner_id, t_attachable_items.from_outer member_profile_image_from_outer, t_attachable_items.alias member_profile_image_alias,
t_attachable_items.url member_profile_image_url, t_attachable_items.size member_profile_image_size, t_attachable_items.mime_type_id member_profile_image_mime_type_id FROM t_messages
LEFT JOIN m_members ON t_messages.sender_id = m_members.member_id
LEFT JOIN t_images ON m_members.profile_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE chat_room_action_id = ANY(@chat_room_action_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'posted_at' THEN posted_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_posted_at' THEN posted_at END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'last_edited_at' THEN last_edited_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_last_edited_at' THEN last_edited_at END DESC NULLS LAST,
	t_messages_pkey ASC;

-- name: GetPluralMessagesWithSenderByChatRoomActionIDsUseNumberedPaginate :many
SELECT t_messages.*, m_members.name member_name, m_members.first_name member_first_name, m_members.last_name member_last_name, m_members.email member_email, m_members.grade_id member_grade_id, m_members.group_id member_group_id,
m_members.profile_image_id member_profile_image_id, t_images.height member_profile_image_height,
t_images.width member_profile_image_width, t_images.attachable_item_id member_profile_image_attachable_item_id,
t_attachable_items.owner_id member_profile_image_owner_id, t_attachable_items.from_outer member_profile_image_from_outer, t_attachable_items.alias member_profile_image_alias,
t_attachable_items.url member_profile_image_url, t_attachable_items.size member_profile_image_size, t_attachable_items.mime_type_id member_profile_image_mime_type_id FROM t_messages
LEFT JOIN m_members ON t_messages.sender_id = m_members.member_id
LEFT JOIN t_images ON m_members.profile_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE chat_room_action_id = ANY(@chat_room_action_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'posted_at' THEN posted_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_posted_at' THEN posted_at END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'last_edited_at' THEN last_edited_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_last_edited_at' THEN last_edited_at END DESC NULLS LAST,
	t_messages_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountMessages :one
SELECT COUNT(*) FROM t_messages
WHERE
	CASE WHEN @where_in_chat_room::boolean = true THEN t_chat_room_actions in (SELECT chat_room_action_id FROM t_chat_room_actions WHERE chat_room_id = ANY(@in_chat_room::uuid[])) ELSE TRUE END
AND
	CASE WHEN @where_in_sender::boolean = true THEN sender_id = ANY(@in_sender::uuid[]) ELSE TRUE END
AND
	CASE WHEN @where_like_body::boolean = true THEN body LIKE '%' || @search_body::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_earlier_posted_at::boolean = true THEN posted_at >= @earlier_posted_at ELSE TRUE END
AND
	CASE WHEN @where_later_posted_at::boolean = true THEN posted_at <= @later_posted_at ELSE TRUE END
AND
	CASE WHEN @where_earlier_last_edited_at::boolean = true THEN last_edited_at >= @earlier_last_edited_at ELSE TRUE END
AND
	CASE WHEN @where_later_last_edited_at::boolean = true THEN last_edited_at <= @later_last_edited_at ELSE TRUE END;
