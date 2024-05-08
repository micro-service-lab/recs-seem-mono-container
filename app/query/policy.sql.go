// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: policy.sql

package query

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const countPolicies = `-- name: CountPolicies :one
SELECT COUNT(*) FROM m_policies
WHERE
	CASE WHEN $1::boolean = true THEN name LIKE '%' || $2::text || '%' ELSE TRUE END
AND
	CASE WHEN $3::boolean = true THEN policy_category_id = ANY($4::uuid[]) ELSE TRUE END
`

type CountPoliciesParams struct {
	WhereLikeName   bool        `json:"where_like_name"`
	SearchName      string      `json:"search_name"`
	WhereInCategory bool        `json:"where_in_category"`
	InCategories    []uuid.UUID `json:"in_categories"`
}

func (q *Queries) CountPolicies(ctx context.Context, arg CountPoliciesParams) (int64, error) {
	row := q.db.QueryRow(ctx, countPolicies,
		arg.WhereLikeName,
		arg.SearchName,
		arg.WhereInCategory,
		arg.InCategories,
	)
	var count int64
	err := row.Scan(&count)
	return count, err
}

type CreatePoliciesParams struct {
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Key              string    `json:"key"`
	PolicyCategoryID uuid.UUID `json:"policy_category_id"`
}

const createPolicy = `-- name: CreatePolicy :one
INSERT INTO m_policies (name, description, key, policy_category_id) VALUES ($1, $2, $3, $4) RETURNING m_policies_pkey, policy_id, name, description, key, policy_category_id
`

type CreatePolicyParams struct {
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Key              string    `json:"key"`
	PolicyCategoryID uuid.UUID `json:"policy_category_id"`
}

func (q *Queries) CreatePolicy(ctx context.Context, arg CreatePolicyParams) (Policy, error) {
	row := q.db.QueryRow(ctx, createPolicy,
		arg.Name,
		arg.Description,
		arg.Key,
		arg.PolicyCategoryID,
	)
	var i Policy
	err := row.Scan(
		&i.MPoliciesPkey,
		&i.PolicyID,
		&i.Name,
		&i.Description,
		&i.Key,
		&i.PolicyCategoryID,
	)
	return i, err
}

const deletePolicy = `-- name: DeletePolicy :exec
DELETE FROM m_policies WHERE policy_id = $1
`

func (q *Queries) DeletePolicy(ctx context.Context, policyID uuid.UUID) error {
	_, err := q.db.Exec(ctx, deletePolicy, policyID)
	return err
}

const deletePolicyByKey = `-- name: DeletePolicyByKey :exec
DELETE FROM m_policies WHERE key = $1
`

func (q *Queries) DeletePolicyByKey(ctx context.Context, key string) error {
	_, err := q.db.Exec(ctx, deletePolicyByKey, key)
	return err
}

const findPolicyByID = `-- name: FindPolicyByID :one
SELECT m_policies_pkey, policy_id, name, description, key, policy_category_id FROM m_policies WHERE policy_id = $1
`

func (q *Queries) FindPolicyByID(ctx context.Context, policyID uuid.UUID) (Policy, error) {
	row := q.db.QueryRow(ctx, findPolicyByID, policyID)
	var i Policy
	err := row.Scan(
		&i.MPoliciesPkey,
		&i.PolicyID,
		&i.Name,
		&i.Description,
		&i.Key,
		&i.PolicyCategoryID,
	)
	return i, err
}

const findPolicyByIDWithCategory = `-- name: FindPolicyByIDWithCategory :one
SELECT m_policies.m_policies_pkey, m_policies.policy_id, m_policies.name, m_policies.description, m_policies.key, m_policies.policy_category_id, m_policy_categories.m_policy_categories_pkey, m_policy_categories.policy_category_id, m_policy_categories.name, m_policy_categories.description, m_policy_categories.key FROM m_policies
JOIN m_policy_categories ON m_policies.policy_category_id = m_policy_categories.policy_category_id
WHERE m_policies.policy_id = $1
`

