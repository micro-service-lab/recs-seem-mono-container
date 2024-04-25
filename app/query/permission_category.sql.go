// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: permission_category.sql

package query

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const countPermissionCategories = `-- name: CountPermissionCategories :one
SELECT COUNT(*) FROM m_permission_categories
WHERE
	CASE WHEN $1::boolean = true THEN name LIKE '%' || $2::text || '%' ELSE TRUE END
`

type CountPermissionCategoriesParams struct {
	WhereLikeName bool   `json:"where_like_name"`
	SearchName    string `json:"search_name"`
}

func (q *Queries) CountPermissionCategories(ctx context.Context, arg CountPermissionCategoriesParams) (int64, error) {
	row := q.db.QueryRow(ctx, countPermissionCategories, arg.WhereLikeName, arg.SearchName)
	var count int64
	err := row.Scan(&count)
	return count, err
}

type CreatePermissionCategoriesParams struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Key         string `json:"key"`
}

const createPermissionCategory = `-- name: CreatePermissionCategory :one
INSERT INTO m_permission_categories (name, description, key) VALUES ($1, $2, $3) RETURNING m_permission_categories_pkey, permission_category_id, name, description, key
`

type CreatePermissionCategoryParams struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Key         string `json:"key"`
}

func (q *Queries) CreatePermissionCategory(ctx context.Context, arg CreatePermissionCategoryParams) (PermissionCategory, error) {
	row := q.db.QueryRow(ctx, createPermissionCategory, arg.Name, arg.Description, arg.Key)
	var i PermissionCategory
	err := row.Scan(
		&i.MPermissionCategoriesPkey,
		&i.PermissionCategoryID,
		&i.Name,
		&i.Description,
		&i.Key,
	)
	return i, err
}

const deletePermissionCategory = `-- name: DeletePermissionCategory :exec
DELETE FROM m_permission_categories WHERE permission_category_id = $1
`

func (q *Queries) DeletePermissionCategory(ctx context.Context, permissionCategoryID uuid.UUID) error {
	_, err := q.db.Exec(ctx, deletePermissionCategory, permissionCategoryID)
	return err
}

const deletePermissionCategoryByKey = `-- name: DeletePermissionCategoryByKey :exec
DELETE FROM m_permission_categories WHERE key = $1
`

func (q *Queries) DeletePermissionCategoryByKey(ctx context.Context, key string) error {
	_, err := q.db.Exec(ctx, deletePermissionCategoryByKey, key)
	return err
}

const findPermissionCategoryByID = `-- name: FindPermissionCategoryByID :one
SELECT m_permission_categories_pkey, permission_category_id, name, description, key FROM m_permission_categories WHERE permission_category_id = $1
`

func (q *Queries) FindPermissionCategoryByID(ctx context.Context, permissionCategoryID uuid.UUID) (PermissionCategory, error) {
	row := q.db.QueryRow(ctx, findPermissionCategoryByID, permissionCategoryID)
	var i PermissionCategory
	err := row.Scan(
		&i.MPermissionCategoriesPkey,
		&i.PermissionCategoryID,
		&i.Name,
		&i.Description,
		&i.Key,
	)
	return i, err
}

const findPermissionCategoryByKey = `-- name: FindPermissionCategoryByKey :one
SELECT m_permission_categories_pkey, permission_category_id, name, description, key FROM m_permission_categories WHERE key = $1
`

func (q *Queries) FindPermissionCategoryByKey(ctx context.Context, key string) (PermissionCategory, error) {
	row := q.db.QueryRow(ctx, findPermissionCategoryByKey, key)
	var i PermissionCategory
	err := row.Scan(
		&i.MPermissionCategoriesPkey,
		&i.PermissionCategoryID,
		&i.Name,
		&i.Description,
		&i.Key,
	)
	return i, err
}

const getPermissionCategories = `-- name: GetPermissionCategories :many
SELECT m_permission_categories_pkey, permission_category_id, name, description, key FROM m_permission_categories
WHERE
	CASE WHEN $1::boolean = true THEN m_permission_categories.name LIKE '%' || $2::text || '%' ELSE TRUE END
ORDER BY
	CASE WHEN $3::text = 'name' THEN m_permission_categories.name END ASC,
	CASE WHEN $3::text = 'r_name' THEN m_permission_categories.name END DESC,
	m_permission_categories_pkey DESC
`

type GetPermissionCategoriesParams struct {
	WhereLikeName bool   `json:"where_like_name"`
	SearchName    string `json:"search_name"`
	OrderMethod   string `json:"order_method"`
}

