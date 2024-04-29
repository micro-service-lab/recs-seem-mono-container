// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: role.sql

package query

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const countRoles = `-- name: CountRoles :one
SELECT COUNT(*) FROM m_roles
WHERE
	CASE WHEN $1::boolean = true THEN name LIKE '%' || $2::text || '%' ELSE TRUE END
`

type CountRolesParams struct {
	WhereLikeName bool   `json:"where_like_name"`
	SearchName    string `json:"search_name"`
}

func (q *Queries) CountRoles(ctx context.Context, arg CountRolesParams) (int64, error) {
	row := q.db.QueryRow(ctx, countRoles, arg.WhereLikeName, arg.SearchName)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createRole = `-- name: CreateRole :one
INSERT INTO m_roles (name, description, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING m_roles_pkey, role_id, name, description, created_at, updated_at
`

type CreateRoleParams struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (q *Queries) CreateRole(ctx context.Context, arg CreateRoleParams) (Role, error) {
	row := q.db.QueryRow(ctx, createRole,
		arg.Name,
		arg.Description,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i Role
	err := row.Scan(
		&i.MRolesPkey,
		&i.RoleID,
		&i.Name,
		&i.Description,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

type CreateRolesParams struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

const deleteRole = `-- name: DeleteRole :exec
DELETE FROM m_roles WHERE role_id = $1
`

func (q *Queries) DeleteRole(ctx context.Context, roleID uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteRole, roleID)
	return err
}

const findRoleByID = `-- name: FindRoleByID :one
SELECT m_roles_pkey, role_id, name, description, created_at, updated_at FROM m_roles WHERE role_id = $1
`

func (q *Queries) FindRoleByID(ctx context.Context, roleID uuid.UUID) (Role, error) {
	row := q.db.QueryRow(ctx, findRoleByID, roleID)
	var i Role
	err := row.Scan(
		&i.MRolesPkey,
		&i.RoleID,
		&i.Name,
		&i.Description,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getPluckRoles = `-- name: GetPluckRoles :many
SELECT role_id, name FROM m_roles
WHERE
	role_id = ANY($3::uuid[])
ORDER BY
	m_roles_pkey ASC
LIMIT $1 OFFSET $2
`

type GetPluckRolesParams struct {
	Limit   int32       `json:"limit"`
	Offset  int32       `json:"offset"`
	RoleIds []uuid.UUID `json:"role_ids"`
}

type GetPluckRolesRow struct {
	RoleID uuid.UUID `json:"role_id"`
	Name   string    `json:"name"`
}

func (q *Queries) GetPluckRoles(ctx context.Context, arg GetPluckRolesParams) ([]GetPluckRolesRow, error) {
	rows, err := q.db.Query(ctx, getPluckRoles, arg.Limit, arg.Offset, arg.RoleIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPluckRolesRow{}
	for rows.Next() {
		var i GetPluckRolesRow
		if err := rows.Scan(&i.RoleID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRoles = `-- name: GetRoles :many
SELECT m_roles_pkey, role_id, name, description, created_at, updated_at FROM m_roles
WHERE
	CASE WHEN $1::boolean = true THEN m_roles.name LIKE '%' || $2::text || '%' ELSE TRUE END
ORDER BY
	CASE WHEN $3::text = 'name' THEN m_roles.name END ASC,
	CASE WHEN $3::text = 'r_name' THEN m_roles.name END DESC,
	m_roles_pkey ASC
`

type GetRolesParams struct {
	WhereLikeName bool   `json:"where_like_name"`
	SearchName    string `json:"search_name"`
	OrderMethod   string `json:"order_method"`
}

func (q *Queries) GetRoles(ctx context.Context, arg GetRolesParams) ([]Role, error) {
	rows, err := q.db.Query(ctx, getRoles, arg.WhereLikeName, arg.SearchName, arg.OrderMethod)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Role{}
	for rows.Next() {
		var i Role
		if err := rows.Scan(
			&i.MRolesPkey,
			&i.RoleID,
			&i.Name,
			&i.Description,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const getRolesUseKeysetPaginate = `-- name: GetRolesUseKeysetPaginate :many
SELECT m_roles_pkey, role_id, name, description, created_at, updated_at FROM m_roles
WHERE
	CASE WHEN $2::boolean = true THEN m_roles.name LIKE '%' || $3::text || '%' ELSE TRUE END
AND
	CASE $4::text
		WHEN 'next' THEN
			CASE $5::text
				WHEN 'name' THEN name > $6 OR (name = $6 AND m_roles_pkey > $7::int)
				WHEN 'r_name' THEN name < $6 OR (name = $6 AND m_roles_pkey > $7::int)
				ELSE m_roles_pkey > $7::int
			END
		WHEN 'prev' THEN
			CASE $5::text
				WHEN 'name' THEN name < $6 OR (name = $6 AND m_roles_pkey < $7::int)
				WHEN 'r_name' THEN name > $6 OR (name = $6 AND m_roles_pkey < $7::int)
				ELSE m_roles_pkey < $7::int
			END
	END
ORDER BY
	CASE WHEN $5::text = 'name' THEN m_roles.name END ASC,
	CASE WHEN $5::text = 'r_name' THEN m_roles.name END DESC,
	m_roles_pkey ASC
LIMIT $1
`

type GetRolesUseKeysetPaginateParams struct {
	Limit           int32  `json:"limit"`
	WhereLikeName   bool   `json:"where_like_name"`
	SearchName      string `json:"search_name"`
	CursorDirection string `json:"cursor_direction"`
	OrderMethod     string `json:"order_method"`
	NameCursor      string `json:"name_cursor"`
	Cursor          int32  `json:"cursor"`
}

func (q *Queries) GetRolesUseKeysetPaginate(ctx context.Context, arg GetRolesUseKeysetPaginateParams) ([]Role, error) {
	rows, err := q.db.Query(ctx, getRolesUseKeysetPaginate,
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
	items := []Role{}
	for rows.Next() {
		var i Role
		if err := rows.Scan(
			&i.MRolesPkey,
			&i.RoleID,
			&i.Name,
			&i.Description,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const getRolesUseNumberedPaginate = `-- name: GetRolesUseNumberedPaginate :many
SELECT m_roles_pkey, role_id, name, description, created_at, updated_at FROM m_roles
WHERE
	CASE WHEN $3::boolean = true THEN m_roles.name LIKE '%' || $4::text || '%' ELSE TRUE END
ORDER BY
	CASE WHEN $5::text = 'name' THEN m_roles.name END ASC,
	CASE WHEN $5::text = 'r_name' THEN m_roles.name END DESC,
	m_roles_pkey ASC
LIMIT $1 OFFSET $2
`

type GetRolesUseNumberedPaginateParams struct {
	Limit         int32  `json:"limit"`
	Offset        int32  `json:"offset"`
	WhereLikeName bool   `json:"where_like_name"`
	SearchName    string `json:"search_name"`
	OrderMethod   string `json:"order_method"`
}

func (q *Queries) GetRolesUseNumberedPaginate(ctx context.Context, arg GetRolesUseNumberedPaginateParams) ([]Role, error) {
	rows, err := q.db.Query(ctx, getRolesUseNumberedPaginate,
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
	items := []Role{}
	for rows.Next() {
		var i Role
		if err := rows.Scan(
			&i.MRolesPkey,
			&i.RoleID,
			&i.Name,
			&i.Description,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const updateRole = `-- name: UpdateRole :one
UPDATE m_roles SET name = $2, description = $3, updated_at = $4 WHERE role_id = $1 RETURNING m_roles_pkey, role_id, name, description, created_at, updated_at
`

type UpdateRoleParams struct {
	RoleID      uuid.UUID `json:"role_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (q *Queries) UpdateRole(ctx context.Context, arg UpdateRoleParams) (Role, error) {
	row := q.db.QueryRow(ctx, updateRole,
		arg.RoleID,
		arg.Name,
		arg.Description,
		arg.UpdatedAt,
	)
	var i Role
	err := row.Scan(
		&i.MRolesPkey,
		&i.RoleID,
		&i.Name,
		&i.Description,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
