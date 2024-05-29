-- name: CreateChatRoomRemovedMembers :copyfrom
INSERT INTO t_chat_room_removed_members (member_id, chat_room_remove_member_action_id) VALUES ($1, $2);

-- name: CreateChatRoomRemovedMember :one
INSERT INTO t_chat_room_removed_members (member_id, chat_room_remove_member_action_id) VALUES ($1, $2) RETURNING *;

-- name: DeleteChatRoomRemovedMember :execrows
DELETE FROM t_chat_room_removed_members WHERE member_id = $1 AND chat_room_remove_member_action_id = $2;

-- name: DeleteChatRoomRemovedMembersOnMember :execrows
DELETE FROM t_chat_room_removed_members WHERE member_id = $1;

-- name: DeleteChatRoomRemovedMembersOnMembers :execrows
DELETE FROM t_chat_room_removed_members WHERE member_id = ANY(@member_ids::uuid[]);

-- name: DeleteChatRoomRemovedMembersOnChatRoomRemoveMemberAction :execrows
DELETE FROM t_chat_room_removed_members WHERE chat_room_remove_member_action_id = $1;

-- name: DeleteChatRoomRemovedMembersOnChatRoomRemoveMemberActions :execrows
DELETE FROM t_chat_room_removed_members WHERE chat_room_remove_member_action_id = ANY(@chat_room_remove_member_action_ids::uuid[]);

-- name: GetMembersOnChatRoomRemoveMemberAction :many
SELECT t_chat_room_removed_members.*, m_members.name member_name, m_members.first_name member_first_name, m_members.last_name member_last_name, m_members.email member_email,
m_members.profile_image_id member_profile_image_id
FROM t_chat_room_removed_members
LEFT JOIN m_members ON t_chat_room_removed_members.member_id = m_members.member_id
WHERE chat_room_remove_member_action_id = $1
ORDER BY
	t_chat_room_removed_members_pkey ASC;

-- name: GetMembersOnChatRoomRemoveMemberActionUseNumberedPaginate :many
SELECT t_chat_room_removed_members.*, m_members.name member_name, m_members.first_name member_first_name, m_members.last_name member_last_name, m_members.email member_email,
m_members.profile_image_id member_profile_image_id
FROM t_chat_room_removed_members
LEFT JOIN m_members ON t_chat_room_removed_members.member_id = m_members.member_id
WHERE chat_room_remove_member_action_id = $1
ORDER BY
	t_chat_room_removed_members_pkey ASC
LIMIT $2 OFFSET $3;

-- name: GetMembersOnChatRoomRemoveMemberActionUseKeysetPaginate :many
SELECT t_chat_room_removed_members.*, m_members.name member_name, m_members.first_name member_first_name, m_members.last_name member_last_name, m_members.email member_email,
m_members.profile_image_id member_profile_image_id
FROM t_chat_room_removed_members
LEFT JOIN m_members ON t_chat_room_removed_members.member_id = m_members.member_id
WHERE chat_room_remove_member_action_id = $1
AND CASE @cursor_direction::text
	WHEN 'next' THEN
			t_chat_room_removed_members_pkey > @cursor::int
	WHEN 'prev' THEN
			t_chat_room_removed_members_pkey < @cursor::int
END
ORDER BY
	CASE WHEN @cursor_direction::text = 'next' THEN t_chat_room_removed_members_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN t_chat_room_removed_members_pkey END DESC
LIMIT $2;

-- name: GetPluralMembersOnChatRoomRemoveMemberAction :many
SELECT t_chat_room_removed_members.*, m_members.name member_name, m_members.first_name member_first_name, m_members.last_name member_last_name, m_members.email member_email,
m_members.profile_image_id member_profile_image_id
FROM t_chat_room_removed_members
LEFT JOIN m_members ON t_chat_room_removed_members.member_id = m_members.member_id
WHERE chat_room_remove_member_action_id = ANY(@chat_room_remove_member_action_ids::uuid[])
ORDER BY
	t_chat_room_removed_members_pkey ASC;

-- name: GetPluralMembersOnChatRoomRemoveMemberActionUseNumberedPaginate :many
SELECT t_chat_room_removed_members.*, m_members.name member_name, m_members.first_name member_first_name, m_members.last_name member_last_name, m_members.email member_email,
m_members.profile_image_id member_profile_image_id
FROM t_chat_room_removed_members
LEFT JOIN m_members ON t_chat_room_removed_members.member_id = m_members.member_id
WHERE chat_room_remove_member_action_id = ANY(@chat_room_remove_member_action_ids::uuid[])
ORDER BY
	t_chat_room_removed_members_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountMembersOnChatRoomRemoveMemberAction :one
SELECT COUNT(*) FROM t_chat_room_removed_members WHERE chat_room_remove_member_action_id = $1;
