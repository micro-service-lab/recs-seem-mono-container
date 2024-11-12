CREATE TABLE t_chat_room_withdraw_actions (
	t_chat_room_withdraw_actions_pkey BIGSERIAL,
    chat_room_withdraw_action_id UUID NOT NULL DEFAULT uuid_generate_v4(),
	chat_room_action_id UUID NOT NULL,
	member_id UUID
);
ALTER TABLE t_chat_room_withdraw_actions ADD CONSTRAINT t_chat_room_withdraw_actions_pkey PRIMARY KEY (t_chat_room_withdraw_actions_pkey);
ALTER TABLE t_chat_room_withdraw_actions ADD CONSTRAINT fk_t_chat_room_withdraw_actions_chat_room_action_id FOREIGN KEY (chat_room_action_id) REFERENCES t_chat_room_actions(chat_room_action_id) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE t_chat_room_withdraw_actions ADD CONSTRAINT fk_t_chat_room_withdraw_actions_member_id FOREIGN KEY (member_id) REFERENCES m_members(member_id) ON DELETE SET NULL ON UPDATE SET NULL;
CREATE UNIQUE INDEX idx_t_chat_room_withdraw_actions_id ON t_chat_room_withdraw_actions(chat_room_withdraw_action_id);
