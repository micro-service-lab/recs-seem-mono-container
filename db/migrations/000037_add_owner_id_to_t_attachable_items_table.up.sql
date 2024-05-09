ALTER TABLE t_attachable_items ADD COLUMN owner_id UUID;

ALTER TABLE t_attachable_items ADD CONSTRAINT fk_t_attachable_items_owner_id FOREIGN KEY (owner_id) REFERENCES m_members(member_id) ON DELETE SET NULL ON UPDATE SET NULL;
