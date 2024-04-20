CREATE TABLE m_permission_associations (
	m_permission_associations_pkey BIGSERIAL,
	permission_id UUID NOT NULL,
	work_position_id UUID NOT NULL
);
ALTER TABLE m_permission_associations ADD CONSTRAINT m_permission_associations_pkey PRIMARY KEY (m_permission_associations_pkey);
ALTER TABLE m_permission_associations ADD CONSTRAINT fk_m_permission_associations_permission_id FOREIGN KEY (permission_id) REFERENCES m_permissions(permission_id) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE m_permission_associations ADD CONSTRAINT fk_m_permission_associations_work_position_id FOREIGN KEY (work_position_id) REFERENCES m_work_positions(work_position_id) ON DELETE RESTRICT ON UPDATE RESTRICT;
CREATE UNIQUE INDEX idx_m_permission_associations_id ON m_permission_associations(permission_id, work_position_id);
