-- name: CreateLabIoHistories :copyfrom
INSERT INTO t_lab_io_histories (member_id, entered_at, exited_at) VALUES ($1, $2, $3);

-- name: CreateLabIoHistory :one
INSERT INTO t_lab_io_histories (member_id, entered_at, exited_at) VALUES ($1, $2, $3) RETURNING *;

-- name: UpdateLabIoHistory :one
UPDATE t_lab_io_histories SET member_id = $2, entered_at = $3, exited_at = $4 WHERE lab_io_history_id = $1 RETURNING *;

-- name: DeleteLabIoHistory :exec
DELETE FROM t_lab_io_histories WHERE lab_io_history_id = $1;

-- name: FindLabIoHistoryByID :one
SELECT * FROM t_lab_io_histories WHERE lab_io_history_id = $1;

-- name: GetLabIoHistories :many
SELECT * FROM t_lab_io_histories
WHERE
	CASE WHEN @where_member::boolean = true THEN member_id = @member_id ELSE TRUE END
AND
	CASE WHEN @where_earlier_entered_at::boolean = true THEN entered_at >= @earlier_entered_at ELSE TRUE END
AND
	CASE WHEN @where_later_entered_at::boolean = true THEN entered_at <= @later_entered_at ELSE TRUE END
AND
	CASE WHEN @where_earlier_exited_at::boolean = true THEN exited_at >= @earlier_exited_at ELSE TRUE END
AND
	CASE WHEN @where_later_exited_at::boolean = true THEN exited_at <= @later_exited_at ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'entered_at' THEN entered_at END ASC,
	CASE WHEN @order_method::text = 'exited_at' THEN exited_at END ASC,
	t_lab_io_histories_pkey DESC
LIMIT $1 OFFSET $2;

-- name: CountLabIoHistories :one
SELECT COUNT(*) FROM t_lab_io_histories
WHERE
	CASE WHEN @where_member::boolean = true THEN member_id = @member_id ELSE TRUE END
AND
	CASE WHEN @where_earlier_entered_at::boolean = true THEN entered_at >= @earlier_entered_at ELSE TRUE END
AND
	CASE WHEN @where_later_entered_at::boolean = true THEN entered_at <= @later_entered_at ELSE TRUE END
AND
	CASE WHEN @where_earlier_exited_at::boolean = true THEN exited_at >= @earlier_exited_at ELSE TRUE END
AND
	CASE WHEN @where_later_exited_at::boolean = true THEN exited_at <= @later_exited_at ELSE TRUE END;