type FindPolicyByIDWithCategoryRow struct {
	MPoliciesPkey         pgtype.Int8 `json:"m_policies_pkey"`
	PolicyID              uuid.UUID   `json:"policy_id"`
	Name                  string      `json:"name"`
	Description           string      `json:"description"`
	Key                   string      `json:"key"`
	PolicyCategoryID      uuid.UUID   `json:"policy_category_id"`
	MPolicyCategoriesPkey pgtype.Int8 `json:"m_policy_categories_pkey"`
	PolicyCategoryID_2    uuid.UUID   `json:"policy_category_id_2"`
	Name_2                string      `json:"name_2"`
	Description_2         string      `json:"description_2"`
	Key_2                 string      `json:"key_2"`
}

func (q *Queries) FindPolicyByIDWithCategory(ctx context.Context, policyID uuid.UUID) (FindPolicyByIDWithCategoryRow, error) {
	row := q.db.QueryRow(ctx, findPolicyByIDWithCategory, policyID)
	var i FindPolicyByIDWithCategoryRow
	err := row.Scan(
		&i.MPoliciesPkey,
		&i.PolicyID,
		&i.Name,
		&i.Description,
		&i.Key,
		&i.PolicyCategoryID,
		&i.MPolicyCategoriesPkey,
		&i.PolicyCategoryID_2,
		&i.Name_2,
		&i.Description_2,
		&i.Key_2,
	)
	return i, err
}

const findPolicyByKey = `-- name: FindPolicyByKey :one
SELECT m_policies_pkey, policy_id, name, description, key, policy_category_id FROM m_policies WHERE key = $1
`

func (q *Queries) FindPolicyByKey(ctx context.Context, key string) (Policy, error) {
	row := q.db.QueryRow(ctx, findPolicyByKey, key)
	var i Policy
	err := row.Scan(
		&i.MPoliciesPkey,
		&i.PolicyID,
		&i.Name,
		&i.Description,
		&i.Key,
		&i.PolicyCategoryID,
	)
	return i, err
}

const getPluralPolicies = `-- name: GetPluralPolicies :many
SELECT m_policies_pkey, policy_id, name, description, key, policy_category_id FROM m_policies WHERE policy_id = ANY($3::uuid[])
ORDER BY
	m_policies_pkey ASC
LIMIT $1 OFFSET $2
`

type GetPluralPoliciesParams struct {
	Limit     int32       `json:"limit"`
	Offset    int32       `json:"offset"`
	PolicyIds []uuid.UUID `json:"policy_ids"`
}

