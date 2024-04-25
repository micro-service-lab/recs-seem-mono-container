-- name: CreateLabIoHistories :copyfrom
INSERT INTO t_lab_io_histories (member_id, entered_at, exited_at) VALUES ($1, $2, $3);

-- name: CreateLabIoHistory :one
INSERT INTO t_lab_io_histories (member_id, entered_at, exited_at) VALUES ($1, $2, $3) RETURNING *;

-- name: UpdateLabIoHistory :one
UPDATE t_lab_io_histories SET member_id = $2, entered_at = $3, exited_at = $4 WHERE lab_io_history_id = $1 RETURNING *;

-- name: ExitLabIoHistory :one
UPDATE t_lab_io_histories SET exited_at = $2 WHERE lab_io_history_id = $1 RETURNING *;

-- name: DeleteLabIoHistory :exec
DELETE FROM t_lab_io_histories WHERE lab_io_history_id = $1;

-- name: FindLabIoHistoryByID :one
SELECT * FROM t_lab_io_histories WHERE lab_io_history_id = $1;

-- name: FindLabIoHistoryWithMember :one
SELECT sqlc.embed(t_lab_io_histories), sqlc.embed(m_members) FROM t_lab_io_histories
LEFT JOIN m_members ON t_lab_io_histories.member_id = m_members.member_id
WHERE lab_io_history_id = $1;

-- name: GetLabIoHistories :many
SELECT * FROM t_lab_io_histories
WHERE
	CASE WHEN @where_in_member::boolean = true THEN member_id = ANY(@in_member) ELSE TRUE END
AND
	CASE WHEN @where_earlier_entered_at::boolean = true THEN entered_at >= @earlier_entered_at ELSE TRUE END
AND
	CASE WHEN @where_later_entered_at::boolean = true THEN entered_at <= @later_entered_at ELSE TRUE END
AND
	CASE WHEN @where_earlier_exited_at::boolean = true THEN exited_at >= @earlier_exited_at ELSE TRUE END
AND
	CASE WHEN @where_later_exited_at::boolean = true THEN exited_at <= @later_exited_at ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'old_enter' THEN entered_at END ASC,
	CASE WHEN @order_method::text = 'late_enter' THEN entered_at END DESC,
	CASE WHEN @order_method::text = 'old_exit' THEN exited_at END ASC,
	CASE WHEN @order_method::text = 'late_exit' THEN exited_at END DESC,
	t_lab_io_histories_pkey DESC;

-- name: GetLabIoHistoriesUseNumberedPaginate :many
SELECT * FROM t_lab_io_histories
WHERE
	CASE WHEN @where_in_member::boolean = true THEN member_id = ANY(@in_member) ELSE TRUE END
AND
	CASE WHEN @where_earlier_entered_at::boolean = true THEN entered_at >= @earlier_entered_at ELSE TRUE END
AND
	CASE WHEN @where_later_entered_at::boolean = true THEN entered_at <= @later_entered_at ELSE TRUE END
AND
	CASE WHEN @where_earlier_exited_at::boolean = true THEN exited_at >= @earlier_exited_at ELSE TRUE END
AND
	CASE WHEN @where_later_exited_at::boolean = true THEN exited_at <= @later_exited_at ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'old_enter' THEN entered_at END ASC,
	CASE WHEN @order_method::text = 'late_enter' THEN entered_at END DESC,
	CASE WHEN @order_method::text = 'old_exit' THEN exited_at END ASC,
	CASE WHEN @order_method::text = 'late_exit' THEN exited_at END DESC,
	t_lab_io_histories_pkey DESC
LIMIT $1 OFFSET $2;

-- name: GetLabIoHistoriesUseKeysetPaginate :many
SELECT * FROM t_lab_io_histories
WHERE
	CASE WHEN @where_in_member::boolean = true THEN member_id = ANY(@in_member) ELSE TRUE END
AND
	CASE WHEN @where_earlier_entered_at::boolean = true THEN entered_at >= @earlier_entered_at ELSE TRUE END
AND
	CASE WHEN @where_later_entered_at::boolean = true THEN entered_at <= @later_entered_at ELSE TRUE END
