-- name: CreateReadReceipts :copyfrom
INSERT INTO t_read_receipts (member_id, message_id, read_at) VALUES ($1, $2, $3);

-- name: CreateReadReceipt :one
INSERT INTO t_read_receipts (member_id, message_id, read_at) VALUES ($1, $2, $3) RETURNING *;

-- name: ReadReceipt :one
UPDATE t_read_receipts SET read_at = $3 WHERE member_id = $1 AND message_id = $2 RETURNING *;

-- name: ReadReceipts :execrows
UPDATE t_read_receipts SET read_at = $2 WHERE member_id = $1 AND message_id = ANY(@message_ids::uuid[]);

-- name: GetReadableMembersOnMessage :many
SELECT m_members.*, t_read_receipts.read_at read_at, t_images.height profile_image_height,
t_images.width profile_image_width, t_images.attachable_item_id profile_image_attachable_item_id,
t_attachable_items.owner_id profile_image_owner_id, t_attachable_items.from_outer profile_image_from_outer,
t_attachable_items.url profile_image_url, t_attachable_items.size profile_image_size, t_attachable_items.mime_type_id profile_image_mime_type_id
FROM m_members
LEFT JOIN t_read_receipts ON m_members.member_id = t_read_receipts.member_id
LEFT JOIN t_images ON m_members.profile_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE message_id = $1
AND
	CASE WHEN @where_like_name::boolean = true THEN m_members.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_is_read::boolean = true THEN t_read_receipts.read_at IS NOT NULL ELSE TRUE END
AND
	CASE WHEN @where_is_not_read::boolean = true THEN t_read_receipts.read_at IS NULL ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_members.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_members.name END DESC,
	CASE WHEN @order_method::text = 'read_at' THEN t_read_receipts.read_at END ASC,
	CASE WHEN @order_method::text = 'r_read_at' THEN t_read_receipts.read_at END DESC,
	m_members_pkey ASC;

-- name: GetReadableMembersOnMessageUseNumberedPaginate :many
SELECT m_members.*, t_read_receipts.read_at read_at, t_images.height profile_image_height,
t_images.width profile_image_width, t_images.attachable_item_id profile_image_attachable_item_id,
t_attachable_items.owner_id profile_image_owner_id, t_attachable_items.from_outer profile_image_from_outer,
t_attachable_items.url profile_image_url, t_attachable_items.size profile_image_size, t_attachable_items.mime_type_id profile_image_mime_type_id
FROM m_members
LEFT JOIN t_read_receipts ON m_members.member_id = t_read_receipts.member_id
LEFT JOIN t_images ON m_members.profile_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE message_id = $1
AND
	CASE WHEN @where_like_name::boolean = true THEN m_members.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_is_read::boolean = true THEN t_read_receipts.read_at IS NOT NULL ELSE TRUE END
AND
	CASE WHEN @where_is_not_read::boolean = true THEN t_read_receipts.read_at IS NULL ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_members.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_members.name END DESC,
	CASE WHEN @order_method::text = 'read_at' THEN t_read_receipts.read_at END ASC,
	CASE WHEN @order_method::text = 'r_read_at' THEN t_read_receipts.read_at END DESC,
	m_members_pkey ASC
LIMIT $2 OFFSET $3;

-- name: GetReadableMembersOnMessageUseKeysetPaginate :many
SELECT m_members.*, t_read_receipts.read_at read_at, t_images.height profile_image_height,
t_images.width profile_image_width, t_images.attachable_item_id profile_image_attachable_item_id,
t_attachable_items.owner_id profile_image_owner_id, t_attachable_items.from_outer profile_image_from_outer,
t_attachable_items.url profile_image_url, t_attachable_items.size profile_image_size, t_attachable_items.mime_type_id profile_image_mime_type_id
FROM m_members
LEFT JOIN t_read_receipts ON m_members.member_id = t_read_receipts.member_id
LEFT JOIN t_images ON m_members.profile_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE message_id = $1
AND
	CASE WHEN @where_like_name::boolean = true THEN m_members.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE WHEN @where_is_read::boolean = true THEN t_read_receipts.read_at IS NOT NULL ELSE TRUE END
