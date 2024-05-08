-- name: CreatePermissionAssociations :copyfrom
INSERT INTO m_permission_associations (permission_id, work_position_id) VALUES ($1, $2);

-- name: CreatePermissionAssociation :one
INSERT INTO m_permission_associations (permission_id, work_position_id) VALUES ($1, $2) RETURNING *;

-- name: DeletePermissionAssociation :exec
DELETE FROM m_permission_associations WHERE permission_id = $1 AND work_position_id = $2;

-- name: DeletePermissionOnPermission :exec
DELETE FROM m_permission_associations WHERE permission_id = $1;

-- name: DeletePermissionOnPermissions :exec
DELETE FROM m_permission_associations WHERE permission_id = ANY($1::uuid[]);

-- name: PluralDeletePermissionAssociationsOnPermission :exec
DELETE FROM m_permission_associations WHERE permission_id = $1 AND work_position_id = ANY($2::uuid[]);

-- name: DeletePermissionOnWorkPosition :exec
DELETE FROM m_permission_associations WHERE work_position_id = $1;

-- name: DeletePermissionOnWorkPositions :exec
DELETE FROM m_permission_associations WHERE work_position_id = ANY($1::uuid[]);

-- name: PluralDeletePermissionAssociationsOnWorkPosition :exec
DELETE FROM m_permission_associations WHERE work_position_id = $1 AND permission_id = ANY($2::uuid[]);

-- name: GetWorkPositionsOnPermission :many
SELECT sqlc.embed(m_permission_associations), sqlc.embed(m_work_positions) FROM m_permission_associations
LEFT JOIN m_work_positions ON m_permission_associations.work_position_id = m_work_positions.work_position_id
WHERE permission_id = $1
AND
	CASE WHEN @where_like_name::boolean = true THEN m_work_positions.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_work_positions.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_work_positions.name END DESC,
	m_permission_associations_pkey ASC;

-- name: GetWorkPositionsOnPermissionUseNumberedPaginate :many
SELECT sqlc.embed(m_permission_associations), sqlc.embed(m_work_positions) FROM m_permission_associations
LEFT JOIN m_work_positions ON m_permission_associations.work_position_id = m_work_positions.work_position_id
WHERE permission_id = $1
AND
	CASE WHEN @where_like_name::boolean = true THEN m_work_positions.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_work_positions.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_work_positions.name END DESC,
	m_permission_associations_pkey ASC
LIMIT $2 OFFSET $3;

-- name: GetWorkPositionsOnPermissionUseKeysetPaginate :many
SELECT sqlc.embed(m_permission_associations), sqlc.embed(m_work_positions) FROM m_permission_associations
LEFT JOIN m_work_positions ON m_permission_associations.work_position_id = m_work_positions.work_position_id
WHERE permission_id = $1
AND
	CASE WHEN @where_like_name::boolean = true THEN m_work_positions.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'name' THEN m_work_positions.name > @name_cursor OR (m_work_positions.name = @name_cursor AND m_permission_associations_pkey > @cursor::int)
				WHEN 'r_name' THEN m_work_positions.name < @name_cursor OR (m_work_positions.name = @name_cursor AND m_permission_associations_pkey > @cursor::int)
				ELSE m_permission_associations_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'name' THEN m_work_positions.name < @name_cursor OR (m_work_positions.name = @name_cursor AND m_permission_associations_pkey < @cursor::int)
				WHEN 'r_name' THEN m_work_positions.name > @name_cursor OR (m_work_positions.name = @name_cursor AND m_permission_associations_pkey < @cursor::int)
				ELSE m_permission_associations_pkey < @cursor::int
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'next' THEN m_work_positions.name END ASC,
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'prev' THEN m_work_positions.name END DESC,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'next' THEN m_work_positions.name END ASC,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'prev' THEN m_work_positions.name END DESC,
	CASE WHEN @cursor_direction::text = 'next' THEN m_permission_associations_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN m_permission_associations_pkey END DESC
LIMIT $2;

-- name: GetPluralWorkPositionsOnPermission :many
SELECT sqlc.embed(m_permission_associations), sqlc.embed(m_work_positions) FROM m_permission_associations
LEFT JOIN m_work_positions ON m_permission_associations.work_position_id = m_work_positions.work_position_id
WHERE permission_id = ANY(@permission_ids::uuid[])
ORDER BY
	m_permission_associations_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountWorkPositionsOnPermission :one
