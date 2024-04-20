CREATE TABLE t_attendances (
	t_attendances_pkey BIGSERIAL,
    attendance_id UUID NOT NULL DEFAULT uuid_generate_v4(),
	attendance_type_id UUID NOT NULL,
	member_id UUID NOT NULL,
	description TEXT NOT NULL,
	date DATE NOT NULL,
	mail_send_flag BOOLEAN NOT NULL,
	send_organization_id UUID,
	posted_at TIMESTAMPTZ NOT NULL,
	last_edited_at TIMESTAMPTZ NOT NULL
);
ALTER TABLE t_attendances ADD CONSTRAINT t_attendances_pkey PRIMARY KEY (t_attendances_pkey);
ALTER TABLE t_attendances ADD CONSTRAINT fk_t_attendances_attendance_type_id FOREIGN KEY (attendance_type_id) REFERENCES m_attendance_types(attendance_type_id) ON DELETE RESTRICT ON UPDATE RESTRICT;
ALTER TABLE t_attendances ADD CONSTRAINT fk_t_attendances_member_id FOREIGN KEY (member_id) REFERENCES m_members(member_id) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE t_attendances ADD CONSTRAINT fk_t_attendances_send_organization_id FOREIGN KEY (send_organization_id) REFERENCES m_organizations(organization_id) ON DELETE SET NULL ON UPDATE SET NULL;
CREATE UNIQUE INDEX idx_t_attendances_id ON t_attendances(attendance_id);
