CREATE TABLE t_records (
	t_records_pkey BIGSERIAL,
    record_id UUID NOT NULL DEFAULT uuid_generate_v4(),
	record_type_id UUID NOT NULL,
	title VARCHAR(255) NOT NULL,
	body TEXT,
	organization_id UUID,
	posted_by UUID,
	last_edited_by UUID,
	posted_at TIMESTAMPTZ NOT NULL,
	last_edited_at TIMESTAMPTZ NOT NULL
);
ALTER TABLE t_records ADD CONSTRAINT t_records_pkey PRIMARY KEY (t_records_pkey);
ALTER TABLE t_records ADD CONSTRAINT fk_t_records_record_type_id FOREIGN KEY (record_type_id) REFERENCES m_record_types(record_type_id) ON DELETE RESTRICT ON UPDATE RESTRICT;
ALTER TABLE t_records ADD CONSTRAINT fk_t_records_organization_id FOREIGN KEY (organization_id) REFERENCES m_organizations(organization_id) ON DELETE SET NULL ON UPDATE SET NULL;
ALTER TABLE t_records ADD CONSTRAINT fk_t_records_posted_by FOREIGN KEY (posted_by) REFERENCES m_members(member_id) ON DELETE SET NULL ON UPDATE SET NULL;
ALTER TABLE t_records ADD CONSTRAINT fk_t_records_last_edited_by FOREIGN KEY (last_edited_by) REFERENCES m_members(member_id) ON DELETE SET NULL ON UPDATE SET NULL;
CREATE UNIQUE INDEX idx_t_records_id ON t_records(record_id);
