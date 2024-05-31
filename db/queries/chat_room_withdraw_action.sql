-- name: CreateChatRoomWithdrawActions :copyfrom
INSERT INTO t_chat_room_withdraw_actions (chat_room_action_id, member_id) VALUES ($1, $2);

-- name: CreateChatRoomWithdrawAction :one
INSERT INTO t_chat_room_withdraw_actions (chat_room_action_id, member_id) VALUES ($1, $2) RETURNING *;

-- name: DeleteChatRoomWithdrawAction :execrows
DELETE FROM t_chat_room_withdraw_actions WHERE chat_room_withdraw_action_id = $1;

-- name: PluralDeleteChatRoomWithdrawActions :execrows
DELETE FROM t_chat_room_withdraw_actions WHERE chat_room_withdraw_action_id = ANY(@chat_room_withdraw_action_ids::uuid[]);

-- name: GetChatRoomWithdrawActionsOnChatRoom :many
SELECT t_chat_room_withdraw_actions.*,
m_members.name withdraw_member_name, m_members.first_name withdraw_member_first_name, m_members.last_name withdraw_member_last_name, m_members.email withdraw_member_email,
m_members.profile_image_id withdraw_member_profile_image_id, m_members.grade_id withdraw_member_grade_id, m_members.group_id withdraw_member_group_id
FROM t_chat_room_withdraw_actions
LEFT JOIN m_members ON t_chat_room_withdraw_actions.member_id = m_members.member_id
WHERE EXISTS (
	SELECT 1 FROM t_chat_room_actions WHERE chat_room_action_id = t_chat_room_withdraw_actions.chat_room_action_id AND t_chat_room_actions.chat_room_id = $1
)
ORDER BY
	t_chat_room_withdraw_actions_pkey ASC;

-- name: GetChatRoomWithdrawActionsOnChatRoomUseNumberedPaginate :many
SELECT t_chat_room_withdraw_actions.*,
m_members.name withdraw_member_name, m_members.first_name withdraw_member_first_name, m_members.last_name withdraw_member_last_name, m_members.email withdraw_member_email,
m_members.profile_image_id withdraw_member_profile_image_id, m_members.grade_id withdraw_member_grade_id, m_members.group_id withdraw_member_group_id
FROM t_chat_room_withdraw_actions
LEFT JOIN m_members ON t_chat_room_withdraw_actions.member_id = m_members.member_id
WHERE EXISTS (
	SELECT 1 FROM t_chat_room_actions WHERE chat_room_action_id = t_chat_room_withdraw_actions.chat_room_action_id AND chat_room_id = $1
)
ORDER BY
	t_chat_room_withdraw_actions_pkey ASC
LIMIT $2 OFFSET $3;

-- name: GetChatRoomWithdrawActionsOnChatRoomUseKeysetPaginate :many
SELECT t_chat_room_withdraw_actions.*,
m_members.name withdraw_member_name, m_members.first_name withdraw_member_first_name, m_members.last_name withdraw_member_last_name, m_members.email withdraw_member_email,
m_members.profile_image_id withdraw_member_profile_image_id, m_members.grade_id withdraw_member_grade_id, m_members.group_id withdraw_member_group_id
FROM t_chat_room_withdraw_actions
LEFT JOIN m_members ON t_chat_room_withdraw_actions.member_id = m_members.member_id
WHERE EXISTS (
	SELECT 1 FROM t_chat_room_actions WHERE chat_room_action_id = t_chat_room_withdraw_actions.chat_room_action_id AND chat_room_id = $1
)
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			t_chat_room_withdraw_actions_pkey > @cursor::int
		WHEN 'prev' THEN
			t_chat_room_withdraw_actions_pkey < @cursor::int
	END
ORDER BY
	CASE WHEN @cursor_direction::text = 'next' THEN t_chat_room_withdraw_actions_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN t_chat_room_withdraw_actions_pkey END DESC
LIMIT $2;

-- name: GetPluralChatRoomWithdrawActions :many
SELECT t_chat_room_withdraw_actions.*,
m_members.name withdraw_member_name, m_members.first_name withdraw_member_first_name, m_members.last_name withdraw_member_last_name, m_members.email withdraw_member_email,
m_members.profile_image_id withdraw_member_profile_image_id, m_members.grade_id withdraw_member_grade_id, m_members.group_id withdraw_member_group_id
FROM t_chat_room_withdraw_actions
LEFT JOIN m_members ON t_chat_room_withdraw_actions.member_id = m_members.member_id
WHERE chat_room_withdraw_action_id = ANY(@chat_room_withdraw_action_ids::uuid[])
ORDER BY
	t_chat_room_withdraw_actions_pkey ASC;

-- name: GetPluralChatRoomWithdrawActionsUseNumberedPaginate :many
SELECT t_chat_room_withdraw_actions.*,
m_members.name withdraw_member_name, m_members.first_name withdraw_member_first_name, m_members.last_name withdraw_member_last_name, m_members.email withdraw_member_email,
m_members.profile_image_id withdraw_member_profile_image_id, m_members.grade_id withdraw_member_grade_id, m_members.group_id withdraw_member_group_id
FROM t_chat_room_withdraw_actions
LEFT JOIN m_members ON t_chat_room_withdraw_actions.member_id = m_members.member_id
WHERE chat_room_withdraw_action_id = ANY(@chat_room_withdraw_action_ids::uuid[])
ORDER BY
	t_chat_room_withdraw_actions_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountChatRoomWithdrawActions :one
SELECT COUNT(*) FROM t_chat_room_withdraw_actions;
