CREATE TABLE t_files (
	t_files_pkey BIGSERIAL,
    file_id UUID NOT NULL DEFAULT uuid_generate_v4(),
	attachable_item_id UUID NOT NULL
);
ALTER TABLE t_files ADD CONSTRAINT t_files_pkey PRIMARY KEY (t_files_pkey);
ALTER TABLE t_files ADD CONSTRAINT fk_t_files_attachable_item_id FOREIGN KEY (attachable_item_id) REFERENCES t_attachable_items(attachable_item_id) ON DELETE RESTRICT ON UPDATE RESTRICT;
CREATE UNIQUE INDEX idx_t_files_id ON t_files(file_id);
