CREATE TABLE t_chat_room_added_members (
	t_chat_room_added_members_pkey BIGSERIAL,
	chat_room_add_member_action_id UUID NOT NULL,
	member_id UUID
);
ALTER TABLE t_chat_room_added_members ADD CONSTRAINT t_chat_room_added_members_pkey PRIMARY KEY (t_chat_room_added_members_pkey);
ALTER TABLE t_chat_room_added_members ADD CONSTRAINT fk_t_chat_room_added_members_chat_room_add_member_action_id FOREIGN KEY (chat_room_add_member_action_id) REFERENCES t_chat_room_add_member_actions(chat_room_add_member_action_id) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE t_chat_room_added_members ADD CONSTRAINT fk_t_chat_room_added_members_member_id FOREIGN KEY (member_id) REFERENCES m_members(member_id) ON DELETE SET NULL ON UPDATE SET NULL;
CREATE UNIQUE INDEX idx_t_chat_room_added_members_id ON t_chat_room_added_members(chat_room_add_member_action_id, member_id);
