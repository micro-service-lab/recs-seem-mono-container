-- name: CreateLabIOHistories :copyfrom
INSERT INTO t_lab_io_histories (member_id, entered_at, exited_at) VALUES ($1, $2, $3);

-- name: CreateLabIOHistory :one
INSERT INTO t_lab_io_histories (member_id, entered_at, exited_at) VALUES ($1, $2, $3) RETURNING *;

-- name: UpdateLabIOHistory :one
UPDATE t_lab_io_histories SET member_id = $2, entered_at = $3, exited_at = $4 WHERE lab_io_history_id = $1 RETURNING *;

-- name: ExitLabIOHistory :one
UPDATE t_lab_io_histories SET exited_at = $2 WHERE lab_io_history_id = $1 RETURNING *;

-- name: DeleteLabIOHistory :exec
DELETE FROM t_lab_io_histories WHERE lab_io_history_id = $1;

-- name: FindLabIOHistoryByID :one
SELECT * FROM t_lab_io_histories WHERE lab_io_history_id = $1;

-- name: FindLabIOHistoryWithMember :one
SELECT sqlc.embed(t_lab_io_histories), sqlc.embed(m_members) FROM t_lab_io_histories
LEFT JOIN m_members ON t_lab_io_histories.member_id = m_members.member_id
WHERE lab_io_history_id = $1;

-- name: GetLabIOHistories :many
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
	t_lab_io_histories_pkey ASC;

-- name: GetLabIOHistoriesUseNumberedPaginate :many
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
	t_lab_io_histories_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetLabIOHistoriesUseKeysetPaginate :many
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
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'old_enter' THEN entered_at > @enter_cursor OR (entered_at = @enter_cursor AND t_lab_io_histories_pkey > @cursor::int)
				WHEN 'late_enter' THEN entered_at < @enter_cursor OR (entered_at = @enter_cursor AND t_lab_io_histories_pkey > @cursor::int)
				WHEN 'old_exit' THEN exited_at > @exit_cursor OR (exited_at = @exit_cursor AND t_lab_io_histories_pkey > @cursor::int)
				WHEN 'late_exit' THEN exited_at < @exit_cursor OR (exited_at = @exit_cursor AND t_lab_io_histories_pkey > @cursor::int)
				ELSE t_lab_io_histories_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'old_enter' THEN entered_at < @enter_cursor OR (entered_at = @enter_cursor AND t_lab_io_histories_pkey < @cursor::int)
				WHEN 'late_enter' THEN entered_at > @enter_cursor OR (entered_at = @enter_cursor AND t_lab_io_histories_pkey < @cursor::int)
				WHEN 'old_exit' THEN exited_at < @exit_cursor OR (exited_at = @exit_cursor AND t_lab_io_histories_pkey < @cursor::int)
				WHEN 'late_exit' THEN exited_at > @exit_cursor OR (exited_at = @exit_cursor AND t_lab_io_histories_pkey < @cursor::int)
				ELSE t_lab_io_histories_pkey < @cursor::int
		END
	END
ORDER BY
	CASE WHEN @order_method::text = 'old_enter' THEN entered_at END ASC,
	CASE WHEN @order_method::text = 'late_enter' THEN entered_at END DESC,
	CASE WHEN @order_method::text = 'old_exit' THEN exited_at END ASC,
	CASE WHEN @order_method::text = 'late_exit' THEN exited_at END DESC,
	t_lab_io_histories_pkey ASC
LIMIT $1;

-- name: GetPluralLabIOHistories :many
SELECT * FROM t_lab_io_histories WHERE lab_io_history_id = ANY(@lab_io_history_ids::uuid[])
ORDER BY
	t_lab_io_histories_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetLabIOHistoriesWithMember :many
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
	t_lab_io_histories_pkey ASC;

-- name: GetLabIOHistoriesWithMemberUseNumberedPaginate :many
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
	t_lab_io_histories_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetLabIOHistoriesWithMemberUseKeysetPaginate :many
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
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'old_enter' THEN entered_at > @enter_cursor OR (entered_at = @enter_cursor AND t_lab_io_histories_pkey > @cursor::int)
				WHEN 'late_enter' THEN entered_at < @enter_cursor OR (entered_at = @enter_cursor AND t_lab_io_histories_pkey > @cursor::int)
				WHEN 'old_exit' THEN exited_at > @exit_cursor OR (exited_at = @exit_cursor AND t_lab_io_histories_pkey > @cursor::int)
				WHEN 'late_exit' THEN exited_at < @exit_cursor OR (exited_at = @exit_cursor AND t_lab_io_histories_pkey > @cursor::int)
				ELSE t_lab_io_histories_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'old_enter' THEN entered_at < @enter_cursor OR (entered_at = @enter_cursor AND t_lab_io_histories_pkey < @cursor::int)
				WHEN 'late_enter' THEN entered_at > @enter_cursor OR (entered_at = @enter_cursor AND t_lab_io_histories_pkey < @cursor::int)
				WHEN 'old_exit' THEN exited_at < @exit_cursor OR (exited_at = @exit_cursor AND t_lab_io_histories_pkey < @cursor::int)
				WHEN 'late_exit' THEN exited_at > @exit_cursor OR (exited_at = @exit_cursor AND t_lab_io_histories_pkey < @cursor::int)
				ELSE t_lab_io_histories_pkey < @cursor::int
		END
	END
ORDER BY
	CASE WHEN @order_method::text = 'old_enter' THEN entered_at END ASC,
	CASE WHEN @order_method::text = 'late_enter' THEN entered_at END DESC,
	CASE WHEN @order_method::text = 'old_exit' THEN exited_at END ASC,
	CASE WHEN @order_method::text = 'late_exit' THEN exited_at END DESC,
	t_lab_io_histories_pkey ASC
LIMIT $1;

-- name: GetPluralLabIOHistoriesWithMember :many
SELECT sqlc.embed(t_lab_io_histories), sqlc.embed(m_members) FROM t_lab_io_histories
LEFT JOIN m_members ON t_lab_io_histories.member_id = m_members.member_id
WHERE lab_io_history_id = ANY(@lab_io_history_ids::uuid[])
ORDER BY
	t_lab_io_histories_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountLabIOHistories :one
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
