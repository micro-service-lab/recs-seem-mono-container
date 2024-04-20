CREATE TABLE m_policies (
	m_policies_pkey BIGSERIAL,
    policy_id UUID NOT NULL DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
	key VARCHAR(255) NOT NULL,
	policy_category_id UUID NOT NULL
);
ALTER TABLE m_policies ADD CONSTRAINT m_policies_pkey PRIMARY KEY (m_policies_pkey);
ALTER TABLE m_policies ADD CONSTRAINT fk_m_policies_policy_category_id FOREIGN KEY (policy_category_id) REFERENCES m_policy_categories(policy_category_id) ON DELETE RESTRICT ON UPDATE RESTRICT;
CREATE UNIQUE INDEX idx_m_policies_id ON m_policies(policy_id);
CREATE UNIQUE INDEX idx_m_policies_key ON m_policies(key);
