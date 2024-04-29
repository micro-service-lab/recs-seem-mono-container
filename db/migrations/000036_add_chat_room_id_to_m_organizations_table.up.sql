ALTER TABLE m_organizations ADD COLUMN chat_room_id UUID;

ALTER TABLE m_organizations ADD CONSTRAINT fk_m_organizations_chat_room_id FOREIGN KEY (chat_room_id) REFERENCES m_chat_rooms(chat_room_id) ON DELETE RESTRICT ON UPDATE RESTRICT;
