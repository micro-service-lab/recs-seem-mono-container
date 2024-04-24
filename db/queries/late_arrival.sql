-- name: CreateLateArrivals :copyfrom
INSERT INTO t_late_arrivals (attendance_id, arrive_time) VALUES ($1, $2);

-- name: CreateLateArrival :one
INSERT INTO t_late_arrivals (attendance_id, arrive_time) VALUES ($1, $2) RETURNING *;

-- name: DeleteLateArrival :exec
DELETE FROM t_late_arrivals WHERE late_arrival_id = $1;

-- name: FindLateArrivalByID :one
SELECT * FROM t_late_arrivals WHERE late_arrival_id = $1;

-- name: GetLateArrivals :many
SELECT * FROM t_late_arrivals
ORDER BY
	t_late_arrivals_pkey DESC
LIMIT $1 OFFSET $2;

-- name: CountLateArrivals :one
SELECT COUNT(*) FROM t_late_arrivals;
