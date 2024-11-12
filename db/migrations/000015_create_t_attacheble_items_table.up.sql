CREATE TABLE t_attachable_items (
	t_attachable_items_pkey BIGSERIAL,
    attachable_item_id UUID NOT NULL DEFAULT uuid_generate_v4(),
	url TEXT NOT NULL,
	size DOUBLE PRECISION,
	alias VARCHAR(255) NOT NULL,
	mime_type_id UUID NOT NULL
);
ALTER TABLE t_attachable_items ADD CONSTRAINT t_attachable_items_pkey PRIMARY KEY (t_attachable_items_pkey);
ALTER TABLE t_attachable_items ADD CONSTRAINT fk_t_attachable_items_mime_type_id FOREIGN KEY (mime_type_id) REFERENCES m_mime_types(mime_type_id) ON DELETE RESTRICT ON UPDATE RESTRICT;
CREATE UNIQUE INDEX idx_t_attachable_items_id ON t_attachable_items(attachable_item_id);
CREATE INDEX idx_t_attachable_items_url ON t_attachable_items(url);
