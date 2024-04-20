CREATE TABLE m_attendance_types (
	m_attendance_types_pkey BIGSERIAL,
    attendance_type_id UUID NOT NULL DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
	key VARCHAR(255) NOT NULL,
	color VARCHAR(15) NOT NULL
);
ALTER TABLE m_attendance_types ADD CONSTRAINT m_attendance_types_pkey PRIMARY KEY (m_attendance_types_pkey);
CREATE UNIQUE INDEX idx_m_attendance_types_id ON m_attendance_types(attendance_type_id);
CREATE UNIQUE INDEX idx_m_attendance_types_key ON m_attendance_types(key);
