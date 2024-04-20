CREATE TABLE m_students (
	m_students_pkey BIGSERIAL,
    student_id UUID NOT NULL DEFAULT uuid_generate_v4(),
	member_id UUID NOT NULL
);
ALTER TABLE m_students ADD CONSTRAINT m_students_pkey PRIMARY KEY (m_students_pkey);
ALTER TABLE m_students ADD CONSTRAINT fk_m_students_member_id FOREIGN KEY (member_id) REFERENCES m_members(member_id) ON DELETE RESTRICT ON UPDATE RESTRICT;
CREATE UNIQUE INDEX idx_m_students_id ON m_students(student_id);
