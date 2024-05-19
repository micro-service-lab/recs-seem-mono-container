// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: image.sql

package query

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const countImages = `-- name: CountImages :one
SELECT COUNT(*) FROM t_images
`

func (q *Queries) CountImages(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, countImages)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createImage = `-- name: CreateImage :one
INSERT INTO t_images (height, width, attachable_item_id) VALUES ($1, $2, $3) RETURNING t_images_pkey, image_id, height, width, attachable_item_id
`

type CreateImageParams struct {
	Height           pgtype.Float8 `json:"height"`
	Width            pgtype.Float8 `json:"width"`
	AttachableItemID uuid.UUID     `json:"attachable_item_id"`
}

func (q *Queries) CreateImage(ctx context.Context, arg CreateImageParams) (Image, error) {
	row := q.db.QueryRow(ctx, createImage, arg.Height, arg.Width, arg.AttachableItemID)
	var i Image
	err := row.Scan(
		&i.TImagesPkey,
		&i.ImageID,
		&i.Height,
		&i.Width,
		&i.AttachableItemID,
	)
	return i, err
}

type CreateImagesParams struct {
	Height           pgtype.Float8 `json:"height"`
	Width            pgtype.Float8 `json:"width"`
	AttachableItemID uuid.UUID     `json:"attachable_item_id"`
}

const deleteImage = `-- name: DeleteImage :execrows
DELETE FROM t_images WHERE image_id = $1
`

