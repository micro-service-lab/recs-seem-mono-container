// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: attendance_type.sql

package query

import (
	"context"

	"github.com/google/uuid"
)

const countAttendanceTypes = `-- name: CountAttendanceTypes :one
SELECT COUNT(*) FROM m_attendance_types
WHERE
	CASE WHEN $1::boolean = true THEN name LIKE '%' || $2::text || '%' ELSE TRUE END
`

type CountAttendanceTypesParams struct {
	WhereLikeName bool   `json:"where_like_name"`
	SearchName    string `json:"search_name"`
}

func (q *Queries) CountAttendanceTypes(ctx context.Context, arg CountAttendanceTypesParams) (int64, error) {
	row := q.db.QueryRow(ctx, countAttendanceTypes, arg.WhereLikeName, arg.SearchName)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createAttendanceType = `-- name: CreateAttendanceType :one
INSERT INTO m_attendance_types (name, key, color) VALUES ($1, $2, $3) RETURNING m_attendance_types_pkey, attendance_type_id, name, key, color
`

type CreateAttendanceTypeParams struct {
	Name  string `json:"name"`
	Key   string `json:"key"`
	Color string `json:"color"`
}

func (q *Queries) CreateAttendanceType(ctx context.Context, arg CreateAttendanceTypeParams) (AttendanceType, error) {
	row := q.db.QueryRow(ctx, createAttendanceType, arg.Name, arg.Key, arg.Color)
	var i AttendanceType
	err := row.Scan(
		&i.MAttendanceTypesPkey,
		&i.AttendanceTypeID,
		&i.Name,
		&i.Key,
		&i.Color,
	)
	return i, err
}

type CreateAttendanceTypesParams struct {
	Name  string `json:"name"`
	Key   string `json:"key"`
	Color string `json:"color"`
}

const deleteAttendanceType = `-- name: DeleteAttendanceType :exec
DELETE FROM m_attendance_types WHERE attendance_type_id = $1
`

func (q *Queries) DeleteAttendanceType(ctx context.Context, attendanceTypeID uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteAttendanceType, attendanceTypeID)
	return err
}

const deleteAttendanceTypeByKey = `-- name: DeleteAttendanceTypeByKey :exec
DELETE FROM m_attendance_types WHERE key = $1
`

func (q *Queries) DeleteAttendanceTypeByKey(ctx context.Context, key string) error {
	_, err := q.db.Exec(ctx, deleteAttendanceTypeByKey, key)
	return err
}

const findAttendanceTypeByID = `-- name: FindAttendanceTypeByID :one
SELECT m_attendance_types_pkey, attendance_type_id, name, key, color FROM m_attendance_types WHERE attendance_type_id = $1
`

func (q *Queries) FindAttendanceTypeByID(ctx context.Context, attendanceTypeID uuid.UUID) (AttendanceType, error) {
	row := q.db.QueryRow(ctx, findAttendanceTypeByID, attendanceTypeID)
	var i AttendanceType
	err := row.Scan(
		&i.MAttendanceTypesPkey,
		&i.AttendanceTypeID,
		&i.Name,
		&i.Key,
		&i.Color,
	)
	return i, err
}

const findAttendanceTypeByKey = `-- name: FindAttendanceTypeByKey :one
SELECT m_attendance_types_pkey, attendance_type_id, name, key, color FROM m_attendance_types WHERE key = $1
`

func (q *Queries) FindAttendanceTypeByKey(ctx context.Context, key string) (AttendanceType, error) {
	row := q.db.QueryRow(ctx, findAttendanceTypeByKey, key)
	var i AttendanceType
	err := row.Scan(
		&i.MAttendanceTypesPkey,
		&i.AttendanceTypeID,
		&i.Name,
		&i.Key,
		&i.Color,
	)
	return i, err
}

const getAttendanceTypes = `-- name: GetAttendanceTypes :many
SELECT m_attendance_types_pkey, attendance_type_id, name, key, color FROM m_attendance_types
WHERE
	CASE WHEN $1::boolean = true THEN m_attendance_types.name LIKE '%' || $2::text || '%' ELSE TRUE END
ORDER BY
	CASE WHEN $3::text = 'name' THEN name END ASC,
	CASE WHEN $3::text = 'r_name' THEN name END DESC,
	m_attendance_types_pkey ASC
`

