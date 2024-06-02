-- name: CreateChatRoomDeleteMessageActions :copyfrom
INSERT INTO t_chat_room_delete_message_actions (chat_room_action_id, deleted_by) VALUES ($1, $2);

-- name: CreateChatRoomDeleteMessageAction :one
INSERT INTO t_chat_room_delete_message_actions (chat_room_action_id, deleted_by) VALUES ($1, $2) RETURNING *;

-- name: DeleteChatRoomDeleteMessageAction :execrows
DELETE FROM t_chat_room_delete_message_actions WHERE chat_room_delete_message_action_id = $1;

-- name: PluralDeleteChatRoomDeleteMessageActions :execrows
DELETE FROM t_chat_room_delete_message_actions WHERE chat_room_delete_message_action_id = ANY(@chat_room_delete_message_action_ids::uuid[]);

-- name: GetChatRoomDeleteMessageActionsOnChatRoom :many
SELECT t_chat_room_delete_message_actions.*,
m_members.name delete_message_member_name, m_members.first_name delete_message_member_first_name, m_members.last_name delete_message_member_last_name, m_members.email delete_message_member_email,
m_members.profile_image_id delete_message_member_profile_image_id, m_members.grade_id delete_message_member_grade_id, m_members.group_id delete_message_member_group_id
FROM t_chat_room_delete_message_actions
LEFT JOIN m_members ON t_chat_room_delete_message_actions.deleted_by = m_members.member_id
WHERE EXISTS (
	SELECT 1 FROM t_chat_room_actions WHERE chat_room_action_id = t_chat_room_delete_message_actions.chat_room_action_id AND t_chat_room_actions.chat_room_id = $1
)
ORDER BY
	t_chat_room_delete_message_actions_pkey ASC;

-- name: GetChatRoomDeleteMessageActionsOnChatRoomUseNumberedPaginate :many
SELECT t_chat_room_delete_message_actions.*,
m_members.name delete_message_member_name, m_members.first_name delete_message_member_first_name, m_members.last_name delete_message_member_last_name, m_members.email delete_message_member_email,
m_members.profile_image_id delete_message_member_profile_image_id, m_members.grade_id delete_message_member_grade_id, m_members.group_id delete_message_member_group_id
FROM t_chat_room_delete_message_actions
LEFT JOIN m_members ON t_chat_room_delete_message_actions.deleted_by = m_members.member_id
WHERE EXISTS (
	SELECT 1 FROM t_chat_room_actions WHERE chat_room_action_id = t_chat_room_delete_message_actions.chat_room_action_id AND chat_room_id = $1
)
ORDER BY
	t_chat_room_delete_message_actions_pkey ASC
LIMIT $2 OFFSET $3;

-- name: GetChatRoomDeleteMessageActionsOnChatRoomUseKeysetPaginate :many
SELECT t_chat_room_delete_message_actions.*,
m_members.name delete_message_member_name, m_members.first_name delete_message_member_first_name, m_members.last_name delete_message_member_last_name, m_members.email delete_message_member_email,
m_members.profile_image_id delete_message_member_profile_image_id, m_members.grade_id delete_message_member_grade_id, m_members.group_id delete_message_member_group_id
FROM t_chat_room_delete_message_actions
LEFT JOIN m_members ON t_chat_room_delete_message_actions.deleted_by = m_members.member_id
WHERE EXISTS (
	SELECT 1 FROM t_chat_room_actions WHERE chat_room_action_id = t_chat_room_delete_message_actions.chat_room_action_id AND chat_room_id = $1
)
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			t_chat_room_delete_message_actions_pkey > @cursor::int
		WHEN 'prev' THEN
			t_chat_room_delete_message_actions_pkey < @cursor::int
	END
ORDER BY
	CASE WHEN @cursor_direction::text = 'next' THEN t_chat_room_delete_message_actions_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN t_chat_room_delete_message_actions_pkey END DESC
LIMIT $2;

-- name: GetPluralChatRoomDeleteMessageActions :many
SELECT t_chat_room_delete_message_actions.*,
m_members.name delete_message_member_name, m_members.first_name delete_message_member_first_name, m_members.last_name delete_message_member_last_name, m_members.email delete_message_member_email,
m_members.profile_image_id delete_message_member_profile_image_id, m_members.grade_id delete_message_member_grade_id, m_members.group_id delete_message_member_group_id
FROM t_chat_room_delete_message_actions
LEFT JOIN m_members ON t_chat_room_delete_message_actions.deleted_by = m_members.member_id
WHERE chat_room_delete_message_action_id = ANY(@chat_room_delete_message_action_ids::uuid[])
ORDER BY
	t_chat_room_delete_message_actions_pkey ASC;

-- name: GetPluralChatRoomDeleteMessageActionsUseNumberedPaginate :many
SELECT t_chat_room_delete_message_actions.*,
m_members.name delete_message_member_name, m_members.first_name delete_message_member_first_name, m_members.last_name delete_message_member_last_name, m_members.email delete_message_member_email,
m_members.profile_image_id delete_message_member_profile_image_id, m_members.grade_id delete_message_member_grade_id, m_members.group_id delete_message_member_group_id
FROM t_chat_room_delete_message_actions
LEFT JOIN m_members ON t_chat_room_delete_message_actions.deleted_by = m_members.member_id
WHERE chat_room_delete_message_action_id = ANY(@chat_room_delete_message_action_ids::uuid[])
ORDER BY
	t_chat_room_delete_message_actions_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetPluralChatRoomDeleteMessageActionsByChatRoomActionIDs :many
SELECT t_chat_room_delete_message_actions.*,
m_members.name delete_message_member_name, m_members.first_name delete_message_member_first_name, m_members.last_name delete_message_member_last_name, m_members.email delete_message_member_email,
m_members.profile_image_id delete_message_member_profile_image_id, m_members.grade_id delete_message_member_grade_id, m_members.group_id delete_message_member_group_id
FROM t_chat_room_delete_message_actions
LEFT JOIN m_members ON t_chat_room_delete_message_actions.deleted_by = m_members.member_id
WHERE chat_room_action_id = ANY(@chat_room_action_ids::uuid[])
ORDER BY
	t_chat_room_delete_message_actions_pkey ASC;

-- name: GetPluralChatRoomDeleteMessageActionsByChatRoomActionIDsUseNumberedPaginate :many
SELECT t_chat_room_delete_message_actions.*,
m_members.name delete_message_member_name, m_members.first_name delete_message_member_first_name, m_members.last_name delete_message_member_last_name, m_members.email delete_message_member_email,
m_members.profile_image_id delete_message_member_profile_image_id, m_members.grade_id delete_message_member_grade_id, m_members.group_id delete_message_member_group_id
FROM t_chat_room_delete_message_actions
LEFT JOIN m_members ON t_chat_room_delete_message_actions.deleted_by = m_members.member_id
WHERE chat_room_action_id = ANY(@chat_room_action_ids::uuid[])
ORDER BY
	t_chat_room_delete_message_actions_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountChatRoomDeleteMessageActions :one
SELECT COUNT(*) FROM t_chat_room_delete_message_actions;
