-- name: CreateChatRoomActions :copyfrom
INSERT INTO t_chat_room_actions (chat_room_id, chat_room_action_type_id, acted_at) VALUES ($1, $2, $3);

-- name: CreateChatRoomAction :one
INSERT INTO t_chat_room_actions (chat_room_id, chat_room_action_type_id, acted_at) VALUES ($1, $2, $3) RETURNING *;

-- name: UpdateChatRoomAction :one
UPDATE t_chat_room_actions SET chat_room_action_type_id = $2 WHERE chat_room_action_id = $1 RETURNING *;

-- name: DeleteChatRoomAction :execrows
DELETE FROM t_chat_room_actions WHERE chat_room_action_id = $1;

-- name: PluralDeleteChatRoomActions :execrows
DELETE FROM t_chat_room_actions WHERE chat_room_action_id = ANY(@chat_room_action_ids::uuid[]);

-- name: GetChatRoomActionsOnChatRoom :many
SELECT * FROM t_chat_room_actions
WHERE chat_room_id = $1
AND
	CASE WHEN @where_in_chat_room_action_type_ids::boolean = true THEN chat_room_action_type_id = ANY(@in_chat_room_action_type_ids::uuid[]) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'acted_at' THEN acted_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_acted_at' THEN acted_at END DESC NULLS LAST,
	t_chat_room_actions_pkey ASC;

-- name: GetChatRoomActionsOnChatRoomUseNumberedPaginate :many
SELECT * FROM t_chat_room_actions
WHERE chat_room_id = $1
AND
	CASE WHEN @where_in_chat_room_action_type_ids::boolean = true THEN chat_room_action_type_id = ANY(@in_chat_room_action_type_ids::uuid[]) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'acted_at' THEN acted_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_acted_at' THEN acted_at END DESC NULLS LAST,
	t_chat_room_actions_pkey ASC
LIMIT $2 OFFSET $3;

-- name: GetChatRoomActionsOnChatRoomUseKeysetPaginate :many
SELECT * FROM t_chat_room_actions
WHERE chat_room_id = $1
AND
	CASE WHEN @where_in_chat_room_action_type_ids::boolean = true THEN chat_room_action_type_id = ANY(@in_chat_room_action_type_ids::uuid[]) ELSE TRUE END
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'acted_at' THEN acted_at > @acted_at_cursor OR (acted_at = @acted_at_cursor AND t_chat_room_actions_pkey > @cursor::int)
				WHEN 'r_acted_at' THEN acted_at < @acted_at_cursor OR (acted_at = @acted_at_cursor AND t_chat_room_actions_pkey > @cursor::int)
				ELSE t_chat_room_actions_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'acted_at' THEN acted_at < @acted_at_cursor OR (acted_at = @acted_at_cursor AND t_chat_room_actions_pkey < @cursor::int)
				WHEN 'r_acted_at' THEN acted_at > @acted_at_cursor OR (acted_at = @acted_at_cursor AND t_chat_room_actions_pkey < @cursor::int)
				ELSE t_chat_room_actions_pkey < @cursor::int
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'acted_at' THEN acted_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_acted_at' THEN acted_at END DESC NULLS LAST,
	t_chat_room_actions_pkey ASC
LIMIT $2;

