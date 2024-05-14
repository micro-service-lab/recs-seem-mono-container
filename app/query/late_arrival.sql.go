// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: late_arrival.sql

package query

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const countLateArrivals = `-- name: CountLateArrivals :one
SELECT COUNT(*) FROM t_late_arrivals
`

func (q *Queries) CountLateArrivals(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, countLateArrivals)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createLateArrival = `-- name: CreateLateArrival :one
INSERT INTO t_late_arrivals (attendance_id, arrive_time) VALUES ($1, $2) RETURNING t_late_arrivals_pkey, late_arrival_id, attendance_id, arrive_time
`

type CreateLateArrivalParams struct {
	AttendanceID uuid.UUID `json:"attendance_id"`
	ArriveTime   time.Time `json:"arrive_time"`
}

func (q *Queries) CreateLateArrival(ctx context.Context, arg CreateLateArrivalParams) (LateArrival, error) {
	row := q.db.QueryRow(ctx, createLateArrival, arg.AttendanceID, arg.ArriveTime)
	var i LateArrival
	err := row.Scan(
		&i.TLateArrivalsPkey,
		&i.LateArrivalID,
		&i.AttendanceID,
		&i.ArriveTime,
	)
	return i, err
}

type CreateLateArrivalsParams struct {
	AttendanceID uuid.UUID `json:"attendance_id"`
	ArriveTime   time.Time `json:"arrive_time"`
}

const deleteLateArrival = `-- name: DeleteLateArrival :execrows
DELETE FROM t_late_arrivals WHERE late_arrival_id = $1
`

func (q *Queries) DeleteLateArrival(ctx context.Context, lateArrivalID uuid.UUID) (int64, error) {
	result, err := q.db.Exec(ctx, deleteLateArrival, lateArrivalID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const findLateArrivalByID = `-- name: FindLateArrivalByID :one
SELECT t_late_arrivals_pkey, late_arrival_id, attendance_id, arrive_time FROM t_late_arrivals WHERE late_arrival_id = $1
`

func (q *Queries) FindLateArrivalByID(ctx context.Context, lateArrivalID uuid.UUID) (LateArrival, error) {
	row := q.db.QueryRow(ctx, findLateArrivalByID, lateArrivalID)
	var i LateArrival
	err := row.Scan(
		&i.TLateArrivalsPkey,
		&i.LateArrivalID,
		&i.AttendanceID,
		&i.ArriveTime,
	)
	return i, err
}

const getLateArrivals = `-- name: GetLateArrivals :many
SELECT t_late_arrivals_pkey, late_arrival_id, attendance_id, arrive_time FROM t_late_arrivals
ORDER BY
	t_late_arrivals_pkey ASC
`

func (q *Queries) GetLateArrivals(ctx context.Context) ([]LateArrival, error) {
	rows, err := q.db.Query(ctx, getLateArrivals)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []LateArrival{}
	for rows.Next() {
		var i LateArrival
		if err := rows.Scan(
			&i.TLateArrivalsPkey,
			&i.LateArrivalID,
			&i.AttendanceID,
			&i.ArriveTime,
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

const getLateArrivalsUseKeysetPaginate = `-- name: GetLateArrivalsUseKeysetPaginate :many
SELECT t_late_arrivals_pkey, late_arrival_id, attendance_id, arrive_time FROM t_late_arrivals
WHERE
	CASE $2::text
		WHEN 'next' THEN
			t_late_arrivals_pkey > $3::int
		WHEN 'prev' THEN
			t_late_arrivals_pkey < $3::int
	END
ORDER BY
	CASE WHEN $2::text = 'next' THEN t_late_arrivals_pkey END ASC,
	CASE WHEN $2::text = 'prev' THEN t_late_arrivals_pkey END DESC
LIMIT $1
`

type GetLateArrivalsUseKeysetPaginateParams struct {
	Limit           int32  `json:"limit"`
	CursorDirection string `json:"cursor_direction"`
	Cursor          int32  `json:"cursor"`
}

func (q *Queries) GetLateArrivalsUseKeysetPaginate(ctx context.Context, arg GetLateArrivalsUseKeysetPaginateParams) ([]LateArrival, error) {
	rows, err := q.db.Query(ctx, getLateArrivalsUseKeysetPaginate, arg.Limit, arg.CursorDirection, arg.Cursor)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []LateArrival{}
	for rows.Next() {
		var i LateArrival
		if err := rows.Scan(
			&i.TLateArrivalsPkey,
			&i.LateArrivalID,
			&i.AttendanceID,
			&i.ArriveTime,
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

const getLateArrivalsUseNumberedPaginate = `-- name: GetLateArrivalsUseNumberedPaginate :many
SELECT t_late_arrivals_pkey, late_arrival_id, attendance_id, arrive_time FROM t_late_arrivals
ORDER BY
	t_late_arrivals_pkey ASC
LIMIT $1 OFFSET $2
`

type GetLateArrivalsUseNumberedPaginateParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) GetLateArrivalsUseNumberedPaginate(ctx context.Context, arg GetLateArrivalsUseNumberedPaginateParams) ([]LateArrival, error) {
	rows, err := q.db.Query(ctx, getLateArrivalsUseNumberedPaginate, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []LateArrival{}
	for rows.Next() {
		var i LateArrival
		if err := rows.Scan(
			&i.TLateArrivalsPkey,
			&i.LateArrivalID,
			&i.AttendanceID,
			&i.ArriveTime,
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

const getPluralLateArrivals = `-- name: GetPluralLateArrivals :many
SELECT t_late_arrivals_pkey, late_arrival_id, attendance_id, arrive_time FROM t_late_arrivals
WHERE attendance_id = ANY($3::uuid[])
ORDER BY
	t_late_arrivals_pkey ASC
LIMIT $1 OFFSET $2
`

type GetPluralLateArrivalsParams struct {
	Limit         int32       `json:"limit"`
	Offset        int32       `json:"offset"`
	AttendanceIds []uuid.UUID `json:"attendance_ids"`
}

func (q *Queries) GetPluralLateArrivals(ctx context.Context, arg GetPluralLateArrivalsParams) ([]LateArrival, error) {
	rows, err := q.db.Query(ctx, getPluralLateArrivals, arg.Limit, arg.Offset, arg.AttendanceIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []LateArrival{}
	for rows.Next() {
		var i LateArrival
		if err := rows.Scan(
			&i.TLateArrivalsPkey,
			&i.LateArrivalID,
			&i.AttendanceID,
			&i.ArriveTime,
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

const pluralDeleteLateArrivals = `-- name: PluralDeleteLateArrivals :execrows
DELETE FROM t_late_arrivals WHERE late_arrival_id = ANY($1::uuid[])
`

func (q *Queries) PluralDeleteLateArrivals(ctx context.Context, dollar_1 []uuid.UUID) (int64, error) {
	result, err := q.db.Exec(ctx, pluralDeleteLateArrivals, dollar_1)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}
