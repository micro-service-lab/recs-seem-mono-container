-- name: CreateEarlyLeavings :copyfrom
INSERT INTO t_early_leavings (attendance_id, leave_time) VALUES ($1, $2);

-- name: CreateEarlyLeaving :one
INSERT INTO t_early_leavings (attendance_id, leave_time) VALUES ($1, $2) RETURNING *;

-- name: DeleteEarlyLeaving :exec
DELETE FROM t_early_leavings WHERE early_leaving_id = $1;

-- name: FindEarlyLeavingByID :one
SELECT * FROM t_early_leavings WHERE early_leaving_id = $1;

-- name: GetEarlyLeavings :many
SELECT * FROM t_early_leavings
ORDER BY
	t_early_leavings_pkey DESC;

-- name: GetEarlyLeavingsUseNumberedPaginate :many
SELECT * FROM t_early_leavings
ORDER BY
	t_early_leavings_pkey DESC
LIMIT $1 OFFSET $2;

-- name: GetEarlyLeavingsUseKeysetPaginate :many
SELECT * FROM t_early_leavings
WHERE
	CASE @cursor_direction
		WHEN 'next' THEN
			t_early_leavings_pkey < @cursor
		WHEN 'prev' THEN
			t_early_leavings_pkey > @cursor
	END
ORDER BY
	t_early_leavings_pkey DESC
LIMIT $1;

-- name: CountEarlyLeavings :one
SELECT COUNT(*) FROM t_early_leavings;
