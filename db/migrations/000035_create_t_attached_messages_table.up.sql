CREATE TABLE t_attached_messages (
	t_attached_messages_pkey BIGSERIAL,
	message_id UUID,
	attachable_item_id UUID NOT NULL
);
ALTER TABLE t_attached_messages ADD CONSTRAINT t_attached_messages_pkey PRIMARY KEY (t_attached_messages_pkey);
ALTER TABLE t_attached_messages ADD CONSTRAINT fk_t_attached_messages_message_id FOREIGN KEY (message_id) REFERENCES t_messages(message_id) ON DELETE SET NULL ON UPDATE SET NULL;
ALTER TABLE t_attached_messages ADD CONSTRAINT fk_t_attached_messages_attachable_item_id FOREIGN KEY (attachable_item_id) REFERENCES t_attachable_items(attachable_item_id) ON DELETE RESTRICT ON UPDATE RESTRICT;
CREATE UNIQUE INDEX idx_t_attached_messages_id ON t_attached_messages(message_id, attachable_item_id);