-- name: GetPluralChatRoomActions :many
SELECT * FROM t_chat_room_actions
WHERE chat_room_action_id = ANY(@chat_room_action_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'acted_at' THEN acted_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_acted_at' THEN acted_at END DESC NULLS LAST,
	t_chat_room_actions_pkey ASC;

-- name: GetPluralChatRoomActionsUseNumberedPaginate :many
SELECT * FROM t_chat_room_actions
WHERE chat_room_action_id = ANY(@chat_room_action_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'acted_at' THEN acted_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_acted_at' THEN acted_at END DESC NULLS LAST,
	t_chat_room_actions_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetChatRoomActionsWithDetailOnChatRoom :many
SELECT t_chat_room_actions.*,
t_chat_room_create_actions.chat_room_create_action_id, t_chat_room_create_actions.name create_name, cm.member_id create_member_id, cm.name create_member_name, cm.first_name create_member_first_name, cm.last_name create_member_last_name, cm.email create_member_email, cm.profile_image_id create_member_profile_image_id, cm.grade_id create_member_grade_id, cm.group_id create_member_group_id,
t_chat_room_update_name_actions.chat_room_update_name_action_id, t_chat_room_update_name_actions.name update_name, um.member_id update_member_id, um.name update_member_name, um.first_name update_member_first_name, um.last_name update_member_last_name, um.email update_member_email, um.profile_image_id update_member_profile_image_id, um.grade_id update_member_grade_id, um.group_id update_member_group_id,
t_chat_room_withdraw_actions.chat_room_withdraw_action_id, wm.member_id withdraw_member_id, wm.name withdraw_member_name, wm.first_name withdraw_member_first_name, wm.last_name withdraw_member_last_name, wm.email withdraw_member_email, wm.profile_image_id withdraw_member_profile_image_id, wm.grade_id withdraw_member_grade_id, wm.group_id withdraw_member_group_id,
t_chat_room_add_member_actions.chat_room_add_member_action_id, am.member_id add_member_id, am.name add_member_name, am.first_name add_member_first_name, am.last_name add_member_last_name, am.email add_member_email, am.profile_image_id add_member_profile_image_id, am.grade_id add_member_grade_id, am.group_id add_member_group_id,
t_chat_room_remove_member_actions.chat_room_remove_member_action_id, rm.member_id remove_member_id, rm.name remove_member_name, rm.first_name remove_member_first_name, rm.last_name remove_member_last_name, rm.email remove_member_email, rm.profile_image_id remove_member_profile_image_id, rm.grade_id remove_member_grade_id, rm.group_id remove_member_group_id,
t_chat_room_delete_message_actions.chat_room_delete_message_action_id, dm.member_id delete_message_member_id, dm.name delete_message_member_name, dm.first_name delete_message_member_first_name, dm.last_name delete_message_member_last_name, dm.email delete_message_member_email, dm.profile_image_id delete_message_member_profile_image_id, dm.grade_id delete_message_member_grade_id, dm.group_id delete_message_member_group_id,
t_messages.message_id, t_messages.sender_id message_sender_id, t_messages.body message_body, t_messages.posted_at message_posted_at, t_messages.last_edited_at message_last_edited_at,
mm.name message_sender_name, mm.first_name message_sender_first_name, mm.last_name message_sender_last_name, mm.email message_sender_email, mm.profile_image_id message_sender_profile_image_id, mm.grade_id message_sender_grade_id, mm.group_id message_sender_group_id,
t_images.height message_sender_profile_image_height, t_images.width message_sender_profile_image_width, t_images.attachable_item_id message_sender_profile_image_attachable_item_id,
t_attachable_items.owner_id message_sender_profile_image_owner_id, t_attachable_items.from_outer message_sender_profile_image_from_outer, t_attachable_items.alias message_sender_profile_image_alias,
t_attachable_items.url message_sender_profile_image_url, t_attachable_items.size message_sender_profile_image_size, t_attachable_items.mime_type_id message_sender_profile_image_mime_type_id
FROM t_chat_room_actions
LEFT JOIN t_chat_room_create_actions ON t_chat_room_actions.chat_room_action_id = t_chat_room_create_actions.chat_room_action_id
LEFT JOIN m_members cm ON t_chat_room_create_actions.created_by = cm.member_id
LEFT JOIN t_chat_room_update_name_actions ON t_chat_room_actions.chat_room_action_id = t_chat_room_update_name_actions.chat_room_action_id
LEFT JOIN m_members um ON t_chat_room_update_name_actions.updated_by = um.member_id
LEFT JOIN t_chat_room_withdraw_actions ON t_chat_room_actions.chat_room_action_id = t_chat_room_withdraw_actions.chat_room_action_id
LEFT JOIN m_members wm ON t_chat_room_withdraw_actions.member_id = wm.member_id
LEFT JOIN t_chat_room_add_member_actions ON t_chat_room_actions.chat_room_action_id = t_chat_room_add_member_actions.chat_room_action_id
LEFT JOIN m_members am ON t_chat_room_add_member_actions.added_by = am.member_id
LEFT JOIN t_chat_room_remove_member_actions ON t_chat_room_actions.chat_room_action_id = t_chat_room_remove_member_actions.chat_room_action_id
LEFT JOIN m_members rm ON t_chat_room_remove_member_actions.removed_by = rm.member_id
LEFT JOIN t_chat_room_delete_message_actions ON t_chat_room_actions.chat_room_action_id = t_chat_room_delete_message_actions.chat_room_action_id
LEFT JOIN m_members dm ON t_chat_room_delete_message_actions.deleted_by = dm.member_id
LEFT JOIN t_messages ON t_chat_room_actions.chat_room_action_id = t_messages.chat_room_action_id
LEFT JOIN m_members mm ON t_messages.sender_id = mm.member_id
LEFT JOIN t_images ON mm.profile_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE t_chat_room_actions.chat_room_id = $1
AND
	CASE WHEN @where_in_chat_room_action_type_ids::boolean = true THEN chat_room_action_type_id = ANY(@in_chat_room_action_type_ids::uuid[]) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'acted_at' THEN acted_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_acted_at' THEN acted_at END DESC NULLS LAST,
	t_chat_room_actions_pkey ASC;

-- name: GetChatRoomActionsWithDetailOnChatRoomUseNumberedPaginate :many
SELECT t_chat_room_actions.*,
t_chat_room_create_actions.chat_room_create_action_id, t_chat_room_create_actions.name create_name, cm.member_id create_member_id, cm.name create_member_name, cm.first_name create_member_first_name, cm.last_name create_member_last_name, cm.email create_member_email, cm.profile_image_id create_member_profile_image_id, cm.grade_id create_member_grade_id, cm.group_id create_member_group_id,
t_chat_room_update_name_actions.chat_room_update_name_action_id, t_chat_room_update_name_actions.name update_name, um.member_id update_member_id, um.name update_member_name, um.first_name update_member_first_name, um.last_name update_member_last_name, um.email update_member_email, um.profile_image_id update_member_profile_image_id, um.grade_id update_member_grade_id, um.group_id update_member_group_id,
t_chat_room_withdraw_actions.chat_room_withdraw_action_id, wm.member_id withdraw_member_id, wm.name withdraw_member_name, wm.first_name withdraw_member_first_name, wm.last_name withdraw_member_last_name, wm.email withdraw_member_email, wm.profile_image_id withdraw_member_profile_image_id, wm.grade_id withdraw_member_grade_id, wm.group_id withdraw_member_group_id,
t_chat_room_add_member_actions.chat_room_add_member_action_id, am.member_id add_member_id, am.name add_member_name, am.first_name add_member_first_name, am.last_name add_member_last_name, am.email add_member_email, am.profile_image_id add_member_profile_image_id, am.grade_id add_member_grade_id, am.group_id add_member_group_id,
t_chat_room_remove_member_actions.chat_room_remove_member_action_id, rm.member_id remove_member_id, rm.name remove_member_name, rm.first_name remove_member_first_name, rm.last_name remove_member_last_name, rm.email remove_member_email, rm.profile_image_id remove_member_profile_image_id, rm.grade_id remove_member_grade_id, rm.group_id remove_member_group_id,
t_chat_room_delete_message_actions.chat_room_delete_message_action_id, dm.member_id delete_message_member_id, dm.name delete_message_member_name, dm.first_name delete_message_member_first_name, dm.last_name delete_message_member_last_name, dm.email delete_message_member_email, dm.profile_image_id delete_message_member_profile_image_id, dm.grade_id delete_message_member_grade_id, dm.group_id delete_message_member_group_id,
t_messages.message_id, t_messages.sender_id message_sender_id, t_messages.body message_body, t_messages.posted_at message_posted_at, t_messages.last_edited_at message_last_edited_at,
mm.name message_sender_name, mm.first_name message_sender_first_name, mm.last_name message_sender_last_name, mm.email message_sender_email, mm.profile_image_id message_sender_profile_image_id, mm.grade_id message_sender_grade_id, mm.group_id message_sender_group_id,
t_images.height message_sender_profile_image_height, t_images.width message_sender_profile_image_width, t_images.attachable_item_id message_sender_profile_image_attachable_item_id,
t_attachable_items.owner_id message_sender_profile_image_owner_id, t_attachable_items.from_outer message_sender_profile_image_from_outer, t_attachable_items.alias message_sender_profile_image_alias,
t_attachable_items.url message_sender_profile_image_url, t_attachable_items.size message_sender_profile_image_size, t_attachable_items.mime_type_id message_sender_profile_image_mime_type_id
FROM t_chat_room_actions
LEFT JOIN t_chat_room_create_actions ON t_chat_room_actions.chat_room_action_id = t_chat_room_create_actions.chat_room_action_id
LEFT JOIN m_members cm ON t_chat_room_create_actions.created_by = cm.member_id
LEFT JOIN t_chat_room_update_name_actions ON t_chat_room_actions.chat_room_action_id = t_chat_room_update_name_actions.chat_room_action_id
LEFT JOIN m_members um ON t_chat_room_update_name_actions.updated_by = um.member_id
LEFT JOIN t_chat_room_withdraw_actions ON t_chat_room_actions.chat_room_action_id = t_chat_room_withdraw_actions.chat_room_action_id
LEFT JOIN m_members wm ON t_chat_room_withdraw_actions.member_id = wm.member_id
LEFT JOIN t_chat_room_add_member_actions ON t_chat_room_actions.chat_room_action_id = t_chat_room_add_member_actions.chat_room_action_id
LEFT JOIN m_members am ON t_chat_room_add_member_actions.added_by = am.member_id
LEFT JOIN t_chat_room_remove_member_actions ON t_chat_room_actions.chat_room_action_id = t_chat_room_remove_member_actions.chat_room_action_id
LEFT JOIN m_members rm ON t_chat_room_remove_member_actions.removed_by = rm.member_id
LEFT JOIN t_chat_room_delete_message_actions ON t_chat_room_actions.chat_room_action_id = t_chat_room_delete_message_actions.chat_room_action_id
LEFT JOIN m_members dm ON t_chat_room_delete_message_actions.deleted_by = dm.member_id
LEFT JOIN t_messages ON t_chat_room_actions.chat_room_action_id = t_messages.chat_room_action_id
LEFT JOIN m_members mm ON t_messages.sender_id = mm.member_id
LEFT JOIN t_images ON mm.profile_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE t_chat_room_actions.chat_room_id = $1
AND
	CASE WHEN @where_in_chat_room_action_type_ids::boolean = true THEN chat_room_action_type_id = ANY(@in_chat_room_action_type_ids::uuid[]) ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'acted_at' THEN acted_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_acted_at' THEN acted_at END DESC NULLS LAST,
	t_chat_room_actions_pkey ASC
LIMIT $2 OFFSET $3;

-- name: GetChatRoomActionsWithDetailOnChatRoomUseKeysetPaginate :many
SELECT t_chat_room_actions.*,
t_chat_room_create_actions.chat_room_create_action_id, t_chat_room_create_actions.name create_name, cm.member_id create_member_id, cm.name create_member_name, cm.first_name create_member_first_name, cm.last_name create_member_last_name, cm.email create_member_email, cm.profile_image_id create_member_profile_image_id, cm.grade_id create_member_grade_id, cm.group_id create_member_group_id,
t_chat_room_update_name_actions.chat_room_update_name_action_id, t_chat_room_update_name_actions.name update_name, um.member_id update_member_id, um.name update_member_name, um.first_name update_member_first_name, um.last_name update_member_last_name, um.email update_member_email, um.profile_image_id update_member_profile_image_id, um.grade_id update_member_grade_id, um.group_id update_member_group_id,
t_chat_room_withdraw_actions.chat_room_withdraw_action_id, wm.member_id withdraw_member_id, wm.name withdraw_member_name, wm.first_name withdraw_member_first_name, wm.last_name withdraw_member_last_name, wm.email withdraw_member_email, wm.profile_image_id withdraw_member_profile_image_id, wm.grade_id withdraw_member_grade_id, wm.group_id withdraw_member_group_id,
t_chat_room_add_member_actions.chat_room_add_member_action_id, am.member_id add_member_id, am.name add_member_name, am.first_name add_member_first_name, am.last_name add_member_last_name, am.email add_member_email, am.profile_image_id add_member_profile_image_id, am.grade_id add_member_grade_id, am.group_id add_member_group_id,
t_chat_room_remove_member_actions.chat_room_remove_member_action_id, rm.member_id remove_member_id, rm.name remove_member_name, rm.first_name remove_member_first_name, rm.last_name remove_member_last_name, rm.email remove_member_email, rm.profile_image_id remove_member_profile_image_id, rm.grade_id remove_member_grade_id, rm.group_id remove_member_group_id,
t_chat_room_delete_message_actions.chat_room_delete_message_action_id, dm.member_id delete_message_member_id, dm.name delete_message_member_name, dm.first_name delete_message_member_first_name, dm.last_name delete_message_member_last_name, dm.email delete_message_member_email, dm.profile_image_id delete_message_member_profile_image_id, dm.grade_id delete_message_member_grade_id, dm.group_id delete_message_member_group_id,
t_messages.message_id, t_messages.sender_id message_sender_id, t_messages.body message_body, t_messages.posted_at message_posted_at, t_messages.last_edited_at message_last_edited_at,
mm.name message_sender_name, mm.first_name message_sender_first_name, mm.last_name message_sender_last_name, mm.email message_sender_email, mm.profile_image_id message_sender_profile_image_id, mm.grade_id message_sender_grade_id, mm.group_id message_sender_group_id,
t_images.height message_sender_profile_image_height, t_images.width message_sender_profile_image_width, t_images.attachable_item_id message_sender_profile_image_attachable_item_id,
t_attachable_items.owner_id message_sender_profile_image_owner_id, t_attachable_items.from_outer message_sender_profile_image_from_outer, t_attachable_items.alias message_sender_profile_image_alias,
t_attachable_items.url message_sender_profile_image_url, t_attachable_items.size message_sender_profile_image_size, t_attachable_items.mime_type_id message_sender_profile_image_mime_type_id
FROM t_chat_room_actions
LEFT JOIN t_chat_room_create_actions ON t_chat_room_actions.chat_room_action_id = t_chat_room_create_actions.chat_room_action_id
LEFT JOIN m_members cm ON t_chat_room_create_actions.created_by = cm.member_id
LEFT JOIN t_chat_room_update_name_actions ON t_chat_room_actions.chat_room_action_id = t_chat_room_update_name_actions.chat_room_action_id
LEFT JOIN m_members um ON t_chat_room_update_name_actions.updated_by = um.member_id
LEFT JOIN t_chat_room_withdraw_actions ON t_chat_room_actions.chat_room_action_id = t_chat_room_withdraw_actions.chat_room_action_id
LEFT JOIN m_members wm ON t_chat_room_withdraw_actions.member_id = wm.member_id
LEFT JOIN t_chat_room_add_member_actions ON t_chat_room_actions.chat_room_action_id = t_chat_room_add_member_actions.chat_room_action_id
LEFT JOIN m_members am ON t_chat_room_add_member_actions.added_by = am.member_id
LEFT JOIN t_chat_room_remove_member_actions ON t_chat_room_actions.chat_room_action_id = t_chat_room_remove_member_actions.chat_room_action_id
LEFT JOIN m_members rm ON t_chat_room_remove_member_actions.removed_by = rm.member_id
LEFT JOIN t_chat_room_delete_message_actions ON t_chat_room_actions.chat_room_action_id = t_chat_room_delete_message_actions.chat_room_action_id
LEFT JOIN m_members dm ON t_chat_room_delete_message_actions.deleted_by = dm.member_id
LEFT JOIN t_messages ON t_chat_room_actions.chat_room_action_id = t_messages.chat_room_action_id
LEFT JOIN m_members mm ON t_messages.sender_id = mm.member_id
LEFT JOIN t_images ON mm.profile_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE t_chat_room_actions.chat_room_id = $1
AND
	CASE WHEN @where_in_chat_room_action_type_ids::boolean = true THEN chat_room_action_type_id = ANY(@in_chat_room_action_type_ids::uuid[]) ELSE TRUE END
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'acted_at' THEN acted_at > @acted_at_cursor OR (acted_at = @acted_at_cursor AND t_chat_room_actions_pkey > @cursor::int)
				WHEN 'r_acted_at' THEN acted_at < @acted_at_cursor OR (acted_at = @acted_at_cursor AND t_chat_room_actions_pkey > @cursor::int)
				ELSE t_chat_room_actions_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'acted_at' THEN acted_at < @acted_at_cursor OR (acted_at = @acted_at_cursor AND t_chat_room_actions_pkey < @cursor::int)
				WHEN 'r_acted_at' THEN acted_at > @acted_at_cursor OR (acted_at = @acted_at_cursor AND t_chat_room_actions_pkey < @cursor::int)
				ELSE t_chat_room_actions_pkey < @cursor::int
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'acted_at' AND @cursor_direction::text = 'next' THEN acted_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'acted_at' AND @cursor_direction::text = 'prev' THEN acted_at END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'r_acted_at' AND @cursor_direction::text = 'next' THEN acted_at END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'r_acted_at' AND @cursor_direction::text = 'prev' THEN acted_at END ASC NULLS LAST,
	t_chat_room_actions_pkey ASC
LIMIT $2;

-- name: GetPluralChatRoomActionsWithDetail :many
SELECT t_chat_room_actions.*,
t_chat_room_create_actions.chat_room_create_action_id, t_chat_room_create_actions.name create_name, cm.member_id create_member_id, cm.name create_member_name, cm.first_name create_member_first_name, cm.last_name create_member_last_name, cm.email create_member_email, cm.profile_image_id create_member_profile_image_id, cm.grade_id create_member_grade_id, cm.group_id create_member_group_id,
t_chat_room_update_name_actions.chat_room_update_name_action_id, t_chat_room_update_name_actions.name update_name, um.member_id update_member_id, um.name update_member_name, um.first_name update_member_first_name, um.last_name update_member_last_name, um.email update_member_email, um.profile_image_id update_member_profile_image_id, um.grade_id update_member_grade_id, um.group_id update_member_group_id,
t_chat_room_withdraw_actions.chat_room_withdraw_action_id, wm.member_id withdraw_member_id, wm.name withdraw_member_name, wm.first_name withdraw_member_first_name, wm.last_name withdraw_member_last_name, wm.email withdraw_member_email, wm.profile_image_id withdraw_member_profile_image_id, wm.grade_id withdraw_member_grade_id, wm.group_id withdraw_member_group_id,
t_chat_room_add_member_actions.chat_room_add_member_action_id, am.member_id add_member_id, am.name add_member_name, am.first_name add_member_first_name, am.last_name add_member_last_name, am.email add_member_email, am.profile_image_id add_member_profile_image_id, am.grade_id add_member_grade_id, am.group_id add_member_group_id,
t_chat_room_remove_member_actions.chat_room_remove_member_action_id, rm.member_id remove_member_id, rm.name remove_member_name, rm.first_name remove_member_first_name, rm.last_name remove_member_last_name, rm.email remove_member_email, rm.profile_image_id remove_member_profile_image_id, rm.grade_id remove_member_grade_id, rm.group_id remove_member_group_id,
t_chat_room_delete_message_actions.chat_room_delete_message_action_id, dm.member_id delete_message_member_id, dm.name delete_message_member_name, dm.first_name delete_message_member_first_name, dm.last_name delete_message_member_last_name, dm.email delete_message_member_email, dm.profile_image_id delete_message_member_profile_image_id, dm.grade_id delete_message_member_grade_id, dm.group_id delete_message_member_group_id,
t_messages.message_id, t_messages.sender_id message_sender_id, t_messages.body message_body, t_messages.posted_at message_posted_at, t_messages.last_edited_at message_last_edited_at,
mm.name message_sender_name, mm.first_name message_sender_first_name, mm.last_name message_sender_last_name, mm.email message_sender_email, mm.profile_image_id message_sender_profile_image_id, mm.grade_id message_sender_grade_id, mm.group_id message_sender_group_id,
t_images.height message_sender_profile_image_height, t_images.width message_sender_profile_image_width, t_images.attachable_item_id message_sender_profile_image_attachable_item_id,
t_attachable_items.owner_id message_sender_profile_image_owner_id, t_attachable_items.from_outer message_sender_profile_image_from_outer, t_attachable_items.alias message_sender_profile_image_alias,
t_attachable_items.url message_sender_profile_image_url, t_attachable_items.size message_sender_profile_image_size, t_attachable_items.mime_type_id message_sender_profile_image_mime_type_id
FROM t_chat_room_actions
LEFT JOIN t_chat_room_create_actions ON t_chat_room_actions.chat_room_action_id = t_chat_room_create_actions.chat_room_action_id
LEFT JOIN m_members cm ON t_chat_room_create_actions.created_by = cm.member_id
LEFT JOIN t_chat_room_update_name_actions ON t_chat_room_actions.chat_room_action_id = t_chat_room_update_name_actions.chat_room_action_id
LEFT JOIN m_members um ON t_chat_room_update_name_actions.updated_by = um.member_id
LEFT JOIN t_chat_room_withdraw_actions ON t_chat_room_actions.chat_room_action_id = t_chat_room_withdraw_actions.chat_room_action_id
LEFT JOIN m_members wm ON t_chat_room_withdraw_actions.member_id = wm.member_id
LEFT JOIN t_chat_room_add_member_actions ON t_chat_room_actions.chat_room_action_id = t_chat_room_add_member_actions.chat_room_action_id
LEFT JOIN m_members am ON t_chat_room_add_member_actions.added_by = am.member_id
LEFT JOIN t_chat_room_remove_member_actions ON t_chat_room_actions.chat_room_action_id = t_chat_room_remove_member_actions.chat_room_action_id
LEFT JOIN m_members rm ON t_chat_room_remove_member_actions.removed_by = rm.member_id
LEFT JOIN t_chat_room_delete_message_actions ON t_chat_room_actions.chat_room_action_id = t_chat_room_delete_message_actions.chat_room_action_id
LEFT JOIN m_members dm ON t_chat_room_delete_message_actions.deleted_by = dm.member_id
LEFT JOIN t_messages ON t_chat_room_actions.chat_room_action_id = t_messages.chat_room_action_id
LEFT JOIN m_members mm ON t_messages.sender_id = mm.member_id
LEFT JOIN t_images ON mm.profile_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE chat_room_action_id = ANY(@chat_room_action_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'acted_at' THEN acted_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_acted_at' THEN acted_at END DESC NULLS LAST,
	t_chat_room_actions_pkey ASC;

-- name: GetPluralChatRoomActionsWithDetailUseNumberedPaginate :many
SELECT t_chat_room_actions.*,
t_chat_room_create_actions.chat_room_create_action_id, t_chat_room_create_actions.name create_name, cm.member_id create_member_id, cm.name create_member_name, cm.first_name create_member_first_name, cm.last_name create_member_last_name, cm.email create_member_email, cm.profile_image_id create_member_profile_image_id, cm.grade_id create_member_grade_id, cm.group_id create_member_group_id,
t_chat_room_update_name_actions.chat_room_update_name_action_id, t_chat_room_update_name_actions.name update_name, um.member_id update_member_id, um.name update_member_name, um.first_name update_member_first_name, um.last_name update_member_last_name, um.email update_member_email, um.profile_image_id update_member_profile_image_id, um.grade_id update_member_grade_id, um.group_id update_member_group_id,
t_chat_room_withdraw_actions.chat_room_withdraw_action_id, wm.member_id withdraw_member_id, wm.name withdraw_member_name, wm.first_name withdraw_member_first_name, wm.last_name withdraw_member_last_name, wm.email withdraw_member_email, wm.profile_image_id withdraw_member_profile_image_id, wm.grade_id withdraw_member_grade_id, wm.group_id withdraw_member_group_id,
t_chat_room_add_member_actions.chat_room_add_member_action_id, am.member_id add_member_id, am.name add_member_name, am.first_name add_member_first_name, am.last_name add_member_last_name, am.email add_member_email, am.profile_image_id add_member_profile_image_id, am.grade_id add_member_grade_id, am.group_id add_member_group_id,
t_chat_room_remove_member_actions.chat_room_remove_member_action_id, rm.member_id remove_member_id, rm.name remove_member_name, rm.first_name remove_member_first_name, rm.last_name remove_member_last_name, rm.email remove_member_email, rm.profile_image_id remove_member_profile_image_id, rm.grade_id remove_member_grade_id, rm.group_id remove_member_group_id,
t_chat_room_delete_message_actions.chat_room_delete_message_action_id, dm.member_id delete_message_member_id, dm.name delete_message_member_name, dm.first_name delete_message_member_first_name, dm.last_name delete_message_member_last_name, dm.email delete_message_member_email, dm.profile_image_id delete_message_member_profile_image_id, dm.grade_id delete_message_member_grade_id, dm.group_id delete_message_member_group_id,
t_messages.message_id, t_messages.sender_id message_sender_id, t_messages.body message_body, t_messages.posted_at message_posted_at, t_messages.last_edited_at message_last_edited_at,
mm.name message_sender_name, mm.first_name message_sender_first_name, mm.last_name message_sender_last_name, mm.email message_sender_email, mm.profile_image_id message_sender_profile_image_id, mm.grade_id message_sender_grade_id, mm.group_id message_sender_group_id,
t_images.height message_sender_profile_image_height, t_images.width message_sender_profile_image_width, t_images.attachable_item_id message_sender_profile_image_attachable_item_id,
t_attachable_items.owner_id message_sender_profile_image_owner_id, t_attachable_items.from_outer message_sender_profile_image_from_outer, t_attachable_items.alias message_sender_profile_image_alias,
t_attachable_items.url message_sender_profile_image_url, t_attachable_items.size message_sender_profile_image_size, t_attachable_items.mime_type_id message_sender_profile_image_mime_type_id
FROM t_chat_room_actions
LEFT JOIN t_chat_room_create_actions ON t_chat_room_actions.chat_room_action_id = t_chat_room_create_actions.chat_room_action_id
LEFT JOIN m_members cm ON t_chat_room_create_actions.created_by = cm.member_id
LEFT JOIN t_chat_room_update_name_actions ON t_chat_room_actions.chat_room_action_id = t_chat_room_update_name_actions.chat_room_action_id
LEFT JOIN m_members um ON t_chat_room_update_name_actions.updated_by = um.member_id
LEFT JOIN t_chat_room_withdraw_actions ON t_chat_room_actions.chat_room_action_id = t_chat_room_withdraw_actions.chat_room_action_id
LEFT JOIN m_members wm ON t_chat_room_withdraw_actions.member_id = wm.member_id
LEFT JOIN t_chat_room_add_member_actions ON t_chat_room_actions.chat_room_action_id = t_chat_room_add_member_actions.chat_room_action_id
LEFT JOIN m_members am ON t_chat_room_add_member_actions.added_by = am.member_id
LEFT JOIN t_chat_room_remove_member_actions ON t_chat_room_actions.chat_room_action_id = t_chat_room_remove_member_actions.chat_room_action_id
LEFT JOIN m_members rm ON t_chat_room_remove_member_actions.removed_by = rm.member_id
LEFT JOIN t_chat_room_delete_message_actions ON t_chat_room_actions.chat_room_action_id = t_chat_room_delete_message_actions.chat_room_action_id
LEFT JOIN m_members dm ON t_chat_room_delete_message_actions.deleted_by = dm.member_id
LEFT JOIN t_messages ON t_chat_room_actions.chat_room_action_id = t_messages.chat_room_action_id
LEFT JOIN m_members mm ON t_messages.sender_id = mm.member_id
LEFT JOIN t_images ON mm.profile_image_id = t_images.image_id
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
WHERE chat_room_action_id = ANY(@chat_room_action_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'acted_at' THEN acted_at END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_acted_at' THEN acted_at END DESC NULLS LAST,
	t_chat_room_actions_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountChatRoomActions :one
SELECT COUNT(*) FROM t_chat_room_actions
WHERE
	CASE WHEN @where_in_chat_room_action_type_ids::boolean = true THEN chat_room_action_type_id = ANY(@in_chat_room_action_type_ids::uuid[]) ELSE TRUE END;

-- name: CountChatRoomActionsOnChatRoom :one
SELECT COUNT(*) FROM t_chat_room_actions
WHERE chat_room_id = $1
AND
	CASE WHEN @where_in_chat_room_action_type_ids::boolean = true THEN chat_room_action_type_id = ANY(@in_chat_room_action_type_ids::uuid[]) ELSE TRUE END;
