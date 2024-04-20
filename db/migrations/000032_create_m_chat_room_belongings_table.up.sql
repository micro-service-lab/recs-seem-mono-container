CREATE TABLE m_chat_room_belongings (
	m_chat_room_belongings_pkey BIGSERIAL,
	member_id UUID NOT NULL,
	chat_room_id UUID NOT NULL,
	added_at TIMESTAMPTZ NOT NULL
);
ALTER TABLE m_chat_room_belongings ADD CONSTRAINT m_chat_room_belongings_pkey PRIMARY KEY (m_chat_room_belongings_pkey);
ALTER TABLE m_chat_room_belongings ADD CONSTRAINT fk_m_chat_room_belongings_member_id FOREIGN KEY (member_id) REFERENCES m_members(member_id) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE m_chat_room_belongings ADD CONSTRAINT fk_m_chat_room_belongings_chat_room_id FOREIGN KEY (chat_room_id) REFERENCES m_chat_rooms(chat_room_id) ON DELETE CASCADE ON UPDATE CASCADE;
CREATE UNIQUE INDEX idx_m_chat_room_belongings_id ON m_chat_room_belongings(member_id, chat_room_id);
