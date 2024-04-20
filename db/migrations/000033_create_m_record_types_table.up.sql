CREATE TABLE m_record_types (
	m_record_types_pkey BIGSERIAL,
    record_type_id UUID NOT NULL DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
	key VARCHAR(255) NOT NULL
);
ALTER TABLE m_record_types ADD CONSTRAINT m_record_types_pkey PRIMARY KEY (m_record_types_pkey);
CREATE UNIQUE INDEX idx_m_record_types_id ON m_record_types(record_type_id);
CREATE UNIQUE INDEX idx_m_record_types_key ON m_record_types(key);
