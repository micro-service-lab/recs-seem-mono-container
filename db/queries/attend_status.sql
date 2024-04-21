-- name: CreateAttendStatuses :copyfrom
INSERT INTO m_attend_statuses (name, key) VALUES ($1, $2);

-- name: CreateAttendStatus :one
INSERT INTO m_attend_statuses (name, key) VALUES ($1, $2) RETURNING *;

-- name: UpdateAttendStatus :one
UPDATE m_attend_statuses SET name = $2 WHERE attend_status_id = $1 RETURNING *;

-- name: UpdateAttendStatusKey :one
UPDATE m_attend_statuses SET key = $2 WHERE attend_status_id = $1 RETURNING *;

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
WHERE CASE
	WHEN @where_like_name::boolean = true THEN m_attend_statuses.name LIKE '%' || @search_name::text || '%'
END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_attend_statuses.name END ASC,
	m_attend_statuses_pkey DESC
LIMIT $1 OFFSET $2;

-- name: GetAttendStatusesByKeys :many
SELECT * FROM m_attend_statuses WHERE key = ANY(@keys::varchar[])
AND CASE
	WHEN @where_like_name::boolean = true THEN m_attend_statuses.name LIKE '%' || @search_name::text || '%'
END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_attend_statuses.name END ASC,
	m_attend_statuses_pkey DESC
LIMIT $1 OFFSET $2;

-- name: CountAttendStatuses :one
SELECT COUNT(*) FROM m_attend_statuses;
