CREATE TABLE m_chat_room_action_types (
	m_chat_room_action_types_pkey BIGSERIAL,
    chat_room_action_type_id UUID NOT NULL DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
	key VARCHAR(255) NOT NULL
);
ALTER TABLE m_chat_room_action_types ADD CONSTRAINT m_chat_room_action_types_pkey PRIMARY KEY (m_chat_room_action_types_pkey);
CREATE UNIQUE INDEX idx_m_chat_room_action_types_id ON m_chat_room_action_types(chat_room_action_type_id);
CREATE UNIQUE INDEX idx_m_chat_room_action_types_key ON m_chat_room_action_types(key);
