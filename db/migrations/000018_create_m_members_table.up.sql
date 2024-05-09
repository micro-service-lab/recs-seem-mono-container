CREATE TABLE m_members (
	m_members_pkey BIGSERIAL,
    member_id UUID NOT NULL DEFAULT uuid_generate_v4(),
	login_id VARCHAR(255) NOT NULL,
	password VARCHAR(255) NOT NULL,
	email VARCHAR(255) NOT NULL,
	name VARCHAR(255) NOT NULL,
	attend_status_id UUID NOT NULL,
	profile_image_url TEXT,
	grade_id UUID NOT NULL,
	group_id UUID NOT NULL,
	personal_organization_id UUID NOT NULL,
	role_id UUID,
	created_at TIMESTAMPTZ NOT NULL,
	updated_at TIMESTAMPTZ NOT NULL
);
ALTER TABLE m_members ADD CONSTRAINT m_members_pkey PRIMARY KEY (m_members_pkey);
ALTER TABLE m_members ADD CONSTRAINT fk_m_members_attend_status_id FOREIGN KEY (attend_status_id) REFERENCES m_attend_statuses(attend_status_id) ON DELETE RESTRICT ON UPDATE RESTRICT;
ALTER TABLE m_members ADD CONSTRAINT fk_m_members_grade_id FOREIGN KEY (grade_id) REFERENCES m_grades(grade_id) ON DELETE RESTRICT ON UPDATE RESTRICT;
ALTER TABLE m_members ADD CONSTRAINT fk_m_members_group_id FOREIGN KEY (group_id) REFERENCES m_groups(group_id) ON DELETE RESTRICT ON UPDATE RESTRICT;
ALTER TABLE m_members ADD CONSTRAINT fk_m_members_personal_organization_id FOREIGN KEY (personal_organization_id) REFERENCES m_organizations(organization_id) ON DELETE RESTRICT ON UPDATE RESTRICT;
ALTER TABLE m_members ADD CONSTRAINT fk_m_members_role_id FOREIGN KEY (role_id) REFERENCES m_roles(role_id) ON DELETE SET NULL ON UPDATE SET NULL;
CREATE UNIQUE INDEX idx_m_members_id ON m_members(member_id);
CREATE UNIQUE INDEX idx_m_members_login_id ON m_members(login_id);
