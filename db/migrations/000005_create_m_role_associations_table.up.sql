CREATE TABLE m_role_associations (
	m_role_associations_pkey BIGSERIAL,
	role_id UUID NOT NULL,
	policy_id UUID NOT NULL
);
ALTER TABLE m_role_associations ADD CONSTRAINT m_role_associations_pkey PRIMARY KEY (m_role_associations_pkey);
ALTER TABLE m_role_associations ADD CONSTRAINT fk_m_role_associations_role_id FOREIGN KEY (role_id) REFERENCES m_roles(role_id) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE m_role_associations ADD CONSTRAINT fk_m_role_associations_policy_id FOREIGN KEY (policy_id) REFERENCES m_policies(policy_id) ON DELETE RESTRICT ON UPDATE RESTRICT;
CREATE UNIQUE INDEX idx_m_role_associations_id ON m_role_associations(role_id, policy_id);
