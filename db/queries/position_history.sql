-- name: CreatePositionHistories :copyfrom
INSERT INTO t_position_histories (member_id, x_pos, y_pos, sent_at) VALUES ($1, $2, $3, $4);

-- name: CreatePositionHistory :one
INSERT INTO t_position_histories (member_id, x_pos, y_pos, sent_at) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: UpdatePositionHistory :one
UPDATE t_position_histories SET member_id = $2, x_pos = $3, y_pos = $4, sent_at = $5 WHERE position_history_id = $1 RETURNING *;

-- name: DeletePositionHistory :exec
DELETE FROM t_position_histories WHERE position_history_id = $1;

-- name: FindPositionHistoryByID :one
SELECT * FROM t_position_histories WHERE position_history_id = $1;

-- name: FindPositionHistoryByIDWithMember :one
SELECT sqlc.embed(t_position_histories), sqlc.embed(m_members) FROM t_position_histories
LEFT JOIN m_members ON t_position_histories.member_id = m_members.member_id
WHERE position_history_id = $1;

-- name: GetPositionHistories :many
SELECT * FROM t_position_histories
WHERE
	CASE WHEN @where_in_member::boolean = true THEN member_id = ANY(@in_member_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @where_earlier_sent_at::boolean = true THEN sent_at >= @earlier_sent_at ELSE TRUE END
AND
	CASE WHEN @where_later_sent_at::boolean = true THEN sent_at <= @later_sent_at ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'old_send' THEN sent_at END ASC,
	CASE WHEN @order_method::text = 'late_send' THEN sent_at END DESC,
	t_position_histories_pkey DESC;

-- name: GetPositionHistoriesUseNumberedPaginate :many
SELECT * FROM t_position_histories
WHERE
	CASE WHEN @where_in_member::boolean = true THEN member_id = ANY(@in_member_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @where_earlier_sent_at::boolean = true THEN sent_at >= @earlier_sent_at ELSE TRUE END
AND
	CASE WHEN @where_later_sent_at::boolean = true THEN sent_at <= @later_sent_at ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'old_send' THEN sent_at END ASC,
	CASE WHEN @order_method::text = 'late_send' THEN sent_at END DESC,
	t_position_histories_pkey DESC
LIMIT $1 OFFSET $2;

-- name: GetPositionHistoriesUseKeysetPaginate :many
SELECT * FROM t_position_histories
WHERE
	CASE WHEN @where_in_member::boolean = true THEN member_id = ANY(@in_member_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @where_earlier_sent_at::boolean = true THEN sent_at >= @earlier_sent_at ELSE TRUE END
AND
	CASE WHEN @where_later_sent_at::boolean = true THEN sent_at <= @later_sent_at ELSE TRUE END
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'old_send' THEN sent_at > @send_cursor OR (sent_at = @send_cursor AND t_position_histories_pkey < @cursor::int)
				WHEN 'late_send' THEN sent_at < @send_cursor OR (sent_at = @send_cursor AND t_position_histories_pkey < @cursor::int)
				ELSE t_position_histories_pkey < @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'old_send' THEN sent_at < @send_cursor OR (sent_at = @send_cursor AND t_position_histories_pkey > @cursor::int)
				WHEN 'late_send' THEN sent_at > @send_cursor OR (sent_at = @send_cursor AND t_position_histories_pkey > @cursor::int)
				ELSE t_position_histories_pkey > @cursor::int
		END
	END
ORDER BY
	CASE WHEN @order_method::text = 'old_send' THEN sent_at END ASC,
	CASE WHEN @order_method::text = 'late_send' THEN sent_at END DESC,
	t_position_histories_pkey DESC
LIMIT $1;

-- name: GetPluralPositionHistories :many
SELECT * FROM t_position_histories WHERE position_history_id = ANY(@position_history_ids::uuid[])
ORDER BY
	t_position_histories_pkey DESC
LIMIT $1 OFFSET $2;

-- name: GetPositionHistoriesWithMember :many
SELECT sqlc.embed(t_position_histories), sqlc.embed(m_members) FROM t_position_histories
LEFT JOIN m_members ON t_position_histories.member_id = m_members.member_id
WHERE
	CASE WHEN @where_in_member::boolean = true THEN member_id = ANY(@in_member_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @where_earlier_sent_at::boolean = true THEN sent_at >= @earlier_sent_at ELSE TRUE END
AND
	CASE WHEN @where_later_sent_at::boolean = true THEN sent_at <= @later_sent_at ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'old_send' THEN sent_at END ASC,
	CASE WHEN @order_method::text = 'late_send' THEN sent_at END DESC,
	t_position_histories_pkey DESC;

-- name: GetPositionHistoriesWithMemberUseNumberedPaginate :many
SELECT sqlc.embed(t_position_histories), sqlc.embed(m_members) FROM t_position_histories
LEFT JOIN m_members ON t_position_histories.member_id = m_members.member_id
WHERE
	CASE WHEN @where_in_member::boolean = true THEN member_id = ANY(@in_member_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @where_earlier_sent_at::boolean = true THEN sent_at >= @earlier_sent_at ELSE TRUE END
AND
	CASE WHEN @where_later_sent_at::boolean = true THEN sent_at <= @later_sent_at ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'old_send' THEN sent_at END ASC,
	CASE WHEN @order_method::text = 'late_send' THEN sent_at END DESC,
	t_position_histories_pkey DESC
LIMIT $1 OFFSET $2;

-- name: GetPositionHistoriesWithMemberUseKeysetPaginate :many
SELECT sqlc.embed(t_position_histories), sqlc.embed(m_members) FROM t_position_histories
LEFT JOIN m_members ON t_position_histories.member_id = m_members.member_id
WHERE
	CASE WHEN @where_in_member::boolean = true THEN member_id = ANY(@in_member_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @where_earlier_sent_at::boolean = true THEN sent_at >= @earlier_sent_at ELSE TRUE END
AND
	CASE WHEN @where_later_sent_at::boolean = true THEN sent_at <= @later_sent_at ELSE TRUE END
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'old_send' THEN sent_at > @send_cursor OR (sent_at = @send_cursor AND t_position_histories_pkey < @cursor::int)
				WHEN 'late_send' THEN sent_at < @send_cursor OR (sent_at = @send_cursor AND t_position_histories_pkey < @cursor::int)
				ELSE t_position_histories_pkey < @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'old_send' THEN sent_at < @send_cursor OR (sent_at = @send_cursor AND t_position_histories_pkey > @cursor::int)
				WHEN 'late_send' THEN sent_at > @send_cursor OR (sent_at = @send_cursor AND t_position_histories_pkey > @cursor::int)
				ELSE t_position_histories_pkey > @cursor::int
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'old_send' THEN sent_at END ASC,
	CASE WHEN @order_method::text = 'late_send' THEN sent_at END DESC,
	t_position_histories_pkey DESC
LIMIT $1;

-- name: GetPluralPositionHistoriesWithMember :many
SELECT sqlc.embed(t_position_histories), sqlc.embed(m_members) FROM t_position_histories
LEFT JOIN m_members ON t_position_histories.member_id = m_members.member_id
WHERE position_history_id = ANY(@position_history_ids::uuid[])
ORDER BY
	t_position_histories_pkey DESC
LIMIT $1 OFFSET $2;

-- name: CountPositionHistories :one
SELECT COUNT(*) FROM t_position_histories
WHERE
	CASE WHEN @where_in_member::boolean = true THEN member_id = ANY(@in_member_ids::uuid[]) ELSE TRUE END
AND
	CASE WHEN @where_earlier_sent_at::boolean = true THEN sent_at >= @earlier_sent_at ELSE TRUE END
AND
	CASE WHEN @where_later_sent_at::boolean = true THEN sent_at <= @later_sent_at ELSE TRUE END;
