// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: mime_type.sql

package query

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const countMimeTypes = `-- name: CountMimeTypes :one
SELECT COUNT(*) FROM m_mime_types
WHERE
	CASE WHEN $1::boolean = true THEN name LIKE '%' || $2::text || '%' ELSE TRUE END
`

type CountMimeTypesParams struct {
	WhereLikeName bool   `json:"where_like_name"`
	SearchName    string `json:"search_name"`
}

func (q *Queries) CountMimeTypes(ctx context.Context, arg CountMimeTypesParams) (int64, error) {
	row := q.db.QueryRow(ctx, countMimeTypes, arg.WhereLikeName, arg.SearchName)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createMimeType = `-- name: CreateMimeType :one
INSERT INTO m_mime_types (name, key) VALUES ($1, $2) RETURNING m_mime_types_pkey, mime_type_id, name, key
`

type CreateMimeTypeParams struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

func (q *Queries) CreateMimeType(ctx context.Context, arg CreateMimeTypeParams) (MimeType, error) {
	row := q.db.QueryRow(ctx, createMimeType, arg.Name, arg.Key)
	var i MimeType
	err := row.Scan(
		&i.MMimeTypesPkey,
		&i.MimeTypeID,
		&i.Name,
		&i.Key,
	)
	return i, err
}

type CreateMimeTypesParams struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

const deleteMimeType = `-- name: DeleteMimeType :exec
DELETE FROM m_mime_types WHERE mime_type_id = $1
`

func (q *Queries) DeleteMimeType(ctx context.Context, mimeTypeID uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteMimeType, mimeTypeID)
	return err
}

const deleteMimeTypeByKey = `-- name: DeleteMimeTypeByKey :exec
DELETE FROM m_mime_types WHERE key = $1
`

func (q *Queries) DeleteMimeTypeByKey(ctx context.Context, key string) error {
	_, err := q.db.Exec(ctx, deleteMimeTypeByKey, key)
	return err
}

const findMimeTypeByID = `-- name: FindMimeTypeByID :one
SELECT m_mime_types_pkey, mime_type_id, name, key FROM m_mime_types WHERE mime_type_id = $1
`

func (q *Queries) FindMimeTypeByID(ctx context.Context, mimeTypeID uuid.UUID) (MimeType, error) {
	row := q.db.QueryRow(ctx, findMimeTypeByID, mimeTypeID)
	var i MimeType
	err := row.Scan(
		&i.MMimeTypesPkey,
		&i.MimeTypeID,
		&i.Name,
		&i.Key,
	)
	return i, err
}

const findMimeTypeByKey = `-- name: FindMimeTypeByKey :one
SELECT m_mime_types_pkey, mime_type_id, name, key FROM m_mime_types WHERE key = $1
`

func (q *Queries) FindMimeTypeByKey(ctx context.Context, key string) (MimeType, error) {
	row := q.db.QueryRow(ctx, findMimeTypeByKey, key)
	var i MimeType
	err := row.Scan(
		&i.MMimeTypesPkey,
		&i.MimeTypeID,
		&i.Name,
		&i.Key,
	)
	return i, err
}

const getMimeTypes = `-- name: GetMimeTypes :many
SELECT m_mime_types_pkey, mime_type_id, name, key FROM m_mime_types
ORDER BY
	CASE WHEN $1::text = 'name' THEN name END ASC,
	CASE WHEN $1::text = 'r_name' THEN name END DESC,
	m_mime_types_pkey DESC
`

func (q *Queries) GetMimeTypes(ctx context.Context, orderMethod string) ([]MimeType, error) {
	rows, err := q.db.Query(ctx, getMimeTypes, orderMethod)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []MimeType{}
	for rows.Next() {
		var i MimeType
		if err := rows.Scan(
			&i.MMimeTypesPkey,
			&i.MimeTypeID,
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

const getMimeTypesUseKeysetPaginate = `-- name: GetMimeTypesUseKeysetPaginate :many
SELECT m_mime_types_pkey, mime_type_id, name, key FROM m_mime_types
WHERE
	CASE $1
		WHEN 'next' THEN
			CASE $2::text
				WHEN 'name' THEN name > $3 OR (name = $3 AND m_mime_types_pkey < $4)
				WHEN 'r_name' THEN name < $3 OR (name = $3 AND m_mime_types_pkey < $4)
				ELSE m_mime_types_pkey < $4
			END
		WHEN 'prev' THEN
			CASE $2::text
				WHEN 'name' THEN name < $3 OR (name = $3 AND m_mime_types_pkey > $4)
				WHEN 'r_name' THEN name > $3 OR (name = $3 AND m_mime_types_pkey > $4)
				ELSE m_mime_types_pkey > $4
			END
	END
ORDER BY
	CASE WHEN $2::text = 'name' THEN name END ASC,
	CASE WHEN $2::text = 'r_name' THEN name END DESC,
	m_mime_types_pkey DESC
`

type GetMimeTypesUseKeysetPaginateParams struct {
	CursorDirection interface{} `json:"cursor_direction"`
	OrderMethod     string      `json:"order_method"`
	CursorColumn    string      `json:"cursor_column"`
	Cursor          pgtype.Int8 `json:"cursor"`
}

func (q *Queries) GetMimeTypesUseKeysetPaginate(ctx context.Context, arg GetMimeTypesUseKeysetPaginateParams) ([]MimeType, error) {
	rows, err := q.db.Query(ctx, getMimeTypesUseKeysetPaginate,
		arg.CursorDirection,
		arg.OrderMethod,
		arg.CursorColumn,
		arg.Cursor,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []MimeType{}
	for rows.Next() {
		var i MimeType
		if err := rows.Scan(
			&i.MMimeTypesPkey,
			&i.MimeTypeID,
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

const getMimeTypesUseNumberedPaginate = `-- name: GetMimeTypesUseNumberedPaginate :many
SELECT m_mime_types_pkey, mime_type_id, name, key FROM m_mime_types
ORDER BY
	CASE WHEN $3::text = 'name' THEN name END ASC,
	CASE WHEN $3::text = 'r_name' THEN name END DESC,
	m_mime_types_pkey DESC
LIMIT $1 OFFSET $2
`

type GetMimeTypesUseNumberedPaginateParams struct {
	Limit       int32  `json:"limit"`
	Offset      int32  `json:"offset"`
	OrderMethod string `json:"order_method"`
}

func (q *Queries) GetMimeTypesUseNumberedPaginate(ctx context.Context, arg GetMimeTypesUseNumberedPaginateParams) ([]MimeType, error) {
	rows, err := q.db.Query(ctx, getMimeTypesUseNumberedPaginate, arg.Limit, arg.Offset, arg.OrderMethod)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []MimeType{}
	for rows.Next() {
		var i MimeType
		if err := rows.Scan(
			&i.MMimeTypesPkey,
			&i.MimeTypeID,
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
