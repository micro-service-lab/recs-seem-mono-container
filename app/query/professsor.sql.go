// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: professsor.sql

package query

import (
	"context"

	"github.com/google/uuid"
)

const countProfessors = `-- name: CountProfessors :one
SELECT COUNT(*) FROM m_professors
`

func (q *Queries) CountProfessors(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, countProfessors)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createProfessor = `-- name: CreateProfessor :one
INSERT INTO m_professors (member_id) VALUES ($1) RETURNING m_professors_pkey, professor_id, member_id
`

func (q *Queries) CreateProfessor(ctx context.Context, memberID uuid.UUID) (Professor, error) {
	row := q.db.QueryRow(ctx, createProfessor, memberID)
	var i Professor
	err := row.Scan(&i.MProfessorsPkey, &i.ProfessorID, &i.MemberID)
	return i, err
}

const deleteProfessor = `-- name: DeleteProfessor :exec
DELETE FROM m_professors WHERE professor_id = $1
`

func (q *Queries) DeleteProfessor(ctx context.Context, professorID uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteProfessor, professorID)
	return err
}

const findProfessorByID = `-- name: FindProfessorByID :one
SELECT m_professors_pkey, professor_id, member_id FROM m_professors WHERE professor_id = $1
`

func (q *Queries) FindProfessorByID(ctx context.Context, professorID uuid.UUID) (Professor, error) {
	row := q.db.QueryRow(ctx, findProfessorByID, professorID)
	var i Professor
	err := row.Scan(&i.MProfessorsPkey, &i.ProfessorID, &i.MemberID)
	return i, err
}

const findProfessorByIDWithMember = `-- name: FindProfessorByIDWithMember :one
SELECT m_professors.m_professors_pkey, m_professors.professor_id, m_professors.member_id, m_members.m_members_pkey, m_members.member_id, m_members.login_id, m_members.password, m_members.email, m_members.name, m_members.attend_status_id, m_members.profile_image_url, m_members.grade_id, m_members.group_id, m_members.personal_organization_id, m_members.role_id, m_members.created_at, m_members.updated_at FROM m_professors
LEFT JOIN m_members ON m_professors.member_id = m_members.member_id
WHERE professor_id = $1
`

type FindProfessorByIDWithMemberRow struct {
	Professor Professor `json:"professor"`
	Member    Member    `json:"member"`
}

func (q *Queries) FindProfessorByIDWithMember(ctx context.Context, professorID uuid.UUID) (FindProfessorByIDWithMemberRow, error) {
	row := q.db.QueryRow(ctx, findProfessorByIDWithMember, professorID)
	var i FindProfessorByIDWithMemberRow
	err := row.Scan(
		&i.Professor.MProfessorsPkey,
		&i.Professor.ProfessorID,
		&i.Professor.MemberID,
		&i.Member.MMembersPkey,
		&i.Member.MemberID,
		&i.Member.LoginID,
		&i.Member.Password,
		&i.Member.Email,
		&i.Member.Name,
		&i.Member.AttendStatusID,
		&i.Member.ProfileImageUrl,
		&i.Member.GradeID,
		&i.Member.GroupID,
		&i.Member.PersonalOrganizationID,
		&i.Member.RoleID,
		&i.Member.CreatedAt,
		&i.Member.UpdatedAt,
	)
	return i, err
}

const getPluralProfessors = `-- name: GetPluralProfessors :many
SELECT m_professors_pkey, professor_id, member_id FROM m_professors
WHERE member_id = ANY($3::uuid[])
ORDER BY
	m_professors_pkey ASC
LIMIT $1 OFFSET $2
`

type GetPluralProfessorsParams struct {
	Limit     int32       `json:"limit"`
	Offset    int32       `json:"offset"`
	MemberIds []uuid.UUID `json:"member_ids"`
}

func (q *Queries) GetPluralProfessors(ctx context.Context, arg GetPluralProfessorsParams) ([]Professor, error) {
	rows, err := q.db.Query(ctx, getPluralProfessors, arg.Limit, arg.Offset, arg.MemberIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Professor{}
	for rows.Next() {
		var i Professor
		if err := rows.Scan(&i.MProfessorsPkey, &i.ProfessorID, &i.MemberID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getProfessors = `-- name: GetProfessors :many
SELECT m_professors_pkey, professor_id, member_id FROM m_professors
ORDER BY
	m_professors_pkey ASC
`

func (q *Queries) GetProfessors(ctx context.Context) ([]Professor, error) {
	rows, err := q.db.Query(ctx, getProfessors)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Professor{}
	for rows.Next() {
		var i Professor
		if err := rows.Scan(&i.MProfessorsPkey, &i.ProfessorID, &i.MemberID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getProfessorsUseKeysetPaginate = `-- name: GetProfessorsUseKeysetPaginate :many
SELECT m_professors_pkey, professor_id, member_id FROM m_professors
WHERE
	CASE $2::text
		WHEN 'next' THEN
			m_professors_pkey > $3::int
		WHEN 'prev' THEN
			m_professors_pkey < $3::int
	END
ORDER BY
	CASE WHEN $2::text = 'next' THEN m_professors_pkey END ASC,
	CASE WHEN $2::text = 'prev' THEN m_professors_pkey END DESC
LIMIT $1
`

type GetProfessorsUseKeysetPaginateParams struct {
	Limit           int32  `json:"limit"`
	CursorDirection string `json:"cursor_direction"`
	Cursor          int32  `json:"cursor"`
}

func (q *Queries) GetProfessorsUseKeysetPaginate(ctx context.Context, arg GetProfessorsUseKeysetPaginateParams) ([]Professor, error) {
	rows, err := q.db.Query(ctx, getProfessorsUseKeysetPaginate, arg.Limit, arg.CursorDirection, arg.Cursor)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Professor{}
	for rows.Next() {
		var i Professor
		if err := rows.Scan(&i.MProfessorsPkey, &i.ProfessorID, &i.MemberID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getProfessorsUseNumberedPaginate = `-- name: GetProfessorsUseNumberedPaginate :many
SELECT m_professors_pkey, professor_id, member_id FROM m_professors
ORDER BY
	m_professors_pkey ASC
LIMIT $1 OFFSET $2
`

type GetProfessorsUseNumberedPaginateParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) GetProfessorsUseNumberedPaginate(ctx context.Context, arg GetProfessorsUseNumberedPaginateParams) ([]Professor, error) {
	rows, err := q.db.Query(ctx, getProfessorsUseNumberedPaginate, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Professor{}
	for rows.Next() {
		var i Professor
		if err := rows.Scan(&i.MProfessorsPkey, &i.ProfessorID, &i.MemberID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const pluralDeleteProfessors = `-- name: PluralDeleteProfessors :exec
DELETE FROM m_professors WHERE professor_id = ANY($1::uuid[])
`

func (q *Queries) PluralDeleteProfessors(ctx context.Context, dollar_1 []uuid.UUID) error {
	_, err := q.db.Exec(ctx, pluralDeleteProfessors, dollar_1)
	return err
}
