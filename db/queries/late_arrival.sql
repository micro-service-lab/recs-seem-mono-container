-- name: CreateLateArrivals :copyfrom
INSERT INTO t_late_arrivals (attendance_id, arrive_time) VALUES ($1, $2);

-- name: CreateLateArrival :one
INSERT INTO t_late_arrivals (attendance_id, arrive_time) VALUES ($1, $2) RETURNING *;

-- name: DeleteLateArrival :execrows
DELETE FROM t_late_arrivals WHERE late_arrival_id = $1;

-- name: PluralDeleteLateArrivals :execrows
DELETE FROM t_late_arrivals WHERE late_arrival_id = ANY($1::uuid[]);

-- name: FindLateArrivalByID :one
SELECT * FROM t_late_arrivals WHERE late_arrival_id = $1;

-- name: GetLateArrivals :many
SELECT * FROM t_late_arrivals
ORDER BY
	t_late_arrivals_pkey ASC;

-- name: GetLateArrivalsUseNumberedPaginate :many
SELECT * FROM t_late_arrivals
ORDER BY
	t_late_arrivals_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetLateArrivalsUseKeysetPaginate :many
SELECT * FROM t_late_arrivals
WHERE
	CASE @cursor_direction::text
		WHEN 'next' THEN
			t_late_arrivals_pkey > @cursor::int
		WHEN 'prev' THEN
			t_late_arrivals_pkey < @cursor::int
	END
ORDER BY
	CASE WHEN @cursor_direction::text = 'next' THEN t_late_arrivals_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN t_late_arrivals_pkey END DESC
LIMIT $1;

-- name: GetPluralLateArrivals :many
SELECT * FROM t_late_arrivals
WHERE attendance_id = ANY(@attendance_ids::uuid[])
ORDER BY
	t_late_arrivals_pkey ASC;

-- name: GetPluralLateArrivalsUseNumberedPaginate :many
SELECT * FROM t_late_arrivals
WHERE attendance_id = ANY(@attendance_ids::uuid[])
ORDER BY
	t_late_arrivals_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountLateArrivals :one
SELECT COUNT(*) FROM t_late_arrivals;
