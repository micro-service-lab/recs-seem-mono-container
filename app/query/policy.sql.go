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

const findPolicyById = `-- name: FindPolicyById :one
SELECT m_policies_pkey, policy_id, name, description, key, policy_category_id FROM m_policies WHERE policy_id = $1
`

func (q *Queries) FindPolicyById(ctx context.Context, policyID uuid.UUID) (Policy, error) {
	row := q.db.QueryRow(ctx, findPolicyById, policyID)
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

const getPolicies = `-- name: GetPolicies :many
SELECT m_policies_pkey, policy_id, name, description, key, policy_category_id FROM m_policies
WHERE CASE
	WHEN $3::boolean = true THEN m_policies.name LIKE '%' || $4::text || '%'
END
ORDER BY
	CASE WHEN $5::text = 'name' THEN m_policies.name END ASC,
	m_policies_pkey DESC
LIMIT $1 OFFSET $2
`

type GetPoliciesParams struct {
	Limit         int32  `json:"limit"`
	Offset        int32  `json:"offset"`
	WhereLikeName bool   `json:"where_like_name"`
	SearchName    string `json:"search_name"`
	OrderMethod   string `json:"order_method"`
}

func (q *Queries) GetPolicies(ctx context.Context, arg GetPoliciesParams) ([]Policy, error) {
	rows, err := q.db.Query(ctx, getPolicies,
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

const getPoliciesByCategories = `-- name: GetPoliciesByCategories :many
SELECT m_policies_pkey, policy_id, name, description, key, policy_category_id FROM m_policies WHERE policy_category_id = ANY($3::uuid[])
AND CASE
	WHEN $4::boolean = true THEN m_policies.name LIKE '%' || $5::text || '%'
END
ORDER BY
	CASE WHEN $6::text = 'name' THEN m_policies.name END ASC,
	m_policies_pkey DESC
LIMIT $1 OFFSET $2
`

type GetPoliciesByCategoriesParams struct {
	Limit             int32       `json:"limit"`
	Offset            int32       `json:"offset"`
	PolicyCategoryIds []uuid.UUID `json:"policy_category_ids"`
	WhereLikeName     bool        `json:"where_like_name"`
	SearchName        string      `json:"search_name"`
	OrderMethod       string      `json:"order_method"`
}

func (q *Queries) GetPoliciesByCategories(ctx context.Context, arg GetPoliciesByCategoriesParams) ([]Policy, error) {
	rows, err := q.db.Query(ctx, getPoliciesByCategories,
		arg.Limit,
		arg.Offset,
		arg.PolicyCategoryIds,
		arg.WhereLikeName,
		arg.SearchName,
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

const getPoliciesByCategory = `-- name: GetPoliciesByCategory :many
SELECT m_policies_pkey, policy_id, name, description, key, policy_category_id FROM m_policies WHERE policy_category_id = $1
AND CASE
	WHEN $4::boolean = true THEN m_policies.name LIKE '%' || $5::text || '%'
END
ORDER BY
	CASE WHEN $6::text = 'name' THEN m_policies.name END ASC,
	m_policies_pkey DESC
LIMIT $2 OFFSET $3
`

type GetPoliciesByCategoryParams struct {
	PolicyCategoryID uuid.UUID `json:"policy_category_id"`
	Limit            int32     `json:"limit"`
	Offset           int32     `json:"offset"`
	WhereLikeName    bool      `json:"where_like_name"`
	SearchName       string    `json:"search_name"`
	OrderMethod      string    `json:"order_method"`
}

func (q *Queries) GetPoliciesByCategory(ctx context.Context, arg GetPoliciesByCategoryParams) ([]Policy, error) {
	rows, err := q.db.Query(ctx, getPoliciesByCategory,
		arg.PolicyCategoryID,
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

const getPoliciesByKeys = `-- name: GetPoliciesByKeys :many
SELECT m_policies_pkey, policy_id, name, description, key, policy_category_id FROM m_policies WHERE key = ANY($3::varchar[])
AND CASE
	WHEN $4::boolean = true THEN m_policies.name LIKE '%' || $5::text || '%'
END
ORDER BY
	CASE WHEN $6::text = 'name' THEN m_policies.name END ASC,
	m_policies_pkey DESC
LIMIT $1 OFFSET $2
`

type GetPoliciesByKeysParams struct {
	Limit         int32    `json:"limit"`
	Offset        int32    `json:"offset"`
	Keys          []string `json:"keys"`
	WhereLikeName bool     `json:"where_like_name"`
	SearchName    string   `json:"search_name"`
	OrderMethod   string   `json:"order_method"`
}

func (q *Queries) GetPoliciesByKeys(ctx context.Context, arg GetPoliciesByKeysParams) ([]Policy, error) {
	rows, err := q.db.Query(ctx, getPoliciesByKeys,
		arg.Limit,
		arg.Offset,
		arg.Keys,
		arg.WhereLikeName,
		arg.SearchName,
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

const getPoliciesByKeysWithCategory = `-- name: GetPoliciesByKeysWithCategory :many
SELECT m_policies.m_policies_pkey, m_policies.policy_id, m_policies.name, m_policies.description, m_policies.key, m_policies.policy_category_id, m_policy_categories.m_policy_categories_pkey, m_policy_categories.policy_category_id, m_policy_categories.name, m_policy_categories.description, m_policy_categories.key FROM m_policies
JOIN m_policy_categories ON m_policies.policy_category_id = m_policy_categories.policy_category_id
WHERE m_policies.key = ANY($3::varchar[])
AND CASE
	WHEN $4::boolean = true THEN m_policies.name LIKE '%' || $5::text || '%'
END
ORDER BY
	CASE WHEN $6::text = 'name' THEN m_policies.name END ASC,
	m_policies_pkey DESC
LIMIT $1 OFFSET $2
`

type GetPoliciesByKeysWithCategoryParams struct {
	Limit         int32    `json:"limit"`
	Offset        int32    `json:"offset"`
	Keys          []string `json:"keys"`
	WhereLikeName bool     `json:"where_like_name"`
	SearchName    string   `json:"search_name"`
	OrderMethod   string   `json:"order_method"`
}

type GetPoliciesByKeysWithCategoryRow struct {
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

func (q *Queries) GetPoliciesByKeysWithCategory(ctx context.Context, arg GetPoliciesByKeysWithCategoryParams) ([]GetPoliciesByKeysWithCategoryRow, error) {
	rows, err := q.db.Query(ctx, getPoliciesByKeysWithCategory,
		arg.Limit,
		arg.Offset,
		arg.Keys,
		arg.WhereLikeName,
		arg.SearchName,
		arg.OrderMethod,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPoliciesByKeysWithCategoryRow{}
	for rows.Next() {
		var i GetPoliciesByKeysWithCategoryRow
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

const getPoliciesWithCategory = `-- name: GetPoliciesWithCategory :many
SELECT m_policies.m_policies_pkey, m_policies.policy_id, m_policies.name, m_policies.description, m_policies.key, m_policies.policy_category_id, m_policy_categories.m_policy_categories_pkey, m_policy_categories.policy_category_id, m_policy_categories.name, m_policy_categories.description, m_policy_categories.key FROM m_policies
JOIN m_policy_categories ON m_policies.policy_category_id = m_policy_categories.policy_category_id
WHERE CASE
	WHEN $3::boolean = true THEN m_policies.name LIKE '%' || $4::text || '%'
END
ORDER BY
	CASE WHEN $5::text = 'name' THEN m_policies.name END ASC,
	m_policies_pkey DESC
LIMIT $1 OFFSET $2
`

type GetPoliciesWithCategoryParams struct {
	Limit         int32  `json:"limit"`
	Offset        int32  `json:"offset"`
	WhereLikeName bool   `json:"where_like_name"`
	SearchName    string `json:"search_name"`
	OrderMethod   string `json:"order_method"`
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

const getPolicyByKeyWithCategory = `-- name: GetPolicyByKeyWithCategory :one
SELECT m_policies.m_policies_pkey, m_policies.policy_id, m_policies.name, m_policies.description, m_policies.key, m_policies.policy_category_id, m_policy_categories.m_policy_categories_pkey, m_policy_categories.policy_category_id, m_policy_categories.name, m_policy_categories.description, m_policy_categories.key FROM m_policies
JOIN m_policy_categories ON m_policies.policy_category_id = m_policy_categories.policy_category_id
WHERE m_policies.key = $1
`

type GetPolicyByKeyWithCategoryRow struct {
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

func (q *Queries) GetPolicyByKeyWithCategory(ctx context.Context, key string) (GetPolicyByKeyWithCategoryRow, error) {
	row := q.db.QueryRow(ctx, getPolicyByKeyWithCategory, key)
	var i GetPolicyByKeyWithCategoryRow
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

const updatePolicy = `-- name: UpdatePolicy :one
UPDATE m_policies SET name = $2, description = $3 WHERE policy_id = $1 RETURNING m_policies_pkey, policy_id, name, description, key, policy_category_id
`

type UpdatePolicyParams struct {
	PolicyID    uuid.UUID `json:"policy_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

func (q *Queries) UpdatePolicy(ctx context.Context, arg UpdatePolicyParams) (Policy, error) {
	row := q.db.QueryRow(ctx, updatePolicy, arg.PolicyID, arg.Name, arg.Description)
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
UPDATE m_policies SET name = $2, description = $3 WHERE key = $1 RETURNING m_policies_pkey, policy_id, name, description, key, policy_category_id
`

type UpdatePolicyByKeyParams struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (q *Queries) UpdatePolicyByKey(ctx context.Context, arg UpdatePolicyByKeyParams) (Policy, error) {
	row := q.db.QueryRow(ctx, updatePolicyByKey, arg.Key, arg.Name, arg.Description)
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

const updatePolicyCategoryID = `-- name: UpdatePolicyCategoryID :one
UPDATE m_policies SET policy_category_id = $2 WHERE policy_id = $1 RETURNING m_policies_pkey, policy_id, name, description, key, policy_category_id
`

type UpdatePolicyCategoryIDParams struct {
	PolicyID         uuid.UUID `json:"policy_id"`
	PolicyCategoryID uuid.UUID `json:"policy_category_id"`
}

func (q *Queries) UpdatePolicyCategoryID(ctx context.Context, arg UpdatePolicyCategoryIDParams) (Policy, error) {
	row := q.db.QueryRow(ctx, updatePolicyCategoryID, arg.PolicyID, arg.PolicyCategoryID)
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

const updatePolicyKey = `-- name: UpdatePolicyKey :one
UPDATE m_policies SET key = $2 WHERE policy_id = $1 RETURNING m_policies_pkey, policy_id, name, description, key, policy_category_id
`

type UpdatePolicyKeyParams struct {
	PolicyID uuid.UUID `json:"policy_id"`
	Key      string    `json:"key"`
}

func (q *Queries) UpdatePolicyKey(ctx context.Context, arg UpdatePolicyKeyParams) (Policy, error) {
	row := q.db.QueryRow(ctx, updatePolicyKey, arg.PolicyID, arg.Key)
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
