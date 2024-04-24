-- name: CreateAttendanceTypes :copyfrom
INSERT INTO m_attendance_types (name, key, color) VALUES ($1, $2, $3);

-- name: CreateAttendanceType :one
INSERT INTO m_attendance_types (name, key, color) VALUES ($1, $2, $3) RETURNING *;

-- name: UpdateAttendanceType :one
UPDATE m_attendance_types SET name = $2, color = $3 WHERE attendance_type_id = $1 RETURNING *;

-- name: UpdateAttendanceTypeByKey :one
UPDATE m_attendance_types SET name = $2, color = $3 WHERE key = $1 RETURNING *;

-- name: UpdateAttendanceTypeKey :one
UPDATE m_attendance_types SET key = $2 WHERE attendance_type_id = $1 RETURNING *;

-- name: DeleteAttendanceType :exec
DELETE FROM m_attendance_types WHERE attendance_type_id = $1;

-- name: DeleteAttendanceTypeByKey :exec
DELETE FROM m_attendance_types WHERE key = $1;

-- name: FindAttendanceTypeByID :one
SELECT * FROM m_attendance_types WHERE attendance_type_id = $1;

-- name: FindAttendanceTypeByKey :one
SELECT * FROM m_attendance_types WHERE key = $1;

-- name: GetAttendanceTypes :many
SELECT * FROM m_attendance_types
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN name END ASC,
	m_attendance_types_pkey DESC
LIMIT $1 OFFSET $2;

-- name: CountAttendanceTypes :one
SELECT COUNT(*) FROM m_attendance_types
WHERE
	CASE WHEN @where_like_name::boolean = true THEN name LIKE '%' || @search_name::text || '%' ELSE TRUE END;
