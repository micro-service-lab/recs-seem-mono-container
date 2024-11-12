CREATE TABLE t_read_receipts (
	t_read_receipts_pkey BIGSERIAL,
	member_id UUID NOT NULL,
	message_id UUID NOT NULL,
	read_at TIMESTAMPTZ
);
ALTER TABLE t_read_receipts ADD CONSTRAINT t_read_receipts_pkey PRIMARY KEY (t_read_receipts_pkey);
ALTER TABLE t_read_receipts ADD CONSTRAINT fk_t_read_receipts_member_id FOREIGN KEY (member_id) REFERENCES m_members(member_id) ON DELETE SET NULL ON UPDATE SET NULL;
ALTER TABLE t_read_receipts ADD CONSTRAINT fk_t_read_receipts_message_id FOREIGN KEY (message_id) REFERENCES t_messages(message_id) ON DELETE CASCADE ON UPDATE CASCADE;
CREATE UNIQUE INDEX idx_t_read_receipts_id ON t_read_receipts(member_id, message_id);
