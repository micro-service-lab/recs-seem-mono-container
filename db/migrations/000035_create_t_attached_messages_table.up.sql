CREATE TABLE t_attached_messages (
	t_attached_messages_pkey BIGSERIAL,
	attached_message_id UUID NOT NULL DEFAULT uuid_generate_v4(),
	message_id UUID,
	file_url TEXT NOT NULL
);
ALTER TABLE t_attached_messages ADD CONSTRAINT t_attached_messages_pkey PRIMARY KEY (t_attached_messages_pkey);
ALTER TABLE t_attached_messages ADD CONSTRAINT fk_t_attached_messages_message_id FOREIGN KEY (message_id) REFERENCES t_messages(message_id) ON DELETE SET NULL ON UPDATE SET NULL;
CREATE UNIQUE INDEX idx_t_attached_messages_id ON t_attached_messages(attached_message_id);
