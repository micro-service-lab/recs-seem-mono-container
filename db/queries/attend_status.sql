-- name: CreateAttendStatuses :copyfrom
INSERT INTO m_attend_statuses (name, key) VALUES ($1, $2);

-- name: CreateAttendStatus :one
INSERT INTO m_attend_statuses (name, key) VALUES ($1, $2) RETURNING *;

-- name: UpdateAttendStatus :one
UPDATE m_attend_statuses SET name = $2, key = $3 WHERE attend_status_id = $1 RETURNING *;

-- name: UpdateAttendStatusByKey :one
UPDATE m_attend_statuses SET name = $2, key = $3 WHERE key = $1 RETURNING *;

-- name: DeleteAttendStatus :exec
DELETE FROM m_attend_statuses WHERE attend_status_id = $1;

-- name: DeleteAttendStatusByKey :exec
DELETE FROM m_attend_statuses WHERE key = $1;

-- name: FindAttendStatusById :one
SELECT * FROM m_attend_statuses WHERE attend_status_id = $1;

-- name: FindAttendStatusByKey :one
SELECT * FROM m_attend_statuses WHERE key = $1;

-- name: GetAttendStatuses :many
SELECT * FROM m_attend_statuses
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_attend_statuses.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_attend_statuses.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_attend_statuses.name END DESC,
	m_attend_statuses_pkey DESC
LIMIT $1 OFFSET $2;

-- name: CountAttendStatuses :one
SELECT COUNT(*) FROM m_attend_statuses
WHERE
	CASE WHEN @where_like_name::boolean = true THEN name LIKE '%' || @search_name::text || '%' ELSE TRUE END;