func (q *Queries) GetPermissionCategories(ctx context.Context, arg GetPermissionCategoriesParams) ([]PermissionCategory, error) {
	rows, err := q.db.Query(ctx, getPermissionCategories, arg.WhereLikeName, arg.SearchName, arg.OrderMethod)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []PermissionCategory{}
	for rows.Next() {
		var i PermissionCategory
		if err := rows.Scan(
			&i.MPermissionCategoriesPkey,
			&i.PermissionCategoryID,
			&i.Name,
			&i.Description,
			&i.Key,
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

const getPermissionCategoriesUseKeysetPaginate = `-- name: GetPermissionCategoriesUseKeysetPaginate :many
SELECT m_permission_categories_pkey, permission_category_id, name, description, key FROM m_permission_categories
WHERE
	CASE WHEN $2::boolean = true THEN m_permission_categories.name LIKE '%' || $3::text || '%' ELSE TRUE END
AND
	CASE $4
		WHEN 'next' THEN
			CASE $5::text
				WHEN 'name' THEN name > $6 OR (name = $6 AND m_permission_categories_pkey < $7)
				WHEN 'r_name' THEN name < $6 OR (name = $6 AND m_permission_categories_pkey < $7)
				ELSE m_permission_categories_pkey < $7
			END
		WHEN 'prev' THEN
			CASE $5::text
				WHEN 'name' THEN name < $6 OR (name = $6 AND m_permission_categories_pkey > $7)
				WHEN 'r_name' THEN name > $6 OR (name = $6 AND m_permission_categories_pkey > $7)
				ELSE m_permission_categories_pkey > $7
			END
	END
ORDER BY
	CASE WHEN $5::text = 'name' THEN m_permission_categories.name END ASC,
	CASE WHEN $5::text = 'r_name' THEN m_permission_categories.name END DESC,
	m_permission_categories_pkey DESC
LIMIT $1
`

type GetPermissionCategoriesUseKeysetPaginateParams struct {
	Limit           int32       `json:"limit"`
	WhereLikeName   bool        `json:"where_like_name"`
	SearchName      string      `json:"search_name"`
	CursorDirection interface{} `json:"cursor_direction"`
	OrderMethod     string      `json:"order_method"`
	CursorColumn    string      `json:"cursor_column"`
	Cursor          pgtype.Int8 `json:"cursor"`
}

func (q *Queries) GetPermissionCategoriesUseKeysetPaginate(ctx context.Context, arg GetPermissionCategoriesUseKeysetPaginateParams) ([]PermissionCategory, error) {
	rows, err := q.db.Query(ctx, getPermissionCategoriesUseKeysetPaginate,
		arg.Limit,
		arg.WhereLikeName,
		arg.SearchName,
		arg.CursorDirection,
		arg.OrderMethod,
		arg.CursorColumn,
		arg.Cursor,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []PermissionCategory{}
	for rows.Next() {
		var i PermissionCategory
		if err := rows.Scan(
			&i.MPermissionCategoriesPkey,
			&i.PermissionCategoryID,
			&i.Name,
			&i.Description,
			&i.Key,
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

const getPermissionCategoriesUseNumberedPaginate = `-- name: GetPermissionCategoriesUseNumberedPaginate :many
SELECT m_permission_categories_pkey, permission_category_id, name, description, key FROM m_permission_categories
WHERE
	CASE WHEN $3::boolean = true THEN m_permission_categories.name LIKE '%' || $4::text || '%' ELSE TRUE END
ORDER BY
	CASE WHEN $5::text = 'name' THEN m_permission_categories.name END ASC,
	CASE WHEN $5::text = 'r_name' THEN m_permission_categories.name END DESC,
	m_permission_categories_pkey DESC
LIMIT $1 OFFSET $2
`

type GetPermissionCategoriesUseNumberedPaginateParams struct {
	Limit         int32  `json:"limit"`
	Offset        int32  `json:"offset"`
	WhereLikeName bool   `json:"where_like_name"`
	SearchName    string `json:"search_name"`
	OrderMethod   string `json:"order_method"`
}

func (q *Queries) GetPermissionCategoriesUseNumberedPaginate(ctx context.Context, arg GetPermissionCategoriesUseNumberedPaginateParams) ([]PermissionCategory, error) {
	rows, err := q.db.Query(ctx, getPermissionCategoriesUseNumberedPaginate,
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
	items := []PermissionCategory{}
	for rows.Next() {
		var i PermissionCategory
		if err := rows.Scan(
			&i.MPermissionCategoriesPkey,
			&i.PermissionCategoryID,
			&i.Name,
			&i.Description,
			&i.Key,
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

const updatePermissionCategory = `-- name: UpdatePermissionCategory :one
UPDATE m_permission_categories SET name = $2, description = $3, key = $4 WHERE permission_category_id = $1 RETURNING m_permission_categories_pkey, permission_category_id, name, description, key
`

type UpdatePermissionCategoryParams struct {
	PermissionCategoryID uuid.UUID `json:"permission_category_id"`
	Name                 string    `json:"name"`
	Description          string    `json:"description"`
	Key                  string    `json:"key"`
}

func (q *Queries) UpdatePermissionCategory(ctx context.Context, arg UpdatePermissionCategoryParams) (PermissionCategory, error) {
	row := q.db.QueryRow(ctx, updatePermissionCategory,
		arg.PermissionCategoryID,
		arg.Name,
		arg.Description,
		arg.Key,
	)
	var i PermissionCategory
	err := row.Scan(
		&i.MPermissionCategoriesPkey,
		&i.PermissionCategoryID,
		&i.Name,
		&i.Description,
		&i.Key,
	)
	return i, err
}

const updatePermissionCategoryByKey = `-- name: UpdatePermissionCategoryByKey :one
UPDATE m_permission_categories SET name = $2, description = $3 WHERE key = $1 RETURNING m_permission_categories_pkey, permission_category_id, name, description, key
`

type UpdatePermissionCategoryByKeyParams struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (q *Queries) UpdatePermissionCategoryByKey(ctx context.Context, arg UpdatePermissionCategoryByKeyParams) (PermissionCategory, error) {
	row := q.db.QueryRow(ctx, updatePermissionCategoryByKey, arg.Key, arg.Name, arg.Description)
	var i PermissionCategory
	err := row.Scan(
		&i.MPermissionCategoriesPkey,
		&i.PermissionCategoryID,
		&i.Name,
		&i.Description,
		&i.Key,
	)
	return i, err
}
