// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: grade.sql

package query

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const countGrades = `-- name: CountGrades :one
SELECT COUNT(*) FROM m_grades
`

func (q *Queries) CountGrades(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, countGrades)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createGrade = `-- name: CreateGrade :one
INSERT INTO m_grades (key, organization_id) VALUES ($1, $2) RETURNING m_grades_pkey, grade_id, key, organization_id
`

type CreateGradeParams struct {
	Key            string    `json:"key"`
	OrganizationID uuid.UUID `json:"organization_id"`
}

func (q *Queries) CreateGrade(ctx context.Context, arg CreateGradeParams) (Grade, error) {
	row := q.db.QueryRow(ctx, createGrade, arg.Key, arg.OrganizationID)
	var i Grade
	err := row.Scan(
		&i.MGradesPkey,
		&i.GradeID,
		&i.Key,
		&i.OrganizationID,
	)
	return i, err
}

type CreateGradesParams struct {
	Key            string    `json:"key"`
	OrganizationID uuid.UUID `json:"organization_id"`
}

const deleteGrade = `-- name: DeleteGrade :execrows
DELETE FROM m_grades WHERE grade_id = $1
`

func (q *Queries) DeleteGrade(ctx context.Context, gradeID uuid.UUID) (int64, error) {
	result, err := q.db.Exec(ctx, deleteGrade, gradeID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const deleteGradeByKey = `-- name: DeleteGradeByKey :execrows
DELETE FROM m_grades WHERE key = $1
`

func (q *Queries) DeleteGradeByKey(ctx context.Context, key string) (int64, error) {
	result, err := q.db.Exec(ctx, deleteGradeByKey, key)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const findGradeByID = `-- name: FindGradeByID :one
SELECT m_grades_pkey, grade_id, key, organization_id FROM m_grades WHERE grade_id = $1
`

func (q *Queries) FindGradeByID(ctx context.Context, gradeID uuid.UUID) (Grade, error) {
	row := q.db.QueryRow(ctx, findGradeByID, gradeID)
	var i Grade
	err := row.Scan(
		&i.MGradesPkey,
		&i.GradeID,
		&i.Key,
		&i.OrganizationID,
	)
	return i, err
}

const findGradeByIDWithOrganization = `-- name: FindGradeByIDWithOrganization :one
SELECT m_grades.m_grades_pkey, m_grades.grade_id, m_grades.key, m_grades.organization_id, m_organizations.name organization_name, m_organizations.description organization_description,
m_organizations.color organization_color, m_organizations.is_personal organization_is_personal,
m_organizations.is_whole organization_is_whole, m_organizations.chat_room_id organization_chat_room_id
FROM m_grades
LEFT JOIN m_organizations ON m_grades.organization_id = m_organizations.organization_id
WHERE grade_id = $1
`

type FindGradeByIDWithOrganizationRow struct {
	MGradesPkey             pgtype.Int8 `json:"m_grades_pkey"`
	GradeID                 uuid.UUID   `json:"grade_id"`
	Key                     string      `json:"key"`
	OrganizationID          uuid.UUID   `json:"organization_id"`
	OrganizationName        pgtype.Text `json:"organization_name"`
	OrganizationDescription pgtype.Text `json:"organization_description"`
	OrganizationColor       pgtype.Text `json:"organization_color"`
	OrganizationIsPersonal  pgtype.Bool `json:"organization_is_personal"`
	OrganizationIsWhole     pgtype.Bool `json:"organization_is_whole"`
	OrganizationChatRoomID  pgtype.UUID `json:"organization_chat_room_id"`
}

func (q *Queries) FindGradeByIDWithOrganization(ctx context.Context, gradeID uuid.UUID) (FindGradeByIDWithOrganizationRow, error) {
	row := q.db.QueryRow(ctx, findGradeByIDWithOrganization, gradeID)
	var i FindGradeByIDWithOrganizationRow
	err := row.Scan(
		&i.MGradesPkey,
		&i.GradeID,
		&i.Key,
		&i.OrganizationID,
		&i.OrganizationName,
		&i.OrganizationDescription,
		&i.OrganizationColor,
		&i.OrganizationIsPersonal,
		&i.OrganizationIsWhole,
		&i.OrganizationChatRoomID,
	)
	return i, err
}

const findGradeByKey = `-- name: FindGradeByKey :one
SELECT m_grades_pkey, grade_id, key, organization_id FROM m_grades WHERE key = $1
`

func (q *Queries) FindGradeByKey(ctx context.Context, key string) (Grade, error) {
	row := q.db.QueryRow(ctx, findGradeByKey, key)
	var i Grade
	err := row.Scan(
		&i.MGradesPkey,
		&i.GradeID,
		&i.Key,
		&i.OrganizationID,
	)
	return i, err
}

const findGradeByKeyWithOrganization = `-- name: FindGradeByKeyWithOrganization :one
SELECT m_grades.m_grades_pkey, m_grades.grade_id, m_grades.key, m_grades.organization_id, m_organizations.name organization_name, m_organizations.description organization_description,
m_organizations.color organization_color, m_organizations.is_personal organization_is_personal,
m_organizations.is_whole organization_is_whole, m_organizations.chat_room_id organization_chat_room_id
FROM m_grades
LEFT JOIN m_organizations ON m_grades.organization_id = m_organizations.organization_id
WHERE key = $1
`

type FindGradeByKeyWithOrganizationRow struct {
	MGradesPkey             pgtype.Int8 `json:"m_grades_pkey"`
	GradeID                 uuid.UUID   `json:"grade_id"`
	Key                     string      `json:"key"`
	OrganizationID          uuid.UUID   `json:"organization_id"`
	OrganizationName        pgtype.Text `json:"organization_name"`
	OrganizationDescription pgtype.Text `json:"organization_description"`
	OrganizationColor       pgtype.Text `json:"organization_color"`
	OrganizationIsPersonal  pgtype.Bool `json:"organization_is_personal"`
	OrganizationIsWhole     pgtype.Bool `json:"organization_is_whole"`
	OrganizationChatRoomID  pgtype.UUID `json:"organization_chat_room_id"`
}

func (q *Queries) FindGradeByKeyWithOrganization(ctx context.Context, key string) (FindGradeByKeyWithOrganizationRow, error) {
	row := q.db.QueryRow(ctx, findGradeByKeyWithOrganization, key)
	var i FindGradeByKeyWithOrganizationRow
	err := row.Scan(
		&i.MGradesPkey,
		&i.GradeID,
		&i.Key,
		&i.OrganizationID,
		&i.OrganizationName,
		&i.OrganizationDescription,
		&i.OrganizationColor,
		&i.OrganizationIsPersonal,
		&i.OrganizationIsWhole,
		&i.OrganizationChatRoomID,
	)
	return i, err
}

const getGrades = `-- name: GetGrades :many
SELECT m_grades_pkey, grade_id, key, organization_id FROM m_grades
ORDER BY
	m_grades_pkey ASC
`

func (q *Queries) GetGrades(ctx context.Context) ([]Grade, error) {
	rows, err := q.db.Query(ctx, getGrades)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Grade{}
	for rows.Next() {
		var i Grade
		if err := rows.Scan(
			&i.MGradesPkey,
			&i.GradeID,
			&i.Key,
			&i.OrganizationID,
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

const getGradesUseKeysetPaginate = `-- name: GetGradesUseKeysetPaginate :many
SELECT m_grades_pkey, grade_id, key, organization_id FROM m_grades
WHERE
	CASE $2::text
		WHEN 'next' THEN
			m_grades_pkey > $3::int
		WHEN 'prev' THEN
			m_grades_pkey < $3::int
	END
ORDER BY
	CASE WHEN $2::text = 'next' THEN m_grades_pkey END ASC,
	CASE WHEN $2::text = 'prev' THEN m_grades_pkey END DESC
LIMIT $1
`

type GetGradesUseKeysetPaginateParams struct {
	Limit           int32  `json:"limit"`
	CursorDirection string `json:"cursor_direction"`
	Cursor          int32  `json:"cursor"`
}

func (q *Queries) GetGradesUseKeysetPaginate(ctx context.Context, arg GetGradesUseKeysetPaginateParams) ([]Grade, error) {
	rows, err := q.db.Query(ctx, getGradesUseKeysetPaginate, arg.Limit, arg.CursorDirection, arg.Cursor)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Grade{}
	for rows.Next() {
		var i Grade
		if err := rows.Scan(
			&i.MGradesPkey,
			&i.GradeID,
			&i.Key,
			&i.OrganizationID,
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

const getGradesUseNumberedPaginate = `-- name: GetGradesUseNumberedPaginate :many
SELECT m_grades_pkey, grade_id, key, organization_id FROM m_grades
ORDER BY
	m_grades_pkey ASC
LIMIT $1 OFFSET $2
`

type GetGradesUseNumberedPaginateParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) GetGradesUseNumberedPaginate(ctx context.Context, arg GetGradesUseNumberedPaginateParams) ([]Grade, error) {
	rows, err := q.db.Query(ctx, getGradesUseNumberedPaginate, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Grade{}
	for rows.Next() {
		var i Grade
		if err := rows.Scan(
			&i.MGradesPkey,
			&i.GradeID,
			&i.Key,
			&i.OrganizationID,
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

const getGradesWithOrganization = `-- name: GetGradesWithOrganization :many
SELECT m_grades.m_grades_pkey, m_grades.grade_id, m_grades.key, m_grades.organization_id, m_organizations.name organization_name, m_organizations.description organization_description,
m_organizations.color organization_color, m_organizations.is_personal organization_is_personal,
m_organizations.is_whole organization_is_whole, m_organizations.chat_room_id organization_chat_room_id
FROM m_grades
LEFT JOIN m_organizations ON m_grades.organization_id = m_organizations.organization_id
ORDER BY
	m_grades_pkey ASC
`

type GetGradesWithOrganizationRow struct {
	MGradesPkey             pgtype.Int8 `json:"m_grades_pkey"`
	GradeID                 uuid.UUID   `json:"grade_id"`
	Key                     string      `json:"key"`
	OrganizationID          uuid.UUID   `json:"organization_id"`
	OrganizationName        pgtype.Text `json:"organization_name"`
	OrganizationDescription pgtype.Text `json:"organization_description"`
	OrganizationColor       pgtype.Text `json:"organization_color"`
	OrganizationIsPersonal  pgtype.Bool `json:"organization_is_personal"`
	OrganizationIsWhole     pgtype.Bool `json:"organization_is_whole"`
	OrganizationChatRoomID  pgtype.UUID `json:"organization_chat_room_id"`
}

func (q *Queries) GetGradesWithOrganization(ctx context.Context) ([]GetGradesWithOrganizationRow, error) {
	rows, err := q.db.Query(ctx, getGradesWithOrganization)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetGradesWithOrganizationRow{}
	for rows.Next() {
		var i GetGradesWithOrganizationRow
		if err := rows.Scan(
			&i.MGradesPkey,
			&i.GradeID,
			&i.Key,
			&i.OrganizationID,
			&i.OrganizationName,
			&i.OrganizationDescription,
			&i.OrganizationColor,
			&i.OrganizationIsPersonal,
			&i.OrganizationIsWhole,
			&i.OrganizationChatRoomID,
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

const getGradesWithOrganizationUseKeysetPaginate = `-- name: GetGradesWithOrganizationUseKeysetPaginate :many
SELECT m_grades.m_grades_pkey, m_grades.grade_id, m_grades.key, m_grades.organization_id, m_organizations.name organization_name, m_organizations.description organization_description,
m_organizations.color organization_color, m_organizations.is_personal organization_is_personal,
m_organizations.is_whole organization_is_whole, m_organizations.chat_room_id organization_chat_room_id
FROM m_grades
LEFT JOIN m_organizations ON m_grades.organization_id = m_organizations.organization_id
WHERE
	CASE $2::text
		WHEN 'next' THEN
			m_grades_pkey > $3::int
		WHEN 'prev' THEN
			m_grades_pkey < $3::int
	END
ORDER BY
	CASE WHEN $2::text = 'next' THEN m_grades_pkey END ASC,
	CASE WHEN $2::text = 'prev' THEN m_grades_pkey END DESC
LIMIT $1
`

type GetGradesWithOrganizationUseKeysetPaginateParams struct {
	Limit           int32  `json:"limit"`
	CursorDirection string `json:"cursor_direction"`
	Cursor          int32  `json:"cursor"`
}

type GetGradesWithOrganizationUseKeysetPaginateRow struct {
	MGradesPkey             pgtype.Int8 `json:"m_grades_pkey"`
	GradeID                 uuid.UUID   `json:"grade_id"`
	Key                     string      `json:"key"`
	OrganizationID          uuid.UUID   `json:"organization_id"`
	OrganizationName        pgtype.Text `json:"organization_name"`
	OrganizationDescription pgtype.Text `json:"organization_description"`
	OrganizationColor       pgtype.Text `json:"organization_color"`
	OrganizationIsPersonal  pgtype.Bool `json:"organization_is_personal"`
	OrganizationIsWhole     pgtype.Bool `json:"organization_is_whole"`
	OrganizationChatRoomID  pgtype.UUID `json:"organization_chat_room_id"`
}

func (q *Queries) GetGradesWithOrganizationUseKeysetPaginate(ctx context.Context, arg GetGradesWithOrganizationUseKeysetPaginateParams) ([]GetGradesWithOrganizationUseKeysetPaginateRow, error) {
	rows, err := q.db.Query(ctx, getGradesWithOrganizationUseKeysetPaginate, arg.Limit, arg.CursorDirection, arg.Cursor)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetGradesWithOrganizationUseKeysetPaginateRow{}
	for rows.Next() {
		var i GetGradesWithOrganizationUseKeysetPaginateRow
		if err := rows.Scan(
			&i.MGradesPkey,
			&i.GradeID,
			&i.Key,
			&i.OrganizationID,
			&i.OrganizationName,
			&i.OrganizationDescription,
			&i.OrganizationColor,
			&i.OrganizationIsPersonal,
			&i.OrganizationIsWhole,
			&i.OrganizationChatRoomID,
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

const getGradesWithOrganizationUseNumberedPaginate = `-- name: GetGradesWithOrganizationUseNumberedPaginate :many
SELECT m_grades.m_grades_pkey, m_grades.grade_id, m_grades.key, m_grades.organization_id, m_organizations.name organization_name, m_organizations.description organization_description,
m_organizations.color organization_color, m_organizations.is_personal organization_is_personal,
m_organizations.is_whole organization_is_whole, m_organizations.chat_room_id organization_chat_room_id
FROM m_grades
LEFT JOIN m_organizations ON m_grades.organization_id = m_organizations.organization_id
ORDER BY
	m_grades_pkey ASC
LIMIT $1 OFFSET $2
`

type GetGradesWithOrganizationUseNumberedPaginateParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type GetGradesWithOrganizationUseNumberedPaginateRow struct {
	MGradesPkey             pgtype.Int8 `json:"m_grades_pkey"`
	GradeID                 uuid.UUID   `json:"grade_id"`
	Key                     string      `json:"key"`
	OrganizationID          uuid.UUID   `json:"organization_id"`
	OrganizationName        pgtype.Text `json:"organization_name"`
	OrganizationDescription pgtype.Text `json:"organization_description"`
	OrganizationColor       pgtype.Text `json:"organization_color"`
	OrganizationIsPersonal  pgtype.Bool `json:"organization_is_personal"`
	OrganizationIsWhole     pgtype.Bool `json:"organization_is_whole"`
	OrganizationChatRoomID  pgtype.UUID `json:"organization_chat_room_id"`
}

func (q *Queries) GetGradesWithOrganizationUseNumberedPaginate(ctx context.Context, arg GetGradesWithOrganizationUseNumberedPaginateParams) ([]GetGradesWithOrganizationUseNumberedPaginateRow, error) {
	rows, err := q.db.Query(ctx, getGradesWithOrganizationUseNumberedPaginate, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetGradesWithOrganizationUseNumberedPaginateRow{}
	for rows.Next() {
		var i GetGradesWithOrganizationUseNumberedPaginateRow
		if err := rows.Scan(
			&i.MGradesPkey,
			&i.GradeID,
			&i.Key,
			&i.OrganizationID,
			&i.OrganizationName,
			&i.OrganizationDescription,
			&i.OrganizationColor,
			&i.OrganizationIsPersonal,
			&i.OrganizationIsWhole,
			&i.OrganizationChatRoomID,
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

const getPluralGrades = `-- name: GetPluralGrades :many
SELECT m_grades_pkey, grade_id, key, organization_id FROM m_grades
WHERE grade_id = ANY($1::uuid[])
ORDER BY
	m_grades_pkey ASC
`

func (q *Queries) GetPluralGrades(ctx context.Context, gradeIds []uuid.UUID) ([]Grade, error) {
	rows, err := q.db.Query(ctx, getPluralGrades, gradeIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Grade{}
	for rows.Next() {
		var i Grade
		if err := rows.Scan(
			&i.MGradesPkey,
			&i.GradeID,
			&i.Key,
			&i.OrganizationID,
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

const getPluralGradesUseNumberedPaginate = `-- name: GetPluralGradesUseNumberedPaginate :many
SELECT m_grades_pkey, grade_id, key, organization_id FROM m_grades
WHERE grade_id = ANY($3::uuid[])
ORDER BY
	m_grades_pkey ASC
LIMIT $1 OFFSET $2
`

type GetPluralGradesUseNumberedPaginateParams struct {
	Limit    int32       `json:"limit"`
	Offset   int32       `json:"offset"`
	GradeIds []uuid.UUID `json:"grade_ids"`
}

func (q *Queries) GetPluralGradesUseNumberedPaginate(ctx context.Context, arg GetPluralGradesUseNumberedPaginateParams) ([]Grade, error) {
	rows, err := q.db.Query(ctx, getPluralGradesUseNumberedPaginate, arg.Limit, arg.Offset, arg.GradeIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Grade{}
	for rows.Next() {
		var i Grade
		if err := rows.Scan(
			&i.MGradesPkey,
			&i.GradeID,
			&i.Key,
			&i.OrganizationID,
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

const getPluralGradesWithOrganization = `-- name: GetPluralGradesWithOrganization :many
SELECT m_grades.m_grades_pkey, m_grades.grade_id, m_grades.key, m_grades.organization_id, m_organizations.name organization_name, m_organizations.description organization_description,
m_organizations.color organization_color, m_organizations.is_personal organization_is_personal,
m_organizations.is_whole organization_is_whole, m_organizations.chat_room_id organization_chat_room_id
FROM m_grades
LEFT JOIN m_organizations ON m_grades.organization_id = m_organizations.organization_id
WHERE grade_id = ANY($1::uuid[])
ORDER BY
	m_grades_pkey ASC
`

type GetPluralGradesWithOrganizationRow struct {
	MGradesPkey             pgtype.Int8 `json:"m_grades_pkey"`
	GradeID                 uuid.UUID   `json:"grade_id"`
	Key                     string      `json:"key"`
	OrganizationID          uuid.UUID   `json:"organization_id"`
	OrganizationName        pgtype.Text `json:"organization_name"`
	OrganizationDescription pgtype.Text `json:"organization_description"`
	OrganizationColor       pgtype.Text `json:"organization_color"`
	OrganizationIsPersonal  pgtype.Bool `json:"organization_is_personal"`
	OrganizationIsWhole     pgtype.Bool `json:"organization_is_whole"`
	OrganizationChatRoomID  pgtype.UUID `json:"organization_chat_room_id"`
}

func (q *Queries) GetPluralGradesWithOrganization(ctx context.Context, gradeIds []uuid.UUID) ([]GetPluralGradesWithOrganizationRow, error) {
	rows, err := q.db.Query(ctx, getPluralGradesWithOrganization, gradeIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPluralGradesWithOrganizationRow{}
	for rows.Next() {
		var i GetPluralGradesWithOrganizationRow
		if err := rows.Scan(
			&i.MGradesPkey,
			&i.GradeID,
			&i.Key,
			&i.OrganizationID,
			&i.OrganizationName,
			&i.OrganizationDescription,
			&i.OrganizationColor,
			&i.OrganizationIsPersonal,
			&i.OrganizationIsWhole,
			&i.OrganizationChatRoomID,
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

const getPluralGradesWithOrganizationUseNumberedPaginate = `-- name: GetPluralGradesWithOrganizationUseNumberedPaginate :many
SELECT m_grades.m_grades_pkey, m_grades.grade_id, m_grades.key, m_grades.organization_id, m_organizations.name organization_name, m_organizations.description organization_description,
m_organizations.color organization_color, m_organizations.is_personal organization_is_personal,
m_organizations.is_whole organization_is_whole, m_organizations.chat_room_id organization_chat_room_id
FROM m_grades
LEFT JOIN m_organizations ON m_grades.organization_id = m_organizations.organization_id
WHERE grade_id = ANY($3::uuid[])
ORDER BY
	m_grades_pkey ASC
LIMIT $1 OFFSET $2
`

type GetPluralGradesWithOrganizationUseNumberedPaginateParams struct {
	Limit    int32       `json:"limit"`
	Offset   int32       `json:"offset"`
	GradeIds []uuid.UUID `json:"grade_ids"`
}

type GetPluralGradesWithOrganizationUseNumberedPaginateRow struct {
	MGradesPkey             pgtype.Int8 `json:"m_grades_pkey"`
	GradeID                 uuid.UUID   `json:"grade_id"`
	Key                     string      `json:"key"`
	OrganizationID          uuid.UUID   `json:"organization_id"`
	OrganizationName        pgtype.Text `json:"organization_name"`
	OrganizationDescription pgtype.Text `json:"organization_description"`
	OrganizationColor       pgtype.Text `json:"organization_color"`
	OrganizationIsPersonal  pgtype.Bool `json:"organization_is_personal"`
	OrganizationIsWhole     pgtype.Bool `json:"organization_is_whole"`
	OrganizationChatRoomID  pgtype.UUID `json:"organization_chat_room_id"`
}

func (q *Queries) GetPluralGradesWithOrganizationUseNumberedPaginate(ctx context.Context, arg GetPluralGradesWithOrganizationUseNumberedPaginateParams) ([]GetPluralGradesWithOrganizationUseNumberedPaginateRow, error) {
	rows, err := q.db.Query(ctx, getPluralGradesWithOrganizationUseNumberedPaginate, arg.Limit, arg.Offset, arg.GradeIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPluralGradesWithOrganizationUseNumberedPaginateRow{}
	for rows.Next() {
		var i GetPluralGradesWithOrganizationUseNumberedPaginateRow
		if err := rows.Scan(
			&i.MGradesPkey,
			&i.GradeID,
			&i.Key,
			&i.OrganizationID,
			&i.OrganizationName,
			&i.OrganizationDescription,
			&i.OrganizationColor,
			&i.OrganizationIsPersonal,
			&i.OrganizationIsWhole,
			&i.OrganizationChatRoomID,
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

const pluralDeleteGrades = `-- name: PluralDeleteGrades :execrows
DELETE FROM m_grades WHERE grade_id = ANY($1::uuid[])
`

func (q *Queries) PluralDeleteGrades(ctx context.Context, gradeIds []uuid.UUID) (int64, error) {
	result, err := q.db.Exec(ctx, pluralDeleteGrades, gradeIds)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}
