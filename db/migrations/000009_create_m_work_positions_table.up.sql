CREATE TABLE m_work_positions (
	m_work_positions_pkey BIGSERIAL,
    work_position_id UUID NOT NULL DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);
ALTER TABLE m_work_positions ADD CONSTRAINT m_work_positions_pkey PRIMARY KEY (m_work_positions_pkey);
CREATE UNIQUE INDEX idx_m_work_positions_id ON m_work_positions(work_position_id);
CREATE INDEX idx_m_work_positions_name ON m_work_positions(name);
