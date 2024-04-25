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

const deleteFile = `-- name: DeleteFile :exec
DELETE FROM t_files WHERE file_id = $1
`

func (q *Queries) DeleteFile(ctx context.Context, fileID uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteFile, fileID)
	return err
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
SELECT t_files.t_files_pkey, t_files.file_id, t_files.attachable_item_id, t_attachable_items.t_attachable_items_pkey, t_attachable_items.attachable_item_id, t_attachable_items.url, t_attachable_items.size, t_attachable_items.mime_type_id, m_mime_types.m_mime_types_pkey, m_mime_types.mime_type_id, m_mime_types.name, m_mime_types.key FROM t_files
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
		&i.MimeType.MMimeTypesPkey,
		&i.MimeType.MimeTypeID,
		&i.MimeType.Name,
		&i.MimeType.Key,
	)
	return i, err
}

const getFiles = `-- name: GetFiles :many
SELECT t_files_pkey, file_id, attachable_item_id FROM t_files
ORDER BY
	t_files_pkey DESC
LIMIT $1 OFFSET $2
`

type GetFilesParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) GetFiles(ctx context.Context, arg GetFilesParams) ([]File, error) {
	rows, err := q.db.Query(ctx, getFiles, arg.Limit, arg.Offset)
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
SELECT t_files.t_files_pkey, t_files.file_id, t_files.attachable_item_id, t_attachable_items.t_attachable_items_pkey, t_attachable_items.attachable_item_id, t_attachable_items.url, t_attachable_items.size, t_attachable_items.mime_type_id, m_mime_types.m_mime_types_pkey, m_mime_types.mime_type_id, m_mime_types.name, m_mime_types.key FROM t_files
LEFT JOIN t_attachable_items ON t_files.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
ORDER BY
	t_files_pkey DESC
LIMIT $1 OFFSET $2
`

type GetFilesWithAttachableItemParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type GetFilesWithAttachableItemRow struct {
	File           File           `json:"file"`
	AttachableItem AttachableItem `json:"attachable_item"`
	MimeType       MimeType       `json:"mime_type"`
}

func (q *Queries) GetFilesWithAttachableItem(ctx context.Context, arg GetFilesWithAttachableItemParams) ([]GetFilesWithAttachableItemRow, error) {
	rows, err := q.db.Query(ctx, getFilesWithAttachableItem, arg.Limit, arg.Offset)
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
			&i.MimeType.MMimeTypesPkey,
			&i.MimeType.MimeTypeID,
			&i.MimeType.Name,
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