type GetAttendanceTypesParams struct {
	WhereLikeName bool   `json:"where_like_name"`
	SearchName    string `json:"search_name"`
	OrderMethod   string `json:"order_method"`
}

func (q *Queries) GetAttendanceTypes(ctx context.Context, arg GetAttendanceTypesParams) ([]AttendanceType, error) {
	rows, err := q.db.Query(ctx, getAttendanceTypes, arg.WhereLikeName, arg.SearchName, arg.OrderMethod)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []AttendanceType{}
	for rows.Next() {
		var i AttendanceType
		if err := rows.Scan(
			&i.MAttendanceTypesPkey,
			&i.AttendanceTypeID,
			&i.Name,
			&i.Key,
			&i.Color,
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

const getAttendanceTypesUseKeysetPaginate = `-- name: GetAttendanceTypesUseKeysetPaginate :many
SELECT m_attendance_types_pkey, attendance_type_id, name, key, color FROM m_attendance_types
WHERE
	CASE WHEN $2::boolean = true THEN m_attendance_types.name LIKE '%' || $3::text || '%' ELSE TRUE END
AND
	CASE $4::text
		WHEN 'next' THEN
			CASE $5::text
				WHEN 'name' THEN name > $6 OR (name = $6 AND m_attendance_types_pkey > $7::int)
				WHEN 'r_name' THEN name < $6 OR (name = $6 AND m_attendance_types_pkey > $7::int)
				ELSE m_attendance_types_pkey > $7::int
			END
		WHEN 'prev' THEN
			CASE $5::text
				WHEN 'name' THEN name < $6 OR (name = $6 AND m_attendance_types_pkey < $7::int)
				WHEN 'r_name' THEN name > $6 OR (name = $6 AND m_attendance_types_pkey < $7::int)
				ELSE m_attendance_types_pkey < $7::int
			END
	END
ORDER BY
	CASE WHEN $5::text = 'name' AND $4::text = 'next' THEN name END ASC,
	CASE WHEN $5::text = 'name' AND $4::text = 'prev' THEN name END DESC,
	CASE WHEN $5::text = 'r_name' AND $4::text = 'next' THEN name END ASC,
	CASE WHEN $5::text = 'r_name' AND $4::text = 'prev' THEN name END DESC,
	CASE WHEN $4::text = 'next' THEN m_attendance_types_pkey END ASC,
	CASE WHEN $4::text = 'prev' THEN m_attendance_types_pkey END DESC
LIMIT $1
`

type GetAttendanceTypesUseKeysetPaginateParams struct {
	Limit           int32  `json:"limit"`
	WhereLikeName   bool   `json:"where_like_name"`
	SearchName      string `json:"search_name"`
	CursorDirection string `json:"cursor_direction"`
	OrderMethod     string `json:"order_method"`
	NameCursor      string `json:"name_cursor"`
	Cursor          int32  `json:"cursor"`
}

func (q *Queries) GetAttendanceTypesUseKeysetPaginate(ctx context.Context, arg GetAttendanceTypesUseKeysetPaginateParams) ([]AttendanceType, error) {
	rows, err := q.db.Query(ctx, getAttendanceTypesUseKeysetPaginate,
		arg.Limit,
		arg.WhereLikeName,
		arg.SearchName,
		arg.CursorDirection,
		arg.OrderMethod,
		arg.NameCursor,
		arg.Cursor,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []AttendanceType{}
	for rows.Next() {
		var i AttendanceType
		if err := rows.Scan(
			&i.MAttendanceTypesPkey,
			&i.AttendanceTypeID,
			&i.Name,
			&i.Key,
			&i.Color,
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

const getAttendanceTypesUseNumberedPaginate = `-- name: GetAttendanceTypesUseNumberedPaginate :many
SELECT m_attendance_types_pkey, attendance_type_id, name, key, color FROM m_attendance_types
WHERE
	CASE WHEN $3::boolean = true THEN m_attendance_types.name LIKE '%' || $4::text || '%' ELSE TRUE END
ORDER BY
	CASE WHEN $5::text = 'name' THEN name END ASC,
	CASE WHEN $5::text = 'r_name' THEN name END DESC,
	m_attendance_types_pkey ASC
LIMIT $1 OFFSET $2
`

type GetAttendanceTypesUseNumberedPaginateParams struct {
	Limit         int32  `json:"limit"`
	Offset        int32  `json:"offset"`
	WhereLikeName bool   `json:"where_like_name"`
	SearchName    string `json:"search_name"`
	OrderMethod   string `json:"order_method"`
}

func (q *Queries) GetAttendanceTypesUseNumberedPaginate(ctx context.Context, arg GetAttendanceTypesUseNumberedPaginateParams) ([]AttendanceType, error) {
	rows, err := q.db.Query(ctx, getAttendanceTypesUseNumberedPaginate,
		arg.Limit,
		arg.Offset,
		arg.WhereLikeName,
		arg.SearchName,
		arg.OrderMethod,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []AttendanceType{}
	for rows.Next() {
		var i AttendanceType
		if err := rows.Scan(
			&i.MAttendanceTypesPkey,
			&i.AttendanceTypeID,
			&i.Name,
			&i.Key,
			&i.Color,
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

const getPluralAttendanceTypes = `-- name: GetPluralAttendanceTypes :many
SELECT m_attendance_types_pkey, attendance_type_id, name, key, color FROM m_attendance_types
WHERE attendance_type_id = ANY($3::uuid[])
ORDER BY
	m_attendance_types_pkey ASC
LIMIT $1 OFFSET $2
`

type GetPluralAttendanceTypesParams struct {
	Limit             int32       `json:"limit"`
	Offset            int32       `json:"offset"`
	AttendanceTypeIds []uuid.UUID `json:"attendance_type_ids"`
}

func (q *Queries) GetPluralAttendanceTypes(ctx context.Context, arg GetPluralAttendanceTypesParams) ([]AttendanceType, error) {
	rows, err := q.db.Query(ctx, getPluralAttendanceTypes, arg.Limit, arg.Offset, arg.AttendanceTypeIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []AttendanceType{}
	for rows.Next() {
		var i AttendanceType
		if err := rows.Scan(
			&i.MAttendanceTypesPkey,
			&i.AttendanceTypeID,
			&i.Name,
			&i.Key,
			&i.Color,
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

const pluralDeleteAttendanceTypes = `-- name: PluralDeleteAttendanceTypes :exec
DELETE FROM m_attendance_types WHERE attendance_type_id = ANY($1::uuid[])
`

func (q *Queries) PluralDeleteAttendanceTypes(ctx context.Context, dollar_1 []uuid.UUID) error {
	_, err := q.db.Exec(ctx, pluralDeleteAttendanceTypes, dollar_1)
	return err
}

const updateAttendanceType = `-- name: UpdateAttendanceType :one
UPDATE m_attendance_types SET name = $2, key = $3, color = $4 WHERE attendance_type_id = $1 RETURNING m_attendance_types_pkey, attendance_type_id, name, key, color
`

type UpdateAttendanceTypeParams struct {
	AttendanceTypeID uuid.UUID `json:"attendance_type_id"`
	Name             string    `json:"name"`
	Key              string    `json:"key"`
	Color            string    `json:"color"`
}

func (q *Queries) UpdateAttendanceType(ctx context.Context, arg UpdateAttendanceTypeParams) (AttendanceType, error) {
	row := q.db.QueryRow(ctx, updateAttendanceType,
		arg.AttendanceTypeID,
		arg.Name,
		arg.Key,
		arg.Color,
	)
	var i AttendanceType
	err := row.Scan(
		&i.MAttendanceTypesPkey,
		&i.AttendanceTypeID,
		&i.Name,
		&i.Key,
		&i.Color,
	)
	return i, err
}

const updateAttendanceTypeByKey = `-- name: UpdateAttendanceTypeByKey :one
UPDATE m_attendance_types SET name = $2, color = $3 WHERE key = $1 RETURNING m_attendance_types_pkey, attendance_type_id, name, key, color
`

type UpdateAttendanceTypeByKeyParams struct {
	Key   string `json:"key"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

func (q *Queries) UpdateAttendanceTypeByKey(ctx context.Context, arg UpdateAttendanceTypeByKeyParams) (AttendanceType, error) {
	row := q.db.QueryRow(ctx, updateAttendanceTypeByKey, arg.Key, arg.Name, arg.Color)
	var i AttendanceType
	err := row.Scan(
		&i.MAttendanceTypesPkey,
		&i.AttendanceTypeID,
		&i.Name,
		&i.Key,
		&i.Color,
	)
	return i, err
}
