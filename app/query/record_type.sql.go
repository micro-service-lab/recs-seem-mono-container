// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: record_type.sql

package query

import (
	"context"

	"github.com/google/uuid"
)

const countRecordTypes = `-- name: CountRecordTypes :one
SELECT COUNT(*) FROM m_record_types
WHERE
	CASE WHEN $1::boolean = true THEN name LIKE '%' || $2::text || '%' ELSE TRUE END
`

type CountRecordTypesParams struct {
	WhereLikeName bool   `json:"where_like_name"`
	SearchName    string `json:"search_name"`
}

func (q *Queries) CountRecordTypes(ctx context.Context, arg CountRecordTypesParams) (int64, error) {
	row := q.db.QueryRow(ctx, countRecordTypes, arg.WhereLikeName, arg.SearchName)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createRecordType = `-- name: CreateRecordType :one
INSERT INTO m_record_types (name, key) VALUES ($1, $2) RETURNING m_record_types_pkey, record_type_id, name, key
`

type CreateRecordTypeParams struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

func (q *Queries) CreateRecordType(ctx context.Context, arg CreateRecordTypeParams) (RecordType, error) {
	row := q.db.QueryRow(ctx, createRecordType, arg.Name, arg.Key)
	var i RecordType
	err := row.Scan(
		&i.MRecordTypesPkey,
		&i.RecordTypeID,
		&i.Name,
		&i.Key,
	)
	return i, err
}

type CreateRecordTypesParams struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

const deleteRecordType = `-- name: DeleteRecordType :exec
DELETE FROM m_record_types WHERE record_type_id = $1
`

func (q *Queries) DeleteRecordType(ctx context.Context, recordTypeID uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteRecordType, recordTypeID)
	return err
}

const deleteRecordTypeByKey = `-- name: DeleteRecordTypeByKey :exec
DELETE FROM m_record_types WHERE key = $1
`

func (q *Queries) DeleteRecordTypeByKey(ctx context.Context, key string) error {
	_, err := q.db.Exec(ctx, deleteRecordTypeByKey, key)
	return err
}

const findRecordTypeByID = `-- name: FindRecordTypeByID :one
SELECT m_record_types_pkey, record_type_id, name, key FROM m_record_types WHERE record_type_id = $1
`

func (q *Queries) FindRecordTypeByID(ctx context.Context, recordTypeID uuid.UUID) (RecordType, error) {
	row := q.db.QueryRow(ctx, findRecordTypeByID, recordTypeID)
	var i RecordType
	err := row.Scan(
		&i.MRecordTypesPkey,
		&i.RecordTypeID,
		&i.Name,
		&i.Key,
	)
	return i, err
}

const findRecordTypeByKey = `-- name: FindRecordTypeByKey :one
SELECT m_record_types_pkey, record_type_id, name, key FROM m_record_types WHERE key = $1
`

func (q *Queries) FindRecordTypeByKey(ctx context.Context, key string) (RecordType, error) {
	row := q.db.QueryRow(ctx, findRecordTypeByKey, key)
	var i RecordType
	err := row.Scan(
		&i.MRecordTypesPkey,
		&i.RecordTypeID,
		&i.Name,
		&i.Key,
	)
	return i, err
}

const getRecordTypes = `-- name: GetRecordTypes :many
SELECT m_record_types_pkey, record_type_id, name, key FROM m_record_types
WHERE
	CASE WHEN $3::boolean = true THEN m_record_types.name LIKE '%' || $4::text || '%' ELSE TRUE END
ORDER BY
	CASE WHEN $5::text = 'name' THEN name END ASC,
	m_record_types_pkey DESC
LIMIT $1 OFFSET $2
`

type GetRecordTypesParams struct {
	Limit         int32  `json:"limit"`
	Offset        int32  `json:"offset"`
	WhereLikeName bool   `json:"where_like_name"`
	SearchName    string `json:"search_name"`
	OrderMethod   string `json:"order_method"`
}

func (q *Queries) GetRecordTypes(ctx context.Context, arg GetRecordTypesParams) ([]RecordType, error) {
	rows, err := q.db.Query(ctx, getRecordTypes,
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
	items := []RecordType{}
	for rows.Next() {
		var i RecordType
		if err := rows.Scan(
			&i.MRecordTypesPkey,
			&i.RecordTypeID,
			&i.Name,
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

const updateRecordType = `-- name: UpdateRecordType :one
UPDATE m_record_types SET name = $2, key = $3 WHERE record_type_id = $1 RETURNING m_record_types_pkey, record_type_id, name, key
`

type UpdateRecordTypeParams struct {
	RecordTypeID uuid.UUID `json:"record_type_id"`
	Name         string    `json:"name"`
	Key          string    `json:"key"`
}

func (q *Queries) UpdateRecordType(ctx context.Context, arg UpdateRecordTypeParams) (RecordType, error) {
	row := q.db.QueryRow(ctx, updateRecordType, arg.RecordTypeID, arg.Name, arg.Key)
	var i RecordType
	err := row.Scan(
		&i.MRecordTypesPkey,
		&i.RecordTypeID,
		&i.Name,
		&i.Key,
	)
	return i, err
}

const updateRecordTypeByKey = `-- name: UpdateRecordTypeByKey :one
UPDATE m_record_types SET name = $2 WHERE key = $1 RETURNING m_record_types_pkey, record_type_id, name, key
`

type UpdateRecordTypeByKeyParams struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

func (q *Queries) UpdateRecordTypeByKey(ctx context.Context, arg UpdateRecordTypeByKeyParams) (RecordType, error) {
	row := q.db.QueryRow(ctx, updateRecordTypeByKey, arg.Key, arg.Name)
	var i RecordType
	err := row.Scan(
		&i.MRecordTypesPkey,
		&i.RecordTypeID,
		&i.Name,
		&i.Key,
	)
	return i, err
}
