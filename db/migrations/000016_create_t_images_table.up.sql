CREATE TABLE t_images (
	t_images_pkey BIGSERIAL,
    image_id UUID NOT NULL DEFAULT uuid_generate_v4(),
	height DOUBLE PRECISION,
	width DOUBLE PRECISION,
	attachable_item_id UUID NOT NULL
);
ALTER TABLE t_images ADD CONSTRAINT t_images_pkey PRIMARY KEY (t_images_pkey);
ALTER TABLE t_images ADD CONSTRAINT fk_t_images_attachable_item_id FOREIGN KEY (attachable_item_id) REFERENCES t_attachable_items(attachable_item_id) ON DELETE RESTRICT ON UPDATE RESTRICT;
CREATE UNIQUE INDEX idx_t_images_id ON t_images(image_id);
