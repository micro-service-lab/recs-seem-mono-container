// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: file.sql

package query

import (
	"context"

	"github.com/google/uuid"
)

const countFiles = `-- name: CountFiles :one
SELECT COUNT(*) FROM t_files
`

func (q *Queries) CountFiles(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, countFiles)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createFile = `-- name: CreateFile :one
INSERT INTO t_files (attachable_item_id) VALUES ($1) RETURNING t_files_pkey, file_id, attachable_item_id
`

func (q *Queries) CreateFile(ctx context.Context, attachableItemID uuid.UUID) (File, error) {
	row := q.db.QueryRow(ctx, createFile, attachableItemID)
	var i File
	err := row.Scan(&i.TFilesPkey, &i.FileID, &i.AttachableItemID)
	return i, err
}

const deleteFile = `-- name: DeleteFile :execrows
DELETE FROM t_files WHERE file_id = $1
`

func (q *Queries) DeleteFile(ctx context.Context, fileID uuid.UUID) (int64, error) {
	result, err := q.db.Exec(ctx, deleteFile, fileID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const findFileByID = `-- name: FindFileByID :one
SELECT t_files_pkey, file_id, attachable_item_id FROM t_files WHERE file_id = $1
`

func (q *Queries) FindFileByID(ctx context.Context, fileID uuid.UUID) (File, error) {
	row := q.db.QueryRow(ctx, findFileByID, fileID)
	var i File
	err := row.Scan(&i.TFilesPkey, &i.FileID, &i.AttachableItemID)
	return i, err
}

const findFileByIDWithAttachableItem = `-- name: FindFileByIDWithAttachableItem :one
SELECT t_files.t_files_pkey, t_files.file_id, t_files.attachable_item_id, t_attachable_items.t_attachable_items_pkey, t_attachable_items.attachable_item_id, t_attachable_items.url, t_attachable_items.size, t_attachable_items.mime_type_id, t_attachable_items.owner_id, t_attachable_items.from_outer, m_mime_types.m_mime_types_pkey, m_mime_types.mime_type_id, m_mime_types.name, m_mime_types.kind, m_mime_types.key FROM t_files
LEFT JOIN t_attachable_items ON t_files.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
WHERE file_id = $1
`

type FindFileByIDWithAttachableItemRow struct {
	File           File           `json:"file"`
	AttachableItem AttachableItem `json:"attachable_item"`
	MimeType       MimeType       `json:"mime_type"`
}

func (q *Queries) FindFileByIDWithAttachableItem(ctx context.Context, fileID uuid.UUID) (FindFileByIDWithAttachableItemRow, error) {
	row := q.db.QueryRow(ctx, findFileByIDWithAttachableItem, fileID)
	var i FindFileByIDWithAttachableItemRow
	err := row.Scan(
		&i.File.TFilesPkey,
		&i.File.FileID,
		&i.File.AttachableItemID,
		&i.AttachableItem.TAttachableItemsPkey,
		&i.AttachableItem.AttachableItemID,
		&i.AttachableItem.Url,
		&i.AttachableItem.Size,
		&i.AttachableItem.MimeTypeID,
		&i.AttachableItem.OwnerID,
		&i.AttachableItem.FromOuter,
		&i.MimeType.MMimeTypesPkey,
		&i.MimeType.MimeTypeID,
		&i.MimeType.Name,
		&i.MimeType.Kind,
		&i.MimeType.Key,
	)
	return i, err
}

const getFiles = `-- name: GetFiles :many
SELECT t_files_pkey, file_id, attachable_item_id FROM t_files
ORDER BY
	t_files_pkey ASC
`

func (q *Queries) GetFiles(ctx context.Context) ([]File, error) {
	rows, err := q.db.Query(ctx, getFiles)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []File{}
	for rows.Next() {
		var i File
		if err := rows.Scan(&i.TFilesPkey, &i.FileID, &i.AttachableItemID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFilesUseKeysetPaginate = `-- name: GetFilesUseKeysetPaginate :many
SELECT t_files_pkey, file_id, attachable_item_id FROM t_files
WHERE
	CASE $2::text
		WHEN 'next' THEN
			t_files_pkey > $3::int
		WHEN 'prev' THEN
			t_files_pkey < $3::int
	END
ORDER BY
	CASE WHEN $2::text = 'next' THEN t_files_pkey END ASC,
	CASE WHEN $2::text = 'prev' THEN t_files_pkey END DESC
LIMIT $1
`

type GetFilesUseKeysetPaginateParams struct {
	Limit           int32  `json:"limit"`
	CursorDirection string `json:"cursor_direction"`
	Cursor          int32  `json:"cursor"`
}

func (q *Queries) GetFilesUseKeysetPaginate(ctx context.Context, arg GetFilesUseKeysetPaginateParams) ([]File, error) {
	rows, err := q.db.Query(ctx, getFilesUseKeysetPaginate, arg.Limit, arg.CursorDirection, arg.Cursor)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []File{}
	for rows.Next() {
		var i File
		if err := rows.Scan(&i.TFilesPkey, &i.FileID, &i.AttachableItemID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFilesUseNumberedPaginate = `-- name: GetFilesUseNumberedPaginate :many
SELECT t_files_pkey, file_id, attachable_item_id FROM t_files
ORDER BY
	t_files_pkey ASC
LIMIT $1 OFFSET $2
`

type GetFilesUseNumberedPaginateParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) GetFilesUseNumberedPaginate(ctx context.Context, arg GetFilesUseNumberedPaginateParams) ([]File, error) {
	rows, err := q.db.Query(ctx, getFilesUseNumberedPaginate, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []File{}
	for rows.Next() {
		var i File
		if err := rows.Scan(&i.TFilesPkey, &i.FileID, &i.AttachableItemID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFilesWithAttachableItem = `-- name: GetFilesWithAttachableItem :many
SELECT t_files.t_files_pkey, t_files.file_id, t_files.attachable_item_id, t_attachable_items.t_attachable_items_pkey, t_attachable_items.attachable_item_id, t_attachable_items.url, t_attachable_items.size, t_attachable_items.mime_type_id, t_attachable_items.owner_id, t_attachable_items.from_outer, m_mime_types.m_mime_types_pkey, m_mime_types.mime_type_id, m_mime_types.name, m_mime_types.kind, m_mime_types.key FROM t_files
LEFT JOIN t_attachable_items ON t_files.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
ORDER BY
	t_files_pkey ASC
`

type GetFilesWithAttachableItemRow struct {
	File           File           `json:"file"`
	AttachableItem AttachableItem `json:"attachable_item"`
	MimeType       MimeType       `json:"mime_type"`
}

func (q *Queries) GetFilesWithAttachableItem(ctx context.Context) ([]GetFilesWithAttachableItemRow, error) {
	rows, err := q.db.Query(ctx, getFilesWithAttachableItem)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetFilesWithAttachableItemRow{}
	for rows.Next() {
		var i GetFilesWithAttachableItemRow
		if err := rows.Scan(
			&i.File.TFilesPkey,
			&i.File.FileID,
			&i.File.AttachableItemID,
			&i.AttachableItem.TAttachableItemsPkey,
			&i.AttachableItem.AttachableItemID,
			&i.AttachableItem.Url,
			&i.AttachableItem.Size,
			&i.AttachableItem.MimeTypeID,
			&i.AttachableItem.OwnerID,
			&i.AttachableItem.FromOuter,
			&i.MimeType.MMimeTypesPkey,
			&i.MimeType.MimeTypeID,
			&i.MimeType.Name,
			&i.MimeType.Kind,
			&i.MimeType.Key,
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

const getFilesWithAttachableItemUseKeysetPaginate = `-- name: GetFilesWithAttachableItemUseKeysetPaginate :many
SELECT t_files.t_files_pkey, t_files.file_id, t_files.attachable_item_id, t_attachable_items.t_attachable_items_pkey, t_attachable_items.attachable_item_id, t_attachable_items.url, t_attachable_items.size, t_attachable_items.mime_type_id, t_attachable_items.owner_id, t_attachable_items.from_outer, m_mime_types.m_mime_types_pkey, m_mime_types.mime_type_id, m_mime_types.name, m_mime_types.kind, m_mime_types.key FROM t_files
LEFT JOIN t_attachable_items ON t_files.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
WHERE
	CASE $2::text
		WHEN 'next' THEN
			t_files_pkey > $3::int
		WHEN 'prev' THEN
			t_files_pkey < $3::int
	END
ORDER BY
	CASE WHEN $2::text = 'next' THEN t_files_pkey END ASC,
	CASE WHEN $2::text = 'prev' THEN t_files_pkey END DESC
LIMIT $1
`

type GetFilesWithAttachableItemUseKeysetPaginateParams struct {
	Limit           int32  `json:"limit"`
	CursorDirection string `json:"cursor_direction"`
	Cursor          int32  `json:"cursor"`
}

type GetFilesWithAttachableItemUseKeysetPaginateRow struct {
	File           File           `json:"file"`
	AttachableItem AttachableItem `json:"attachable_item"`
	MimeType       MimeType       `json:"mime_type"`
}

func (q *Queries) GetFilesWithAttachableItemUseKeysetPaginate(ctx context.Context, arg GetFilesWithAttachableItemUseKeysetPaginateParams) ([]GetFilesWithAttachableItemUseKeysetPaginateRow, error) {
	rows, err := q.db.Query(ctx, getFilesWithAttachableItemUseKeysetPaginate, arg.Limit, arg.CursorDirection, arg.Cursor)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetFilesWithAttachableItemUseKeysetPaginateRow{}
	for rows.Next() {
		var i GetFilesWithAttachableItemUseKeysetPaginateRow
		if err := rows.Scan(
			&i.File.TFilesPkey,
			&i.File.FileID,
			&i.File.AttachableItemID,
			&i.AttachableItem.TAttachableItemsPkey,
			&i.AttachableItem.AttachableItemID,
			&i.AttachableItem.Url,
			&i.AttachableItem.Size,
			&i.AttachableItem.MimeTypeID,
			&i.AttachableItem.OwnerID,
			&i.AttachableItem.FromOuter,
			&i.MimeType.MMimeTypesPkey,
			&i.MimeType.MimeTypeID,
			&i.MimeType.Name,
			&i.MimeType.Kind,
			&i.MimeType.Key,
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

const getFilesWithAttachableItemUseNumberedPaginate = `-- name: GetFilesWithAttachableItemUseNumberedPaginate :many
SELECT t_files.t_files_pkey, t_files.file_id, t_files.attachable_item_id, t_attachable_items.t_attachable_items_pkey, t_attachable_items.attachable_item_id, t_attachable_items.url, t_attachable_items.size, t_attachable_items.mime_type_id, t_attachable_items.owner_id, t_attachable_items.from_outer, m_mime_types.m_mime_types_pkey, m_mime_types.mime_type_id, m_mime_types.name, m_mime_types.kind, m_mime_types.key FROM t_files
LEFT JOIN t_attachable_items ON t_files.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
ORDER BY
	t_files_pkey ASC
LIMIT $1 OFFSET $2
`

type GetFilesWithAttachableItemUseNumberedPaginateParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type GetFilesWithAttachableItemUseNumberedPaginateRow struct {
	File           File           `json:"file"`
	AttachableItem AttachableItem `json:"attachable_item"`
	MimeType       MimeType       `json:"mime_type"`
}

func (q *Queries) GetFilesWithAttachableItemUseNumberedPaginate(ctx context.Context, arg GetFilesWithAttachableItemUseNumberedPaginateParams) ([]GetFilesWithAttachableItemUseNumberedPaginateRow, error) {
	rows, err := q.db.Query(ctx, getFilesWithAttachableItemUseNumberedPaginate, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetFilesWithAttachableItemUseNumberedPaginateRow{}
	for rows.Next() {
		var i GetFilesWithAttachableItemUseNumberedPaginateRow
		if err := rows.Scan(
			&i.File.TFilesPkey,
			&i.File.FileID,
			&i.File.AttachableItemID,
			&i.AttachableItem.TAttachableItemsPkey,
			&i.AttachableItem.AttachableItemID,
			&i.AttachableItem.Url,
			&i.AttachableItem.Size,
			&i.AttachableItem.MimeTypeID,
			&i.AttachableItem.OwnerID,
			&i.AttachableItem.FromOuter,
			&i.MimeType.MMimeTypesPkey,
			&i.MimeType.MimeTypeID,
			&i.MimeType.Name,
			&i.MimeType.Kind,
			&i.MimeType.Key,
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

const getPluralFiles = `-- name: GetPluralFiles :many
SELECT t_files_pkey, file_id, attachable_item_id FROM t_files
WHERE attachable_item_id = ANY($1::uuid[])
ORDER BY
	t_files_pkey ASC
`

func (q *Queries) GetPluralFiles(ctx context.Context, attachableItemIds []uuid.UUID) ([]File, error) {
	rows, err := q.db.Query(ctx, getPluralFiles, attachableItemIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []File{}
	for rows.Next() {
		var i File
		if err := rows.Scan(&i.TFilesPkey, &i.FileID, &i.AttachableItemID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPluralFilesUseNumberedPaginate = `-- name: GetPluralFilesUseNumberedPaginate :many
SELECT t_files_pkey, file_id, attachable_item_id FROM t_files
WHERE attachable_item_id = ANY($3::uuid[])
ORDER BY
	t_files_pkey ASC
LIMIT $1 OFFSET $2
`

type GetPluralFilesUseNumberedPaginateParams struct {
	Limit             int32       `json:"limit"`
	Offset            int32       `json:"offset"`
	AttachableItemIds []uuid.UUID `json:"attachable_item_ids"`
}

func (q *Queries) GetPluralFilesUseNumberedPaginate(ctx context.Context, arg GetPluralFilesUseNumberedPaginateParams) ([]File, error) {
	rows, err := q.db.Query(ctx, getPluralFilesUseNumberedPaginate, arg.Limit, arg.Offset, arg.AttachableItemIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []File{}
	for rows.Next() {
		var i File
		if err := rows.Scan(&i.TFilesPkey, &i.FileID, &i.AttachableItemID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPluralFilesWithAttachableItem = `-- name: GetPluralFilesWithAttachableItem :many
SELECT t_files.t_files_pkey, t_files.file_id, t_files.attachable_item_id, t_attachable_items.t_attachable_items_pkey, t_attachable_items.attachable_item_id, t_attachable_items.url, t_attachable_items.size, t_attachable_items.mime_type_id, t_attachable_items.owner_id, t_attachable_items.from_outer, m_mime_types.m_mime_types_pkey, m_mime_types.mime_type_id, m_mime_types.name, m_mime_types.kind, m_mime_types.key FROM t_files
LEFT JOIN t_attachable_items ON t_files.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
WHERE attachable_item_id = ANY($1::uuid[])
ORDER BY
	t_files_pkey ASC
`

type GetPluralFilesWithAttachableItemRow struct {
	File           File           `json:"file"`
	AttachableItem AttachableItem `json:"attachable_item"`
	MimeType       MimeType       `json:"mime_type"`
}

func (q *Queries) GetPluralFilesWithAttachableItem(ctx context.Context, attachableItemIds []uuid.UUID) ([]GetPluralFilesWithAttachableItemRow, error) {
	rows, err := q.db.Query(ctx, getPluralFilesWithAttachableItem, attachableItemIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPluralFilesWithAttachableItemRow{}
	for rows.Next() {
		var i GetPluralFilesWithAttachableItemRow
		if err := rows.Scan(
			&i.File.TFilesPkey,
			&i.File.FileID,
			&i.File.AttachableItemID,
			&i.AttachableItem.TAttachableItemsPkey,
			&i.AttachableItem.AttachableItemID,
			&i.AttachableItem.Url,
			&i.AttachableItem.Size,
			&i.AttachableItem.MimeTypeID,
			&i.AttachableItem.OwnerID,
			&i.AttachableItem.FromOuter,
			&i.MimeType.MMimeTypesPkey,
			&i.MimeType.MimeTypeID,
			&i.MimeType.Name,
			&i.MimeType.Kind,
			&i.MimeType.Key,
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

const getPluralFilesWithAttachableItemUseNumberedPaginate = `-- name: GetPluralFilesWithAttachableItemUseNumberedPaginate :many
SELECT t_files.t_files_pkey, t_files.file_id, t_files.attachable_item_id, t_attachable_items.t_attachable_items_pkey, t_attachable_items.attachable_item_id, t_attachable_items.url, t_attachable_items.size, t_attachable_items.mime_type_id, t_attachable_items.owner_id, t_attachable_items.from_outer, m_mime_types.m_mime_types_pkey, m_mime_types.mime_type_id, m_mime_types.name, m_mime_types.kind, m_mime_types.key FROM t_files
LEFT JOIN t_attachable_items ON t_files.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
WHERE attachable_item_id = ANY($3::uuid[])
ORDER BY
	t_files_pkey ASC
LIMIT $1 OFFSET $2
`

type GetPluralFilesWithAttachableItemUseNumberedPaginateParams struct {
	Limit             int32       `json:"limit"`
	Offset            int32       `json:"offset"`
	AttachableItemIds []uuid.UUID `json:"attachable_item_ids"`
}

type GetPluralFilesWithAttachableItemUseNumberedPaginateRow struct {
	File           File           `json:"file"`
	AttachableItem AttachableItem `json:"attachable_item"`
	MimeType       MimeType       `json:"mime_type"`
}

func (q *Queries) GetPluralFilesWithAttachableItemUseNumberedPaginate(ctx context.Context, arg GetPluralFilesWithAttachableItemUseNumberedPaginateParams) ([]GetPluralFilesWithAttachableItemUseNumberedPaginateRow, error) {
	rows, err := q.db.Query(ctx, getPluralFilesWithAttachableItemUseNumberedPaginate, arg.Limit, arg.Offset, arg.AttachableItemIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPluralFilesWithAttachableItemUseNumberedPaginateRow{}
	for rows.Next() {
		var i GetPluralFilesWithAttachableItemUseNumberedPaginateRow
		if err := rows.Scan(
			&i.File.TFilesPkey,
			&i.File.FileID,
			&i.File.AttachableItemID,
			&i.AttachableItem.TAttachableItemsPkey,
			&i.AttachableItem.AttachableItemID,
			&i.AttachableItem.Url,
			&i.AttachableItem.Size,
			&i.AttachableItem.MimeTypeID,
			&i.AttachableItem.OwnerID,
			&i.AttachableItem.FromOuter,
			&i.MimeType.MMimeTypesPkey,
			&i.MimeType.MimeTypeID,
			&i.MimeType.Name,
			&i.MimeType.Kind,
			&i.MimeType.Key,
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

const pluralDeleteFiles = `-- name: PluralDeleteFiles :execrows
DELETE FROM t_files WHERE file_id = ANY($1::uuid[])
`

func (q *Queries) PluralDeleteFiles(ctx context.Context, fileIds []uuid.UUID) (int64, error) {
	result, err := q.db.Exec(ctx, pluralDeleteFiles, fileIds)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}
