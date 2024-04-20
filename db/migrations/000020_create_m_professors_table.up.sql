CREATE TABLE m_professors (
	m_professors_pkey BIGSERIAL,
    professor_id UUID NOT NULL DEFAULT uuid_generate_v4(),
	member_id UUID NOT NULL
);
ALTER TABLE m_professors ADD CONSTRAINT m_professors_pkey PRIMARY KEY (m_professors_pkey);
ALTER TABLE m_professors ADD CONSTRAINT fk_m_professors_member_id FOREIGN KEY (member_id) REFERENCES m_members(member_id) ON DELETE RESTRICT ON UPDATE RESTRICT;
CREATE UNIQUE INDEX idx_m_professors_id ON m_professors(professor_id);