AND
	CASE WHEN @where_earlier_exited_at::boolean = true THEN exited_at >= @earlier_exited_at ELSE TRUE END
AND
	CASE WHEN @where_later_exited_at::boolean = true THEN exited_at <= @later_exited_at ELSE TRUE END
WHERE
	CASE @cursor_direction
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'old_enter' THEN entered_at > @cursor_column OR (entered_at = @cursor_column AND t_lab_io_histories_pkey < @cursor)
				WHEN 'late_enter' THEN entered_at < @cursor_column OR (entered_at = @cursor_column AND t_lab_io_histories_pkey < @cursor)
				WHEN 'old_exit' THEN exited_at > @cursor_column OR (exited_at = @cursor_column AND t_lab_io_histories_pkey < @cursor)
				WHEN 'late_exit' THEN exited_at < @cursor_column OR (exited_at = @cursor_column AND t_lab_io_histories_pkey < @cursor)
				ELSE t_lab_io_histories_pkey < @cursor
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'old_enter' THEN entered_at < @cursor_column OR (entered_at = @cursor_column AND t_lab_io_histories_pkey > @cursor)
				WHEN 'late_enter' THEN entered_at > @cursor_column OR (entered_at = @cursor_column AND t_lab_io_histories_pkey > @cursor)
				WHEN 'old_exit' THEN exited_at < @cursor_column OR (exited_at = @cursor_column AND t_lab_io_histories_pkey > @cursor)
				WHEN 'late_exit' THEN exited_at > @cursor_column OR (exited_at = @cursor_column AND t_lab_io_histories_pkey > @cursor)
				ELSE t_lab_io_histories_pkey > @cursor
	END
ORDER BY
	CASE WHEN @order_method::text = 'old_enter' THEN entered_at END ASC,
	CASE WHEN @order_method::text = 'late_enter' THEN entered_at END DESC,
	CASE WHEN @order_method::text = 'old_exit' THEN exited_at END ASC,
	CASE WHEN @order_method::text = 'late_exit' THEN exited_at END DESC,
	t_lab_io_histories_pkey DESC
LIMIT $1;

-- name: GetLabIoHistoriesWithMember :many
SELECT sqlc.embed(t_lab_io_histories), sqlc.embed(m_members) FROM t_lab_io_histories
LEFT JOIN m_members ON t_lab_io_histories.member_id = m_members.member_id
WHERE
	CASE WHEN @where_in_member::boolean = true THEN t_lab_io_histories.member_id = ANY(@in_member) ELSE TRUE END
AND
	CASE WHEN @where_earlier_entered_at::boolean = true THEN entered_at >= @earlier_entered_at ELSE TRUE END
AND
	CASE WHEN @where_later_entered_at::boolean = true THEN entered_at <= @later_entered_at ELSE TRUE END
AND
	CASE WHEN @where_earlier_exited_at::boolean = true THEN exited_at >= @earlier_exited_at ELSE TRUE END
AND
	CASE WHEN @where_later_exited_at::boolean = true THEN exited_at <= @later_exited_at ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'old_enter' THEN entered_at END ASC,
	CASE WHEN @order_method::text = 'late_enter' THEN entered_at END DESC,
	CASE WHEN @order_method::text = 'old_exit' THEN exited_at END ASC,
	CASE WHEN @order_method::text = 'late_exit' THEN exited_at END DESC,
	t_lab_io_histories_pkey DESC;

-- name: GetLabIoHistoriesWithMemberUseNumberedPaginate :many
SELECT sqlc.embed(t_lab_io_histories), sqlc.embed(m_members) FROM t_lab_io_histories
LEFT JOIN m_members ON t_lab_io_histories.member_id = m_members.member_id
WHERE
	CASE WHEN @where_in_member::boolean = true THEN t_lab_io_histories.member_id = ANY(@in_member) ELSE TRUE END
AND
	CASE WHEN @where_earlier_entered_at::boolean = true THEN entered_at >= @earlier_entered_at ELSE TRUE END
AND
	CASE WHEN @where_later_entered_at::boolean = true THEN entered_at <= @later_entered_at ELSE TRUE END
AND
	CASE WHEN @where_earlier_exited_at::boolean = true THEN exited_at >= @earlier_exited_at ELSE TRUE END
