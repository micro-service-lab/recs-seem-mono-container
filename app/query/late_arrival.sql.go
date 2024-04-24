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

const deleteLateArrival = `-- name: DeleteLateArrival :exec
DELETE FROM t_late_arrivals WHERE late_arrival_id = $1
`

func (q *Queries) DeleteLateArrival(ctx context.Context, lateArrivalID uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteLateArrival, lateArrivalID)
	return err
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
	t_late_arrivals_pkey DESC
LIMIT $1 OFFSET $2
`

type GetLateArrivalsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) GetLateArrivals(ctx context.Context, arg GetLateArrivalsParams) ([]LateArrival, error) {
	rows, err := q.db.Query(ctx, getLateArrivals, arg.Limit, arg.Offset)
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
