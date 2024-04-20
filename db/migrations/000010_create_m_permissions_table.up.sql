CREATE TABLE m_permissions (
	m_permissions_pkey BIGSERIAL,
    permission_id UUID NOT NULL DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
	key VARCHAR(255) NOT NULL,
	permission_category_id UUID NOT NULL
);
ALTER TABLE m_permissions ADD CONSTRAINT m_permissions_pkey PRIMARY KEY (m_permissions_pkey);
ALTER TABLE m_permissions ADD CONSTRAINT fk_m_permissions_permission_category_id FOREIGN KEY (permission_category_id) REFERENCES m_permission_categories(permission_category_id) ON DELETE RESTRICT ON UPDATE RESTRICT;
CREATE UNIQUE INDEX idx_m_permissions_id ON m_permissions(permission_id);
CREATE UNIQUE INDEX idx_m_permissions_key ON m_permissions(key);
