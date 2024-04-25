CREATE TABLE t_position_histories (
	t_position_histories_pkey BIGSERIAL,
    position_history_id UUID NOT NULL DEFAULT uuid_generate_v4(),
	member_id UUID NOT NULL,
	x_pos DOUBLE PRECISION NOT NULL,
	y_pos DOUBLE PRECISION NOT NULL,
	sent_at TIMESTAMPTZ NOT NULL
);
ALTER TABLE t_position_histories ADD CONSTRAINT t_position_histories_pkey PRIMARY KEY (t_position_histories_pkey);
ALTER TABLE t_position_histories ADD CONSTRAINT fk_t_position_histories_member_id FOREIGN KEY (member_id) REFERENCES m_members(member_id) ON DELETE CASCADE ON UPDATE CASCADE;
CREATE UNIQUE INDEX idx_t_position_histories_id ON t_position_histories(position_history_id);
