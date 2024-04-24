-- name: CreatePositionHistories :copyfrom
INSERT INTO t_position_histories (member_id, x_pos, y_pos, send_at) VALUES ($1, $2, $3, $4);

-- name: CreatePositionHistory :one
INSERT INTO t_position_histories (member_id, x_pos, y_pos, send_at) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: UpdatePositionHistory :one
UPDATE t_position_histories SET member_id = $2, x_pos = $3, y_pos = $4, send_at = $5 WHERE position_history_id = $1 RETURNING *;

-- name: DeletePositionHistory :exec
DELETE FROM t_position_histories WHERE position_history_id = $1;

-- name: FindPositionHistoryByID :one
SELECT * FROM t_position_histories WHERE position_history_id = $1;

-- name: FindPositionHistoryByIDWithMember :one
SELECT sqlc.embed(t_position_histories), sqlc.embed(m_members) FROM t_position_histories
INNER JOIN m_members ON t_position_histories.member_id = m_members.member_id
WHERE position_history_id = $1;

-- name: GetPositionHistories :many
SELECT * FROM t_position_histories
WHERE
	CASE WHEN @where_member::boolean = true THEN member_id = @member_id ELSE TRUE END
AND
	CASE WHEN @where_earlier_send_at::boolean = true THEN send_at >= @earlier_send_at ELSE TRUE END
AND
	CASE WHEN @where_later_send_at::boolean = true THEN send_at <= @later_send_at ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'send_at' THEN send_at END ASC,
	t_position_histories_pkey DESC
LIMIT $1 OFFSET $2;

-- name: GetPositionHistoriesWithMember :many
SELECT sqlc.embed(t_position_histories), sqlc.embed(m_members) FROM t_position_histories
INNER JOIN m_members ON t_position_histories.member_id = m_members.member_id
WHERE
	CASE WHEN @where_member::boolean = true THEN t_position_histories.member_id = @member_id ELSE TRUE END
AND
	CASE WHEN @where_earlier_send_at::boolean = true THEN send_at >= @earlier_send_at ELSE TRUE END
AND
	CASE WHEN @where_later_send_at::boolean = true THEN send_at <= @later_send_at ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'send_at' THEN send_at END ASC,
	t_position_histories_pkey DESC
LIMIT $1 OFFSET $2;

-- name: CountPositionHistories :one
SELECT COUNT(*) FROM t_position_histories
WHERE
	CASE WHEN @where_member::boolean = true THEN member_id = @member_id ELSE TRUE END
AND
	CASE WHEN @where_earlier_send_at::boolean = true THEN send_at >= @earlier_send_at ELSE TRUE END
AND
	CASE WHEN @where_later_send_at::boolean = true THEN send_at <= @later_send_at ELSE TRUE END;
