CREATE TABLE m_grades (
	m_grades_pkey BIGSERIAL,
    grade_id UUID NOT NULL DEFAULT uuid_generate_v4(),
	key VARCHAR(255) NOT NULL,
	organization_id UUID NOT NULL
);
ALTER TABLE m_grades ADD CONSTRAINT m_grades_pkey PRIMARY KEY (m_grades_pkey);
ALTER TABLE m_grades ADD CONSTRAINT fk_m_grades_organization_id FOREIGN KEY (organization_id) REFERENCES m_organizations(organization_id) ON DELETE RESTRICT ON UPDATE RESTRICT;
CREATE UNIQUE INDEX idx_m_grades_id ON m_grades(grade_id);
CREATE UNIQUE INDEX idx_m_grades_key ON m_grades(key);
