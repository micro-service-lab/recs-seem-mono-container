-- name: CreateRecordTypes :copyfrom
INSERT INTO m_record_types (name, key) VALUES ($1, $2);

-- name: CreateRecordType :one
INSERT INTO m_record_types (name, key) VALUES ($1, $2) RETURNING *;

-- name: UpdateRecordType :one
UPDATE m_record_types SET name = $2, key = $3 WHERE record_type_id = $1 RETURNING *;

-- name: UpdateRecordTypeByKey :one
UPDATE m_record_types SET name = $2 WHERE key = $1 RETURNING *;

-- name: DeleteRecordType :execrows
DELETE FROM m_record_types WHERE record_type_id = $1;

-- name: DeleteRecordTypeByKey :execrows
DELETE FROM m_record_types WHERE key = $1;

-- name: PluralDeleteRecordTypes :execrows
DELETE FROM m_record_types WHERE record_type_id = ANY($1::uuid[]);

-- name: FindRecordTypeByID :one
SELECT * FROM m_record_types WHERE record_type_id = $1;

-- name: FindRecordTypeByKey :one
SELECT * FROM m_record_types WHERE key = $1;

-- name: GetRecordTypes :many
SELECT * FROM m_record_types
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_record_types.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN name END DESC,
	m_record_types_pkey ASC;

-- name: GetRecordTypesUseNumberedPaginate :many
SELECT * FROM m_record_types
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_record_types.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN name END DESC,
	m_record_types_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetRecordTypesUseKeysetPaginate :many
SELECT * FROM m_record_types
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_record_types.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'name' THEN name > @name_cursor OR (name = @name_cursor AND m_record_types_pkey > @cursor::int)
				WHEN 'r_name' THEN name < @name_cursor OR (name = @name_cursor AND m_record_types_pkey > @cursor::int)
				ELSE m_record_types_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'name' THEN name < @name_cursor OR (name = @name_cursor AND m_record_types_pkey < @cursor::int)
				WHEN 'r_name' THEN name > @name_cursor OR (name = @name_cursor AND m_record_types_pkey < @cursor::int)
				ELSE m_record_types_pkey < @cursor::int
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'next' THEN name END ASC,
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'prev' THEN name END DESC,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'next' THEN name END DESC,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'prev' THEN name END ASC,
	CASE WHEN @cursor_direction::text = 'next' THEN m_record_types_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN m_record_types_pkey END DESC
LIMIT $1;

-- name: GetPluralRecordTypes :many
SELECT * FROM m_record_types
WHERE
	record_type_id = ANY(@record_type_ids::uuid[])
ORDER BY
	m_record_types_pkey ASC;

-- name: GetPluralRecordTypesUseNumberedPaginate :many
SELECT * FROM m_record_types
WHERE
	record_type_id = ANY(@record_type_ids::uuid[])
ORDER BY
	m_record_types_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountRecordTypes :one
SELECT COUNT(*) FROM m_record_types
WHERE
	CASE WHEN @where_like_name::boolean = true THEN name LIKE '%' || @search_name::text || '%' ELSE TRUE END;