func (q *Queries) GetPluralPolicies(ctx context.Context, arg GetPluralPoliciesParams) ([]Policy, error) {
	rows, err := q.db.Query(ctx, getPluralPolicies, arg.Limit, arg.Offset, arg.PolicyIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Policy{}
	for rows.Next() {
		var i Policy
		if err := rows.Scan(
			&i.MPoliciesPkey,
			&i.PolicyID,
			&i.Name,
			&i.Description,
			&i.Key,
			&i.PolicyCategoryID,
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

const getPluralPoliciesWithCategory = `-- name: GetPluralPoliciesWithCategory :many
SELECT m_policies.m_policies_pkey, m_policies.policy_id, m_policies.name, m_policies.description, m_policies.key, m_policies.policy_category_id, m_policy_categories.m_policy_categories_pkey, m_policy_categories.policy_category_id, m_policy_categories.name, m_policy_categories.description, m_policy_categories.key FROM m_policies
JOIN m_policy_categories ON m_policies.policy_category_id = m_policy_categories.policy_category_id
WHERE policy_id = ANY($3::uuid[])
ORDER BY
	m_policies_pkey ASC
LIMIT $1 OFFSET $2
`

type GetPluralPoliciesWithCategoryParams struct {
	Limit     int32       `json:"limit"`
	Offset    int32       `json:"offset"`
	PolicyIds []uuid.UUID `json:"policy_ids"`
}

type GetPluralPoliciesWithCategoryRow struct {
	MPoliciesPkey         pgtype.Int8 `json:"m_policies_pkey"`
	PolicyID              uuid.UUID   `json:"policy_id"`
	Name                  string      `json:"name"`
	Description           string      `json:"description"`
	Key                   string      `json:"key"`
	PolicyCategoryID      uuid.UUID   `json:"policy_category_id"`
	MPolicyCategoriesPkey pgtype.Int8 `json:"m_policy_categories_pkey"`
	PolicyCategoryID_2    uuid.UUID   `json:"policy_category_id_2"`
	Name_2                string      `json:"name_2"`
	Description_2         string      `json:"description_2"`
	Key_2                 string      `json:"key_2"`
}

func (q *Queries) GetPluralPoliciesWithCategory(ctx context.Context, arg GetPluralPoliciesWithCategoryParams) ([]GetPluralPoliciesWithCategoryRow, error) {
	rows, err := q.db.Query(ctx, getPluralPoliciesWithCategory, arg.Limit, arg.Offset, arg.PolicyIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPluralPoliciesWithCategoryRow{}
	for rows.Next() {
		var i GetPluralPoliciesWithCategoryRow
		if err := rows.Scan(
			&i.MPoliciesPkey,
			&i.PolicyID,
			&i.Name,
			&i.Description,
			&i.Key,
			&i.PolicyCategoryID,
			&i.MPolicyCategoriesPkey,
			&i.PolicyCategoryID_2,
			&i.Name_2,
			&i.Description_2,
			&i.Key_2,
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

const getPolicies = `-- name: GetPolicies :many
SELECT m_policies_pkey, policy_id, name, description, key, policy_category_id FROM m_policies
WHERE
	CASE WHEN $1::boolean = true THEN m_policies.name LIKE '%' || $2::text || '%' ELSE TRUE END
AND
	CASE WHEN $3::boolean = true THEN policy_category_id = ANY($4::uuid[]) ELSE TRUE END
ORDER BY
	CASE WHEN $5::text = 'name' THEN m_policies.name END ASC,
	CASE WHEN $5::text = 'r_name' THEN m_policies.name END DESC,
	m_policies_pkey ASC
`

type GetPoliciesParams struct {
	WhereLikeName   bool        `json:"where_like_name"`
	SearchName      string      `json:"search_name"`
	WhereInCategory bool        `json:"where_in_category"`
	InCategories    []uuid.UUID `json:"in_categories"`
	OrderMethod     string      `json:"order_method"`
}

func (q *Queries) GetPolicies(ctx context.Context, arg GetPoliciesParams) ([]Policy, error) {
	rows, err := q.db.Query(ctx, getPolicies,
		arg.WhereLikeName,
		arg.SearchName,
		arg.WhereInCategory,
		arg.InCategories,
		arg.OrderMethod,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Policy{}
	for rows.Next() {
		var i Policy
		if err := rows.Scan(
			&i.MPoliciesPkey,
			&i.PolicyID,
			&i.Name,
			&i.Description,
			&i.Key,
			&i.PolicyCategoryID,
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

const getPoliciesUseKeysetPaginate = `-- name: GetPoliciesUseKeysetPaginate :many
SELECT m_policies_pkey, policy_id, name, description, key, policy_category_id FROM m_policies
WHERE
	CASE WHEN $2::boolean = true THEN m_policies.name LIKE '%' || $3::text || '%' ELSE TRUE END
AND
	CASE WHEN $4::boolean = true THEN policy_category_id = ANY($5::uuid[]) ELSE TRUE END
AND
	CASE $6::text
		WHEN 'next' THEN
			CASE $7::text
				WHEN 'name' THEN m_policies.name > $8 OR (m_policies.name = $8 AND m_policies_pkey > $9::int)
				WHEN 'r_name' THEN m_policies.name < $8 OR (m_policies.name = $8 AND m_policies_pkey > $9::int)
				ELSE m_policies_pkey > $9::int
			END
		WHEN 'prev' THEN
			CASE $7::text
				WHEN 'name' THEN m_policies.name < $8 OR (m_policies.name = $8 AND m_policies_pkey < $9::int)
				WHEN 'r_name' THEN m_policies.name > $8 OR (m_policies.name = $8 AND m_policies_pkey < $9::int)
				ELSE m_policies_pkey < $9::int
			END
	END
ORDER BY
	CASE WHEN $7::text = 'name' AND $6::text = 'next' THEN m_policies.name END ASC,
	CASE WHEN $7::text = 'name' AND $6::text = 'prev' THEN m_policies.name END DESC,
	CASE WHEN $7::text = 'r_name' AND $6::text = 'next' THEN m_policies.name END ASC,
	CASE WHEN $7::text = 'r_name' AND $6::text = 'prev' THEN m_policies.name END DESC,
	CASE WHEN $6::text = 'next' THEN m_policies_pkey END ASC,
	CASE WHEN $6::text = 'prev' THEN m_policies_pkey END DESC
LIMIT $1
`

type GetPoliciesUseKeysetPaginateParams struct {
	Limit           int32       `json:"limit"`
	WhereLikeName   bool        `json:"where_like_name"`
	SearchName      string      `json:"search_name"`
	WhereInCategory bool        `json:"where_in_category"`
	InCategories    []uuid.UUID `json:"in_categories"`
	CursorDirection string      `json:"cursor_direction"`
	OrderMethod     string      `json:"order_method"`
	NameCursor      string      `json:"name_cursor"`
	Cursor          int32       `json:"cursor"`
}

func (q *Queries) GetPoliciesUseKeysetPaginate(ctx context.Context, arg GetPoliciesUseKeysetPaginateParams) ([]Policy, error) {
	rows, err := q.db.Query(ctx, getPoliciesUseKeysetPaginate,
		arg.Limit,
		arg.WhereLikeName,
		arg.SearchName,
		arg.WhereInCategory,
		arg.InCategories,
		arg.CursorDirection,
		arg.OrderMethod,
		arg.NameCursor,
		arg.Cursor,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Policy{}
	for rows.Next() {
		var i Policy
		if err := rows.Scan(
			&i.MPoliciesPkey,
			&i.PolicyID,
			&i.Name,
			&i.Description,
			&i.Key,
			&i.PolicyCategoryID,
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

const getPoliciesUseNumberedPaginate = `-- name: GetPoliciesUseNumberedPaginate :many
SELECT m_policies_pkey, policy_id, name, description, key, policy_category_id FROM m_policies
WHERE
	CASE WHEN $3::boolean = true THEN m_policies.name LIKE '%' || $4::text || '%' ELSE TRUE END
AND
	CASE WHEN $5::boolean = true THEN policy_category_id = ANY($6::uuid[]) ELSE TRUE END
ORDER BY
	CASE WHEN $7::text = 'name' THEN m_policies.name END ASC,
	CASE WHEN $7::text = 'r_name' THEN m_policies.name END DESC,
	m_policies_pkey ASC
LIMIT $1 OFFSET $2
`

type GetPoliciesUseNumberedPaginateParams struct {
	Limit           int32       `json:"limit"`
	Offset          int32       `json:"offset"`
	WhereLikeName   bool        `json:"where_like_name"`
	SearchName      string      `json:"search_name"`
	WhereInCategory bool        `json:"where_in_category"`
	InCategories    []uuid.UUID `json:"in_categories"`
	OrderMethod     string      `json:"order_method"`
}

func (q *Queries) GetPoliciesUseNumberedPaginate(ctx context.Context, arg GetPoliciesUseNumberedPaginateParams) ([]Policy, error) {
	rows, err := q.db.Query(ctx, getPoliciesUseNumberedPaginate,
		arg.Limit,
		arg.Offset,
		arg.WhereLikeName,
		arg.SearchName,
		arg.WhereInCategory,
		arg.InCategories,
		arg.OrderMethod,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Policy{}
	for rows.Next() {
		var i Policy
		if err := rows.Scan(
			&i.MPoliciesPkey,
			&i.PolicyID,
			&i.Name,
			&i.Description,
			&i.Key,
			&i.PolicyCategoryID,
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

const getPoliciesWithCategory = `-- name: GetPoliciesWithCategory :many
SELECT m_policies.m_policies_pkey, m_policies.policy_id, m_policies.name, m_policies.description, m_policies.key, m_policies.policy_category_id, m_policy_categories.m_policy_categories_pkey, m_policy_categories.policy_category_id, m_policy_categories.name, m_policy_categories.description, m_policy_categories.key FROM m_policies
JOIN m_policy_categories ON m_policies.policy_category_id = m_policy_categories.policy_category_id
WHERE
	CASE WHEN $1::boolean = true THEN m_policies.name LIKE '%' || $2::text || '%' ELSE TRUE END
AND
	CASE WHEN $3::boolean = true THEN policy_category_id = ANY($4::uuid[]) ELSE TRUE END
ORDER BY
	CASE WHEN $5::text = 'name' THEN m_policies.name END ASC,
	CASE WHEN $5::text = 'r_name' THEN m_policies.name END DESC,
	m_policies_pkey ASC
`

type GetPoliciesWithCategoryParams struct {
	WhereLikeName   bool        `json:"where_like_name"`
	SearchName      string      `json:"search_name"`
	WhereInCategory bool        `json:"where_in_category"`
	InCategories    []uuid.UUID `json:"in_categories"`
	OrderMethod     string      `json:"order_method"`
}

type GetPoliciesWithCategoryRow struct {
	MPoliciesPkey         pgtype.Int8 `json:"m_policies_pkey"`
	PolicyID              uuid.UUID   `json:"policy_id"`
	Name                  string      `json:"name"`
	Description           string      `json:"description"`
	Key                   string      `json:"key"`
	PolicyCategoryID      uuid.UUID   `json:"policy_category_id"`
	MPolicyCategoriesPkey pgtype.Int8 `json:"m_policy_categories_pkey"`
	PolicyCategoryID_2    uuid.UUID   `json:"policy_category_id_2"`
	Name_2                string      `json:"name_2"`
	Description_2         string      `json:"description_2"`
	Key_2                 string      `json:"key_2"`
}

func (q *Queries) GetPoliciesWithCategory(ctx context.Context, arg GetPoliciesWithCategoryParams) ([]GetPoliciesWithCategoryRow, error) {
	rows, err := q.db.Query(ctx, getPoliciesWithCategory,
		arg.WhereLikeName,
		arg.SearchName,
		arg.WhereInCategory,
		arg.InCategories,
		arg.OrderMethod,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPoliciesWithCategoryRow{}
	for rows.Next() {
		var i GetPoliciesWithCategoryRow
		if err := rows.Scan(
			&i.MPoliciesPkey,
			&i.PolicyID,
			&i.Name,
			&i.Description,
			&i.Key,
			&i.PolicyCategoryID,
			&i.MPolicyCategoriesPkey,
			&i.PolicyCategoryID_2,
			&i.Name_2,
			&i.Description_2,
			&i.Key_2,
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

const getPoliciesWithCategoryUseKeysetPaginate = `-- name: GetPoliciesWithCategoryUseKeysetPaginate :many
SELECT m_policies.m_policies_pkey, m_policies.policy_id, m_policies.name, m_policies.description, m_policies.key, m_policies.policy_category_id, m_policy_categories.m_policy_categories_pkey, m_policy_categories.policy_category_id, m_policy_categories.name, m_policy_categories.description, m_policy_categories.key FROM m_policies
JOIN m_policy_categories ON m_policies.policy_category_id = m_policy_categories.policy_category_id
WHERE
	CASE WHEN $2::boolean = true THEN m_policies.name LIKE '%' || $3::text || '%' ELSE TRUE END
AND
	CASE WHEN $4::boolean = true THEN policy_category_id = ANY($5::uuid[]) ELSE TRUE END
AND
	CASE $6::text
		WHEN 'next' THEN
			CASE $7::text
				WHEN 'name' THEN m_policies.name > $8 OR (m_policies.name = $8 AND m_policies_pkey > $9::int)
				WHEN 'r_name' THEN m_policies.name < $8 OR (m_policies.name = $8 AND m_policies_pkey > $9::int)
				ELSE m_policies_pkey > $9::int
			END
		WHEN 'prev' THEN
			CASE $7::text
				WHEN 'name' THEN m_policies.name < $8 OR (m_policies.name = $8 AND m_policies_pkey < $9::int)
				WHEN 'r_name' THEN m_policies.name > $8 OR (m_policies.name = $8 AND m_policies_pkey < $9::int)
				ELSE m_policies_pkey < $9::int
			END
	END
ORDER BY
	CASE WHEN $7::text = 'name' AND $6::text = 'next' THEN m_policies.name END ASC,
	CASE WHEN $7::text = 'name' AND $6::text = 'prev' THEN m_policies.name END DESC,
	CASE WHEN $7::text = 'r_name' AND $6::text = 'next' THEN m_policies.name END ASC,
	CASE WHEN $7::text = 'r_name' AND $6::text = 'prev' THEN m_policies.name END DESC,
	CASE WHEN $6::text = 'next' THEN m_policies_pkey END ASC,
	CASE WHEN $6::text = 'prev' THEN m_policies_pkey END DESC
LIMIT $1
`

type GetPoliciesWithCategoryUseKeysetPaginateParams struct {
	Limit           int32       `json:"limit"`
	WhereLikeName   bool        `json:"where_like_name"`
	SearchName      string      `json:"search_name"`
	WhereInCategory bool        `json:"where_in_category"`
	InCategories    []uuid.UUID `json:"in_categories"`
	CursorDirection string      `json:"cursor_direction"`
	OrderMethod     string      `json:"order_method"`
	NameCursor      string      `json:"name_cursor"`
	Cursor          int32       `json:"cursor"`
}

type GetPoliciesWithCategoryUseKeysetPaginateRow struct {
	MPoliciesPkey         pgtype.Int8 `json:"m_policies_pkey"`
	PolicyID              uuid.UUID   `json:"policy_id"`
	Name                  string      `json:"name"`
	Description           string      `json:"description"`
	Key                   string      `json:"key"`
	PolicyCategoryID      uuid.UUID   `json:"policy_category_id"`
	MPolicyCategoriesPkey pgtype.Int8 `json:"m_policy_categories_pkey"`
	PolicyCategoryID_2    uuid.UUID   `json:"policy_category_id_2"`
	Name_2                string      `json:"name_2"`
	Description_2         string      `json:"description_2"`
	Key_2                 string      `json:"key_2"`
}

func (q *Queries) GetPoliciesWithCategoryUseKeysetPaginate(ctx context.Context, arg GetPoliciesWithCategoryUseKeysetPaginateParams) ([]GetPoliciesWithCategoryUseKeysetPaginateRow, error) {
	rows, err := q.db.Query(ctx, getPoliciesWithCategoryUseKeysetPaginate,
		arg.Limit,
		arg.WhereLikeName,
		arg.SearchName,
		arg.WhereInCategory,
		arg.InCategories,
		arg.CursorDirection,
		arg.OrderMethod,
		arg.NameCursor,
		arg.Cursor,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPoliciesWithCategoryUseKeysetPaginateRow{}
	for rows.Next() {
		var i GetPoliciesWithCategoryUseKeysetPaginateRow
		if err := rows.Scan(
			&i.MPoliciesPkey,
			&i.PolicyID,
			&i.Name,
			&i.Description,
			&i.Key,
			&i.PolicyCategoryID,
			&i.MPolicyCategoriesPkey,
			&i.PolicyCategoryID_2,
			&i.Name_2,
			&i.Description_2,
			&i.Key_2,
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

const getPoliciesWithCategoryUseNumberedPaginate = `-- name: GetPoliciesWithCategoryUseNumberedPaginate :many
SELECT m_policies.m_policies_pkey, m_policies.policy_id, m_policies.name, m_policies.description, m_policies.key, m_policies.policy_category_id, m_policy_categories.m_policy_categories_pkey, m_policy_categories.policy_category_id, m_policy_categories.name, m_policy_categories.description, m_policy_categories.key FROM m_policies
JOIN m_policy_categories ON m_policies.policy_category_id = m_policy_categories.policy_category_id
WHERE
	CASE WHEN $3::boolean = true THEN m_policies.name LIKE '%' || $4::text || '%' ELSE TRUE END
AND
	CASE WHEN $5::boolean = true THEN policy_category_id = ANY($6::uuid[]) ELSE TRUE END
ORDER BY
	CASE WHEN $7::text = 'name' THEN m_policies.name END ASC,
	CASE WHEN $7::text = 'r_name' THEN m_policies.name END DESC,
	m_policies_pkey ASC
LIMIT $1 OFFSET $2
`

type GetPoliciesWithCategoryUseNumberedPaginateParams struct {
	Limit           int32       `json:"limit"`
	Offset          int32       `json:"offset"`
	WhereLikeName   bool        `json:"where_like_name"`
	SearchName      string      `json:"search_name"`
	WhereInCategory bool        `json:"where_in_category"`
	InCategories    []uuid.UUID `json:"in_categories"`
	OrderMethod     string      `json:"order_method"`
}

type GetPoliciesWithCategoryUseNumberedPaginateRow struct {
	MPoliciesPkey         pgtype.Int8 `json:"m_policies_pkey"`
	PolicyID              uuid.UUID   `json:"policy_id"`
	Name                  string      `json:"name"`
	Description           string      `json:"description"`
	Key                   string      `json:"key"`
	PolicyCategoryID      uuid.UUID   `json:"policy_category_id"`
	MPolicyCategoriesPkey pgtype.Int8 `json:"m_policy_categories_pkey"`
	PolicyCategoryID_2    uuid.UUID   `json:"policy_category_id_2"`
	Name_2                string      `json:"name_2"`
	Description_2         string      `json:"description_2"`
	Key_2                 string      `json:"key_2"`
}

func (q *Queries) GetPoliciesWithCategoryUseNumberedPaginate(ctx context.Context, arg GetPoliciesWithCategoryUseNumberedPaginateParams) ([]GetPoliciesWithCategoryUseNumberedPaginateRow, error) {
	rows, err := q.db.Query(ctx, getPoliciesWithCategoryUseNumberedPaginate,
		arg.Limit,
		arg.Offset,
		arg.WhereLikeName,
		arg.SearchName,
		arg.WhereInCategory,
		arg.InCategories,
		arg.OrderMethod,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPoliciesWithCategoryUseNumberedPaginateRow{}
	for rows.Next() {
		var i GetPoliciesWithCategoryUseNumberedPaginateRow
		if err := rows.Scan(
			&i.MPoliciesPkey,
			&i.PolicyID,
			&i.Name,
			&i.Description,
			&i.Key,
			&i.PolicyCategoryID,
			&i.MPolicyCategoriesPkey,
			&i.PolicyCategoryID_2,
			&i.Name_2,
			&i.Description_2,
			&i.Key_2,
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

const pluralDeletePolicies = `-- name: PluralDeletePolicies :exec
DELETE FROM m_policies WHERE policy_id = ANY($1::uuid[])
`

func (q *Queries) PluralDeletePolicies(ctx context.Context, dollar_1 []uuid.UUID) error {
	_, err := q.db.Exec(ctx, pluralDeletePolicies, dollar_1)
	return err
}

const updatePolicy = `-- name: UpdatePolicy :one
UPDATE m_policies SET name = $2, description = $3, key = $4, policy_category_id = $5 WHERE policy_id = $1 RETURNING m_policies_pkey, policy_id, name, description, key, policy_category_id
`

type UpdatePolicyParams struct {
	PolicyID         uuid.UUID `json:"policy_id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	Key              string    `json:"key"`
	PolicyCategoryID uuid.UUID `json:"policy_category_id"`
}

func (q *Queries) UpdatePolicy(ctx context.Context, arg UpdatePolicyParams) (Policy, error) {
	row := q.db.QueryRow(ctx, updatePolicy,
		arg.PolicyID,
		arg.Name,
		arg.Description,
		arg.Key,
		arg.PolicyCategoryID,
	)
	var i Policy
	err := row.Scan(
		&i.MPoliciesPkey,
		&i.PolicyID,
		&i.Name,
		&i.Description,
		&i.Key,
		&i.PolicyCategoryID,
	)
	return i, err
}

const updatePolicyByKey = `-- name: UpdatePolicyByKey :one
UPDATE m_policies SET name = $2, description = $3, policy_category_id = $4 WHERE key = $1 RETURNING m_policies_pkey, policy_id, name, description, key, policy_category_id
`

type UpdatePolicyByKeyParams struct {
	Key              string    `json:"key"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	PolicyCategoryID uuid.UUID `json:"policy_category_id"`
}

func (q *Queries) UpdatePolicyByKey(ctx context.Context, arg UpdatePolicyByKeyParams) (Policy, error) {
	row := q.db.QueryRow(ctx, updatePolicyByKey,
		arg.Key,
		arg.Name,
		arg.Description,
		arg.PolicyCategoryID,
	)
	var i Policy
	err := row.Scan(
		&i.MPoliciesPkey,
		&i.PolicyID,
		&i.Name,
		&i.Description,
		&i.Key,
		&i.PolicyCategoryID,
	)
	return i, err
}
