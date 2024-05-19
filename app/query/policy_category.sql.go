// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: policy_category.sql

package query

import (
	"context"

	"github.com/google/uuid"
)

const countPolicyCategories = `-- name: CountPolicyCategories :one
SELECT COUNT(*) FROM m_policy_categories
WHERE
	CASE WHEN $1::boolean = true THEN name LIKE '%' || $2::text || '%' ELSE TRUE END
`

type CountPolicyCategoriesParams struct {
	WhereLikeName bool   `json:"where_like_name"`
	SearchName    string `json:"search_name"`
}

func (q *Queries) CountPolicyCategories(ctx context.Context, arg CountPolicyCategoriesParams) (int64, error) {
	row := q.db.QueryRow(ctx, countPolicyCategories, arg.WhereLikeName, arg.SearchName)
	var count int64
	err := row.Scan(&count)
	return count, err
}

type CreatePolicyCategoriesParams struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Key         string `json:"key"`
}

const createPolicyCategory = `-- name: CreatePolicyCategory :one
INSERT INTO m_policy_categories (name, description, key) VALUES ($1, $2, $3) RETURNING m_policy_categories_pkey, policy_category_id, name, description, key
`

type CreatePolicyCategoryParams struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Key         string `json:"key"`
}

func (q *Queries) CreatePolicyCategory(ctx context.Context, arg CreatePolicyCategoryParams) (PolicyCategory, error) {
	row := q.db.QueryRow(ctx, createPolicyCategory, arg.Name, arg.Description, arg.Key)
	var i PolicyCategory
	err := row.Scan(
		&i.MPolicyCategoriesPkey,
		&i.PolicyCategoryID,
		&i.Name,
		&i.Description,
		&i.Key,
	)
	return i, err
}

const deletePolicyCategory = `-- name: DeletePolicyCategory :execrows
DELETE FROM m_policy_categories WHERE policy_category_id = $1
`