AND
	CASE WHEN @where_later_exited_at::boolean = true THEN exited_at <= @later_exited_at ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'old_enter' THEN entered_at END ASC,
	CASE WHEN @order_method::text = 'late_enter' THEN entered_at END DESC,
	CASE WHEN @order_method::text = 'old_exit' THEN exited_at END ASC,
	CASE WHEN @order_method::text = 'late_exit' THEN exited_at END DESC,
	t_lab_io_histories_pkey DESC
LIMIT $1 OFFSET $2;

-- name: GetLabIoHistoriesWithMemberUseKeysetPaginate :many
SELECT sqlc.embed(t_lab_io_histories), sqlc.embed(m_members) FROM t_lab_io_histories
LEFT JOIN m_members ON t_lab_io_histories.member_id = m_members.member_id
WHERE
	CASE WHEN @where_in_member::boolean = true THEN t_lab_io_histories.member_id = ANY(@in_member) ELSE TRUE END
AND
	CASE WHEN @where_earlier_entered_at::boolean = true THEN entered_at >= @earlier_entered_at ELSE TRUE END
AND
	CASE WHEN @where_later_entered_at::boolean = true THEN entered_at <= @later_entered_at ELSE TRUE END
AND
	CASE WHEN @where_earlier_exited_at::boolean = true THEN exited_at >= @earlier_exited_at ELSE TRUE END
AND
	CASE WHEN @where_later_exited_at::boolean = true THEN exited_at <= @later_exited_at ELSE TRUE END
WHERE
	CASE @cursor_direction
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'old_enter' THEN entered_at > @cursor_column OR (entered_at = @cursor_column AND t_lab_io_histories_pkey < @cursor)
				WHEN 'late_enter' THEN entered_at < @cursor_column OR (entered_at = @cursor_column AND t_lab_io_histories_pkey < @cursor)
				WHEN 'old_exit' THEN exited_at > @cursor_column OR (exited_at = @cursor_column AND t_lab_io_histories_pkey < @cursor)
				WHEN 'late_exit' THEN exited_at < @cursor_column OR (exited_at = @cursor_column AND t_lab_io_histories_pkey < @cursor)
				ELSE t_lab_io_histories_pkey < @cursor
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'old_enter' THEN entered_at < @cursor_column OR (entered_at = @cursor_column AND t_lab_io_histories_pkey > @cursor)
				WHEN 'late_enter' THEN entered_at > @cursor_column OR (entered_at = @cursor_column AND t_lab_io_histories_pkey > @cursor)
				WHEN 'old_exit' THEN exited_at < @cursor_column OR (exited_at = @cursor_column AND t_lab_io_histories_pkey > @cursor)
				WHEN 'late_exit' THEN exited_at > @cursor_column OR (exited_at = @cursor_column AND t_lab_io_histories_pkey > @cursor)
				ELSE t_lab_io_histories_pkey > @cursor
	END
ORDER BY
	CASE WHEN @order_method::text = 'old_enter' THEN entered_at END ASC,
	CASE WHEN @order_method::text = 'late_enter' THEN entered_at END DESC,
	CASE WHEN @order_method::text = 'old_exit' THEN exited_at END ASC,
	CASE WHEN @order_method::text = 'late_exit' THEN exited_at END DESC,
	t_lab_io_histories_pkey DESC
LIMIT $1;

-- name: CountLabIoHistories :one
SELECT COUNT(*) FROM t_lab_io_histories
where
	CASE WHEN @where_in_member::boolean = true THEN member_id = ANY(@in_member) ELSE TRUE END
AND
	CASE WHEN @where_earlier_entered_at::boolean = true THEN entered_at >= @earlier_entered_at ELSE TRUE END
AND
	CASE WHEN @where_later_entered_at::boolean = true THEN entered_at <= @later_entered_at ELSE TRUE END
AND
	CASE WHEN @where_earlier_exited_at::boolean = true THEN exited_at >= @earlier_exited_at ELSE TRUE END
AND
	CASE WHEN @where_later_exited_at::boolean = true THEN exited_at <= @later_exited_at ELSE TRUE END;
