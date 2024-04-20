CREATE TABLE t_messages (
	t_messages_pkey BIGSERIAL,
    message_id UUID NOT NULL DEFAULT uuid_generate_v4(),
	chat_room_id UUID NOT NULL,
	sender_id UUID,
	body TEXT NOT NULL,
	posted_at TIMESTAMPTZ NOT NULL,
	last_edited_at TIMESTAMPTZ NOT NULL
);
ALTER TABLE t_messages ADD CONSTRAINT t_messages_pkey PRIMARY KEY (t_messages_pkey);
ALTER TABLE t_messages ADD CONSTRAINT fk_t_messages_chat_room_id FOREIGN KEY (chat_room_id) REFERENCES m_chat_rooms(chat_room_id) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE t_messages ADD CONSTRAINT fk_t_messages_sender_id FOREIGN KEY (sender_id) REFERENCES m_members(member_id) ON DELETE SET NULL ON UPDATE SET NULL;
CREATE UNIQUE INDEX idx_t_messages_id ON t_messages(message_id);
