CREATE TABLE t_late_arrivals (
	t_late_arrivals_pkey BIGSERIAL,
    late_arrival_id UUID NOT NULL DEFAULT uuid_generate_v4(),
	attendance_id UUID NOT NULL,
	arrive_time TIMESTAMPTZ NOT NULL
);
ALTER TABLE t_late_arrivals ADD CONSTRAINT t_late_arrivals_pkey PRIMARY KEY (t_late_arrivals_pkey);
ALTER TABLE t_late_arrivals ADD CONSTRAINT fk_t_late_arrivals_attendance_id FOREIGN KEY (attendance_id) REFERENCES t_attendances(attendance_id) ON DELETE RESTRICT ON UPDATE RESTRICT;
CREATE UNIQUE INDEX idx_t_late_arrivals_id ON t_late_arrivals(late_arrival_id);
