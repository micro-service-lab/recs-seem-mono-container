CREATE TABLE m_permission_categories (
	m_permission_categories_pkey BIGSERIAL,
    permission_category_id UUID NOT NULL DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
	key VARCHAR(255) NOT NULL
);
ALTER TABLE m_permission_categories ADD CONSTRAINT m_permission_categories_pkey PRIMARY KEY (m_permission_categories_pkey);
CREATE UNIQUE INDEX idx_m_permission_categories_id ON m_permission_categories(permission_category_id);
CREATE UNIQUE INDEX idx_m_permission_categories_key ON m_permission_categories(key);
