-- name: CreateMimeTypes :copyfrom
INSERT INTO m_mime_types (name, key, kind) VALUES ($1, $2, $3);

-- name: CreateMimeType :one
INSERT INTO m_mime_types (name, key, kind) VALUES ($1, $2, $3) RETURNING *;

-- name: UpdateMimeType :one
UPDATE m_mime_types SET name = $2, key = $3, kind = $4 WHERE mime_type_id = $1 RETURNING *;

-- name: UpdateMimeTypeByKey :one
UPDATE m_mime_types SET name = $2, kind = $3 WHERE key = $1 RETURNING *;

-- name: DeleteMimeType :execrows
DELETE FROM m_mime_types WHERE mime_type_id = $1;

-- name: DeleteMimeTypeByKey :execrows
DELETE FROM m_mime_types WHERE key = $1;

-- name: PluralDeleteMimeTypes :execrows
DELETE FROM m_mime_types WHERE mime_type_id = ANY(@mime_type_ids::uuid[]);

-- name: FindMimeTypeByID :one
SELECT * FROM m_mime_types WHERE mime_type_id = $1;

-- name: FindMimeTypeByKey :one
SELECT * FROM m_mime_types WHERE key = $1;

-- name: FindMimeTypeByKind :one
SELECT * FROM m_mime_types WHERE kind = $1;

-- name: GetMimeTypes :many
SELECT * FROM m_mime_types
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_mime_types.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN name END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' THEN name END DESC NULLS LAST,
	m_mime_types_pkey ASC;

-- name: GetMimeTypesUseNumberedPaginate :many
SELECT * FROM m_mime_types
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_mime_types.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN name END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' THEN name END DESC NULLS LAST,
	m_mime_types_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetMimeTypesUseKeysetPaginate :many
SELECT * FROM m_mime_types
WHERE
	CASE WHEN @where_like_name::boolean = true THEN m_mime_types.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'name' THEN name > @name_cursor OR (name = @name_cursor AND m_mime_types_pkey > @cursor::int)
				WHEN 'r_name' THEN name < @name_cursor OR (name = @name_cursor AND m_mime_types_pkey > @cursor::int)
				ELSE m_mime_types_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'name' THEN name < @name_cursor OR (name = @name_cursor AND m_mime_types_pkey < @cursor::int)
				WHEN 'r_name' THEN name > @name_cursor OR (name = @name_cursor AND m_mime_types_pkey < @cursor::int)
				ELSE m_mime_types_pkey < @cursor::int
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'next' THEN name END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'prev' THEN name END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'next' THEN name END DESC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'prev' THEN name END ASC NULLS LAST,
	CASE WHEN @cursor_direction::text = 'next' THEN m_mime_types_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN m_mime_types_pkey END DESC
LIMIT $1;

-- name: GetPluralMimeTypes :many
SELECT * FROM m_mime_types
WHERE mime_type_id = ANY(@mime_type_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN name END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' THEN name END DESC NULLS LAST,
	m_mime_types_pkey ASC;

-- name: GetPluralMimeTypesUseNumberedPaginate :many
SELECT * FROM m_mime_types
WHERE mime_type_id = ANY(@mime_type_ids::uuid[])
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN name END ASC NULLS LAST,
	CASE WHEN @order_method::text = 'r_name' THEN name END DESC NULLS LAST,
	m_mime_types_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountMimeTypes :one
SELECT COUNT(*) FROM m_mime_types
WHERE
	CASE WHEN @where_like_name::boolean = true THEN name LIKE '%' || @search_name::text || '%' ELSE TRUE END;
