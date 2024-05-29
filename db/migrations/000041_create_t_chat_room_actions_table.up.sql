CREATE TABLE t_chat_room_actions (
	t_chat_room_actions_pkey BIGSERIAL,
    chat_room_action_id UUID NOT NULL DEFAULT uuid_generate_v4(),
	chat_room_id UUID NOT NULL,
	chat_room_action_type_id UUID NOT NULL,
	acted_at TIMESTAMPTZ NOT NULL
);
ALTER TABLE t_chat_room_actions ADD CONSTRAINT t_chat_room_actions_pkey PRIMARY KEY (t_chat_room_actions_pkey);
ALTER TABLE t_chat_room_actions ADD CONSTRAINT fk_t_chat_room_actions_chat_room_id FOREIGN KEY (chat_room_id) REFERENCES m_chat_rooms(chat_room_id) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE t_chat_room_actions ADD CONSTRAINT fk_t_chat_room_actions_chat_room_action_type_id FOREIGN KEY (chat_room_action_type_id) REFERENCES m_chat_room_action_types(chat_room_action_type_id) ON DELETE RESTRICT ON UPDATE RESTRICT;
CREATE UNIQUE INDEX idx_t_chat_room_actions_id ON t_chat_room_actions(chat_room_action_id);
