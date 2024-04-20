CREATE TABLE m_chat_rooms (
	m_chat_rooms_pkey BIGSERIAL,
    chat_room_id UUID NOT NULL DEFAULT uuid_generate_v4(),
	name VARCHAR(255),
	is_private BOOLEAN NOT NULL,
	cover_image_id UUID,
	owner_id UUID,
	created_at TIMESTAMPTZ NOT NULL,
	updated_at TIMESTAMPTZ NOT NULL
);
ALTER TABLE m_chat_rooms ADD CONSTRAINT m_chat_rooms_pkey PRIMARY KEY (m_chat_rooms_pkey);
ALTER TABLE m_chat_rooms ADD CONSTRAINT fk_m_chat_rooms_cover_image_id FOREIGN KEY (cover_image_id) REFERENCES t_images(image_id) ON DELETE SET NULL ON UPDATE SET NULL;
ALTER TABLE m_chat_rooms ADD CONSTRAINT fk_m_chat_rooms_owner_id FOREIGN KEY (owner_id) REFERENCES m_members(member_id) ON DELETE SET NULL ON UPDATE SET NULL;
CREATE UNIQUE INDEX idx_m_chat_rooms_id ON m_chat_rooms(chat_room_id);
