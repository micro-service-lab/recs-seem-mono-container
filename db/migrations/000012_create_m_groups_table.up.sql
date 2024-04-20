CREATE TABLE m_groups (
	m_groups_pkey BIGSERIAL,
    group_id UUID NOT NULL DEFAULT uuid_generate_v4(),
	key VARCHAR(255) NOT NULL,
	organization_id UUID NOT NULL
);
ALTER TABLE m_groups ADD CONSTRAINT m_groups_pkey PRIMARY KEY (m_groups_pkey);
ALTER TABLE m_groups ADD CONSTRAINT fk_m_groups_organization_id FOREIGN KEY (organization_id) REFERENCES m_organizations(organization_id) ON DELETE RESTRICT ON UPDATE RESTRICT;
CREATE UNIQUE INDEX idx_m_groups_id ON m_groups(group_id);
CREATE UNIQUE INDEX idx_m_groups_key ON m_groups(key);
