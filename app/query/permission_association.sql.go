// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: permission_association.sql

package query

import (
	"context"

	"github.com/google/uuid"
)

const countPermissionsOnWorkPosition = `-- name: CountPermissionsOnWorkPosition :one
SELECT COUNT(*) FROM m_permission_associations
LEFT JOIN m_permissions ON m_permission_associations.permission_id = m_permissions.permission_id
WHERE work_position_id = $1
AND
	CASE WHEN $2::boolean = true THEN m_permissions.name LIKE '%' || $3::text || '%' ELSE TRUE END
`

type CountPermissionsOnWorkPositionParams struct {
	WorkPositionID uuid.UUID `json:"work_position_id"`
	WhereLikeName  bool      `json:"where_like_name"`
	SearchName     string    `json:"search_name"`
}

func (q *Queries) CountPermissionsOnWorkPosition(ctx context.Context, arg CountPermissionsOnWorkPositionParams) (int64, error) {
	row := q.db.QueryRow(ctx, countPermissionsOnWorkPosition, arg.WorkPositionID, arg.WhereLikeName, arg.SearchName)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countWorkPositionsOnPermission = `-- name: CountWorkPositionsOnPermission :one
SELECT COUNT(*) FROM m_permission_associations
LEFT JOIN m_work_positions ON m_permission_associations.work_position_id = m_work_positions.work_position_id
WHERE permission_id = $1
AND
	CASE WHEN $2::boolean = true THEN m_work_positions.name LIKE '%' || $3::text || '%' ELSE TRUE END
`

type CountWorkPositionsOnPermissionParams struct {
	PermissionID  uuid.UUID `json:"permission_id"`
	WhereLikeName bool      `json:"where_like_name"`
	SearchName    string    `json:"search_name"`
}

func (q *Queries) CountWorkPositionsOnPermission(ctx context.Context, arg CountWorkPositionsOnPermissionParams) (int64, error) {
	row := q.db.QueryRow(ctx, countWorkPositionsOnPermission, arg.PermissionID, arg.WhereLikeName, arg.SearchName)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createPermissionAssociation = `-- name: CreatePermissionAssociation :one
INSERT INTO m_permission_associations (permission_id, work_position_id) VALUES ($1, $2) RETURNING m_permission_associations_pkey, permission_id, work_position_id
`

type CreatePermissionAssociationParams struct {
	PermissionID   uuid.UUID `json:"permission_id"`
	WorkPositionID uuid.UUID `json:"work_position_id"`
}

func (q *Queries) CreatePermissionAssociation(ctx context.Context, arg CreatePermissionAssociationParams) (PermissionAssociation, error) {
	row := q.db.QueryRow(ctx, createPermissionAssociation, arg.PermissionID, arg.WorkPositionID)
	var i PermissionAssociation
	err := row.Scan(&i.MPermissionAssociationsPkey, &i.PermissionID, &i.WorkPositionID)
	return i, err
}

type CreatePermissionAssociationsParams struct {
	PermissionID   uuid.UUID `json:"permission_id"`
	WorkPositionID uuid.UUID `json:"work_position_id"`
}

const deletePermissionAssociation = `-- name: DeletePermissionAssociation :exec
DELETE FROM m_permission_associations WHERE permission_id = $1 AND work_position_id = $2
`

type DeletePermissionAssociationParams struct {
	PermissionID   uuid.UUID `json:"permission_id"`
	WorkPositionID uuid.UUID `json:"work_position_id"`
}

func (q *Queries) DeletePermissionAssociation(ctx context.Context, arg DeletePermissionAssociationParams) error {
	_, err := q.db.Exec(ctx, deletePermissionAssociation, arg.PermissionID, arg.WorkPositionID)
	return err
}

const deletePermissionOnPermission = `-- name: DeletePermissionOnPermission :exec
DELETE FROM m_permission_associations WHERE permission_id = $1
`

func (q *Queries) DeletePermissionOnPermission(ctx context.Context, permissionID uuid.UUID) error {
	_, err := q.db.Exec(ctx, deletePermissionOnPermission, permissionID)
	return err
}

const deletePermissionOnPermissions = `-- name: DeletePermissionOnPermissions :exec
DELETE FROM m_permission_associations WHERE permission_id = ANY($1::uuid[])
`

func (q *Queries) DeletePermissionOnPermissions(ctx context.Context, dollar_1 []uuid.UUID) error {
	_, err := q.db.Exec(ctx, deletePermissionOnPermissions, dollar_1)
	return err
}

const deletePermissionOnWorkPosition = `-- name: DeletePermissionOnWorkPosition :exec
DELETE FROM m_permission_associations WHERE work_position_id = $1
`

func (q *Queries) DeletePermissionOnWorkPosition(ctx context.Context, workPositionID uuid.UUID) error {
	_, err := q.db.Exec(ctx, deletePermissionOnWorkPosition, workPositionID)
	return err
}

const deletePermissionOnWorkPositions = `-- name: DeletePermissionOnWorkPositions :exec
DELETE FROM m_permission_associations WHERE work_position_id = ANY($1::uuid[])
`

func (q *Queries) DeletePermissionOnWorkPositions(ctx context.Context, dollar_1 []uuid.UUID) error {
	_, err := q.db.Exec(ctx, deletePermissionOnWorkPositions, dollar_1)
	return err
}

const getPermissionsOnWorkPosition = `-- name: GetPermissionsOnWorkPosition :many
SELECT m_permission_associations.m_permission_associations_pkey, m_permission_associations.permission_id, m_permission_associations.work_position_id, m_permissions.m_permissions_pkey, m_permissions.permission_id, m_permissions.name, m_permissions.description, m_permissions.key, m_permissions.permission_category_id FROM m_permission_associations
LEFT JOIN m_permissions ON m_permission_associations.permission_id = m_permissions.permission_id
WHERE work_position_id = $1
AND
	CASE WHEN $2::boolean = true THEN m_permissions.name LIKE '%' || $3::text || '%' ELSE TRUE END
ORDER BY
	CASE WHEN $4::text = 'name' THEN m_permissions.name END ASC,
	CASE WHEN $4::text = 'r_name' THEN m_permissions.name END DESC,
	m_permission_associations_pkey ASC
`

type GetPermissionsOnWorkPositionParams struct {
	WorkPositionID uuid.UUID `json:"work_position_id"`
	WhereLikeName  bool      `json:"where_like_name"`
	SearchName     string    `json:"search_name"`
	OrderMethod    string    `json:"order_method"`
}

type GetPermissionsOnWorkPositionRow struct {
	PermissionAssociation PermissionAssociation `json:"permission_association"`
	Permission            Permission            `json:"permission"`
}

func (q *Queries) GetPermissionsOnWorkPosition(ctx context.Context, arg GetPermissionsOnWorkPositionParams) ([]GetPermissionsOnWorkPositionRow, error) {
	rows, err := q.db.Query(ctx, getPermissionsOnWorkPosition,
		arg.WorkPositionID,
		arg.WhereLikeName,
		arg.SearchName,
		arg.OrderMethod,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPermissionsOnWorkPositionRow{}
	for rows.Next() {
		var i GetPermissionsOnWorkPositionRow
		if err := rows.Scan(
			&i.PermissionAssociation.MPermissionAssociationsPkey,
			&i.PermissionAssociation.PermissionID,
			&i.PermissionAssociation.WorkPositionID,
			&i.Permission.MPermissionsPkey,
			&i.Permission.PermissionID,
			&i.Permission.Name,
			&i.Permission.Description,
			&i.Permission.Key,
			&i.Permission.PermissionCategoryID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPermissionsOnWorkPositionUseKeysetPaginate = `-- name: GetPermissionsOnWorkPositionUseKeysetPaginate :many
SELECT m_permission_associations.m_permission_associations_pkey, m_permission_associations.permission_id, m_permission_associations.work_position_id, m_permissions.m_permissions_pkey, m_permissions.permission_id, m_permissions.name, m_permissions.description, m_permissions.key, m_permissions.permission_category_id FROM m_permission_associations
LEFT JOIN m_permissions ON m_permission_associations.permission_id = m_permissions.permission_id
WHERE work_position_id = $1
AND
	CASE WHEN $3::boolean = true THEN m_permissions.name LIKE '%' || $4::text || '%' ELSE TRUE END
AND
	CASE $5::text
		WHEN 'next' THEN
			CASE $6::text
				WHEN 'name' THEN m_permissions.name > $7 OR (m_permissions.name = $7 AND m_permission_associations_pkey > $8::int)
				WHEN 'r_name' THEN m_permissions.name < $7 OR (m_permissions.name = $7 AND m_permission_associations_pkey > $8::int)
				ELSE m_permission_associations_pkey > $8::int
			END
		WHEN 'prev' THEN
			CASE $6::text
				WHEN 'name' THEN m_permissions.name < $7 OR (m_permissions.name = $7 AND m_permission_associations_pkey < $8::int)
				WHEN 'r_name' THEN m_permissions.name > $7 OR (m_permissions.name = $7 AND m_permission_associations_pkey < $8::int)
				ELSE m_permission_associations_pkey < $8::int
			END
	END
ORDER BY
	CASE WHEN $6::text = 'name' AND $5::text = 'next' THEN m_permissions.name END ASC,
	CASE WHEN $6::text = 'name' AND $5::text = 'prev' THEN m_permissions.name END DESC,
	CASE WHEN $6::text = 'r_name' AND $5::text = 'next' THEN m_permissions.name END ASC,
	CASE WHEN $6::text = 'r_name' AND $5::text = 'prev' THEN m_permissions.name END DESC,
	CASE WHEN $5::text = 'next' THEN m_permission_associations_pkey END ASC,
	CASE WHEN $5::text = 'prev' THEN m_permission_associations_pkey END DESC
LIMIT $2
`

type GetPermissionsOnWorkPositionUseKeysetPaginateParams struct {
	WorkPositionID  uuid.UUID `json:"work_position_id"`
	Limit           int32     `json:"limit"`
	WhereLikeName   bool      `json:"where_like_name"`
	SearchName      string    `json:"search_name"`
	CursorDirection string    `json:"cursor_direction"`
	OrderMethod     string    `json:"order_method"`
	NameCursor      string    `json:"name_cursor"`
	Cursor          int32     `json:"cursor"`
}

type GetPermissionsOnWorkPositionUseKeysetPaginateRow struct {
	PermissionAssociation PermissionAssociation `json:"permission_association"`
	Permission            Permission            `json:"permission"`
}

func (q *Queries) GetPermissionsOnWorkPositionUseKeysetPaginate(ctx context.Context, arg GetPermissionsOnWorkPositionUseKeysetPaginateParams) ([]GetPermissionsOnWorkPositionUseKeysetPaginateRow, error) {
	rows, err := q.db.Query(ctx, getPermissionsOnWorkPositionUseKeysetPaginate,
		arg.WorkPositionID,
		arg.Limit,
		arg.WhereLikeName,
		arg.SearchName,
		arg.CursorDirection,
		arg.OrderMethod,
		arg.NameCursor,
		arg.Cursor,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPermissionsOnWorkPositionUseKeysetPaginateRow{}
	for rows.Next() {
		var i GetPermissionsOnWorkPositionUseKeysetPaginateRow
		if err := rows.Scan(
			&i.PermissionAssociation.MPermissionAssociationsPkey,
			&i.PermissionAssociation.PermissionID,
			&i.PermissionAssociation.WorkPositionID,
			&i.Permission.MPermissionsPkey,
			&i.Permission.PermissionID,
			&i.Permission.Name,
			&i.Permission.Description,
			&i.Permission.Key,
			&i.Permission.PermissionCategoryID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPermissionsOnWorkPositionUseNumberedPaginate = `-- name: GetPermissionsOnWorkPositionUseNumberedPaginate :many
SELECT m_permission_associations.m_permission_associations_pkey, m_permission_associations.permission_id, m_permission_associations.work_position_id, m_permissions.m_permissions_pkey, m_permissions.permission_id, m_permissions.name, m_permissions.description, m_permissions.key, m_permissions.permission_category_id FROM m_permission_associations
LEFT JOIN m_permissions ON m_permission_associations.permission_id = m_permissions.permission_id
WHERE work_position_id = $1
AND
	CASE WHEN $4::boolean = true THEN m_permissions.name LIKE '%' || $5::text || '%' ELSE TRUE END
ORDER BY
	CASE WHEN $6::text = 'name' THEN m_permissions.name END ASC,
	CASE WHEN $6::text = 'r_name' THEN m_permissions.name END DESC,
	m_permission_associations_pkey ASC
LIMIT $2 OFFSET $3
`

type GetPermissionsOnWorkPositionUseNumberedPaginateParams struct {
	WorkPositionID uuid.UUID `json:"work_position_id"`
	Limit          int32     `json:"limit"`
	Offset         int32     `json:"offset"`
	WhereLikeName  bool      `json:"where_like_name"`
	SearchName     string    `json:"search_name"`
	OrderMethod    string    `json:"order_method"`
}

type GetPermissionsOnWorkPositionUseNumberedPaginateRow struct {
	PermissionAssociation PermissionAssociation `json:"permission_association"`
	Permission            Permission            `json:"permission"`
}

func (q *Queries) GetPermissionsOnWorkPositionUseNumberedPaginate(ctx context.Context, arg GetPermissionsOnWorkPositionUseNumberedPaginateParams) ([]GetPermissionsOnWorkPositionUseNumberedPaginateRow, error) {
	rows, err := q.db.Query(ctx, getPermissionsOnWorkPositionUseNumberedPaginate,
		arg.WorkPositionID,
		arg.Limit,
		arg.Offset,
		arg.WhereLikeName,
		arg.SearchName,
		arg.OrderMethod,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPermissionsOnWorkPositionUseNumberedPaginateRow{}
	for rows.Next() {
		var i GetPermissionsOnWorkPositionUseNumberedPaginateRow
		if err := rows.Scan(
			&i.PermissionAssociation.MPermissionAssociationsPkey,
			&i.PermissionAssociation.PermissionID,
			&i.PermissionAssociation.WorkPositionID,
			&i.Permission.MPermissionsPkey,
			&i.Permission.PermissionID,
			&i.Permission.Name,
			&i.Permission.Description,
			&i.Permission.Key,
			&i.Permission.PermissionCategoryID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPluralPermissionsOnWorkPosition = `-- name: GetPluralPermissionsOnWorkPosition :many
SELECT m_permission_associations.m_permission_associations_pkey, m_permission_associations.permission_id, m_permission_associations.work_position_id, m_permissions.m_permissions_pkey, m_permissions.permission_id, m_permissions.name, m_permissions.description, m_permissions.key, m_permissions.permission_category_id FROM m_permission_associations
LEFT JOIN m_permissions ON m_permission_associations.permission_id = m_permissions.permission_id
WHERE work_position_id = ANY($3::uuid[])
ORDER BY
	m_permission_associations_pkey ASC
LIMIT $1 OFFSET $2
`

type GetPluralPermissionsOnWorkPositionParams struct {
	Limit           int32       `json:"limit"`
	Offset          int32       `json:"offset"`
	WorkPositionIds []uuid.UUID `json:"work_position_ids"`
}

type GetPluralPermissionsOnWorkPositionRow struct {
	PermissionAssociation PermissionAssociation `json:"permission_association"`
	Permission            Permission            `json:"permission"`
}

func (q *Queries) GetPluralPermissionsOnWorkPosition(ctx context.Context, arg GetPluralPermissionsOnWorkPositionParams) ([]GetPluralPermissionsOnWorkPositionRow, error) {
	rows, err := q.db.Query(ctx, getPluralPermissionsOnWorkPosition, arg.Limit, arg.Offset, arg.WorkPositionIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPluralPermissionsOnWorkPositionRow{}
	for rows.Next() {
		var i GetPluralPermissionsOnWorkPositionRow
		if err := rows.Scan(
			&i.PermissionAssociation.MPermissionAssociationsPkey,
			&i.PermissionAssociation.PermissionID,
			&i.PermissionAssociation.WorkPositionID,
			&i.Permission.MPermissionsPkey,
			&i.Permission.PermissionID,
			&i.Permission.Name,
			&i.Permission.Description,
			&i.Permission.Key,
			&i.Permission.PermissionCategoryID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPluralWorkPositionsOnPermission = `-- name: GetPluralWorkPositionsOnPermission :many
SELECT m_permission_associations.m_permission_associations_pkey, m_permission_associations.permission_id, m_permission_associations.work_position_id, m_work_positions.m_work_positions_pkey, m_work_positions.work_position_id, m_work_positions.name, m_work_positions.description, m_work_positions.created_at, m_work_positions.updated_at FROM m_permission_associations
LEFT JOIN m_work_positions ON m_permission_associations.work_position_id = m_work_positions.work_position_id
WHERE permission_id = ANY($3::uuid[])
ORDER BY
	m_permission_associations_pkey ASC
LIMIT $1 OFFSET $2
`

type GetPluralWorkPositionsOnPermissionParams struct {
	Limit         int32       `json:"limit"`
	Offset        int32       `json:"offset"`
	PermissionIds []uuid.UUID `json:"permission_ids"`
}

type GetPluralWorkPositionsOnPermissionRow struct {
	PermissionAssociation PermissionAssociation `json:"permission_association"`
	WorkPosition          WorkPosition          `json:"work_position"`
}

func (q *Queries) GetPluralWorkPositionsOnPermission(ctx context.Context, arg GetPluralWorkPositionsOnPermissionParams) ([]GetPluralWorkPositionsOnPermissionRow, error) {
	rows, err := q.db.Query(ctx, getPluralWorkPositionsOnPermission, arg.Limit, arg.Offset, arg.PermissionIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPluralWorkPositionsOnPermissionRow{}
	for rows.Next() {
		var i GetPluralWorkPositionsOnPermissionRow
		if err := rows.Scan(
			&i.PermissionAssociation.MPermissionAssociationsPkey,
			&i.PermissionAssociation.PermissionID,
			&i.PermissionAssociation.WorkPositionID,
			&i.WorkPosition.MWorkPositionsPkey,
			&i.WorkPosition.WorkPositionID,
			&i.WorkPosition.Name,
			&i.WorkPosition.Description,
			&i.WorkPosition.CreatedAt,
			&i.WorkPosition.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getWorkPositionsOnPermission = `-- name: GetWorkPositionsOnPermission :many
SELECT m_permission_associations.m_permission_associations_pkey, m_permission_associations.permission_id, m_permission_associations.work_position_id, m_work_positions.m_work_positions_pkey, m_work_positions.work_position_id, m_work_positions.name, m_work_positions.description, m_work_positions.created_at, m_work_positions.updated_at FROM m_permission_associations
LEFT JOIN m_work_positions ON m_permission_associations.work_position_id = m_work_positions.work_position_id
WHERE permission_id = $1
AND
	CASE WHEN $2::boolean = true THEN m_work_positions.name LIKE '%' || $3::text || '%' ELSE TRUE END
ORDER BY
	CASE WHEN $4::text = 'name' THEN m_work_positions.name END ASC,
	CASE WHEN $4::text = 'r_name' THEN m_work_positions.name END DESC,
	m_permission_associations_pkey ASC
`

type GetWorkPositionsOnPermissionParams struct {
	PermissionID  uuid.UUID `json:"permission_id"`
	WhereLikeName bool      `json:"where_like_name"`
	SearchName    string    `json:"search_name"`
	OrderMethod   string    `json:"order_method"`
}

type GetWorkPositionsOnPermissionRow struct {
	PermissionAssociation PermissionAssociation `json:"permission_association"`
	WorkPosition          WorkPosition          `json:"work_position"`
}

func (q *Queries) GetWorkPositionsOnPermission(ctx context.Context, arg GetWorkPositionsOnPermissionParams) ([]GetWorkPositionsOnPermissionRow, error) {
	rows, err := q.db.Query(ctx, getWorkPositionsOnPermission,
		arg.PermissionID,
		arg.WhereLikeName,
		arg.SearchName,
		arg.OrderMethod,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetWorkPositionsOnPermissionRow{}
	for rows.Next() {
		var i GetWorkPositionsOnPermissionRow
		if err := rows.Scan(
			&i.PermissionAssociation.MPermissionAssociationsPkey,
			&i.PermissionAssociation.PermissionID,
			&i.PermissionAssociation.WorkPositionID,
			&i.WorkPosition.MWorkPositionsPkey,
			&i.WorkPosition.WorkPositionID,
			&i.WorkPosition.Name,
			&i.WorkPosition.Description,
			&i.WorkPosition.CreatedAt,
			&i.WorkPosition.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getWorkPositionsOnPermissionUseKeysetPaginate = `-- name: GetWorkPositionsOnPermissionUseKeysetPaginate :many
SELECT m_permission_associations.m_permission_associations_pkey, m_permission_associations.permission_id, m_permission_associations.work_position_id, m_work_positions.m_work_positions_pkey, m_work_positions.work_position_id, m_work_positions.name, m_work_positions.description, m_work_positions.created_at, m_work_positions.updated_at FROM m_permission_associations
LEFT JOIN m_work_positions ON m_permission_associations.work_position_id = m_work_positions.work_position_id
WHERE permission_id = $1
AND
	CASE WHEN $3::boolean = true THEN m_work_positions.name LIKE '%' || $4::text || '%' ELSE TRUE END
AND
	CASE $5::text
		WHEN 'next' THEN
			CASE $6::text
				WHEN 'name' THEN m_work_positions.name > $7 OR (m_work_positions.name = $7 AND m_permission_associations_pkey > $8::int)
				WHEN 'r_name' THEN m_work_positions.name < $7 OR (m_work_positions.name = $7 AND m_permission_associations_pkey > $8::int)
				ELSE m_permission_associations_pkey > $8::int
			END
		WHEN 'prev' THEN
			CASE $6::text
				WHEN 'name' THEN m_work_positions.name < $7 OR (m_work_positions.name = $7 AND m_permission_associations_pkey < $8::int)
				WHEN 'r_name' THEN m_work_positions.name > $7 OR (m_work_positions.name = $7 AND m_permission_associations_pkey < $8::int)
				ELSE m_permission_associations_pkey < $8::int
			END
	END
ORDER BY
	CASE WHEN $6::text = 'name' AND $5::text = 'next' THEN m_work_positions.name END ASC,
	CASE WHEN $6::text = 'name' AND $5::text = 'prev' THEN m_work_positions.name END DESC,
	CASE WHEN $6::text = 'r_name' AND $5::text = 'next' THEN m_work_positions.name END ASC,
	CASE WHEN $6::text = 'r_name' AND $5::text = 'prev' THEN m_work_positions.name END DESC,
	CASE WHEN $5::text = 'next' THEN m_permission_associations_pkey END ASC,
	CASE WHEN $5::text = 'prev' THEN m_permission_associations_pkey END DESC
LIMIT $2
`

type GetWorkPositionsOnPermissionUseKeysetPaginateParams struct {
	PermissionID    uuid.UUID `json:"permission_id"`
	Limit           int32     `json:"limit"`
	WhereLikeName   bool      `json:"where_like_name"`
	SearchName      string    `json:"search_name"`
	CursorDirection string    `json:"cursor_direction"`
	OrderMethod     string    `json:"order_method"`
	NameCursor      string    `json:"name_cursor"`
	Cursor          int32     `json:"cursor"`
}

type GetWorkPositionsOnPermissionUseKeysetPaginateRow struct {
	PermissionAssociation PermissionAssociation `json:"permission_association"`
	WorkPosition          WorkPosition          `json:"work_position"`
}

func (q *Queries) GetWorkPositionsOnPermissionUseKeysetPaginate(ctx context.Context, arg GetWorkPositionsOnPermissionUseKeysetPaginateParams) ([]GetWorkPositionsOnPermissionUseKeysetPaginateRow, error) {
	rows, err := q.db.Query(ctx, getWorkPositionsOnPermissionUseKeysetPaginate,
		arg.PermissionID,
		arg.Limit,
		arg.WhereLikeName,
		arg.SearchName,
		arg.CursorDirection,
		arg.OrderMethod,
		arg.NameCursor,
		arg.Cursor,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetWorkPositionsOnPermissionUseKeysetPaginateRow{}
	for rows.Next() {
		var i GetWorkPositionsOnPermissionUseKeysetPaginateRow
		if err := rows.Scan(
			&i.PermissionAssociation.MPermissionAssociationsPkey,
			&i.PermissionAssociation.PermissionID,
			&i.PermissionAssociation.WorkPositionID,
			&i.WorkPosition.MWorkPositionsPkey,
			&i.WorkPosition.WorkPositionID,
			&i.WorkPosition.Name,
			&i.WorkPosition.Description,
			&i.WorkPosition.CreatedAt,
			&i.WorkPosition.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getWorkPositionsOnPermissionUseNumberedPaginate = `-- name: GetWorkPositionsOnPermissionUseNumberedPaginate :many
SELECT m_permission_associations.m_permission_associations_pkey, m_permission_associations.permission_id, m_permission_associations.work_position_id, m_work_positions.m_work_positions_pkey, m_work_positions.work_position_id, m_work_positions.name, m_work_positions.description, m_work_positions.created_at, m_work_positions.updated_at FROM m_permission_associations
LEFT JOIN m_work_positions ON m_permission_associations.work_position_id = m_work_positions.work_position_id
WHERE permission_id = $1
AND
	CASE WHEN $4::boolean = true THEN m_work_positions.name LIKE '%' || $5::text || '%' ELSE TRUE END
ORDER BY
	CASE WHEN $6::text = 'name' THEN m_work_positions.name END ASC,
	CASE WHEN $6::text = 'r_name' THEN m_work_positions.name END DESC,
	m_permission_associations_pkey ASC
LIMIT $2 OFFSET $3
`

type GetWorkPositionsOnPermissionUseNumberedPaginateParams struct {
	PermissionID  uuid.UUID `json:"permission_id"`
	Limit         int32     `json:"limit"`
	Offset        int32     `json:"offset"`
	WhereLikeName bool      `json:"where_like_name"`
	SearchName    string    `json:"search_name"`
	OrderMethod   string    `json:"order_method"`
}

type GetWorkPositionsOnPermissionUseNumberedPaginateRow struct {
	PermissionAssociation PermissionAssociation `json:"permission_association"`
	WorkPosition          WorkPosition          `json:"work_position"`
}

func (q *Queries) GetWorkPositionsOnPermissionUseNumberedPaginate(ctx context.Context, arg GetWorkPositionsOnPermissionUseNumberedPaginateParams) ([]GetWorkPositionsOnPermissionUseNumberedPaginateRow, error) {
	rows, err := q.db.Query(ctx, getWorkPositionsOnPermissionUseNumberedPaginate,
		arg.PermissionID,
		arg.Limit,
		arg.Offset,
		arg.WhereLikeName,
		arg.SearchName,
		arg.OrderMethod,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetWorkPositionsOnPermissionUseNumberedPaginateRow{}
	for rows.Next() {
		var i GetWorkPositionsOnPermissionUseNumberedPaginateRow
		if err := rows.Scan(
			&i.PermissionAssociation.MPermissionAssociationsPkey,
			&i.PermissionAssociation.PermissionID,
			&i.PermissionAssociation.WorkPositionID,
			&i.WorkPosition.MWorkPositionsPkey,
			&i.WorkPosition.WorkPositionID,
			&i.WorkPosition.Name,
			&i.WorkPosition.Description,
			&i.WorkPosition.CreatedAt,
			&i.WorkPosition.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const pluralDeletePermissionAssociationsOnPermission = `-- name: PluralDeletePermissionAssociationsOnPermission :exec
DELETE FROM m_permission_associations WHERE permission_id = $1 AND work_position_id = ANY($2::uuid[])
`

type PluralDeletePermissionAssociationsOnPermissionParams struct {
	PermissionID uuid.UUID   `json:"permission_id"`
	Column2      []uuid.UUID `json:"column_2"`
}

func (q *Queries) PluralDeletePermissionAssociationsOnPermission(ctx context.Context, arg PluralDeletePermissionAssociationsOnPermissionParams) error {
	_, err := q.db.Exec(ctx, pluralDeletePermissionAssociationsOnPermission, arg.PermissionID, arg.Column2)
	return err
}

const pluralDeletePermissionAssociationsOnWorkPosition = `-- name: PluralDeletePermissionAssociationsOnWorkPosition :exec
DELETE FROM m_permission_associations WHERE work_position_id = $1 AND permission_id = ANY($2::uuid[])
`

type PluralDeletePermissionAssociationsOnWorkPositionParams struct {
	WorkPositionID uuid.UUID   `json:"work_position_id"`
	Column2        []uuid.UUID `json:"column_2"`
}

func (q *Queries) PluralDeletePermissionAssociationsOnWorkPosition(ctx context.Context, arg PluralDeletePermissionAssociationsOnWorkPositionParams) error {
	_, err := q.db.Exec(ctx, pluralDeletePermissionAssociationsOnWorkPosition, arg.WorkPositionID, arg.Column2)
	return err
}
