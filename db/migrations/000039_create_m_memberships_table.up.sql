CREATE TABLE m_memberships (
	m_memberships_pkey BIGSERIAL,
	member_id UUID NOT NULL,
	organization_id UUID NOT NULL,
	work_position_id UUID,
	added_at TIMESTAMPTZ NOT NULL
);
ALTER TABLE m_memberships ADD CONSTRAINT m_memberships_pkey PRIMARY KEY (m_memberships_pkey);
ALTER TABLE m_memberships ADD CONSTRAINT fk_m_memberships_member_id FOREIGN KEY (member_id) REFERENCES m_members(member_id) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE m_memberships ADD CONSTRAINT fk_m_memberships_organization_id FOREIGN KEY (organization_id) REFERENCES m_organizations(organization_id) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE m_memberships ADD CONSTRAINT fk_m_memberships_work_position_id FOREIGN KEY (work_position_id) REFERENCES m_work_positions(work_position_id) ON DELETE SET NULL ON UPDATE SET NULL;
CREATE UNIQUE INDEX idx_m_memberships_id ON m_memberships(member_id, organization_id);