func (q *Queries) DeletePolicyCategory(ctx context.Context, policyCategoryID uuid.UUID) (int64, error) {
	result, err := q.db.Exec(ctx, deletePolicyCategory, policyCategoryID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const deletePolicyCategoryByKey = `-- name: DeletePolicyCategoryByKey :execrows
DELETE FROM m_policy_categories WHERE key = $1
`

func (q *Queries) DeletePolicyCategoryByKey(ctx context.Context, key string) (int64, error) {
	result, err := q.db.Exec(ctx, deletePolicyCategoryByKey, key)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const findPolicyCategoryByID = `-- name: FindPolicyCategoryByID :one
SELECT m_policy_categories_pkey, policy_category_id, name, description, key FROM m_policy_categories WHERE policy_category_id = $1
`

func (q *Queries) FindPolicyCategoryByID(ctx context.Context, policyCategoryID uuid.UUID) (PolicyCategory, error) {
	row := q.db.QueryRow(ctx, findPolicyCategoryByID, policyCategoryID)
	var i PolicyCategory
	err := row.Scan(
		&i.MPolicyCategoriesPkey,
		&i.PolicyCategoryID,
		&i.Name,
		&i.Description,
		&i.Key,
	)
	return i, err
}

const findPolicyCategoryByKey = `-- name: FindPolicyCategoryByKey :one
SELECT m_policy_categories_pkey, policy_category_id, name, description, key FROM m_policy_categories WHERE key = $1
`

func (q *Queries) FindPolicyCategoryByKey(ctx context.Context, key string) (PolicyCategory, error) {
	row := q.db.QueryRow(ctx, findPolicyCategoryByKey, key)
	var i PolicyCategory
	err := row.Scan(
		&i.MPolicyCategoriesPkey,
		&i.PolicyCategoryID,
		&i.Name,
		&i.Description,
		&i.Key,
	)
	return i, err
}

const getPluralPolicyCategories = `-- name: GetPluralPolicyCategories :many
SELECT m_policy_categories_pkey, policy_category_id, name, description, key FROM m_policy_categories
WHERE policy_category_id = ANY($1::uuid[])
ORDER BY
	CASE WHEN $2::text = 'name' THEN m_policy_categories.name END ASC,
	CASE WHEN $2::text = 'r_name' THEN m_policy_categories.name END DESC,
	m_policy_categories_pkey ASC
`

type GetPluralPolicyCategoriesParams struct {
	PolicyCategoryIds []uuid.UUID `json:"policy_category_ids"`
	OrderMethod       string      `json:"order_method"`
}

func (q *Queries) GetPluralPolicyCategories(ctx context.Context, arg GetPluralPolicyCategoriesParams) ([]PolicyCategory, error) {
	rows, err := q.db.Query(ctx, getPluralPolicyCategories, arg.PolicyCategoryIds, arg.OrderMethod)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []PolicyCategory{}
	for rows.Next() {
		var i PolicyCategory
		if err := rows.Scan(
			&i.MPolicyCategoriesPkey,
			&i.PolicyCategoryID,
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

const getPluralPolicyCategoriesUseNumberedPaginate = `-- name: GetPluralPolicyCategoriesUseNumberedPaginate :many
SELECT m_policy_categories_pkey, policy_category_id, name, description, key FROM m_policy_categories
WHERE policy_category_id = ANY($3::uuid[])
ORDER BY
	CASE WHEN $4::text = 'name' THEN m_policy_categories.name END ASC,
	CASE WHEN $4::text = 'r_name' THEN m_policy_categories.name END DESC,
	m_policy_categories_pkey ASC
LIMIT $1 OFFSET $2
`

type GetPluralPolicyCategoriesUseNumberedPaginateParams struct {
	Limit             int32       `json:"limit"`
	Offset            int32       `json:"offset"`
	PolicyCategoryIds []uuid.UUID `json:"policy_category_ids"`
	OrderMethod       string      `json:"order_method"`
}

func (q *Queries) GetPluralPolicyCategoriesUseNumberedPaginate(ctx context.Context, arg GetPluralPolicyCategoriesUseNumberedPaginateParams) ([]PolicyCategory, error) {
	rows, err := q.db.Query(ctx, getPluralPolicyCategoriesUseNumberedPaginate,
		arg.Limit,
		arg.Offset,
		arg.PolicyCategoryIds,
		arg.OrderMethod,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []PolicyCategory{}
	for rows.Next() {
		var i PolicyCategory
		if err := rows.Scan(
			&i.MPolicyCategoriesPkey,
			&i.PolicyCategoryID,
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

const getPolicyCategories = `-- name: GetPolicyCategories :many
SELECT m_policy_categories_pkey, policy_category_id, name, description, key FROM m_policy_categories
WHERE
	CASE WHEN $1::boolean = true THEN m_policy_categories.name LIKE '%' || $2::text || '%' ELSE TRUE END
ORDER BY
	CASE WHEN $3::text = 'name' THEN m_policy_categories.name END ASC,
	CASE WHEN $3::text = 'r_name' THEN m_policy_categories.name END DESC,
	m_policy_categories_pkey ASC
`

type GetPolicyCategoriesParams struct {
	WhereLikeName bool   `json:"where_like_name"`
	SearchName    string `json:"search_name"`
	OrderMethod   string `json:"order_method"`
}

func (q *Queries) GetPolicyCategories(ctx context.Context, arg GetPolicyCategoriesParams) ([]PolicyCategory, error) {
	rows, err := q.db.Query(ctx, getPolicyCategories, arg.WhereLikeName, arg.SearchName, arg.OrderMethod)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []PolicyCategory{}
	for rows.Next() {
		var i PolicyCategory
		if err := rows.Scan(
			&i.MPolicyCategoriesPkey,
			&i.PolicyCategoryID,
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

const getPolicyCategoriesUseKeysetPaginate = `-- name: GetPolicyCategoriesUseKeysetPaginate :many
SELECT m_policy_categories_pkey, policy_category_id, name, description, key FROM m_policy_categories
WHERE
	CASE WHEN $2::boolean = true THEN m_policy_categories.name LIKE '%' || $3::text || '%' ELSE TRUE END
AND
	CASE $4::text
		WHEN 'next' THEN
			CASE $5::text
				WHEN 'name' THEN name > $6 OR (name = $6 AND m_policy_categories_pkey > $7::int)
				WHEN 'r_name' THEN name < $6 OR (name = $6 AND m_policy_categories_pkey > $7::int)
				ELSE m_policy_categories_pkey > $7::int
			END
		WHEN 'prev' THEN
			CASE $5::text
				WHEN 'name' THEN name < $6 OR (name = $6 AND m_policy_categories_pkey < $7::int)
				WHEN 'r_name' THEN name > $6 OR (name = $6 AND m_policy_categories_pkey < $7::int)
				ELSE m_policy_categories_pkey < $7::int
			END
	END
ORDER BY
	CASE WHEN $5::text = 'name' AND $4::text = 'next' THEN m_policy_categories.name END ASC,
	CASE WHEN $5::text = 'name' AND $4::text = 'prev' THEN m_policy_categories.name END DESC,
	CASE WHEN $5::text = 'r_name' AND $4::text = 'next' THEN m_policy_categories.name END DESC,
	CASE WHEN $5::text = 'r_name' AND $4::text = 'prev' THEN m_policy_categories.name END ASC,
	CASE WHEN $4::text = 'next' THEN m_policy_categories_pkey END ASC,
	CASE WHEN $4::text = 'prev' THEN m_policy_categories_pkey END DESC
LIMIT $1
`

type GetPolicyCategoriesUseKeysetPaginateParams struct {
	Limit           int32  `json:"limit"`
	WhereLikeName   bool   `json:"where_like_name"`
	SearchName      string `json:"search_name"`
	CursorDirection string `json:"cursor_direction"`
	OrderMethod     string `json:"order_method"`
	NameCursor      string `json:"name_cursor"`
	Cursor          int32  `json:"cursor"`
}

func (q *Queries) GetPolicyCategoriesUseKeysetPaginate(ctx context.Context, arg GetPolicyCategoriesUseKeysetPaginateParams) ([]PolicyCategory, error) {
	rows, err := q.db.Query(ctx, getPolicyCategoriesUseKeysetPaginate,
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
	items := []PolicyCategory{}
	for rows.Next() {
		var i PolicyCategory
		if err := rows.Scan(
			&i.MPolicyCategoriesPkey,
			&i.PolicyCategoryID,
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

const getPolicyCategoriesUseNumberedPaginate = `-- name: GetPolicyCategoriesUseNumberedPaginate :many
SELECT m_policy_categories_pkey, policy_category_id, name, description, key FROM m_policy_categories
WHERE
	CASE WHEN $3::boolean = true THEN m_policy_categories.name LIKE '%' || $4::text || '%' ELSE TRUE END
ORDER BY
	CASE WHEN $5::text = 'name' THEN m_policy_categories.name END ASC,
	CASE WHEN $5::text = 'r_name' THEN m_policy_categories.name END DESC,
	m_policy_categories_pkey ASC
LIMIT $1 OFFSET $2
`

type GetPolicyCategoriesUseNumberedPaginateParams struct {
	Limit         int32  `json:"limit"`
	Offset        int32  `json:"offset"`
	WhereLikeName bool   `json:"where_like_name"`
	SearchName    string `json:"search_name"`
	OrderMethod   string `json:"order_method"`
}

func (q *Queries) GetPolicyCategoriesUseNumberedPaginate(ctx context.Context, arg GetPolicyCategoriesUseNumberedPaginateParams) ([]PolicyCategory, error) {
	rows, err := q.db.Query(ctx, getPolicyCategoriesUseNumberedPaginate,
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
	items := []PolicyCategory{}
	for rows.Next() {
		var i PolicyCategory
		if err := rows.Scan(
			&i.MPolicyCategoriesPkey,
			&i.PolicyCategoryID,
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

const pluralDeletePolicyCategories = `-- name: PluralDeletePolicyCategories :execrows
DELETE FROM m_policy_categories WHERE policy_category_id = ANY($1::uuid[])
`

func (q *Queries) PluralDeletePolicyCategories(ctx context.Context, dollar_1 []uuid.UUID) (int64, error) {
	result, err := q.db.Exec(ctx, pluralDeletePolicyCategories, dollar_1)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const updatePolicyCategory = `-- name: UpdatePolicyCategory :one
UPDATE m_policy_categories SET name = $2, description = $3, key = $4 WHERE policy_category_id = $1 RETURNING m_policy_categories_pkey, policy_category_id, name, description, key
`

type UpdatePolicyCategoryParams struct {
	PolicyCategoryID uuid.UUID `json:"policy_category_id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Key              string    `json:"key"`
}

func (q *Queries) UpdatePolicyCategory(ctx context.Context, arg UpdatePolicyCategoryParams) (PolicyCategory, error) {
	row := q.db.QueryRow(ctx, updatePolicyCategory,
		arg.PolicyCategoryID,
		arg.Name,
		arg.Description,
		arg.Key,
	)
	var i PolicyCategory
	err := row.Scan(
		&i.MPolicyCategoriesPkey,
		&i.PolicyCategoryID,
		&i.Name,
		&i.Description,
		&i.Key,
	)
	return i, err
}

const updatePolicyCategoryByKey = `-- name: UpdatePolicyCategoryByKey :one
UPDATE m_policy_categories SET name = $2, description = $3 WHERE key = $1 RETURNING m_policy_categories_pkey, policy_category_id, name, description, key
`

type UpdatePolicyCategoryByKeyParams struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (q *Queries) UpdatePolicyCategoryByKey(ctx context.Context, arg UpdatePolicyCategoryByKeyParams) (PolicyCategory, error) {
	row := q.db.QueryRow(ctx, updatePolicyCategoryByKey, arg.Key, arg.Name, arg.Description)
	var i PolicyCategory
	err := row.Scan(
		&i.MPolicyCategoriesPkey,
		&i.PolicyCategoryID,
		&i.Name,
		&i.Description,
		&i.Key,
	)
	return i, err
}
