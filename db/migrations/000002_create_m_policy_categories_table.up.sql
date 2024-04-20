CREATE TABLE m_policy_categories (
	m_policy_categories_pkey BIGSERIAL,
    policy_category_id UUID NOT NULL DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
	key VARCHAR(255) NOT NULL
);
ALTER TABLE m_policy_categories ADD CONSTRAINT m_policy_categories_pkey PRIMARY KEY (m_policy_categories_pkey);
CREATE UNIQUE INDEX idx_m_policy_categories_id ON m_policy_categories(policy_category_id);
CREATE UNIQUE INDEX idx_m_policy_categories_key ON m_policy_categories(key);
