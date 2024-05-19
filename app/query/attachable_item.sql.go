// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: attachable_item.sql

package query

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const countAttachableItems = `-- name: CountAttachableItems :one
SELECT COUNT(*) FROM t_attachable_items
WHERE
	CASE WHEN $1::boolean = true THEN mime_type_id = ANY($2::uuid[]) ELSE TRUE END
AND
	CASE WHEN $3::boolean = true THEN owner_id = ANY($4::uuid[]) ELSE TRUE END
`

type CountAttachableItemsParams struct {
	WhereInMimeTypeIds bool        `json:"where_in_mime_type_ids"`
	InMimeTypeIds      []uuid.UUID `json:"in_mime_type_ids"`
	WhereInOwnerIds    bool        `json:"where_in_owner_ids"`
	InOwnerIds         []uuid.UUID `json:"in_owner_ids"`
}

func (q *Queries) CountAttachableItems(ctx context.Context, arg CountAttachableItemsParams) (int64, error) {
	row := q.db.QueryRow(ctx, countAttachableItems,
		arg.WhereInMimeTypeIds,
		arg.InMimeTypeIds,
		arg.WhereInOwnerIds,
		arg.InOwnerIds,
	)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createAttachableItem = `-- name: CreateAttachableItem :one
INSERT INTO t_attachable_items (url, size, owner_id, from_outer, mime_type_id) VALUES ($1, $2, $3, $4, $5) RETURNING t_attachable_items_pkey, attachable_item_id, url, size, mime_type_id, owner_id, from_outer
`

type CreateAttachableItemParams struct {
	Url        string        `json:"url"`
	Size       pgtype.Float8 `json:"size"`
	OwnerID    pgtype.UUID   `json:"owner_id"`
	FromOuter  bool          `json:"from_outer"`
	MimeTypeID uuid.UUID     `json:"mime_type_id"`
}

func (q *Queries) CreateAttachableItem(ctx context.Context, arg CreateAttachableItemParams) (AttachableItem, error) {
	row := q.db.QueryRow(ctx, createAttachableItem,
		arg.Url,
		arg.Size,
		arg.OwnerID,
		arg.FromOuter,
		arg.MimeTypeID,
	)
	var i AttachableItem
	err := row.Scan(
		&i.TAttachableItemsPkey,
		&i.AttachableItemID,
		&i.Url,
		&i.Size,
		&i.MimeTypeID,
		&i.OwnerID,
		&i.FromOuter,
	)
	return i, err
}

type CreateAttachableItemsParams struct {
	Url        string        `json:"url"`
	Size       pgtype.Float8 `json:"size"`
	OwnerID    pgtype.UUID   `json:"owner_id"`
	FromOuter  bool          `json:"from_outer"`
	MimeTypeID uuid.UUID     `json:"mime_type_id"`
}

const deleteAttachableItem = `-- name: DeleteAttachableItem :execrows
DELETE FROM t_attachable_items WHERE attachable_item_id = $1
`

func (q *Queries) DeleteAttachableItem(ctx context.Context, attachableItemID uuid.UUID) (int64, error) {
	result, err := q.db.Exec(ctx, deleteAttachableItem, attachableItemID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const findAttachableItemByID = `-- name: FindAttachableItemByID :one
SELECT t_attachable_items.t_attachable_items_pkey, t_attachable_items.attachable_item_id, t_attachable_items.url, t_attachable_items.size, t_attachable_items.mime_type_id, t_attachable_items.owner_id, t_attachable_items.from_outer, t_images.image_id, t_images.height image_height,
t_images.width image_width, t_files.file_id
FROM t_attachable_items
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE t_attachable_items.attachable_item_id = $1
`

type FindAttachableItemByIDRow struct {
	TAttachableItemsPkey pgtype.Int8   `json:"t_attachable_items_pkey"`
	AttachableItemID     uuid.UUID     `json:"attachable_item_id"`
	Url                  string        `json:"url"`
	Size                 pgtype.Float8 `json:"size"`
	MimeTypeID           uuid.UUID     `json:"mime_type_id"`
	OwnerID              pgtype.UUID   `json:"owner_id"`
	FromOuter            bool          `json:"from_outer"`
	ImageID              pgtype.UUID   `json:"image_id"`
	ImageHeight          pgtype.Float8 `json:"image_height"`
	ImageWidth           pgtype.Float8 `json:"image_width"`
	FileID               pgtype.UUID   `json:"file_id"`
}

func (q *Queries) FindAttachableItemByID(ctx context.Context, attachableItemID uuid.UUID) (FindAttachableItemByIDRow, error) {
	row := q.db.QueryRow(ctx, findAttachableItemByID, attachableItemID)
	var i FindAttachableItemByIDRow
	err := row.Scan(
		&i.TAttachableItemsPkey,
		&i.AttachableItemID,
		&i.Url,
		&i.Size,
		&i.MimeTypeID,
		&i.OwnerID,
		&i.FromOuter,
		&i.ImageID,
		&i.ImageHeight,
		&i.ImageWidth,
		&i.FileID,
	)
	return i, err
}

const findAttachableItemByIDWithMimeType = `-- name: FindAttachableItemByIDWithMimeType :one
SELECT t_attachable_items.t_attachable_items_pkey, t_attachable_items.attachable_item_id, t_attachable_items.url, t_attachable_items.size, t_attachable_items.mime_type_id, t_attachable_items.owner_id, t_attachable_items.from_outer, m_mime_types.name mime_type_name, m_mime_types.kind mime_type_kind,
m_mime_types.key mime_type_key, t_images.height image_height, t_images.width image_width, t_images.image_id, t_files.file_id
FROM t_attachable_items
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE t_attachable_items.attachable_item_id = $1
`

type FindAttachableItemByIDWithMimeTypeRow struct {
	TAttachableItemsPkey pgtype.Int8   `json:"t_attachable_items_pkey"`
	AttachableItemID     uuid.UUID     `json:"attachable_item_id"`
	Url                  string        `json:"url"`
	Size                 pgtype.Float8 `json:"size"`
	MimeTypeID           uuid.UUID     `json:"mime_type_id"`
	OwnerID              pgtype.UUID   `json:"owner_id"`
	FromOuter            bool          `json:"from_outer"`
	MimeTypeName         pgtype.Text   `json:"mime_type_name"`
	MimeTypeKind         pgtype.Text   `json:"mime_type_kind"`
	MimeTypeKey          pgtype.Text   `json:"mime_type_key"`
	ImageHeight          pgtype.Float8 `json:"image_height"`
	ImageWidth           pgtype.Float8 `json:"image_width"`
	ImageID              pgtype.UUID   `json:"image_id"`
	FileID               pgtype.UUID   `json:"file_id"`
}

func (q *Queries) FindAttachableItemByIDWithMimeType(ctx context.Context, attachableItemID uuid.UUID) (FindAttachableItemByIDWithMimeTypeRow, error) {
	row := q.db.QueryRow(ctx, findAttachableItemByIDWithMimeType, attachableItemID)
	var i FindAttachableItemByIDWithMimeTypeRow
	err := row.Scan(
		&i.TAttachableItemsPkey,
		&i.AttachableItemID,
		&i.Url,
		&i.Size,
		&i.MimeTypeID,
		&i.OwnerID,
		&i.FromOuter,
		&i.MimeTypeName,
		&i.MimeTypeKind,
		&i.MimeTypeKey,
		&i.ImageHeight,
		&i.ImageWidth,
		&i.ImageID,
		&i.FileID,
	)
	return i, err
}

const getAttachableItems = `-- name: GetAttachableItems :many
SELECT t_attachable_items.t_attachable_items_pkey, t_attachable_items.attachable_item_id, t_attachable_items.url, t_attachable_items.size, t_attachable_items.mime_type_id, t_attachable_items.owner_id, t_attachable_items.from_outer, t_images.image_id, t_images.height image_height,
t_images.width image_width, t_files.file_id
FROM t_attachable_items
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE
	CASE WHEN $1::boolean = true THEN mime_type_id = ANY($2::uuid[]) ELSE TRUE END
AND
	CASE WHEN $3::boolean = true THEN owner_id = ANY($4::uuid[]) ELSE TRUE END
ORDER BY
	t_attachable_items_pkey ASC
`

type GetAttachableItemsParams struct {
	WhereInMimeTypeIds bool        `json:"where_in_mime_type_ids"`
	InMimeTypeIds      []uuid.UUID `json:"in_mime_type_ids"`
	WhereInOwnerIds    bool        `json:"where_in_owner_ids"`
	InOwnerIds         []uuid.UUID `json:"in_owner_ids"`
}

type GetAttachableItemsRow struct {
	TAttachableItemsPkey pgtype.Int8   `json:"t_attachable_items_pkey"`
	AttachableItemID     uuid.UUID     `json:"attachable_item_id"`
	Url                  string        `json:"url"`
	Size                 pgtype.Float8 `json:"size"`
	MimeTypeID           uuid.UUID     `json:"mime_type_id"`
	OwnerID              pgtype.UUID   `json:"owner_id"`
	FromOuter            bool          `json:"from_outer"`
	ImageID              pgtype.UUID   `json:"image_id"`
	ImageHeight          pgtype.Float8 `json:"image_height"`
	ImageWidth           pgtype.Float8 `json:"image_width"`
	FileID               pgtype.UUID   `json:"file_id"`
}

func (q *Queries) GetAttachableItems(ctx context.Context, arg GetAttachableItemsParams) ([]GetAttachableItemsRow, error) {
	rows, err := q.db.Query(ctx, getAttachableItems,
		arg.WhereInMimeTypeIds,
		arg.InMimeTypeIds,
		arg.WhereInOwnerIds,
		arg.InOwnerIds,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAttachableItemsRow{}
	for rows.Next() {
		var i GetAttachableItemsRow
		if err := rows.Scan(
			&i.TAttachableItemsPkey,
			&i.AttachableItemID,
			&i.Url,
			&i.Size,
			&i.MimeTypeID,
			&i.OwnerID,
			&i.FromOuter,
			&i.ImageID,
			&i.ImageHeight,
			&i.ImageWidth,
			&i.FileID,
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

const getAttachableItemsUseKeysetPaginate = `-- name: GetAttachableItemsUseKeysetPaginate :many
SELECT t_attachable_items.t_attachable_items_pkey, t_attachable_items.attachable_item_id, t_attachable_items.url, t_attachable_items.size, t_attachable_items.mime_type_id, t_attachable_items.owner_id, t_attachable_items.from_outer, t_images.image_id, t_images.height image_height,
t_images.width image_width, t_files.file_id
FROM t_attachable_items
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE
	CASE WHEN $2::boolean = true THEN mime_type_id = ANY($3::uuid[]) ELSE TRUE END
AND
	CASE WHEN $4::boolean = true THEN owner_id = ANY($5::uuid[]) ELSE TRUE END
AND
	CASE $6::text
		WHEN 'next' THEN
			t_attachable_items_pkey < $7::int
		WHEN 'prev' THEN
			t_attachable_items_pkey > $7::int
	END
ORDER BY
	CASE WHEN $6::text = 'next' THEN t_attachable_items_pkey END ASC,
	CASE WHEN $6::text = 'prev' THEN t_attachable_items_pkey END DESC
LIMIT $1
`

type GetAttachableItemsUseKeysetPaginateParams struct {
	Limit              int32       `json:"limit"`
	WhereInMimeTypeIds bool        `json:"where_in_mime_type_ids"`
	InMimeTypeIds      []uuid.UUID `json:"in_mime_type_ids"`
	WhereInOwnerIds    bool        `json:"where_in_owner_ids"`
	InOwnerIds         []uuid.UUID `json:"in_owner_ids"`
	CursorDirection    string      `json:"cursor_direction"`
	Cursor             int32       `json:"cursor"`
}

type GetAttachableItemsUseKeysetPaginateRow struct {
	TAttachableItemsPkey pgtype.Int8   `json:"t_attachable_items_pkey"`
	AttachableItemID     uuid.UUID     `json:"attachable_item_id"`
	Url                  string        `json:"url"`
	Size                 pgtype.Float8 `json:"size"`
	MimeTypeID           uuid.UUID     `json:"mime_type_id"`
	OwnerID              pgtype.UUID   `json:"owner_id"`
	FromOuter            bool          `json:"from_outer"`
	ImageID              pgtype.UUID   `json:"image_id"`
	ImageHeight          pgtype.Float8 `json:"image_height"`
	ImageWidth           pgtype.Float8 `json:"image_width"`
	FileID               pgtype.UUID   `json:"file_id"`
}

func (q *Queries) GetAttachableItemsUseKeysetPaginate(ctx context.Context, arg GetAttachableItemsUseKeysetPaginateParams) ([]GetAttachableItemsUseKeysetPaginateRow, error) {
	rows, err := q.db.Query(ctx, getAttachableItemsUseKeysetPaginate,
		arg.Limit,
		arg.WhereInMimeTypeIds,
		arg.InMimeTypeIds,
		arg.WhereInOwnerIds,
		arg.InOwnerIds,
		arg.CursorDirection,
		arg.Cursor,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAttachableItemsUseKeysetPaginateRow{}
	for rows.Next() {
		var i GetAttachableItemsUseKeysetPaginateRow
		if err := rows.Scan(
			&i.TAttachableItemsPkey,
			&i.AttachableItemID,
			&i.Url,
			&i.Size,
			&i.MimeTypeID,
			&i.OwnerID,
			&i.FromOuter,
			&i.ImageID,
			&i.ImageHeight,
			&i.ImageWidth,
			&i.FileID,
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

const getAttachableItemsUseNumberedPaginate = `-- name: GetAttachableItemsUseNumberedPaginate :many
SELECT t_attachable_items.t_attachable_items_pkey, t_attachable_items.attachable_item_id, t_attachable_items.url, t_attachable_items.size, t_attachable_items.mime_type_id, t_attachable_items.owner_id, t_attachable_items.from_outer, t_images.image_id, t_images.height image_height,
t_images.width image_width, t_files.file_id
FROM t_attachable_items
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE
	CASE WHEN $3::boolean = true THEN mime_type_id = ANY($4::uuid[]) ELSE TRUE END
AND
	CASE WHEN $5::boolean = true THEN owner_id = ANY($6::uuid[]) ELSE TRUE END
ORDER BY
	t_attachable_items_pkey ASC
LIMIT $1 OFFSET $2
`

type GetAttachableItemsUseNumberedPaginateParams struct {
	Limit              int32       `json:"limit"`
	Offset             int32       `json:"offset"`
	WhereInMimeTypeIds bool        `json:"where_in_mime_type_ids"`
	InMimeTypeIds      []uuid.UUID `json:"in_mime_type_ids"`
	WhereInOwnerIds    bool        `json:"where_in_owner_ids"`
	InOwnerIds         []uuid.UUID `json:"in_owner_ids"`
}

type GetAttachableItemsUseNumberedPaginateRow struct {
	TAttachableItemsPkey pgtype.Int8   `json:"t_attachable_items_pkey"`
	AttachableItemID     uuid.UUID     `json:"attachable_item_id"`
	Url                  string        `json:"url"`
	Size                 pgtype.Float8 `json:"size"`
	MimeTypeID           uuid.UUID     `json:"mime_type_id"`
	OwnerID              pgtype.UUID   `json:"owner_id"`
	FromOuter            bool          `json:"from_outer"`
	ImageID              pgtype.UUID   `json:"image_id"`
	ImageHeight          pgtype.Float8 `json:"image_height"`
	ImageWidth           pgtype.Float8 `json:"image_width"`
	FileID               pgtype.UUID   `json:"file_id"`
}

func (q *Queries) GetAttachableItemsUseNumberedPaginate(ctx context.Context, arg GetAttachableItemsUseNumberedPaginateParams) ([]GetAttachableItemsUseNumberedPaginateRow, error) {
	rows, err := q.db.Query(ctx, getAttachableItemsUseNumberedPaginate,
		arg.Limit,
		arg.Offset,
		arg.WhereInMimeTypeIds,
		arg.InMimeTypeIds,
		arg.WhereInOwnerIds,
		arg.InOwnerIds,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAttachableItemsUseNumberedPaginateRow{}
	for rows.Next() {
		var i GetAttachableItemsUseNumberedPaginateRow
		if err := rows.Scan(
			&i.TAttachableItemsPkey,
			&i.AttachableItemID,
			&i.Url,
			&i.Size,
			&i.MimeTypeID,
			&i.OwnerID,
			&i.FromOuter,
			&i.ImageID,
			&i.ImageHeight,
			&i.ImageWidth,
			&i.FileID,
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

const getAttachableItemsWithMimeType = `-- name: GetAttachableItemsWithMimeType :many
SELECT t_attachable_items.t_attachable_items_pkey, t_attachable_items.attachable_item_id, t_attachable_items.url, t_attachable_items.size, t_attachable_items.mime_type_id, t_attachable_items.owner_id, t_attachable_items.from_outer, m_mime_types.name mime_type_name, m_mime_types.kind mime_type_kind,
m_mime_types.key mime_type_key, t_images.height image_height, t_images.width image_width, t_images.image_id, t_files.file_id
FROM t_attachable_items
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.where_mime_type_id
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE
	CASE WHEN $1::boolean = true THEN mime_type_id = ANY($2::uuid[]) ELSE TRUE END
AND
	CASE WHEN $3::boolean = true THEN owner_id = ANY($4::uuid[]) ELSE TRUE END
ORDER BY
	t_attachable_items_pkey ASC
`

type GetAttachableItemsWithMimeTypeParams struct {
	WhereInMimeTypeIds bool        `json:"where_in_mime_type_ids"`
	InMimeTypeIds      []uuid.UUID `json:"in_mime_type_ids"`
	WhereInOwnerIds    bool        `json:"where_in_owner_ids"`
	InOwnerIds         []uuid.UUID `json:"in_owner_ids"`
}

type GetAttachableItemsWithMimeTypeRow struct {
	TAttachableItemsPkey pgtype.Int8   `json:"t_attachable_items_pkey"`
	AttachableItemID     uuid.UUID     `json:"attachable_item_id"`
	Url                  string        `json:"url"`
	Size                 pgtype.Float8 `json:"size"`
	MimeTypeID           uuid.UUID     `json:"mime_type_id"`
	OwnerID              pgtype.UUID   `json:"owner_id"`
	FromOuter            bool          `json:"from_outer"`
	MimeTypeName         pgtype.Text   `json:"mime_type_name"`
	MimeTypeKind         pgtype.Text   `json:"mime_type_kind"`
	MimeTypeKey          pgtype.Text   `json:"mime_type_key"`
	ImageHeight          pgtype.Float8 `json:"image_height"`
	ImageWidth           pgtype.Float8 `json:"image_width"`
	ImageID              pgtype.UUID   `json:"image_id"`
	FileID               pgtype.UUID   `json:"file_id"`
}

func (q *Queries) GetAttachableItemsWithMimeType(ctx context.Context, arg GetAttachableItemsWithMimeTypeParams) ([]GetAttachableItemsWithMimeTypeRow, error) {
	rows, err := q.db.Query(ctx, getAttachableItemsWithMimeType,
		arg.WhereInMimeTypeIds,
		arg.InMimeTypeIds,
		arg.WhereInOwnerIds,
		arg.InOwnerIds,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAttachableItemsWithMimeTypeRow{}
	for rows.Next() {
		var i GetAttachableItemsWithMimeTypeRow
		if err := rows.Scan(
			&i.TAttachableItemsPkey,
			&i.AttachableItemID,
			&i.Url,
			&i.Size,
			&i.MimeTypeID,
			&i.OwnerID,
			&i.FromOuter,
			&i.MimeTypeName,
			&i.MimeTypeKind,
			&i.MimeTypeKey,
			&i.ImageHeight,
			&i.ImageWidth,
			&i.ImageID,
			&i.FileID,
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

const getAttachableItemsWithMimeTypeUseKeysetPaginate = `-- name: GetAttachableItemsWithMimeTypeUseKeysetPaginate :many
SELECT t_attachable_items.t_attachable_items_pkey, t_attachable_items.attachable_item_id, t_attachable_items.url, t_attachable_items.size, t_attachable_items.mime_type_id, t_attachable_items.owner_id, t_attachable_items.from_outer, m_mime_types.name mime_type_name, m_mime_types.kind mime_type_kind,
m_mime_types.key mime_type_key, t_images.height image_height, t_images.width image_width, t_images.image_id, t_files.file_id
FROM t_attachable_items
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE
	CASE WHEN $2::boolean = true THEN mime_type_id = ANY($3::uuid[]) ELSE TRUE END
AND
	CASE WHEN $4::boolean = true THEN owner_id = ANY($5::uuid[]) ELSE TRUE END
AND
	CASE $6::text
		WHEN 'next' THEN
			t_attachable_items_pkey > $7::int
		WHEN 'prev' THEN
			t_attachable_items_pkey < $7::int
	END
ORDER BY
	CASE WHEN $6::text = 'next' THEN t_attachable_items_pkey END ASC,
	CASE WHEN $6::text = 'prev' THEN t_attachable_items_pkey END DESC
LIMIT $1
`

type GetAttachableItemsWithMimeTypeUseKeysetPaginateParams struct {
	Limit              int32       `json:"limit"`
	WhereInMimeTypeIds bool        `json:"where_in_mime_type_ids"`
	InMimeTypeIds      []uuid.UUID `json:"in_mime_type_ids"`
	WhereInOwnerIds    bool        `json:"where_in_owner_ids"`
	InOwnerIds         []uuid.UUID `json:"in_owner_ids"`
	CursorDirection    string      `json:"cursor_direction"`
	Cursor             int32       `json:"cursor"`
}

type GetAttachableItemsWithMimeTypeUseKeysetPaginateRow struct {
	TAttachableItemsPkey pgtype.Int8   `json:"t_attachable_items_pkey"`
	AttachableItemID     uuid.UUID     `json:"attachable_item_id"`
	Url                  string        `json:"url"`
	Size                 pgtype.Float8 `json:"size"`
	MimeTypeID           uuid.UUID     `json:"mime_type_id"`
	OwnerID              pgtype.UUID   `json:"owner_id"`
	FromOuter            bool          `json:"from_outer"`
	MimeTypeName         pgtype.Text   `json:"mime_type_name"`
	MimeTypeKind         pgtype.Text   `json:"mime_type_kind"`
	MimeTypeKey          pgtype.Text   `json:"mime_type_key"`
	ImageHeight          pgtype.Float8 `json:"image_height"`
	ImageWidth           pgtype.Float8 `json:"image_width"`
	ImageID              pgtype.UUID   `json:"image_id"`
	FileID               pgtype.UUID   `json:"file_id"`
}

func (q *Queries) GetAttachableItemsWithMimeTypeUseKeysetPaginate(ctx context.Context, arg GetAttachableItemsWithMimeTypeUseKeysetPaginateParams) ([]GetAttachableItemsWithMimeTypeUseKeysetPaginateRow, error) {
	rows, err := q.db.Query(ctx, getAttachableItemsWithMimeTypeUseKeysetPaginate,
		arg.Limit,
		arg.WhereInMimeTypeIds,
		arg.InMimeTypeIds,
		arg.WhereInOwnerIds,
		arg.InOwnerIds,
		arg.CursorDirection,
		arg.Cursor,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAttachableItemsWithMimeTypeUseKeysetPaginateRow{}
	for rows.Next() {
		var i GetAttachableItemsWithMimeTypeUseKeysetPaginateRow
		if err := rows.Scan(
			&i.TAttachableItemsPkey,
			&i.AttachableItemID,
			&i.Url,
			&i.Size,
			&i.MimeTypeID,
			&i.OwnerID,
			&i.FromOuter,
			&i.MimeTypeName,
			&i.MimeTypeKind,
			&i.MimeTypeKey,
			&i.ImageHeight,
			&i.ImageWidth,
			&i.ImageID,
			&i.FileID,
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

const getAttachableItemsWithMimeTypeUseNumberedPaginate = `-- name: GetAttachableItemsWithMimeTypeUseNumberedPaginate :many
SELECT t_attachable_items.t_attachable_items_pkey, t_attachable_items.attachable_item_id, t_attachable_items.url, t_attachable_items.size, t_attachable_items.mime_type_id, t_attachable_items.owner_id, t_attachable_items.from_outer, m_mime_types.name mime_type_name, m_mime_types.kind mime_type_kind,
m_mime_types.key mime_type_key, t_images.height image_height, t_images.width image_width, t_images.image_id, t_files.file_id
FROM t_attachable_items
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE
	CASE WHEN $3::boolean = true THEN mime_type_id = ANY($4::uuid[]) ELSE TRUE END
AND
	CASE WHEN $5::boolean = true THEN owner_id = ANY($6::uuid[]) ELSE TRUE END
ORDER BY
	t_attachable_items_pkey ASC
LIMIT $1 OFFSET $2
`

type GetAttachableItemsWithMimeTypeUseNumberedPaginateParams struct {
	Limit              int32       `json:"limit"`
	Offset             int32       `json:"offset"`
	WhereInMimeTypeIds bool        `json:"where_in_mime_type_ids"`
	InMimeTypeIds      []uuid.UUID `json:"in_mime_type_ids"`
	WhereInOwnerIds    bool        `json:"where_in_owner_ids"`
	InOwnerIds         []uuid.UUID `json:"in_owner_ids"`
}

type GetAttachableItemsWithMimeTypeUseNumberedPaginateRow struct {
	TAttachableItemsPkey pgtype.Int8   `json:"t_attachable_items_pkey"`
	AttachableItemID     uuid.UUID     `json:"attachable_item_id"`
	Url                  string        `json:"url"`
	Size                 pgtype.Float8 `json:"size"`
	MimeTypeID           uuid.UUID     `json:"mime_type_id"`
	OwnerID              pgtype.UUID   `json:"owner_id"`
	FromOuter            bool          `json:"from_outer"`
	MimeTypeName         pgtype.Text   `json:"mime_type_name"`
	MimeTypeKind         pgtype.Text   `json:"mime_type_kind"`
	MimeTypeKey          pgtype.Text   `json:"mime_type_key"`
	ImageHeight          pgtype.Float8 `json:"image_height"`
	ImageWidth           pgtype.Float8 `json:"image_width"`
	ImageID              pgtype.UUID   `json:"image_id"`
	FileID               pgtype.UUID   `json:"file_id"`
}

func (q *Queries) GetAttachableItemsWithMimeTypeUseNumberedPaginate(ctx context.Context, arg GetAttachableItemsWithMimeTypeUseNumberedPaginateParams) ([]GetAttachableItemsWithMimeTypeUseNumberedPaginateRow, error) {
	rows, err := q.db.Query(ctx, getAttachableItemsWithMimeTypeUseNumberedPaginate,
		arg.Limit,
		arg.Offset,
		arg.WhereInMimeTypeIds,
		arg.InMimeTypeIds,
		arg.WhereInOwnerIds,
		arg.InOwnerIds,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAttachableItemsWithMimeTypeUseNumberedPaginateRow{}
	for rows.Next() {
		var i GetAttachableItemsWithMimeTypeUseNumberedPaginateRow
		if err := rows.Scan(
			&i.TAttachableItemsPkey,
			&i.AttachableItemID,
			&i.Url,
			&i.Size,
			&i.MimeTypeID,
			&i.OwnerID,
			&i.FromOuter,
			&i.MimeTypeName,
			&i.MimeTypeKind,
			&i.MimeTypeKey,
			&i.ImageHeight,
			&i.ImageWidth,
			&i.ImageID,
			&i.FileID,
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

const getPluralAttachableItems = `-- name: GetPluralAttachableItems :many
SELECT t_attachable_items.t_attachable_items_pkey, t_attachable_items.attachable_item_id, t_attachable_items.url, t_attachable_items.size, t_attachable_items.mime_type_id, t_attachable_items.owner_id, t_attachable_items.from_outer, t_images.image_id, t_images.height image_height,
t_images.width image_width, t_files.file_id
FROM t_attachable_items
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE attachable_item_id = ANY($1::uuid[])
ORDER BY
	t_attachable_items_pkey ASC
`

type GetPluralAttachableItemsRow struct {
	TAttachableItemsPkey pgtype.Int8   `json:"t_attachable_items_pkey"`
	AttachableItemID     uuid.UUID     `json:"attachable_item_id"`
	Url                  string        `json:"url"`
	Size                 pgtype.Float8 `json:"size"`
	MimeTypeID           uuid.UUID     `json:"mime_type_id"`
	OwnerID              pgtype.UUID   `json:"owner_id"`
	FromOuter            bool          `json:"from_outer"`
	ImageID              pgtype.UUID   `json:"image_id"`
	ImageHeight          pgtype.Float8 `json:"image_height"`
	ImageWidth           pgtype.Float8 `json:"image_width"`
	FileID               pgtype.UUID   `json:"file_id"`
}

func (q *Queries) GetPluralAttachableItems(ctx context.Context, attachableItemIds []uuid.UUID) ([]GetPluralAttachableItemsRow, error) {
	rows, err := q.db.Query(ctx, getPluralAttachableItems, attachableItemIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPluralAttachableItemsRow{}
	for rows.Next() {
		var i GetPluralAttachableItemsRow
		if err := rows.Scan(
			&i.TAttachableItemsPkey,
			&i.AttachableItemID,
			&i.Url,
			&i.Size,
			&i.MimeTypeID,
			&i.OwnerID,
			&i.FromOuter,
			&i.ImageID,
			&i.ImageHeight,
			&i.ImageWidth,
			&i.FileID,
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

const getPluralAttachableItemsUseNumberedPaginate = `-- name: GetPluralAttachableItemsUseNumberedPaginate :many
SELECT t_attachable_items.t_attachable_items_pkey, t_attachable_items.attachable_item_id, t_attachable_items.url, t_attachable_items.size, t_attachable_items.mime_type_id, t_attachable_items.owner_id, t_attachable_items.from_outer, t_images.image_id, t_images.height image_height,
t_images.width image_width, t_files.file_id
FROM t_attachable_items
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE attachable_item_id = ANY($3::uuid[])
ORDER BY
	t_attachable_items_pkey ASC
LIMIT $1 OFFSET $2
`

type GetPluralAttachableItemsUseNumberedPaginateParams struct {
	Limit             int32       `json:"limit"`
	Offset            int32       `json:"offset"`
	AttachableItemIds []uuid.UUID `json:"attachable_item_ids"`
}

type GetPluralAttachableItemsUseNumberedPaginateRow struct {
	TAttachableItemsPkey pgtype.Int8   `json:"t_attachable_items_pkey"`
	AttachableItemID     uuid.UUID     `json:"attachable_item_id"`
	Url                  string        `json:"url"`
	Size                 pgtype.Float8 `json:"size"`
	MimeTypeID           uuid.UUID     `json:"mime_type_id"`
	OwnerID              pgtype.UUID   `json:"owner_id"`
	FromOuter            bool          `json:"from_outer"`
	ImageID              pgtype.UUID   `json:"image_id"`
	ImageHeight          pgtype.Float8 `json:"image_height"`
	ImageWidth           pgtype.Float8 `json:"image_width"`
	FileID               pgtype.UUID   `json:"file_id"`
}

func (q *Queries) GetPluralAttachableItemsUseNumberedPaginate(ctx context.Context, arg GetPluralAttachableItemsUseNumberedPaginateParams) ([]GetPluralAttachableItemsUseNumberedPaginateRow, error) {
	rows, err := q.db.Query(ctx, getPluralAttachableItemsUseNumberedPaginate, arg.Limit, arg.Offset, arg.AttachableItemIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPluralAttachableItemsUseNumberedPaginateRow{}
	for rows.Next() {
		var i GetPluralAttachableItemsUseNumberedPaginateRow
		if err := rows.Scan(
			&i.TAttachableItemsPkey,
			&i.AttachableItemID,
			&i.Url,
			&i.Size,
			&i.MimeTypeID,
			&i.OwnerID,
			&i.FromOuter,
			&i.ImageID,
			&i.ImageHeight,
			&i.ImageWidth,
			&i.FileID,
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

const getPluralAttachableItemsWithMimeType = `-- name: GetPluralAttachableItemsWithMimeType :many
SELECT t_attachable_items.t_attachable_items_pkey, t_attachable_items.attachable_item_id, t_attachable_items.url, t_attachable_items.size, t_attachable_items.mime_type_id, t_attachable_items.owner_id, t_attachable_items.from_outer, m_mime_types.name mime_type_name, m_mime_types.kind mime_type_kind,
m_mime_types.key mime_type_key, t_images.height image_height, t_images.width image_width, t_images.image_id, t_files.file_id
FROM t_attachable_items
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE attachable_item_id = ANY($1::uuid[])
ORDER BY
	t_attachable_items_pkey ASC
`

type GetPluralAttachableItemsWithMimeTypeRow struct {
	TAttachableItemsPkey pgtype.Int8   `json:"t_attachable_items_pkey"`
	AttachableItemID     uuid.UUID     `json:"attachable_item_id"`
	Url                  string        `json:"url"`
	Size                 pgtype.Float8 `json:"size"`
	MimeTypeID           uuid.UUID     `json:"mime_type_id"`
	OwnerID              pgtype.UUID   `json:"owner_id"`
	FromOuter            bool          `json:"from_outer"`
	MimeTypeName         pgtype.Text   `json:"mime_type_name"`
	MimeTypeKind         pgtype.Text   `json:"mime_type_kind"`
	MimeTypeKey          pgtype.Text   `json:"mime_type_key"`
	ImageHeight          pgtype.Float8 `json:"image_height"`
	ImageWidth           pgtype.Float8 `json:"image_width"`
	ImageID              pgtype.UUID   `json:"image_id"`
	FileID               pgtype.UUID   `json:"file_id"`
}

func (q *Queries) GetPluralAttachableItemsWithMimeType(ctx context.Context, attachableItemIds []uuid.UUID) ([]GetPluralAttachableItemsWithMimeTypeRow, error) {
	rows, err := q.db.Query(ctx, getPluralAttachableItemsWithMimeType, attachableItemIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPluralAttachableItemsWithMimeTypeRow{}
	for rows.Next() {
		var i GetPluralAttachableItemsWithMimeTypeRow
		if err := rows.Scan(
			&i.TAttachableItemsPkey,
			&i.AttachableItemID,
			&i.Url,
			&i.Size,
			&i.MimeTypeID,
			&i.OwnerID,
			&i.FromOuter,
			&i.MimeTypeName,
			&i.MimeTypeKind,
			&i.MimeTypeKey,
			&i.ImageHeight,
			&i.ImageWidth,
			&i.ImageID,
			&i.FileID,
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

const getPluralAttachableItemsWithMimeTypeUseNumberedPaginate = `-- name: GetPluralAttachableItemsWithMimeTypeUseNumberedPaginate :many
SELECT t_attachable_items.t_attachable_items_pkey, t_attachable_items.attachable_item_id, t_attachable_items.url, t_attachable_items.size, t_attachable_items.mime_type_id, t_attachable_items.owner_id, t_attachable_items.from_outer, m_mime_types.name mime_type_name, m_mime_types.kind mime_type_kind,
m_mime_types.key mime_type_key, t_images.height image_height, t_images.width image_width, t_images.image_id, t_files.file_id
FROM t_attachable_items
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE attachable_item_id = ANY($3::uuid[])
ORDER BY
	t_attachable_items_pkey ASC
LIMIT $1 OFFSET $2
`

type GetPluralAttachableItemsWithMimeTypeUseNumberedPaginateParams struct {
	Limit             int32       `json:"limit"`
	Offset            int32       `json:"offset"`
	AttachableItemIds []uuid.UUID `json:"attachable_item_ids"`
}

type GetPluralAttachableItemsWithMimeTypeUseNumberedPaginateRow struct {
	TAttachableItemsPkey pgtype.Int8   `json:"t_attachable_items_pkey"`
	AttachableItemID     uuid.UUID     `json:"attachable_item_id"`
	Url                  string        `json:"url"`
	Size                 pgtype.Float8 `json:"size"`
	MimeTypeID           uuid.UUID     `json:"mime_type_id"`
	OwnerID              pgtype.UUID   `json:"owner_id"`
	FromOuter            bool          `json:"from_outer"`
	MimeTypeName         pgtype.Text   `json:"mime_type_name"`
	MimeTypeKind         pgtype.Text   `json:"mime_type_kind"`
	MimeTypeKey          pgtype.Text   `json:"mime_type_key"`
	ImageHeight          pgtype.Float8 `json:"image_height"`
	ImageWidth           pgtype.Float8 `json:"image_width"`
	ImageID              pgtype.UUID   `json:"image_id"`
	FileID               pgtype.UUID   `json:"file_id"`
}

func (q *Queries) GetPluralAttachableItemsWithMimeTypeUseNumberedPaginate(ctx context.Context, arg GetPluralAttachableItemsWithMimeTypeUseNumberedPaginateParams) ([]GetPluralAttachableItemsWithMimeTypeUseNumberedPaginateRow, error) {
	rows, err := q.db.Query(ctx, getPluralAttachableItemsWithMimeTypeUseNumberedPaginate, arg.Limit, arg.Offset, arg.AttachableItemIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPluralAttachableItemsWithMimeTypeUseNumberedPaginateRow{}
	for rows.Next() {
		var i GetPluralAttachableItemsWithMimeTypeUseNumberedPaginateRow
		if err := rows.Scan(
			&i.TAttachableItemsPkey,
			&i.AttachableItemID,
			&i.Url,
			&i.Size,
			&i.MimeTypeID,
			&i.OwnerID,
			&i.FromOuter,
			&i.MimeTypeName,
			&i.MimeTypeKind,
			&i.MimeTypeKey,
			&i.ImageHeight,
			&i.ImageWidth,
			&i.ImageID,
			&i.FileID,
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

const pluralDeleteAttachableItems = `-- name: PluralDeleteAttachableItems :execrows
DELETE FROM t_attachable_items WHERE attachable_item_id = ANY($1::uuid[])
`

func (q *Queries) PluralDeleteAttachableItems(ctx context.Context, attachableItemIds []uuid.UUID) (int64, error) {
	result, err := q.db.Exec(ctx, pluralDeleteAttachableItems, attachableItemIds)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const updateAttachableItem = `-- name: UpdateAttachableItem :one
UPDATE t_attachable_items SET url = $2, size = $3, mime_type_id = $4 WHERE attachable_item_id = $1 RETURNING t_attachable_items_pkey, attachable_item_id, url, size, mime_type_id, owner_id, from_outer
`

type UpdateAttachableItemParams struct {
	AttachableItemID uuid.UUID     `json:"attachable_item_id"`
	Url              string        `json:"url"`
	Size             pgtype.Float8 `json:"size"`
	MimeTypeID       uuid.UUID     `json:"mime_type_id"`
}

func (q *Queries) UpdateAttachableItem(ctx context.Context, arg UpdateAttachableItemParams) (AttachableItem, error) {
	row := q.db.QueryRow(ctx, updateAttachableItem,
		arg.AttachableItemID,
		arg.Url,
		arg.Size,
		arg.MimeTypeID,
	)
	var i AttachableItem
	err := row.Scan(
		&i.TAttachableItemsPkey,
		&i.AttachableItemID,
		&i.Url,
		&i.Size,
		&i.MimeTypeID,
		&i.OwnerID,
		&i.FromOuter,
	)
	return i, err
}
