CREATE TABLE t_events (
	t_events_pkey BIGSERIAL,
    event_id UUID NOT NULL DEFAULT uuid_generate_v4(),
	event_type_id UUID NOT NULL,
	title VARCHAR(255) NOT NULL,
	description TEXT,
	organization_id UUID,
	start_time TIMESTAMPTZ NOT NULL,
	end_time TIMESTAMPTZ NOT NULL,
	mail_send_flag BOOLEAN NOT NULL,
	send_organization_id UUID,
	posted_by UUID,
	last_edited_by UUID,
	posted_at TIMESTAMPTZ NOT NULL,
	last_edited_at TIMESTAMPTZ NOT NULL
);
ALTER TABLE t_events ADD CONSTRAINT t_events_pkey PRIMARY KEY (t_events_pkey);
ALTER TABLE t_events ADD CONSTRAINT fk_t_events_event_type_id FOREIGN KEY (event_type_id) REFERENCES m_event_types(event_type_id) ON DELETE RESTRICT ON UPDATE RESTRICT;
ALTER TABLE t_events ADD CONSTRAINT fk_t_events_organization_id FOREIGN KEY (organization_id) REFERENCES m_organizations(organization_id) ON DELETE SET NULL ON UPDATE SET NULL;
ALTER TABLE t_events ADD CONSTRAINT fk_t_events_send_organization_id FOREIGN KEY (send_organization_id) REFERENCES m_organizations(organization_id) ON DELETE SET NULL ON UPDATE SET NULL;
ALTER TABLE t_events ADD CONSTRAINT fk_t_events_posted_by FOREIGN KEY (posted_by) REFERENCES m_members(member_id) ON DELETE SET NULL ON UPDATE SET NULL;
ALTER TABLE t_events ADD CONSTRAINT fk_t_events_last_edited_by FOREIGN KEY (last_edited_by) REFERENCES m_members(member_id) ON DELETE SET NULL ON UPDATE SET NULL;
CREATE UNIQUE INDEX idx_t_events_id ON t_events(event_id);
