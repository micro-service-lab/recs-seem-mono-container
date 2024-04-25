-- name: CreatePermissionAssociations :copyfrom
INSERT INTO m_permission_associations (permission_id, work_position_id) VALUES ($1, $2);

-- name: CreatePermissionAssociation :one
INSERT INTO m_permission_associations (permission_id, work_position_id) VALUES ($1, $2) RETURNING *;

-- name: DeletePermissionAssociation :exec
DELETE FROM m_permission_associations WHERE permission_id = $1 AND work_position_id = $2;

-- name: GetWorkPositionsOnPermission :many
SELECT sqlc.embed(m_permission_associations), sqlc.embed(m_work_positions) FROM m_permission_associations
LEFT JOIN m_work_positions ON m_permission_associations.work_position_id = m_work_positions.work_position_id
WHERE permission_id = $1
AND
	CASE WHEN @where_like_name::boolean = true THEN m_work_positions.name LIKE '%' || @search_name::text || '%' ELSE TRUE END
ORDER BY
	CASE WHEN @order_method::text = 'name' THEN m_work_positions.name END ASC,
	CASE WHEN @order_method::text = 'r_name' THEN m_work_positions.name END DESC,
	m_permission_associations_pkey DESC
LIMIT $2 OFFSET $3;

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
	m_permission_associations_pkey DESC
LIMIT $2 OFFSET $3;

-- name: CountPermissionsOnWorkPosition :one
SELECT COUNT(*) FROM m_permission_associations
LEFT JOIN m_permissions ON m_permission_associations.permission_id = m_permissions.permission_id
WHERE work_position_id = $1
AND
	CASE WHEN @where_like_name::boolean = true THEN m_permissions.name LIKE '%' || @search_name::text || '%' ELSE TRUE END;
