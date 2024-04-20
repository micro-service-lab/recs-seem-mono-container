CREATE TABLE m_attend_statuses (
	m_attend_statuses_pkey BIGSERIAL,
    attend_status_id UUID NOT NULL DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
	key VARCHAR(255) NOT NULL
);
ALTER TABLE m_attend_statuses ADD CONSTRAINT m_attend_statuses_pkey PRIMARY KEY (m_attend_statuses_pkey);
CREATE UNIQUE INDEX idx_m_attend_statuses_id ON m_attend_statuses(attend_status_id);
CREATE UNIQUE INDEX idx_m_attend_statuses_key ON m_attend_statuses(key);
