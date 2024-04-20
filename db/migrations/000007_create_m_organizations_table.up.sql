CREATE TABLE m_organizations (
	m_organizations_pkey BIGSERIAL,
    organization_id UUID NOT NULL DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
	description TEXT,
	is_personal BOOLEAN NOT NULL DEFAULT FALSE,
	is_whole BOOLEAN NOT NULL DEFAULT FALSE,
	created_at TIMESTAMPTZ NOT NULL,
	updated_at TIMESTAMPTZ NOT NULL
);
ALTER TABLE m_organizations ADD CONSTRAINT m_organizations_pkey PRIMARY KEY (m_organizations_pkey);
CREATE UNIQUE INDEX idx_m_organizations_id ON m_organizations(organization_id);
CREATE UNIQUE INDEX idx_m_organizations_name ON m_organizations(name);
