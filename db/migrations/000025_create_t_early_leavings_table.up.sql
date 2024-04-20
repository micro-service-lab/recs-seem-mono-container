CREATE TABLE t_early_leavings (
	t_early_leavings_pkey BIGSERIAL,
    early_leaving_id UUID NOT NULL DEFAULT uuid_generate_v4(),
	attendance_id UUID NOT NULL,
	leave_time TIMESTAMPTZ NOT NULL
);
ALTER TABLE t_early_leavings ADD CONSTRAINT t_early_leavings_pkey PRIMARY KEY (t_early_leavings_pkey);
ALTER TABLE t_early_leavings ADD CONSTRAINT fk_t_early_leavings_attendance_id FOREIGN KEY (attendance_id) REFERENCES t_attendances(attendance_id) ON DELETE RESTRICT ON UPDATE RESTRICT;
CREATE UNIQUE INDEX idx_t_early_leavings_id ON t_early_leavings(early_leaving_id);
