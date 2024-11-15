// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: early_leaving.sql

package query

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const countEarlyLeavings = `-- name: CountEarlyLeavings :one
SELECT COUNT(*) FROM t_early_leavings
`

func (q *Queries) CountEarlyLeavings(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, countEarlyLeavings)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createEarlyLeaving = `-- name: CreateEarlyLeaving :one
INSERT INTO t_early_leavings (attendance_id, leave_time) VALUES ($1, $2) RETURNING t_early_leavings_pkey, early_leaving_id, attendance_id, leave_time
`

type CreateEarlyLeavingParams struct {
	AttendanceID uuid.UUID `json:"attendance_id"`
	LeaveTime    time.Time `json:"leave_time"`
}

func (q *Queries) CreateEarlyLeaving(ctx context.Context, arg CreateEarlyLeavingParams) (EarlyLeaving, error) {
	row := q.db.QueryRow(ctx, createEarlyLeaving, arg.AttendanceID, arg.LeaveTime)
	var i EarlyLeaving
	err := row.Scan(
		&i.TEarlyLeavingsPkey,
		&i.EarlyLeavingID,
		&i.AttendanceID,
		&i.LeaveTime,
	)
	return i, err
}

type CreateEarlyLeavingsParams struct {
	AttendanceID uuid.UUID `json:"attendance_id"`
	LeaveTime    time.Time `json:"leave_time"`
}

const deleteEarlyLeaving = `-- name: DeleteEarlyLeaving :execrows
DELETE FROM t_early_leavings WHERE early_leaving_id = $1
`

