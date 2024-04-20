CREATE TABLE m_event_types (
	m_event_types_pkey BIGSERIAL,
    event_type_id UUID NOT NULL DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
	key VARCHAR(255) NOT NULL,
	color VARCHAR(15) NOT NULL
);
ALTER TABLE m_event_types ADD CONSTRAINT m_event_types_pkey PRIMARY KEY (m_event_types_pkey);
CREATE UNIQUE INDEX idx_m_event_types_id ON m_event_types(event_type_id);
CREATE UNIQUE INDEX idx_m_event_types_key ON m_event_types(key);
