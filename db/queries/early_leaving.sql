-- name: CreateEarlyLeavings :copyfrom
INSERT INTO t_early_leavings (attendance_id, leave_time) VALUES ($1, $2);

-- name: CreateEarlyLeaving :one
INSERT INTO t_early_leavings (attendance_id, leave_time) VALUES ($1, $2) RETURNING *;

-- name: DeleteEarlyLeaving :execrows
DELETE FROM t_early_leavings WHERE early_leaving_id = $1;

-- name: PluralDeleteEarlyLeavings :execrows
DELETE FROM t_early_leavings WHERE early_leaving_id = ANY(@early_leaving_ids::uuid[]);

-- name: FindEarlyLeavingByID :one
SELECT * FROM t_early_leavings WHERE early_leaving_id = $1;

-- name: GetEarlyLeavings :many
SELECT * FROM t_early_leavings
ORDER BY
	t_early_leavings_pkey ASC;

-- name: GetEarlyLeavingsUseNumberedPaginate :many
SELECT * FROM t_early_leavings
ORDER BY
	t_early_leavings_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetEarlyLeavingsUseKeysetPaginate :many
SELECT * FROM t_early_leavings
WHERE
	CASE @cursor_direction::text
		WHEN 'next' THEN
			t_early_leavings_pkey > @cursor::int
		WHEN 'prev' THEN
			t_early_leavings_pkey < @cursor::int
	END
ORDER BY
	CASE WHEN @cursor_direction::text = 'next' THEN t_early_leavings_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN t_early_leavings_pkey END DESC
LIMIT $1;

-- name: GetPluralEarlyLeavings :many
SELECT * FROM t_early_leavings
WHERE attendance_id = ANY(@attendance_ids::uuid[])
ORDER BY
	t_early_leavings_pkey ASC;

-- name: GetPluralEarlyLeavingsUseNumberedPaginate :many
SELECT * FROM t_early_leavings
WHERE attendance_id = ANY(@attendance_ids::uuid[])
ORDER BY
	t_early_leavings_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountEarlyLeavings :one
SELECT COUNT(*) FROM t_early_leavings;