SELECT COUNT(*) FROM m_permission_associations
LEFT JOIN m_work_positions ON m_permission_associations.work_position_id = m_work_positions.work_position_id
WHERE permission_id = $1
AND
	CASE WHEN @where_like_name::boolean = true THEN m_work_positions.name LIKE '%' || @search_name::text || '%' ELSE TRUE END;

-- name: GetPermissionsOnWorkPosition :many
SELECT sqlc.embed(m_permission_associations), sqlc.embed(m_permissions) FROM m_permission_associations
LEFT JOIN m_permissions ON m_permission_associations.permission_id = m_permissions.permission_id
WHERE work_position_id = $1
AND
	CASE WHEN @where_like_name::boolean = true THEN m_permissions.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_permissions.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_permissions.name END DESC,
	m_permission_associations_pkey ASC;

-- name: GetPermissionsOnWorkPositionUseNumberedPaginate :many
SELECT sqlc.embed(m_permission_associations), sqlc.embed(m_permissions) FROM m_permission_associations
LEFT JOIN m_permissions ON m_permission_associations.permission_id = m_permissions.permission_id
WHERE work_position_id = $1
AND
	CASE WHEN @where_like_name::boolean = true THEN m_permissions.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_permissions.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_permissions.name END DESC,
	m_permission_associations_pkey ASC
LIMIT $2 OFFSET $3;

-- name: GetPermissionsOnWorkPositionUseKeysetPaginate :many
SELECT sqlc.embed(m_permission_associations), sqlc.embed(m_permissions) FROM m_permission_associations
LEFT JOIN m_permissions ON m_permission_associations.permission_id = m_permissions.permission_id
WHERE work_position_id = $1
AND
	CASE WHEN @where_like_name::boolean = true THEN m_permissions.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
AND
	CASE @cursor_direction::text
		WHEN 'next' THEN
			CASE @order_method::text
				WHEN 'name' THEN m_permissions.name > @name_cursor OR (m_permissions.name = @name_cursor AND m_permission_associations_pkey > @cursor::int)
				WHEN 'r_name' THEN m_permissions.name < @name_cursor OR (m_permissions.name = @name_cursor AND m_permission_associations_pkey > @cursor::int)
				ELSE m_permission_associations_pkey > @cursor::int
			END
		WHEN 'prev' THEN
			CASE @order_method::text
				WHEN 'name' THEN m_permissions.name < @name_cursor OR (m_permissions.name = @name_cursor AND m_permission_associations_pkey < @cursor::int)
				WHEN 'r_name' THEN m_permissions.name > @name_cursor OR (m_permissions.name = @name_cursor AND m_permission_associations_pkey < @cursor::int)
				ELSE m_permission_associations_pkey < @cursor::int
			END
	END
ORDER BY
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'next' THEN m_permissions.name END ASC,
	CASE WHEN @order_method::text = 'name' AND @cursor_direction::text = 'prev' THEN m_permissions.name END DESC,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'next' THEN m_permissions.name END ASC,
	CASE WHEN @order_method::text = 'r_name' AND @cursor_direction::text = 'prev' THEN m_permissions.name END DESC,
	CASE WHEN @cursor_direction::text = 'next' THEN m_permission_associations_pkey END ASC,
	CASE WHEN @cursor_direction::text = 'prev' THEN m_permission_associations_pkey END DESC
LIMIT $2;

-- name: GetPluralPermissionsOnWorkPosition :many
SELECT sqlc.embed(m_permission_associations), sqlc.embed(m_permissions) FROM m_permission_associations
LEFT JOIN m_permissions ON m_permission_associations.permission_id = m_permissions.permission_id
WHERE work_position_id = ANY(@work_position_ids::uuid[])
ORDER BY
	m_permission_associations_pkey ASC
LIMIT $1 OFFSET $2;

-- name: CountPermissionsOnWorkPosition :one
SELECT COUNT(*) FROM m_permission_associations
LEFT JOIN m_permissions ON m_permission_associations.permission_id = m_permissions.permission_id
WHERE work_position_id = $1
AND
	CASE WHEN @where_like_name::boolean = true THEN m_permissions.name LIKE '%' || @search_name::text || '%' ELSE TRUE END;
