CREATE TABLE t_chat_room_add_member_actions (
	t_chat_room_add_member_actions_pkey BIGSERIAL,
    chat_room_add_member_action_id UUID NOT NULL DEFAULT uuid_generate_v4(),
	chat_room_action_id UUID NOT NULL,
	added_by UUID
);
ALTER TABLE t_chat_room_add_member_actions ADD CONSTRAINT t_chat_room_add_member_actions_pkey PRIMARY KEY (t_chat_room_add_member_actions_pkey);
ALTER TABLE t_chat_room_add_member_actions ADD CONSTRAINT fk_t_chat_room_add_member_actions_chat_room_action_id FOREIGN KEY (chat_room_action_id) REFERENCES t_chat_room_actions(chat_room_action_id) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE t_chat_room_add_member_actions ADD CONSTRAINT fk_t_chat_room_add_member_actions_added_by FOREIGN KEY (added_by) REFERENCES m_members(member_id) ON DELETE SET NULL ON UPDATE SET NULL;
CREATE UNIQUE INDEX idx_t_chat_room_add_member_actions_id ON t_chat_room_add_member_actions(chat_room_add_member_action_id);
