CREATE TABLE m_mime_types (
	m_mime_types_pkey BIGSERIAL,
    mime_type_id UUID NOT NULL DEFAULT uuid_generate_v4(),
	name VARCHAR(255) NOT NULL,
	key VARCHAR(255) NOT NULL
);
ALTER TABLE m_mime_types ADD CONSTRAINT m_mime_types_pkey PRIMARY KEY (m_mime_types_pkey);
CREATE UNIQUE INDEX idx_m_mime_types_id ON m_mime_types(mime_type_id);
CREATE UNIQUE INDEX idx_m_mime_types_name ON m_mime_types(name);