func (q *Queries) DeleteEarlyLeaving(ctx context.Context, earlyLeavingID uuid.UUID) (int64, error) {
	result, err := q.db.Exec(ctx, deleteEarlyLeaving, earlyLeavingID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const findEarlyLeavingByID = `-- name: FindEarlyLeavingByID :one
SELECT t_early_leavings_pkey, early_leaving_id, attendance_id, leave_time FROM t_early_leavings WHERE early_leaving_id = $1
`

func (q *Queries) FindEarlyLeavingByID(ctx context.Context, earlyLeavingID uuid.UUID) (EarlyLeaving, error) {
	row := q.db.QueryRow(ctx, findEarlyLeavingByID, earlyLeavingID)
	var i EarlyLeaving
	err := row.Scan(
		&i.TEarlyLeavingsPkey,
		&i.EarlyLeavingID,
		&i.AttendanceID,
		&i.LeaveTime,
	)
	return i, err
}

const getEarlyLeavings = `-- name: GetEarlyLeavings :many
SELECT t_early_leavings_pkey, early_leaving_id, attendance_id, leave_time FROM t_early_leavings
ORDER BY
	t_early_leavings_pkey ASC
`

func (q *Queries) GetEarlyLeavings(ctx context.Context) ([]EarlyLeaving, error) {
	rows, err := q.db.Query(ctx, getEarlyLeavings)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []EarlyLeaving{}
	for rows.Next() {
		var i EarlyLeaving
		if err := rows.Scan(
			&i.TEarlyLeavingsPkey,
			&i.EarlyLeavingID,
			&i.AttendanceID,
			&i.LeaveTime,
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

const getEarlyLeavingsUseKeysetPaginate = `-- name: GetEarlyLeavingsUseKeysetPaginate :many
SELECT t_early_leavings_pkey, early_leaving_id, attendance_id, leave_time FROM t_early_leavings
WHERE
	CASE $2::text
		WHEN 'next' THEN
			t_early_leavings_pkey > $3::int
		WHEN 'prev' THEN
			t_early_leavings_pkey < $3::int
	END
ORDER BY
	CASE WHEN $2::text = 'next' THEN t_early_leavings_pkey END ASC,
	CASE WHEN $2::text = 'prev' THEN t_early_leavings_pkey END DESC
LIMIT $1
`

type GetEarlyLeavingsUseKeysetPaginateParams struct {
	Limit           int32  `json:"limit"`
	CursorDirection string `json:"cursor_direction"`
	Cursor          int32  `json:"cursor"`
}

func (q *Queries) GetEarlyLeavingsUseKeysetPaginate(ctx context.Context, arg GetEarlyLeavingsUseKeysetPaginateParams) ([]EarlyLeaving, error) {
	rows, err := q.db.Query(ctx, getEarlyLeavingsUseKeysetPaginate, arg.Limit, arg.CursorDirection, arg.Cursor)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []EarlyLeaving{}
	for rows.Next() {
		var i EarlyLeaving
		if err := rows.Scan(
			&i.TEarlyLeavingsPkey,
			&i.EarlyLeavingID,
			&i.AttendanceID,
			&i.LeaveTime,
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

const getEarlyLeavingsUseNumberedPaginate = `-- name: GetEarlyLeavingsUseNumberedPaginate :many
SELECT t_early_leavings_pkey, early_leaving_id, attendance_id, leave_time FROM t_early_leavings
ORDER BY
	t_early_leavings_pkey ASC
LIMIT $1 OFFSET $2
`

type GetEarlyLeavingsUseNumberedPaginateParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) GetEarlyLeavingsUseNumberedPaginate(ctx context.Context, arg GetEarlyLeavingsUseNumberedPaginateParams) ([]EarlyLeaving, error) {
	rows, err := q.db.Query(ctx, getEarlyLeavingsUseNumberedPaginate, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []EarlyLeaving{}
	for rows.Next() {
		var i EarlyLeaving
		if err := rows.Scan(
			&i.TEarlyLeavingsPkey,
			&i.EarlyLeavingID,
			&i.AttendanceID,
			&i.LeaveTime,
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

const getPluralEarlyLeavings = `-- name: GetPluralEarlyLeavings :many
SELECT t_early_leavings_pkey, early_leaving_id, attendance_id, leave_time FROM t_early_leavings
WHERE attendance_id = ANY($1::uuid[])
ORDER BY
	t_early_leavings_pkey ASC
`

func (q *Queries) GetPluralEarlyLeavings(ctx context.Context, attendanceIds []uuid.UUID) ([]EarlyLeaving, error) {
	rows, err := q.db.Query(ctx, getPluralEarlyLeavings, attendanceIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []EarlyLeaving{}
	for rows.Next() {
		var i EarlyLeaving
		if err := rows.Scan(
			&i.TEarlyLeavingsPkey,
			&i.EarlyLeavingID,
			&i.AttendanceID,
			&i.LeaveTime,
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

const getPluralEarlyLeavingsUseNumberedPaginate = `-- name: GetPluralEarlyLeavingsUseNumberedPaginate :many
SELECT t_early_leavings_pkey, early_leaving_id, attendance_id, leave_time FROM t_early_leavings
WHERE attendance_id = ANY($3::uuid[])
ORDER BY
	t_early_leavings_pkey ASC
LIMIT $1 OFFSET $2
`

type GetPluralEarlyLeavingsUseNumberedPaginateParams struct {
	Limit         int32       `json:"limit"`
	Offset        int32       `json:"offset"`
	AttendanceIds []uuid.UUID `json:"attendance_ids"`
}

func (q *Queries) GetPluralEarlyLeavingsUseNumberedPaginate(ctx context.Context, arg GetPluralEarlyLeavingsUseNumberedPaginateParams) ([]EarlyLeaving, error) {
	rows, err := q.db.Query(ctx, getPluralEarlyLeavingsUseNumberedPaginate, arg.Limit, arg.Offset, arg.AttendanceIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []EarlyLeaving{}
	for rows.Next() {
		var i EarlyLeaving
		if err := rows.Scan(
			&i.TEarlyLeavingsPkey,
			&i.EarlyLeavingID,
			&i.AttendanceID,
			&i.LeaveTime,
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

const pluralDeleteEarlyLeavings = `-- name: PluralDeleteEarlyLeavings :execrows
DELETE FROM t_early_leavings WHERE early_leaving_id = ANY($1::uuid[])
`

func (q *Queries) PluralDeleteEarlyLeavings(ctx context.Context, earlyLeavingIds []uuid.UUID) (int64, error) {
	result, err := q.db.Exec(ctx, pluralDeleteEarlyLeavings, earlyLeavingIds)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}
