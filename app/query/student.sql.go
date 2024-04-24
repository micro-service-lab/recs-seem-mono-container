// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: student.sql

package query

import (
	"context"

	"github.com/google/uuid"
)

const countStudents = `-- name: CountStudents :one
SELECT COUNT(*) FROM m_students
`

func (q *Queries) CountStudents(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, countStudents)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createStudent = `-- name: CreateStudent :one
INSERT INTO m_students (member_id) VALUES ($1) RETURNING m_students_pkey, student_id, member_id
`

func (q *Queries) CreateStudent(ctx context.Context, memberID uuid.UUID) (Student, error) {
	row := q.db.QueryRow(ctx, createStudent, memberID)
	var i Student
	err := row.Scan(&i.MStudentsPkey, &i.StudentID, &i.MemberID)
	return i, err
}

const deleteStudent = `-- name: DeleteStudent :exec
DELETE FROM m_students WHERE student_id = $1
`

func (q *Queries) DeleteStudent(ctx context.Context, studentID uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteStudent, studentID)
	return err
}

const findStudentByID = `-- name: FindStudentByID :one
SELECT m_students_pkey, student_id, member_id FROM m_students WHERE student_id = $1
`

func (q *Queries) FindStudentByID(ctx context.Context, studentID uuid.UUID) (Student, error) {
	row := q.db.QueryRow(ctx, findStudentByID, studentID)
	var i Student
	err := row.Scan(&i.MStudentsPkey, &i.StudentID, &i.MemberID)
	return i, err
}

const findStudentByMemberID = `-- name: FindStudentByMemberID :one
SELECT m_students_pkey, student_id, member_id FROM m_students WHERE member_id = $1
`

func (q *Queries) FindStudentByMemberID(ctx context.Context, memberID uuid.UUID) (Student, error) {
	row := q.db.QueryRow(ctx, findStudentByMemberID, memberID)
	var i Student
	err := row.Scan(&i.MStudentsPkey, &i.StudentID, &i.MemberID)
	return i, err
}

const findStudentWithMember = `-- name: FindStudentWithMember :one
SELECT m_students.m_students_pkey, m_students.student_id, m_students.member_id, m_members.m_members_pkey, m_members.member_id, m_members.login_id, m_members.password, m_members.email, m_members.name, m_members.attend_status_id, m_members.profile_image_id, m_members.grade_id, m_members.group_id, m_members.personal_organization_id, m_members.role_id, m_members.created_at, m_members.updated_at FROM m_students
INNER JOIN m_members ON m_students.member_id = m_members.member_id
WHERE student_id = $1
`

type FindStudentWithMemberRow struct {
	Student Student `json:"student"`
	Member  Member  `json:"member"`
}

func (q *Queries) FindStudentWithMember(ctx context.Context, studentID uuid.UUID) (FindStudentWithMemberRow, error) {
	row := q.db.QueryRow(ctx, findStudentWithMember, studentID)
	var i FindStudentWithMemberRow
	err := row.Scan(
		&i.Student.MStudentsPkey,
		&i.Student.StudentID,
		&i.Student.MemberID,
		&i.Member.MMembersPkey,
		&i.Member.MemberID,
		&i.Member.LoginID,
		&i.Member.Password,
		&i.Member.Email,
		&i.Member.Name,
		&i.Member.AttendStatusID,
		&i.Member.ProfileImageID,
		&i.Member.GradeID,
		&i.Member.GroupID,
		&i.Member.PersonalOrganizationID,
		&i.Member.RoleID,
		&i.Member.CreatedAt,
		&i.Member.UpdatedAt,
	)
	return i, err
}

const getStudents = `-- name: GetStudents :many
SELECT m_students_pkey, student_id, member_id FROM m_students
ORDER BY
	m_students_pkey DESC
LIMIT $1 OFFSET $2
`

type GetStudentsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) GetStudents(ctx context.Context, arg GetStudentsParams) ([]Student, error) {
	rows, err := q.db.Query(ctx, getStudents, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Student{}
	for rows.Next() {
		var i Student
		if err := rows.Scan(&i.MStudentsPkey, &i.StudentID, &i.MemberID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getStudentsWithMember = `-- name: GetStudentsWithMember :many
SELECT m_students.m_students_pkey, m_students.student_id, m_students.member_id, m_members.m_members_pkey, m_members.member_id, m_members.login_id, m_members.password, m_members.email, m_members.name, m_members.attend_status_id, m_members.profile_image_id, m_members.grade_id, m_members.group_id, m_members.personal_organization_id, m_members.role_id, m_members.created_at, m_members.updated_at FROM m_students
INNER JOIN m_members ON m_students.member_id = m_members.member_id
ORDER BY
	m_students_pkey DESC
LIMIT $1 OFFSET $2
`

type GetStudentsWithMemberParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type GetStudentsWithMemberRow struct {
	Student Student `json:"student"`
	Member  Member  `json:"member"`
}

func (q *Queries) GetStudentsWithMember(ctx context.Context, arg GetStudentsWithMemberParams) ([]GetStudentsWithMemberRow, error) {
	rows, err := q.db.Query(ctx, getStudentsWithMember, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetStudentsWithMemberRow{}
	for rows.Next() {
		var i GetStudentsWithMemberRow
		if err := rows.Scan(
			&i.Student.MStudentsPkey,
			&i.Student.StudentID,
			&i.Student.MemberID,
			&i.Member.MMembersPkey,
			&i.Member.MemberID,
			&i.Member.LoginID,
			&i.Member.Password,
			&i.Member.Email,
			&i.Member.Name,
			&i.Member.AttendStatusID,
			&i.Member.ProfileImageID,
			&i.Member.GradeID,
			&i.Member.GroupID,
			&i.Member.PersonalOrganizationID,
			&i.Member.RoleID,
			&i.Member.CreatedAt,
			&i.Member.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
