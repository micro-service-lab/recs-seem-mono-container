-- name: CreateChatRoomAddedMembers :copyfrom
INSERT INTO t_chat_room_added_members (member_id, chat_room_add_member_action_id) VALUES ($1, $2);

-- name: CreateChatRoomAddedMember :one
INSERT INTO t_chat_room_added_members (member_id, chat_room_add_member_action_id) VALUES ($1, $2) RETURNING *;

-- name: DeleteChatRoomAddedMember :execrows
DELETE FROM t_chat_room_added_members WHERE member_id = $1 AND chat_room_add_member_action_id = $2;

-- name: DeleteChatRoomAddedMembersOnMember :execrows
DELETE FROM t_chat_room_added_members WHERE member_id = $1;

-- name: DeleteChatRoomAddedMembersOnMembers :execrows
DELETE FROM t_chat_room_added_members WHERE member_id = ANY(@member_ids::uuid[]);

-- name: DeleteChatRoomAddedMembersOnChatRoomAddMemberAction :execrows
DELETE FROM t_chat_room_added_members WHERE chat_room_add_member_action_id = $1;

-- name: DeleteChatRoomAddedMembersOnChatRoomAddMemberActions :execrows
DELETE FROM t_chat_room_added_members WHERE chat_room_add_member_action_id = ANY(@chat_room_add_member_action_ids::uuid[]);

-- name: GetMembersOnChatRoomAddMemberAction :many
SELECT t_chat_room_added_members.*, m_members.name member_name, m_members.first_name member_first_name, m_members.last_name member_last_name, m_members.email member_email,
m_members.profile_image_id member_profile_image_id
FROM t_chat_room_added_members
LEFT JOIN m_members ON t_chat_room_added_members.member_id = m_members.member_id
WHERE chat_room_add_member_action_id = $1
ORDER BY
	t_chat_room_added_members_pkey ASC;

-- name: GetMembersOnChatRoomAddMemberActionUseNumberedPaginate :many
SELECT t_chat_room_added_members.*, m_members.name member_name, m_members.first_name member_first_name, m_members.last_name member_last_name, m_members.email member_email,
m_members.profile_image_id member_profile_image_id
FROM t_chat_room_added_members
LEFT JOIN m_members ON t_chat_room_added_members.member_id = m_members.member_id
WHERE chat_room_add_member_action_id = $1
ORDER BY
	t_chat_room_added_members_pkey ASC
LIMIT $2 OFFSET $3;

-- name: GetMembersOnChatRoomAddMemberActionUseKeysetPaginate :many
SELECT t_chat_room_added_members.*, m_members.name member_name, m_members.first_name member_first_name, m_members.last_name member_last_name, m_members.email member_email,
m_members.profile_image_id member_profile_image_id
FROM t_chat_room_added_members
LEFT JOIN m_members ON t_chat_room_added_members.member_id = m_members.member_id
WHERE chat_room_add_member_action_id = $1
AND CASE @cursor_direction::text
	WHEN 'next' THEN
			t_chat_room_added_members_pkey > @cursor::int
	WHEN 'prev' THEN
			t_chat_room_added_members_pkey < @cursor::int
END
ORDER BY
	CASE WHEN @cursor_direction::text = 'next' THEN t_chat_room_added_members_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN t_chat_room_added_members_pkey END DESC
LIMIT $2;

-- name: GetPluralMembersOnChatRoomAddMemberAction :many
SELECT t_chat_room_added_members.*, m_members.name member_name, m_members.first_name member_first_name, m_members.last_name member_last_name, m_members.email member_email,
m_members.profile_image_id member_profile_image_id
FROM t_chat_room_added_members
LEFT JOIN m_members ON t_chat_room_added_members.member_id = m_members.member_id
WHERE chat_room_add_member_action_id = ANY(@chat_room_add_member_action_ids::uuid[])
ORDER BY
	t_chat_room_added_members_pkey ASC;

-- name: GetPluralMembersOnChatRoomAddMemberActionUseNumberedPaginate :many
SELECT t_chat_room_added_members.*, m_members.name member_name, m_members.first_name member_first_name, m_members.last_name member_last_name, m_members.email member_email,
m_members.profile_image_id member_profile_image_id
FROM t_chat_room_added_members
LEFT JOIN m_members ON t_chat_room_added_members.member_id = m_members.member_id
WHERE chat_room_add_member_action_id = ANY(@chat_room_add_member_action_ids::uuid[])
ORDER BY
	t_chat_room_added_members_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountMembersOnChatRoomAddMemberAction :one
SELECT COUNT(*) FROM t_chat_room_added_members WHERE chat_room_add_member_action_id = $1;