AND
	CASE WHEN @where_is_not_read::boolean = true THEN t_read_receipts.read_at IS NULL ELSE TRUE END
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'name' THEN m_members.name > @name_cursor OR (m_members.name = @name_cursor AND m_members_pkey > @cursor::int)
				WHEN 'r_name' THEN m_members.name < @name_cursor OR (m_members.name = @name_cursor AND m_members_pkey > @cursor::int)
				WHEN 'read_at' THEN t_read_receipts.read_at > @read_at_cursor OR (t_read_receipts.read_at = @read_at_cursor AND m_members_pkey > @cursor::int)
				WHEN 'r_read_at' THEN t_read_receipts.read_at < @read_at_cursor OR (t_read_receipts.read_at = @read_at_cursor AND m_members_pkey > @cursor::int)
				ELSE m_members_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'name' THEN m_members.name < @name_cursor OR (m_members.name = @name_cursor AND m_members_pkey < @cursor::int)
				WHEN 'r_name' THEN m_members.name > @name_cursor OR (m_members.name = @name_cursor AND m_members_pkey < @cursor::int)
				WHEN 'read_at' THEN t_read_receipts.read_at < @read_at_cursor OR (t_read_receipts.read_at = @read_at_cursor AND m_members_pkey < @cursor::int)
				WHEN 'r_read_at' THEN t_read_receipts.read_at > @read_at_cursor OR (t_read_receipts.read_at = @read_at_cursor AND m_members_pkey < @cursor::int)
				ELSE m_members_pkey < @cursor::int
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'next' THEN m_members.name END ASC,
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'prev' THEN m_members.name END DESC,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'next' THEN m_members.name END DESC,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'prev' THEN m_members.name END ASC,
	CASE WHEN @order_method::text = 'read_at' AND @cursor_direction::text = 'next' THEN t_read_receipts.read_at END ASC,
	CASE WHEN @order_method::text = 'read_at' AND @cursor_direction::text = 'prev' THEN t_read_receipts.read_at END DESC,
	CASE WHEN @order_method::text = 'r_read_at' AND @cursor_direction::text = 'next' THEN t_read_receipts.read_at END DESC,
	CASE WHEN @order_method::text = 'r_read_at' AND @cursor_direction::text = 'prev' THEN t_read_receipts.read_at END ASC,
	CASE WHEN @cursor_direction::text = 'next' THEN m_members_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN m_members_pkey END DESC
LIMIT $2;

