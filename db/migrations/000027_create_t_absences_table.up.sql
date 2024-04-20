CREATE TABLE t_absences (
	t_absences_pkey BIGSERIAL,
    absence_id UUID NOT NULL DEFAULT uuid_generate_v4(),
	attendance_id UUID NOT NULL
);
ALTER TABLE t_absences ADD CONSTRAINT t_absences_pkey PRIMARY KEY (t_absences_pkey);
ALTER TABLE t_absences ADD CONSTRAINT fk_t_absences_attendance_id FOREIGN KEY (attendance_id) REFERENCES t_attendances(attendance_id) ON DELETE RESTRICT ON UPDATE RESTRICT;
CREATE UNIQUE INDEX idx_t_absences_id ON t_absences(absence_id);
