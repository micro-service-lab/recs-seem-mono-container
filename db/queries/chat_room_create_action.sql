-- name: CreateChatRoomCreateActions :copyfrom
INSERT INTO t_chat_room_create_actions (chat_room_action_id, created_by, name) VALUES ($1, $2, $3);

-- name: CreateChatRoomCreateAction :one
INSERT INTO t_chat_room_create_actions (chat_room_action_id, created_by, name) VALUES ($1, $2, $3) RETURNING *;

-- name: DeleteChatRoomCreateAction :execrows
DELETE FROM t_chat_room_create_actions WHERE chat_room_create_action_id = $1;

-- name: PluralDeleteChatRoomCreateActions :execrows
DELETE FROM t_chat_room_create_actions WHERE chat_room_create_action_id = ANY(@chat_room_create_action_ids::uuid[]);

-- name: GetChatRoomCreateActionsOnChatRoom :many
SELECT t_chat_room_create_actions.*,
m_members.name create_member_name, m_members.first_name create_member_first_name, m_members.last_name create_member_last_name, m_members.email create_member_email,
m_members.profile_image_id create_member_profile_image_id
FROM t_chat_room_create_actions
LEFT JOIN m_members ON t_chat_room_create_actions.created_by = m_members.member_id
WHERE EXISTS (
	SELECT 1 FROM t_chat_room_actions WHERE chat_room_action_id = t_chat_room_create_actions.chat_room_action_id AND t_chat_room_actions.chat_room_id = $1
)
ORDER BY
	t_chat_room_create_actions_pkey ASC;

-- name: GetChatRoomCreateActionsOnChatRoomUseNumberedPaginate :many
SELECT t_chat_room_create_actions.*,
m_members.name create_member_name, m_members.first_name create_member_first_name, m_members.last_name create_member_last_name, m_members.email create_member_email,
m_members.profile_image_id create_member_profile_image_id
FROM t_chat_room_create_actions
LEFT JOIN m_members ON t_chat_room_create_actions.created_by = m_members.member_id
WHERE EXISTS (
	SELECT 1 FROM t_chat_room_actions WHERE chat_room_action_id = t_chat_room_create_actions.chat_room_action_id AND chat_room_id = $1
)
ORDER BY
	t_chat_room_create_actions_pkey ASC
LIMIT $2 OFFSET $3;

-- name: GetChatRoomCreateActionsOnChatRoomUseKeysetPaginate :many
SELECT t_chat_room_create_actions.*,
m_members.name create_member_name, m_members.first_name create_member_first_name, m_members.last_name create_member_last_name, m_members.email create_member_email,
m_members.profile_image_id create_member_profile_image_id
FROM t_chat_room_create_actions
LEFT JOIN m_members ON t_chat_room_create_actions.created_by = m_members.member_id
WHERE EXISTS (
	SELECT 1 FROM t_chat_room_actions WHERE chat_room_action_id = t_chat_room_create_actions.chat_room_action_id AND chat_room_id = $1
)
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			t_chat_room_create_actions_pkey > @cursor::int
		WHEN 'prev' THEN
			t_chat_room_create_actions_pkey < @cursor::int
	END
ORDER BY
	CASE WHEN @cursor_direction::text = 'next' THEN t_chat_room_create_actions_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN t_chat_room_create_actions_pkey END DESC
LIMIT $2;

-- name: GetPluralChatRoomCreateActions :many
SELECT t_chat_room_create_actions.*,
m_members.name create_member_name, m_members.first_name create_member_first_name, m_members.last_name create_member_last_name, m_members.email create_member_email,
m_members.profile_image_id create_member_profile_image_id
FROM t_chat_room_create_actions
LEFT JOIN m_members ON t_chat_room_create_actions.created_by = m_members.member_id
WHERE chat_room_create_action_id = ANY(@chat_room_create_action_ids::uuid[])
ORDER BY
	t_chat_room_create_actions_pkey ASC;

-- name: GetPluralChatRoomCreateActionsUseNumberedPaginate :many
SELECT t_chat_room_create_actions.*,
m_members.name create_member_name, m_members.first_name create_member_first_name, m_members.last_name create_member_last_name, m_members.email create_member_email,
m_members.profile_image_id create_member_profile_image_id
FROM t_chat_room_create_actions
LEFT JOIN m_members ON t_chat_room_create_actions.created_by = m_members.member_id
WHERE chat_room_create_action_id = ANY(@chat_room_create_action_ids::uuid[])
ORDER BY
	t_chat_room_create_actions_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountChatRoomCreateActions :one
SELECT COUNT(*) FROM t_chat_room_create_actions;
