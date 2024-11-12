ALTER TABLE t_messages ADD COLUMN chat_room_action_id UUID NOT NULL;

ALTER TABLE t_messages ADD CONSTRAINT fk_t_messages_chat_room_action_id FOREIGN KEY (chat_room_action_id) REFERENCES t_chat_room_actions(chat_room_action_id) ON DELETE CASCADE ON UPDATE CASCADE;
