-- name: CreateChatRoomUpdateNameActions :copyfrom
INSERT INTO t_chat_room_update_name_actions (chat_room_action_id, updated_by, name) VALUES ($1, $2, $3);

-- name: CreateChatRoomUpdateNameAction :one
INSERT INTO t_chat_room_update_name_actions (chat_room_action_id, updated_by, name) VALUES ($1, $2, $3) RETURNING *;

-- name: DeleteChatRoomUpdateNameAction :execrows
DELETE FROM t_chat_room_update_name_actions WHERE chat_room_update_name_action_id = $1;

-- name: PluralDeleteChatRoomUpdateNameActions :execrows
DELETE FROM t_chat_room_update_name_actions WHERE chat_room_update_name_action_id = ANY(@chat_room_update_name_action_ids::uuid[]);

-- name: GetChatRoomUpdateNameActionsOnChatRoom :many
SELECT t_chat_room_update_name_actions.*,
m_members.name update_member_name, m_members.first_name update_member_first_name, m_members.last_name update_member_last_name, m_members.email update_member_email,
m_members.profile_image_id update_member_profile_image_id, m_members.grade_id update_member_grade_id, m_members.group_id update_member_group_id
FROM t_chat_room_update_name_actions
LEFT JOIN m_members ON t_chat_room_update_name_actions.updated_by = m_members.member_id
WHERE EXISTS (
	SELECT 1 FROM t_chat_room_actions WHERE chat_room_action_id = t_chat_room_update_name_actions.chat_room_action_id AND t_chat_room_actions.chat_room_id = $1
)
ORDER BY
	t_chat_room_update_name_actions_pkey ASC;

-- name: GetChatRoomUpdateNameActionsOnChatRoomUseNumberedPaginate :many
SELECT t_chat_room_update_name_actions.*,
m_members.name update_member_name, m_members.first_name update_member_first_name, m_members.last_name update_member_last_name, m_members.email update_member_email,
m_members.profile_image_id update_member_profile_image_id, m_members.grade_id update_member_grade_id, m_members.group_id update_member_group_id
FROM t_chat_room_update_name_actions
LEFT JOIN m_members ON t_chat_room_update_name_actions.updated_by = m_members.member_id
WHERE EXISTS (
	SELECT 1 FROM t_chat_room_actions WHERE chat_room_action_id = t_chat_room_update_name_actions.chat_room_action_id AND chat_room_id = $1
)
ORDER BY
	t_chat_room_update_name_actions_pkey ASC
LIMIT $2 OFFSET $3;

-- name: GetChatRoomUpdateNameActionsOnChatRoomUseKeysetPaginate :many
SELECT t_chat_room_update_name_actions.*,
m_members.name update_member_name, m_members.first_name update_member_first_name, m_members.last_name update_member_last_name, m_members.email update_member_email,
m_members.profile_image_id update_member_profile_image_id, m_members.grade_id update_member_grade_id, m_members.group_id update_member_group_id
FROM t_chat_room_update_name_actions
LEFT JOIN m_members ON t_chat_room_update_name_actions.updated_by = m_members.member_id
WHERE EXISTS (
	SELECT 1 FROM t_chat_room_actions WHERE chat_room_action_id = t_chat_room_update_name_actions.chat_room_action_id AND chat_room_id = $1
)
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			t_chat_room_update_name_actions_pkey > @cursor::int
		WHEN 'prev' THEN
			t_chat_room_update_name_actions_pkey < @cursor::int
	END
ORDER BY
	CASE WHEN @cursor_direction::text = 'next' THEN t_chat_room_update_name_actions_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN t_chat_room_update_name_actions_pkey END DESC
LIMIT $2;

-- name: GetPluralChatRoomUpdateNameActions :many
SELECT t_chat_room_update_name_actions.*,
m_members.name update_member_name, m_members.first_name update_member_first_name, m_members.last_name update_member_last_name, m_members.email update_member_email,
m_members.profile_image_id update_member_profile_image_id, m_members.grade_id update_member_grade_id, m_members.group_id update_member_group_id
FROM t_chat_room_update_name_actions
LEFT JOIN m_members ON t_chat_room_update_name_actions.updated_by = m_members.member_id
WHERE chat_room_update_name_action_id = ANY(@chat_room_update_name_action_ids::uuid[])
ORDER BY
	t_chat_room_update_name_actions_pkey ASC;

-- name: GetPluralChatRoomUpdateNameActionsUseNumberedPaginate :many
SELECT t_chat_room_update_name_actions.*,
m_members.name update_member_name, m_members.first_name update_member_first_name, m_members.last_name update_member_last_name, m_members.email update_member_email,
m_members.profile_image_id update_member_profile_image_id, m_members.grade_id update_member_grade_id, m_members.group_id update_member_group_id
FROM t_chat_room_update_name_actions
LEFT JOIN m_members ON t_chat_room_update_name_actions.updated_by = m_members.member_id
WHERE chat_room_update_name_action_id = ANY(@chat_room_update_name_action_ids::uuid[])
ORDER BY
	t_chat_room_update_name_actions_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountChatRoomUpdateNameActions :one
SELECT COUNT(*) FROM t_chat_room_update_name_actions;
