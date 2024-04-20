CREATE TABLE t_lab_io_histories (
	t_lab_io_histories_pkey BIGSERIAL,
    lab_io_history_id UUID NOT NULL DEFAULT uuid_generate_v4(),
	member_id UUID NOT NULL,
	entered_at TIMESTAMPTZ NOT NULL,
	exited_at TIMESTAMPTZ
);
ALTER TABLE t_lab_io_histories ADD CONSTRAINT t_lab_io_histories_pkey PRIMARY KEY (t_lab_io_histories_pkey);
ALTER TABLE t_lab_io_histories ADD CONSTRAINT fk_t_lab_io_histories_member_id FOREIGN KEY (member_id) REFERENCES m_members(member_id) ON DELETE CASCADE ON UPDATE CASCADE;
CREATE UNIQUE INDEX idx_t_lab_io_histories_id ON t_lab_io_histories(lab_io_history_id);
