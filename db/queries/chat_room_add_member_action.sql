-- name: CreateChatRoomAddMemberActions :copyfrom
INSERT INTO t_chat_room_add_member_actions (chat_room_action_id, added_by) VALUES ($1, $2);

-- name: CreateChatRoomAddMemberAction :one
INSERT INTO t_chat_room_add_member_actions (chat_room_action_id, added_by) VALUES ($1, $2) RETURNING *;

-- name: DeleteChatRoomAddMemberAction :execrows
DELETE FROM t_chat_room_add_member_actions WHERE chat_room_add_member_action_id = $1;

-- name: PluralDeleteChatRoomAddMemberActions :execrows
DELETE FROM t_chat_room_add_member_actions WHERE chat_room_add_member_action_id = ANY(@chat_room_add_member_action_ids::uuid[]);

-- name: GetChatRoomAddMemberActionsOnChatRoom :many
SELECT t_chat_room_add_member_actions.*,
m_members.name add_member_name, m_members.first_name add_member_first_name, m_members.last_name add_member_last_name, m_members.email add_member_email,
m_members.profile_image_id add_member_profile_image_id
FROM t_chat_room_add_member_actions
LEFT JOIN m_members ON t_chat_room_add_member_actions.added_by = m_members.member_id
WHERE EXISTS (
	SELECT 1 FROM t_chat_room_actions WHERE chat_room_action_id = t_chat_room_add_member_actions.chat_room_action_id AND t_chat_room_actions.chat_room_id = $1
)
ORDER BY
	t_chat_room_add_member_actions_pkey ASC;

-- name: GetChatRoomAddMemberActionsOnChatRoomUseNumberedPaginate :many
SELECT t_chat_room_add_member_actions.*,
m_members.name add_member_name, m_members.first_name add_member_first_name, m_members.last_name add_member_last_name, m_members.email add_member_email,
m_members.profile_image_id add_member_profile_image_id
FROM t_chat_room_add_member_actions
LEFT JOIN m_members ON t_chat_room_add_member_actions.added_by = m_members.member_id
WHERE EXISTS (
	SELECT 1 FROM t_chat_room_actions WHERE chat_room_action_id = t_chat_room_add_member_actions.chat_room_action_id AND chat_room_id = $1
)
ORDER BY
	t_chat_room_add_member_actions_pkey ASC
LIMIT $2 OFFSET $3;

-- name: GetChatRoomAddMemberActionsOnChatRoomUseKeysetPaginate :many
SELECT t_chat_room_add_member_actions.*,
m_members.name add_member_name, m_members.first_name add_member_first_name, m_members.last_name add_member_last_name, m_members.email add_member_email,
m_members.profile_image_id add_member_profile_image_id
FROM t_chat_room_add_member_actions
LEFT JOIN m_members ON t_chat_room_add_member_actions.added_by = m_members.member_id
WHERE EXISTS (
	SELECT 1 FROM t_chat_room_actions WHERE chat_room_action_id = t_chat_room_add_member_actions.chat_room_action_id AND chat_room_id = $1
)
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			t_chat_room_add_member_actions_pkey > @cursor::int
		WHEN 'prev' THEN
			t_chat_room_add_member_actions_pkey < @cursor::int
	END
ORDER BY
	CASE WHEN @cursor_direction::text = 'next' THEN t_chat_room_add_member_actions_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN t_chat_room_add_member_actions_pkey END DESC
LIMIT $2;

-- name: GetPluralChatRoomAddMemberActions :many
SELECT t_chat_room_add_member_actions.*,
m_members.name add_member_name, m_members.first_name add_member_first_name, m_members.last_name add_member_last_name, m_members.email add_member_email,
m_members.profile_image_id add_member_profile_image_id
FROM t_chat_room_add_member_actions
LEFT JOIN m_members ON t_chat_room_add_member_actions.added_by = m_members.member_id
WHERE chat_room_add_member_action_id = ANY(@chat_room_add_member_action_ids::uuid[])
ORDER BY
	t_chat_room_add_member_actions_pkey ASC;

-- name: GetPluralChatRoomAddMemberActionsUseNumberedPaginate :many
SELECT t_chat_room_add_member_actions.*,
m_members.name add_member_name, m_members.first_name add_member_first_name, m_members.last_name add_member_last_name, m_members.email add_member_email,
m_members.profile_image_id add_member_profile_image_id
FROM t_chat_room_add_member_actions
LEFT JOIN m_members ON t_chat_room_add_member_actions.added_by = m_members.member_id
WHERE chat_room_add_member_action_id = ANY(@chat_room_add_member_action_ids::uuid[])
ORDER BY
	t_chat_room_add_member_actions_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountChatRoomAddMemberActions :one
SELECT COUNT(*) FROM t_chat_room_add_member_actions;
