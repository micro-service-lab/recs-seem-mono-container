// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: student.sql

package query

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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

const deleteStudent = `-- name: DeleteStudent :execrows
DELETE FROM m_students WHERE student_id = $1
`

func (q *Queries) DeleteStudent(ctx context.Context, studentID uuid.UUID) (int64, error) {
	result, err := q.db.Exec(ctx, deleteStudent, studentID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
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

const findStudentByIDWithMember = `-- name: FindStudentByIDWithMember :one
SELECT m_students.m_students_pkey, m_students.student_id, m_students.member_id, m_members.m_members_pkey, m_members.member_id, m_members.login_id, m_members.password, m_members.email, m_members.name, m_members.attend_status_id, m_members.profile_image_id, m_members.grade_id, m_members.group_id, m_members.personal_organization_id, m_members.role_id, m_members.created_at, m_members.updated_at FROM m_students
LEFT JOIN m_members ON m_students.member_id = m_members.member_id
WHERE student_id = $1
`

type FindStudentByIDWithMemberRow struct {
	MStudentsPkey pgtype.Int8 `json:"m_students_pkey"`
	StudentID     uuid.UUID   `json:"student_id"`
	MemberID      uuid.UUID   `json:"member_id"`
	Member        Member      `json:"member"`
}

func (q *Queries) FindStudentByIDWithMember(ctx context.Context, studentID uuid.UUID) (FindStudentByIDWithMemberRow, error) {
	row := q.db.QueryRow(ctx, findStudentByIDWithMember, studentID)
	var i FindStudentByIDWithMemberRow
	err := row.Scan(
		&i.MStudentsPkey,
		&i.StudentID,
		&i.MemberID,
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

const getPluralStudents = `-- name: GetPluralStudents :many
SELECT m_students_pkey, student_id, member_id FROM m_students
WHERE member_id = ANY($3::uuid[])
ORDER BY
	m_students_pkey ASC
LIMIT $1 OFFSET $2
`

type GetPluralStudentsParams struct {
	Limit     int32       `json:"limit"`
	Offset    int32       `json:"offset"`
	MemberIds []uuid.UUID `json:"member_ids"`
}

func (q *Queries) GetPluralStudents(ctx context.Context, arg GetPluralStudentsParams) ([]Student, error) {
	rows, err := q.db.Query(ctx, getPluralStudents, arg.Limit, arg.Offset, arg.MemberIds)
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

const getStudents = `-- name: GetStudents :many
SELECT m_students_pkey, student_id, member_id FROM m_students
ORDER BY
	m_students_pkey ASC
`

func (q *Queries) GetStudents(ctx context.Context) ([]Student, error) {
	rows, err := q.db.Query(ctx, getStudents)
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

const getStudentsUseKeysetPaginate = `-- name: GetStudentsUseKeysetPaginate :many
SELECT m_students_pkey, student_id, member_id FROM m_students
WHERE
	CASE $2::text
		WHEN 'next' THEN
			m_students_pkey > $3::int
		WHEN 'prev' THEN
			m_students_pkey < $3::int
	END
ORDER BY
	CASE WHEN $2::text = 'next' THEN m_students_pkey END ASC,
	CASE WHEN $2::text = 'prev' THEN m_students_pkey END DESC
LIMIT $1
`

type GetStudentsUseKeysetPaginateParams struct {
	Limit           int32  `json:"limit"`
	CursorDirection string `json:"cursor_direction"`
	Cursor          int32  `json:"cursor"`
}

func (q *Queries) GetStudentsUseKeysetPaginate(ctx context.Context, arg GetStudentsUseKeysetPaginateParams) ([]Student, error) {
	rows, err := q.db.Query(ctx, getStudentsUseKeysetPaginate, arg.Limit, arg.CursorDirection, arg.Cursor)
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

const getStudentsUseNumberedPaginate = `-- name: GetStudentsUseNumberedPaginate :many
SELECT m_students_pkey, student_id, member_id FROM m_students
ORDER BY
	m_students_pkey ASC
LIMIT $1 OFFSET $2
`

type GetStudentsUseNumberedPaginateParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) GetStudentsUseNumberedPaginate(ctx context.Context, arg GetStudentsUseNumberedPaginateParams) ([]Student, error) {
	rows, err := q.db.Query(ctx, getStudentsUseNumberedPaginate, arg.Limit, arg.Offset)
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

const pluralDeleteStudents = `-- name: PluralDeleteStudents :execrows
DELETE FROM m_students WHERE student_id = ANY($1::uuid[])
`

func (q *Queries) PluralDeleteStudents(ctx context.Context, dollar_1 []uuid.UUID) (int64, error) {
	result, err := q.db.Exec(ctx, pluralDeleteStudents, dollar_1)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}
