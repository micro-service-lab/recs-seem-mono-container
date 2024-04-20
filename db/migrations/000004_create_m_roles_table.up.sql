CREATE TABLE m_roles (
	m_roles_pkey BIGSERIAL,
    role_id UUID NOT NULL DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);
ALTER TABLE m_roles ADD CONSTRAINT m_roles_pkey PRIMARY KEY (m_roles_pkey);
CREATE UNIQUE INDEX idx_m_roles_id ON m_roles(role_id);
CREATE INDEX idx_m_roles_name ON m_roles(name);