-- name: GetPluralReadableMembersOnMessage :many
SELECT m_members.*, t_read_receipts.read_at read_at, t_images.height profile_image_height,
t_images.width profile_image_width, t_images.attachable_item_id profile_image_attachable_item_id,
t_attachable_items.owner_id profile_image_owner_id, t_attachable_items.from_outer profile_image_from_outer,
t_attachable_items.url profile_image_url, t_attachable_items.size profile_image_size, t_attachable_items.mime_type_id profile_image_mime_type_id
FROM m_members
LEFT JOIN t_read_receipts ON m_members.member_id = t_read_receipts.member_id
LEFT JOIN t_images ON m_members.profile_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE message_id = ANY(@message_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_members.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_members.name END DESC,
	CASE WHEN @order_method::text = 'read_at' THEN t_read_receipts.read_at END ASC,
	CASE WHEN @order_method::text = 'r_read_at' THEN t_read_receipts.read_at END DESC,
	m_members_pkey ASC;

-- name: GetPluralReadableMembersOnMessageUseNumberedPaginate :many
SELECT m_members.*, t_read_receipts.read_at read_at, t_images.height profile_image_height,
t_images.width profile_image_width, t_images.attachable_item_id profile_image_attachable_item_id,
t_attachable_items.owner_id profile_image_owner_id, t_attachable_items.from_outer profile_image_from_outer,
t_attachable_items.url profile_image_url, t_attachable_items.size profile_image_size, t_attachable_items.mime_type_id profile_image_mime_type_id
FROM m_members
LEFT JOIN t_read_receipts ON m_members.member_id = t_read_receipts.member_id
LEFT JOIN t_images ON m_members.profile_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE message_id = ANY(@message_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_members.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_members.name END DESC,
	CASE WHEN @order_method::text = 'read_at' THEN t_read_receipts.read_at END ASC,
	CASE WHEN @order_method::text = 'r_read_at' THEN t_read_receipts.read_at END DESC,
	m_members_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountReadableMembersOnMessage :one
SELECT COUNT(*) FROM t_read_receipts
LEFT JOIN m_members ON t_read_receipts.member_id = m_members.member_id
WHERE message_id = $1
AND
	CASE WHEN @where_like_name::boolean = true THEN member_id IN (SELECT member_id FROM m_members WHERE name LIKE '%' || @search_name::text || '%') ELSE TRUE END
AND
	CASE WHEN @where_is_read::boolean = true THEN read_at IS NOT NULL ELSE TRUE END
AND
	CASE WHEN @where_is_not_read::boolean = true THEN read_at IS NULL ELSE TRUE END;

-- name: GetReadableMessagesOnMember :many
SELECT t_messages.*, t_read_receipts.read_at read_at FROM t_messages
LEFT JOIN t_read_receipts ON t_messages.message_id = t_read_receipts.message_id
WHERE member_id = $1
AND
	CASE WHEN @where_is_read::boolean = true THEN t_read_receipts.read_at IS NOT NULL ELSE TRUE END
AND
	CASE WHEN @where_is_not_read::boolean = true THEN t_read_receipts.read_at IS NULL ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'read_at' THEN t_read_receipts.read_at END ASC,
	CASE WHEN @order_method::text = 'r_read_at' THEN t_read_receipts.read_at END DESC,
	t_messages_pkey ASC;

-- name: GetReadableMessagesOnMemberUseNumberedPaginate :many
SELECT t_messages.*, t_read_receipts.read_at read_at FROM t_messages
LEFT JOIN t_read_receipts ON t_messages.message_id = t_read_receipts.message_id
WHERE member_id = $1
AND
	CASE WHEN @where_is_read::boolean = true THEN t_read_receipts.read_at IS NOT NULL ELSE TRUE END
AND
	CASE WHEN @where_is_not_read::boolean = true THEN t_read_receipts.read_at IS NULL ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'read_at' THEN t_read_receipts.read_at END ASC,
	CASE WHEN @order_method::text = 'r_read_at' THEN t_read_receipts.read_at END DESC,
	t_messages_pkey ASC
LIMIT $2 OFFSET $3;

-- name: GetReadableMessagesOnMemberUseKeysetPaginate :many
SELECT t_messages.*, t_read_receipts.read_at read_at FROM t_messages
LEFT JOIN t_read_receipts ON t_messages.message_id = t_read_receipts.message_id
WHERE member_id = $1
AND
	CASE WHEN @where_is_read::boolean = true THEN t_read_receipts.read_at IS NOT NULL ELSE TRUE END
AND
	CASE WHEN @where_is_not_read::boolean = true THEN t_read_receipts.read_at IS NULL ELSE TRUE END
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'read_at' THEN t_read_receipts.read_at > @read_at_cursor OR (t_read_receipts.read_at = @read_at_cursor AND t_messages_pkey > @cursor::int)
				WHEN 'r_read_at' THEN t_read_receipts.read_at < @read_at_cursor OR (t_read_receipts.read_at = @read_at_cursor AND t_messages_pkey > @cursor::int)
				ELSE t_messages_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'read_at' THEN t_read_receipts.read_at < @read_at_cursor OR (t_read_receipts.read_at = @read_at_cursor AND t_messages_pkey < @cursor::int)
				WHEN 'r_read_at' THEN t_read_receipts.read_at > @read_at_cursor OR (t_read_receipts.read_at = @read_at_cursor AND t_messages_pkey < @cursor::int)
				ELSE t_messages_pkey < @cursor::int
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'read_at' AND @cursor_direction::text = 'next' THEN t_read_receipts.read_at END ASC,
	CASE WHEN @order_method::text = 'read_at' AND @cursor_direction::text = 'prev' THEN t_read_receipts.read_at END DESC,
	CASE WHEN @order_method::text = 'r_read_at' AND @cursor_direction::text = 'next' THEN t_read_receipts.read_at END DESC,
	CASE WHEN @order_method::text = 'r_read_at' AND @cursor_direction::text = 'prev' THEN t_read_receipts.read_at END ASC,
	CASE WHEN @cursor_direction::text = 'next' THEN t_messages_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN t_messages_pkey END DESC
LIMIT $2;

-- name: GetPluralReadableMessagesOnMember :many
SELECT t_messages.*, t_read_receipts.read_at read_at FROM t_messages
LEFT JOIN t_read_receipts ON t_messages.message_id = t_read_receipts.message_id
WHERE member_id = ANY(@member_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'read_at' THEN t_read_receipts.read_at END ASC,
	CASE WHEN @order_method::text = 'r_read_at' THEN t_read_receipts.read_at END DESC,
	t_messages_pkey ASC;

-- name: GetPluralReadableMessagesOnMemberUseNumberedPaginate :many
SELECT t_messages.*, t_read_receipts.read_at read_at FROM t_messages
LEFT JOIN t_read_receipts ON t_messages.message_id = t_read_receipts.message_id
WHERE member_id = ANY(@member_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'read_at' THEN t_read_receipts.read_at END ASC,
	CASE WHEN @order_method::text = 'r_read_at' THEN t_read_receipts.read_at END DESC,
	t_messages_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountReadableMessagesOnMember :one
SELECT COUNT(*) FROM t_read_receipts
WHERE member_id = $1
AND
	CASE WHEN @where_is_read::boolean = true THEN read_at IS NOT NULL ELSE TRUE END
AND
	CASE WHEN @where_is_not_read::boolean = true THEN read_at IS NULL ELSE TRUE END;

-- name: CountReadableMessagesOnChatRoomAndMember :one
SELECT COUNT(*) FROM t_read_receipts
LEFT JOIN t_messages ON t_read_receipts.message_id = t_messages.message_id
WHERE t_messages.chat_room_id = $1
AND t_read_receipts.member_id = $2
AND
	CASE WHEN @where_is_read::boolean = true THEN t_read_receipts.read_at IS NOT NULL ELSE TRUE END
AND
	CASE WHEN @where_is_not_read::boolean = true THEN t_read_receipts.read_at IS NULL ELSE TRUE END;

-- name: CountReadableMessagesOnChatRoomsAndMember :many
SELECT t_messages.chat_room_id, COUNT(*) FROM t_read_receipts
LEFT JOIN t_messages ON t_read_receipts.message_id = t_messages.message_id
WHERE t_messages.chat_room_id = ANY(@chat_room_ids::uuid[])
AND t_read_receipts.member_id = $1
AND
	CASE WHEN @where_is_read::boolean = true THEN t_read_receipts.read_at IS NOT NULL ELSE TRUE END
AND
	CASE WHEN @where_is_not_read::boolean = true THEN t_read_receipts.read_at IS NULL ELSE TRUE END
GROUP BY t_messages.chat_room_id;