func (q *Queries) DeleteImage(ctx context.Context, imageID uuid.UUID) (int64, error) {
	result, err := q.db.Exec(ctx, deleteImage, imageID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const findImageByID = `-- name: FindImageByID :one
SELECT t_images_pkey, image_id, height, width, attachable_item_id FROM t_images WHERE image_id = $1
`

func (q *Queries) FindImageByID(ctx context.Context, imageID uuid.UUID) (Image, error) {
	row := q.db.QueryRow(ctx, findImageByID, imageID)
	var i Image
	err := row.Scan(
		&i.TImagesPkey,
		&i.ImageID,
		&i.Height,
		&i.Width,
		&i.AttachableItemID,
	)
	return i, err
}

const findImageByIDWithAttachableItem = `-- name: FindImageByIDWithAttachableItem :one
SELECT t_images.t_images_pkey, t_images.image_id, t_images.height, t_images.width, t_images.attachable_item_id, t_attachable_items.t_attachable_items_pkey, t_attachable_items.attachable_item_id, t_attachable_items.url, t_attachable_items.size, t_attachable_items.mime_type_id, t_attachable_items.owner_id, t_attachable_items.from_outer, m_mime_types.m_mime_types_pkey, m_mime_types.mime_type_id, m_mime_types.name, m_mime_types.kind, m_mime_types.key FROM t_images
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
WHERE image_id = $1
`

type FindImageByIDWithAttachableItemRow struct {
	Image          Image          `json:"image"`
	AttachableItem AttachableItem `json:"attachable_item"`
	MimeType       MimeType       `json:"mime_type"`
}

func (q *Queries) FindImageByIDWithAttachableItem(ctx context.Context, imageID uuid.UUID) (FindImageByIDWithAttachableItemRow, error) {
	row := q.db.QueryRow(ctx, findImageByIDWithAttachableItem, imageID)
	var i FindImageByIDWithAttachableItemRow
	err := row.Scan(
		&i.Image.TImagesPkey,
		&i.Image.ImageID,
		&i.Image.Height,
		&i.Image.Width,
		&i.Image.AttachableItemID,
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

const getImages = `-- name: GetImages :many
SELECT t_images_pkey, image_id, height, width, attachable_item_id FROM t_images
ORDER BY
	t_images_pkey ASC
`

func (q *Queries) GetImages(ctx context.Context) ([]Image, error) {
	rows, err := q.db.Query(ctx, getImages)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Image{}
	for rows.Next() {
		var i Image
		if err := rows.Scan(
			&i.TImagesPkey,
			&i.ImageID,
			&i.Height,
			&i.Width,
			&i.AttachableItemID,
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

const getImagesUseKeysetPaginate = `-- name: GetImagesUseKeysetPaginate :many
SELECT t_images_pkey, image_id, height, width, attachable_item_id FROM t_images
WHERE
	CASE $2::text
		WHEN 'next' THEN
			t_images_pkey > $3::int
		WHEN 'prev' THEN
			t_images_pkey < $3::int
	END
ORDER BY
	CASE WHEN $2::text = 'next' THEN t_images_pkey END ASC,
	CASE WHEN $2::text = 'prev' THEN t_images_pkey END DESC
LIMIT $1
`

type GetImagesUseKeysetPaginateParams struct {
	Limit           int32  `json:"limit"`
	CursorDirection string `json:"cursor_direction"`
	Cursor          int32  `json:"cursor"`
}

func (q *Queries) GetImagesUseKeysetPaginate(ctx context.Context, arg GetImagesUseKeysetPaginateParams) ([]Image, error) {
	rows, err := q.db.Query(ctx, getImagesUseKeysetPaginate, arg.Limit, arg.CursorDirection, arg.Cursor)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Image{}
	for rows.Next() {
		var i Image
		if err := rows.Scan(
			&i.TImagesPkey,
			&i.ImageID,
			&i.Height,
			&i.Width,
			&i.AttachableItemID,
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

const getImagesUseNumberedPaginate = `-- name: GetImagesUseNumberedPaginate :many
SELECT t_images_pkey, image_id, height, width, attachable_item_id FROM t_images
ORDER BY
	t_images_pkey ASC
LIMIT $1 OFFSET $2
`

type GetImagesUseNumberedPaginateParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) GetImagesUseNumberedPaginate(ctx context.Context, arg GetImagesUseNumberedPaginateParams) ([]Image, error) {
	rows, err := q.db.Query(ctx, getImagesUseNumberedPaginate, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Image{}
	for rows.Next() {
		var i Image
		if err := rows.Scan(
			&i.TImagesPkey,
			&i.ImageID,
			&i.Height,
			&i.Width,
			&i.AttachableItemID,
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

const getImagesWithAttachableItem = `-- name: GetImagesWithAttachableItem :many
SELECT t_images.t_images_pkey, t_images.image_id, t_images.height, t_images.width, t_images.attachable_item_id, t_attachable_items.t_attachable_items_pkey, t_attachable_items.attachable_item_id, t_attachable_items.url, t_attachable_items.size, t_attachable_items.mime_type_id, t_attachable_items.owner_id, t_attachable_items.from_outer, m_mime_types.m_mime_types_pkey, m_mime_types.mime_type_id, m_mime_types.name, m_mime_types.kind, m_mime_types.key FROM t_images
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
ORDER BY
	t_images_pkey ASC
`

type GetImagesWithAttachableItemRow struct {
	Image          Image          `json:"image"`
	AttachableItem AttachableItem `json:"attachable_item"`
	MimeType       MimeType       `json:"mime_type"`
}

func (q *Queries) GetImagesWithAttachableItem(ctx context.Context) ([]GetImagesWithAttachableItemRow, error) {
	rows, err := q.db.Query(ctx, getImagesWithAttachableItem)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetImagesWithAttachableItemRow{}
	for rows.Next() {
		var i GetImagesWithAttachableItemRow
		if err := rows.Scan(
			&i.Image.TImagesPkey,
			&i.Image.ImageID,
			&i.Image.Height,
			&i.Image.Width,
			&i.Image.AttachableItemID,
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

const getImagesWithAttachableItemUseKeysetPaginate = `-- name: GetImagesWithAttachableItemUseKeysetPaginate :many
SELECT t_images.t_images_pkey, t_images.image_id, t_images.height, t_images.width, t_images.attachable_item_id, t_attachable_items.t_attachable_items_pkey, t_attachable_items.attachable_item_id, t_attachable_items.url, t_attachable_items.size, t_attachable_items.mime_type_id, t_attachable_items.owner_id, t_attachable_items.from_outer, m_mime_types.m_mime_types_pkey, m_mime_types.mime_type_id, m_mime_types.name, m_mime_types.kind, m_mime_types.key FROM t_images
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
WHERE
	CASE $2::text
		WHEN 'next' THEN
			t_images_pkey > $3::int
		WHEN 'prev' THEN
			t_images_pkey < $3::int
	END
ORDER BY
	CASE WHEN $2::text = 'next' THEN t_images_pkey END ASC,
	CASE WHEN $2::text = 'prev' THEN t_images_pkey END DESC
LIMIT $1
`

type GetImagesWithAttachableItemUseKeysetPaginateParams struct {
	Limit           int32  `json:"limit"`
	CursorDirection string `json:"cursor_direction"`
	Cursor          int32  `json:"cursor"`
}

type GetImagesWithAttachableItemUseKeysetPaginateRow struct {
	Image          Image          `json:"image"`
	AttachableItem AttachableItem `json:"attachable_item"`
	MimeType       MimeType       `json:"mime_type"`
}

func (q *Queries) GetImagesWithAttachableItemUseKeysetPaginate(ctx context.Context, arg GetImagesWithAttachableItemUseKeysetPaginateParams) ([]GetImagesWithAttachableItemUseKeysetPaginateRow, error) {
	rows, err := q.db.Query(ctx, getImagesWithAttachableItemUseKeysetPaginate, arg.Limit, arg.CursorDirection, arg.Cursor)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetImagesWithAttachableItemUseKeysetPaginateRow{}
	for rows.Next() {
		var i GetImagesWithAttachableItemUseKeysetPaginateRow
		if err := rows.Scan(
			&i.Image.TImagesPkey,
			&i.Image.ImageID,
			&i.Image.Height,
			&i.Image.Width,
			&i.Image.AttachableItemID,
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

const getImagesWithAttachableItemUseNumberedPaginate = `-- name: GetImagesWithAttachableItemUseNumberedPaginate :many
SELECT t_images.t_images_pkey, t_images.image_id, t_images.height, t_images.width, t_images.attachable_item_id, t_attachable_items.t_attachable_items_pkey, t_attachable_items.attachable_item_id, t_attachable_items.url, t_attachable_items.size, t_attachable_items.mime_type_id, t_attachable_items.owner_id, t_attachable_items.from_outer, m_mime_types.m_mime_types_pkey, m_mime_types.mime_type_id, m_mime_types.name, m_mime_types.kind, m_mime_types.key FROM t_images
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
ORDER BY
	t_images_pkey ASC
LIMIT $1 OFFSET $2
`

type GetImagesWithAttachableItemUseNumberedPaginateParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type GetImagesWithAttachableItemUseNumberedPaginateRow struct {
	Image          Image          `json:"image"`
	AttachableItem AttachableItem `json:"attachable_item"`
	MimeType       MimeType       `json:"mime_type"`
}

func (q *Queries) GetImagesWithAttachableItemUseNumberedPaginate(ctx context.Context, arg GetImagesWithAttachableItemUseNumberedPaginateParams) ([]GetImagesWithAttachableItemUseNumberedPaginateRow, error) {
	rows, err := q.db.Query(ctx, getImagesWithAttachableItemUseNumberedPaginate, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetImagesWithAttachableItemUseNumberedPaginateRow{}
	for rows.Next() {
		var i GetImagesWithAttachableItemUseNumberedPaginateRow
		if err := rows.Scan(
			&i.Image.TImagesPkey,
			&i.Image.ImageID,
			&i.Image.Height,
			&i.Image.Width,
			&i.Image.AttachableItemID,
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

const getPluralImages = `-- name: GetPluralImages :many
SELECT t_images_pkey, image_id, height, width, attachable_item_id FROM t_images
WHERE attachable_item_id = ANY($1::uuid[])
ORDER BY
	t_images_pkey ASC
`

func (q *Queries) GetPluralImages(ctx context.Context, attachableItemIds []uuid.UUID) ([]Image, error) {
	rows, err := q.db.Query(ctx, getPluralImages, attachableItemIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Image{}
	for rows.Next() {
		var i Image
		if err := rows.Scan(
			&i.TImagesPkey,
			&i.ImageID,
			&i.Height,
			&i.Width,
			&i.AttachableItemID,
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

const getPluralImagesUseNumberedPaginate = `-- name: GetPluralImagesUseNumberedPaginate :many
SELECT t_images_pkey, image_id, height, width, attachable_item_id FROM t_images
WHERE attachable_item_id = ANY($3::uuid[])
ORDER BY
	t_images_pkey ASC
LIMIT $1 OFFSET $2
`

type GetPluralImagesUseNumberedPaginateParams struct {
	Limit             int32       `json:"limit"`
	Offset            int32       `json:"offset"`
	AttachableItemIds []uuid.UUID `json:"attachable_item_ids"`
}

func (q *Queries) GetPluralImagesUseNumberedPaginate(ctx context.Context, arg GetPluralImagesUseNumberedPaginateParams) ([]Image, error) {
	rows, err := q.db.Query(ctx, getPluralImagesUseNumberedPaginate, arg.Limit, arg.Offset, arg.AttachableItemIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Image{}
	for rows.Next() {
		var i Image
		if err := rows.Scan(
			&i.TImagesPkey,
			&i.ImageID,
			&i.Height,
			&i.Width,
			&i.AttachableItemID,
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

const getPluralImagesWithAttachableItem = `-- name: GetPluralImagesWithAttachableItem :many
SELECT t_images.t_images_pkey, t_images.image_id, t_images.height, t_images.width, t_images.attachable_item_id, t_attachable_items.t_attachable_items_pkey, t_attachable_items.attachable_item_id, t_attachable_items.url, t_attachable_items.size, t_attachable_items.mime_type_id, t_attachable_items.owner_id, t_attachable_items.from_outer, m_mime_types.m_mime_types_pkey, m_mime_types.mime_type_id, m_mime_types.name, m_mime_types.kind, m_mime_types.key FROM t_images
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
WHERE attachable_item_id = ANY($1::uuid[])
ORDER BY
	t_images_pkey ASC
`

type GetPluralImagesWithAttachableItemRow struct {
	Image          Image          `json:"image"`
	AttachableItem AttachableItem `json:"attachable_item"`
	MimeType       MimeType       `json:"mime_type"`
}

func (q *Queries) GetPluralImagesWithAttachableItem(ctx context.Context, attachableItemIds []uuid.UUID) ([]GetPluralImagesWithAttachableItemRow, error) {
	rows, err := q.db.Query(ctx, getPluralImagesWithAttachableItem, attachableItemIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPluralImagesWithAttachableItemRow{}
	for rows.Next() {
		var i GetPluralImagesWithAttachableItemRow
		if err := rows.Scan(
			&i.Image.TImagesPkey,
			&i.Image.ImageID,
			&i.Image.Height,
			&i.Image.Width,
			&i.Image.AttachableItemID,
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

const getPluralImagesWithAttachableItemUseNumberedPaginate = `-- name: GetPluralImagesWithAttachableItemUseNumberedPaginate :many
SELECT t_images.t_images_pkey, t_images.image_id, t_images.height, t_images.width, t_images.attachable_item_id, t_attachable_items.t_attachable_items_pkey, t_attachable_items.attachable_item_id, t_attachable_items.url, t_attachable_items.size, t_attachable_items.mime_type_id, t_attachable_items.owner_id, t_attachable_items.from_outer, m_mime_types.m_mime_types_pkey, m_mime_types.mime_type_id, m_mime_types.name, m_mime_types.kind, m_mime_types.key FROM t_images
LEFT JOIN t_attachable_items ON t_images.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
WHERE attachable_item_id = ANY($3::uuid[])
ORDER BY
	t_images_pkey ASC
LIMIT $1 OFFSET $2
`

type GetPluralImagesWithAttachableItemUseNumberedPaginateParams struct {
	Limit             int32       `json:"limit"`
	Offset            int32       `json:"offset"`
	AttachableItemIds []uuid.UUID `json:"attachable_item_ids"`
}

type GetPluralImagesWithAttachableItemUseNumberedPaginateRow struct {
	Image          Image          `json:"image"`
	AttachableItem AttachableItem `json:"attachable_item"`
	MimeType       MimeType       `json:"mime_type"`
}

func (q *Queries) GetPluralImagesWithAttachableItemUseNumberedPaginate(ctx context.Context, arg GetPluralImagesWithAttachableItemUseNumberedPaginateParams) ([]GetPluralImagesWithAttachableItemUseNumberedPaginateRow, error) {
	rows, err := q.db.Query(ctx, getPluralImagesWithAttachableItemUseNumberedPaginate, arg.Limit, arg.Offset, arg.AttachableItemIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPluralImagesWithAttachableItemUseNumberedPaginateRow{}
	for rows.Next() {
		var i GetPluralImagesWithAttachableItemUseNumberedPaginateRow
		if err := rows.Scan(
			&i.Image.TImagesPkey,
			&i.Image.ImageID,
			&i.Image.Height,
			&i.Image.Width,
			&i.Image.AttachableItemID,
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

const pluralDeleteImages = `-- name: PluralDeleteImages :execrows
DELETE FROM t_images WHERE image_id = ANY($1::uuid[])
`

func (q *Queries) PluralDeleteImages(ctx context.Context, imageIds []uuid.UUID) (int64, error) {
	result, err := q.db.Exec(ctx, pluralDeleteImages, imageIds)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}
