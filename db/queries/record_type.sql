-- name: CreateRecordTypes :copyfrom
INSERT INTO m_record_types (name, key) VALUES ($1, $2);

-- name: CreateRecordType :one
INSERT INTO m_record_types (name, key) VALUES ($1, $2) RETURNING *;

-- name: UpdateRecordType :one
UPDATE m_record_types SET name = $2 WHERE record_type_id = $1 RETURNING *;

-- name: UpdateRecordTypeByKey :one
UPDATE m_record_types SET name = $2 WHERE key = $1 RETURNING *;

-- name: UpdateRecordTypeKey :one
UPDATE m_record_types SET key = $2 WHERE record_type_id = $1 RETURNING *;

-- name: DeleteRecordType :exec
DELETE FROM m_record_types WHERE record_type_id = $1;

-- name: DeleteRecordTypeByKey :exec
DELETE FROM m_record_types WHERE key = $1;

-- name: FindRecordTypeByID :one
SELECT * FROM m_record_types WHERE record_type_id = $1;

-- name: FindRecordTypeByKey :one
SELECT * FROM m_record_types WHERE key = $1;

-- name: GetRecordTypes :many
SELECT * FROM m_record_types
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN name END ASC,
	m_record_types_pkey DESC
LIMIT $1 OFFSET $2;

-- name: CountRecordTypes :one
SELECT COUNT(*) FROM m_record_types
WHERE
	CASE WHEN @where_like_name::boolean = true THEN name LIKE '%' || @search_name::text || '%' ELSE TRUE END;
