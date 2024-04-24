-- name: CreateEventTypes :copyfrom
INSERT INTO m_event_types (name, key, color) VALUES ($1, $2, $3);

-- name: CreateEventType :one
INSERT INTO m_event_types (name, key, color) VALUES ($1, $2, $3) RETURNING *;

-- name: UpdateEventType :one
UPDATE m_event_types SET name = $2, color = $3 WHERE event_type_id = $1 RETURNING *;

-- name: UpdateEventTypeByKey :one
UPDATE m_event_types SET name = $2, color = $3 WHERE key = $1 RETURNING *;

-- name: UpdateEventTypeKey :one
UPDATE m_event_types SET key = $2 WHERE event_type_id = $1 RETURNING *;

-- name: DeleteEventType :exec
DELETE FROM m_event_types WHERE event_type_id = $1;

-- name: DeleteEventTypeByKey :exec
DELETE FROM m_event_types WHERE key = $1;

-- name: FindEventTypeByID :one
SELECT * FROM m_event_types WHERE event_type_id = $1;

-- name: FindEventTypeByKey :one
SELECT * FROM m_event_types WHERE key = $1;

-- name: GetEventTypes :many
SELECT * FROM m_event_types
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN name END ASC,
	m_event_types_pkey DESC
LIMIT $1 OFFSET $2;

