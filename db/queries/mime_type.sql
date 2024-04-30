-- name: CreateMimeTypes :copyfrom
INSERT INTO m_mime_types (name, key) VALUES ($1, $2);

-- name: CreateMimeType :one
INSERT INTO m_mime_types (name, key) VALUES ($1, $2) RETURNING *;

-- name: DeleteMimeType :exec
DELETE FROM m_mime_types WHERE mime_type_id = $1;

-- name: DeleteMimeTypeByKey :exec
DELETE FROM m_mime_types WHERE key = $1;

-- name: FindMimeTypeByID :one
SELECT * FROM m_mime_types WHERE mime_type_id = $1;

-- name: FindMimeTypeByKey :one
SELECT * FROM m_mime_types WHERE key = $1;

-- name: GetMimeTypes :many
SELECT * FROM m_mime_types
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN name END DESC,
	m_mime_types_pkey ASC;

-- name: GetMimeTypesUseNumberedPaginate :many
SELECT * FROM m_mime_types
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN name END DESC,
	m_mime_types_pkey ASC
LIMIT $1 OFFSET $2;

-- name: GetMimeTypesUseKeysetPaginate :many
SELECT * FROM m_mime_types
WHERE
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
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'next' THEN name END ASC,
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'prev' THEN name END DESC,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'next' THEN name END ASC,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'prev' THEN name END DESC,
	CASE WHEN @cursor_direction::text = 'next' THEN m_mime_types_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN m_mime_types_pkey END DESC
LIMIT $1;

-- name: GetPluralMimeTypes :many
SELECT * FROM m_mime_types
WHERE mime_type_id = ANY(@mime_type_ids::uuid[])
ORDER BY
	m_mime_types_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountMimeTypes :one
SELECT COUNT(*) FROM m_mime_types
WHERE
	CASE WHEN @where_like_name::boolean = true THEN name LIKE '%' || @search_name::text || '%' ELSE TRUE END;
