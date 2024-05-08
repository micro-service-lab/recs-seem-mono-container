// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: attached_message.sql

package query

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const countAttachableItemsOnMessage = `-- name: CountAttachableItemsOnMessage :one
SELECT COUNT(*) FROM t_attached_messages WHERE message_id = $1
`

func (q *Queries) CountAttachableItemsOnMessage(ctx context.Context, messageID pgtype.UUID) (int64, error) {
	row := q.db.QueryRow(ctx, countAttachableItemsOnMessage, messageID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countAttachedMessagesOnChatRoom = `-- name: CountAttachedMessagesOnChatRoom :one
SELECT COUNT(*) FROM t_attached_messages
WHERE message_id IN (
	SELECT message_id FROM t_messages WHERE chat_room_id = $1
)
`

func (q *Queries) CountAttachedMessagesOnChatRoom(ctx context.Context, chatRoomID uuid.UUID) (int64, error) {
	row := q.db.QueryRow(ctx, countAttachedMessagesOnChatRoom, chatRoomID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createAttachedMessage = `-- name: CreateAttachedMessage :one
INSERT INTO t_attached_messages (message_id, attachable_item_id) VALUES ($1, $2) RETURNING t_attached_messages_pkey, message_id, attachable_item_id
`

type CreateAttachedMessageParams struct {
	MessageID        pgtype.UUID `json:"message_id"`
	AttachableItemID uuid.UUID   `json:"attachable_item_id"`
}

func (q *Queries) CreateAttachedMessage(ctx context.Context, arg CreateAttachedMessageParams) (AttachedMessage, error) {
	row := q.db.QueryRow(ctx, createAttachedMessage, arg.MessageID, arg.AttachableItemID)
	var i AttachedMessage
	err := row.Scan(&i.TAttachedMessagesPkey, &i.MessageID, &i.AttachableItemID)
	return i, err
}

type CreateAttachedMessagesParams struct {
	MessageID        pgtype.UUID `json:"message_id"`
	AttachableItemID uuid.UUID   `json:"attachable_item_id"`
}

const deleteAttachedMessage = `-- name: DeleteAttachedMessage :exec
DELETE FROM t_attached_messages WHERE message_id = $1 AND attachable_item_id = $2
`

type DeleteAttachedMessageParams struct {
	MessageID        pgtype.UUID `json:"message_id"`
	AttachableItemID uuid.UUID   `json:"attachable_item_id"`
}

func (q *Queries) DeleteAttachedMessage(ctx context.Context, arg DeleteAttachedMessageParams) error {
	_, err := q.db.Exec(ctx, deleteAttachedMessage, arg.MessageID, arg.AttachableItemID)
	return err
}

const deleteAttachedMessagesOnAttachableItem = `-- name: DeleteAttachedMessagesOnAttachableItem :exec
DELETE FROM t_attached_messages WHERE attachable_item_id = $1
`

func (q *Queries) DeleteAttachedMessagesOnAttachableItem(ctx context.Context, attachableItemID uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteAttachedMessagesOnAttachableItem, attachableItemID)
	return err
}

const deleteAttachedMessagesOnAttachableItems = `-- name: DeleteAttachedMessagesOnAttachableItems :exec
DELETE FROM t_attached_messages WHERE attachable_item_id = ANY($1::uuid[])
`

func (q *Queries) DeleteAttachedMessagesOnAttachableItems(ctx context.Context, dollar_1 []uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteAttachedMessagesOnAttachableItems, dollar_1)
	return err
}

const deleteAttachedMessagesOnMessage = `-- name: DeleteAttachedMessagesOnMessage :exec
DELETE FROM t_attached_messages WHERE message_id = $1
`

func (q *Queries) DeleteAttachedMessagesOnMessage(ctx context.Context, messageID pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteAttachedMessagesOnMessage, messageID)
	return err
}

const deleteAttachedMessagesOnMessages = `-- name: DeleteAttachedMessagesOnMessages :exec
DELETE FROM t_attached_messages WHERE message_id = ANY($1::uuid[])
`

func (q *Queries) DeleteAttachedMessagesOnMessages(ctx context.Context, dollar_1 []uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteAttachedMessagesOnMessages, dollar_1)
	return err
}

const getAttachableItemsOnMessage = `-- name: GetAttachableItemsOnMessage :many
SELECT t_attachable_items.t_attachable_items_pkey, t_attachable_items.attachable_item_id, t_attachable_items.url, t_attachable_items.size, t_attachable_items.mime_type_id, t_images.image_id, t_images.height as image_height, t_images.width as image_width, t_files.file_id FROM t_attached_messages
LEFT JOIN t_attachable_items ON t_attached_messages.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE message_id = $1
ORDER BY
	t_attached_messages_pkey ASC
`

type GetAttachableItemsOnMessageRow struct {
	TAttachableItemsPkey pgtype.Int8   `json:"t_attachable_items_pkey"`
	AttachableItemID     pgtype.UUID   `json:"attachable_item_id"`
	Url                  pgtype.Text   `json:"url"`
	Size                 pgtype.Float8 `json:"size"`
	MimeTypeID           pgtype.UUID   `json:"mime_type_id"`
	ImageID              pgtype.UUID   `json:"image_id"`
	ImageHeight          pgtype.Float8 `json:"image_height"`
	ImageWidth           pgtype.Float8 `json:"image_width"`
	FileID               pgtype.UUID   `json:"file_id"`
}

func (q *Queries) GetAttachableItemsOnMessage(ctx context.Context, messageID pgtype.UUID) ([]GetAttachableItemsOnMessageRow, error) {
	rows, err := q.db.Query(ctx, getAttachableItemsOnMessage, messageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAttachableItemsOnMessageRow{}
	for rows.Next() {
		var i GetAttachableItemsOnMessageRow
		if err := rows.Scan(
			&i.TAttachableItemsPkey,
			&i.AttachableItemID,
			&i.Url,
			&i.Size,
			&i.MimeTypeID,
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

const getAttachableItemsOnMessageUseKeysetPaginate = `-- name: GetAttachableItemsOnMessageUseKeysetPaginate :many
SELECT t_attachable_items.t_attachable_items_pkey, t_attachable_items.attachable_item_id, t_attachable_items.url, t_attachable_items.size, t_attachable_items.mime_type_id, t_images.image_id, t_images.height as image_height, t_images.width as image_width, t_files.file_id FROM t_attached_messages
LEFT JOIN t_attachable_items ON t_attached_messages.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE message_id = $1
AND
	CASE $3::text
		WHEN 'next' THEN
			t_attached_messages_pkey > $4::int
		WHEN 'prev' THEN
			t_attached_messages_pkey < $4::int
	END
ORDER BY
	CASE WHEN $3::text = 'next' THEN t_attached_messages_pkey END ASC,
	CASE WHEN $3::text = 'prev' THEN t_attached_messages_pkey END DESC
LIMIT $2
`

type GetAttachableItemsOnMessageUseKeysetPaginateParams struct {
	MessageID       pgtype.UUID `json:"message_id"`
	Limit           int32       `json:"limit"`
	CursorDirection string      `json:"cursor_direction"`
	Cursor          int32       `json:"cursor"`
}

type GetAttachableItemsOnMessageUseKeysetPaginateRow struct {
	TAttachableItemsPkey pgtype.Int8   `json:"t_attachable_items_pkey"`
	AttachableItemID     pgtype.UUID   `json:"attachable_item_id"`
	Url                  pgtype.Text   `json:"url"`
	Size                 pgtype.Float8 `json:"size"`
	MimeTypeID           pgtype.UUID   `json:"mime_type_id"`
	ImageID              pgtype.UUID   `json:"image_id"`
	ImageHeight          pgtype.Float8 `json:"image_height"`
	ImageWidth           pgtype.Float8 `json:"image_width"`
	FileID               pgtype.UUID   `json:"file_id"`
}

func (q *Queries) GetAttachableItemsOnMessageUseKeysetPaginate(ctx context.Context, arg GetAttachableItemsOnMessageUseKeysetPaginateParams) ([]GetAttachableItemsOnMessageUseKeysetPaginateRow, error) {
	rows, err := q.db.Query(ctx, getAttachableItemsOnMessageUseKeysetPaginate,
		arg.MessageID,
		arg.Limit,
		arg.CursorDirection,
		arg.Cursor,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAttachableItemsOnMessageUseKeysetPaginateRow{}
	for rows.Next() {
		var i GetAttachableItemsOnMessageUseKeysetPaginateRow
		if err := rows.Scan(
			&i.TAttachableItemsPkey,
			&i.AttachableItemID,
			&i.Url,
			&i.Size,
			&i.MimeTypeID,
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

const getAttachableItemsOnMessageUseNumberedPaginate = `-- name: GetAttachableItemsOnMessageUseNumberedPaginate :many
SELECT t_attachable_items.t_attachable_items_pkey, t_attachable_items.attachable_item_id, t_attachable_items.url, t_attachable_items.size, t_attachable_items.mime_type_id, t_images.image_id, t_images.height as image_height, t_images.width as image_width, t_files.file_id FROM t_attached_messages
LEFT JOIN t_attachable_items ON t_attached_messages.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE message_id = $1
ORDER BY
	t_attached_messages_pkey ASC
LIMIT $2 OFFSET $3
`

type GetAttachableItemsOnMessageUseNumberedPaginateParams struct {
	MessageID pgtype.UUID `json:"message_id"`
	Limit     int32       `json:"limit"`
	Offset    int32       `json:"offset"`
}

type GetAttachableItemsOnMessageUseNumberedPaginateRow struct {
	TAttachableItemsPkey pgtype.Int8   `json:"t_attachable_items_pkey"`
	AttachableItemID     pgtype.UUID   `json:"attachable_item_id"`
	Url                  pgtype.Text   `json:"url"`
	Size                 pgtype.Float8 `json:"size"`
	MimeTypeID           pgtype.UUID   `json:"mime_type_id"`
	ImageID              pgtype.UUID   `json:"image_id"`
	ImageHeight          pgtype.Float8 `json:"image_height"`
	ImageWidth           pgtype.Float8 `json:"image_width"`
	FileID               pgtype.UUID   `json:"file_id"`
}

func (q *Queries) GetAttachableItemsOnMessageUseNumberedPaginate(ctx context.Context, arg GetAttachableItemsOnMessageUseNumberedPaginateParams) ([]GetAttachableItemsOnMessageUseNumberedPaginateRow, error) {
	rows, err := q.db.Query(ctx, getAttachableItemsOnMessageUseNumberedPaginate, arg.MessageID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAttachableItemsOnMessageUseNumberedPaginateRow{}
	for rows.Next() {
		var i GetAttachableItemsOnMessageUseNumberedPaginateRow
		if err := rows.Scan(
			&i.TAttachableItemsPkey,
			&i.AttachableItemID,
			&i.Url,
			&i.Size,
			&i.MimeTypeID,
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

const getAttachedMessagesOnChatRoom = `-- name: GetAttachedMessagesOnChatRoom :many
SELECT t_attachable_items.t_attachable_items_pkey, t_attachable_items.attachable_item_id, t_attachable_items.url, t_attachable_items.size, t_attachable_items.mime_type_id, t_images.image_id, t_images.height as image_height, t_images.width as image_width, t_files.file_id FROM t_attached_messages
LEFT JOIN t_attachable_items ON t_attached_messages.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE message_id IN (
	SELECT message_id FROM t_messages WHERE chat_room_id = $1
)
ORDER BY
	t_attached_messages_pkey ASC
`

type GetAttachedMessagesOnChatRoomRow struct {
	TAttachableItemsPkey pgtype.Int8   `json:"t_attachable_items_pkey"`
	AttachableItemID     pgtype.UUID   `json:"attachable_item_id"`
	Url                  pgtype.Text   `json:"url"`
	Size                 pgtype.Float8 `json:"size"`
	MimeTypeID           pgtype.UUID   `json:"mime_type_id"`
	ImageID              pgtype.UUID   `json:"image_id"`
	ImageHeight          pgtype.Float8 `json:"image_height"`
	ImageWidth           pgtype.Float8 `json:"image_width"`
	FileID               pgtype.UUID   `json:"file_id"`
}

func (q *Queries) GetAttachedMessagesOnChatRoom(ctx context.Context, chatRoomID uuid.UUID) ([]GetAttachedMessagesOnChatRoomRow, error) {
	rows, err := q.db.Query(ctx, getAttachedMessagesOnChatRoom, chatRoomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAttachedMessagesOnChatRoomRow{}
	for rows.Next() {
		var i GetAttachedMessagesOnChatRoomRow
		if err := rows.Scan(
			&i.TAttachableItemsPkey,
			&i.AttachableItemID,
			&i.Url,
			&i.Size,
			&i.MimeTypeID,
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

const getAttachedMessagesOnChatRoomUseKeysetPaginate = `-- name: GetAttachedMessagesOnChatRoomUseKeysetPaginate :many
SELECT t_attachable_items.t_attachable_items_pkey, t_attachable_items.attachable_item_id, t_attachable_items.url, t_attachable_items.size, t_attachable_items.mime_type_id, t_images.image_id, t_images.height as image_height, t_images.width as image_width, t_files.file_id FROM t_attached_messages
LEFT JOIN t_attachable_items ON t_attached_messages.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE message_id IN (
	SELECT message_id FROM t_messages WHERE chat_room_id = $1
)
AND
	CASE $3::text
		WHEN 'next' THEN
			t_attached_messages_pkey > $4::int
		WHEN 'prev' THEN
			t_attached_messages_pkey < $4::int
	END
ORDER BY
	CASE WHEN $3::text = 'next' THEN t_attached_messages_pkey END ASC,
	CASE WHEN $3::text = 'prev' THEN t_attached_messages_pkey END DESC
LIMIT $2
`

type GetAttachedMessagesOnChatRoomUseKeysetPaginateParams struct {
	ChatRoomID      uuid.UUID `json:"chat_room_id"`
	Limit           int32     `json:"limit"`
	CursorDirection string    `json:"cursor_direction"`
	Cursor          int32     `json:"cursor"`
}

type GetAttachedMessagesOnChatRoomUseKeysetPaginateRow struct {
	TAttachableItemsPkey pgtype.Int8   `json:"t_attachable_items_pkey"`
	AttachableItemID     pgtype.UUID   `json:"attachable_item_id"`
	Url                  pgtype.Text   `json:"url"`
	Size                 pgtype.Float8 `json:"size"`
	MimeTypeID           pgtype.UUID   `json:"mime_type_id"`
	ImageID              pgtype.UUID   `json:"image_id"`
	ImageHeight          pgtype.Float8 `json:"image_height"`
	ImageWidth           pgtype.Float8 `json:"image_width"`
	FileID               pgtype.UUID   `json:"file_id"`
}

func (q *Queries) GetAttachedMessagesOnChatRoomUseKeysetPaginate(ctx context.Context, arg GetAttachedMessagesOnChatRoomUseKeysetPaginateParams) ([]GetAttachedMessagesOnChatRoomUseKeysetPaginateRow, error) {
	rows, err := q.db.Query(ctx, getAttachedMessagesOnChatRoomUseKeysetPaginate,
		arg.ChatRoomID,
		arg.Limit,
		arg.CursorDirection,
		arg.Cursor,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAttachedMessagesOnChatRoomUseKeysetPaginateRow{}
	for rows.Next() {
		var i GetAttachedMessagesOnChatRoomUseKeysetPaginateRow
		if err := rows.Scan(
			&i.TAttachableItemsPkey,
			&i.AttachableItemID,
			&i.Url,
			&i.Size,
			&i.MimeTypeID,
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

const getAttachedMessagesOnChatRoomUseNumberedPaginate = `-- name: GetAttachedMessagesOnChatRoomUseNumberedPaginate :many
SELECT t_attachable_items.t_attachable_items_pkey, t_attachable_items.attachable_item_id, t_attachable_items.url, t_attachable_items.size, t_attachable_items.mime_type_id, t_images.image_id, t_images.height as image_height, t_images.width as image_width, t_files.file_id FROM t_attached_messages
LEFT JOIN t_attachable_items ON t_attached_messages.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE message_id IN (
	SELECT message_id FROM t_messages WHERE chat_room_id = $1
)
ORDER BY
	t_attached_messages_pkey ASC
LIMIT $2 OFFSET $3
`

type GetAttachedMessagesOnChatRoomUseNumberedPaginateParams struct {
	ChatRoomID uuid.UUID `json:"chat_room_id"`
	Limit      int32     `json:"limit"`
	Offset     int32     `json:"offset"`
}

type GetAttachedMessagesOnChatRoomUseNumberedPaginateRow struct {
	TAttachableItemsPkey pgtype.Int8   `json:"t_attachable_items_pkey"`
	AttachableItemID     pgtype.UUID   `json:"attachable_item_id"`
	Url                  pgtype.Text   `json:"url"`
	Size                 pgtype.Float8 `json:"size"`
	MimeTypeID           pgtype.UUID   `json:"mime_type_id"`
	ImageID              pgtype.UUID   `json:"image_id"`
	ImageHeight          pgtype.Float8 `json:"image_height"`
	ImageWidth           pgtype.Float8 `json:"image_width"`
	FileID               pgtype.UUID   `json:"file_id"`
}

func (q *Queries) GetAttachedMessagesOnChatRoomUseNumberedPaginate(ctx context.Context, arg GetAttachedMessagesOnChatRoomUseNumberedPaginateParams) ([]GetAttachedMessagesOnChatRoomUseNumberedPaginateRow, error) {
	rows, err := q.db.Query(ctx, getAttachedMessagesOnChatRoomUseNumberedPaginate, arg.ChatRoomID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAttachedMessagesOnChatRoomUseNumberedPaginateRow{}
	for rows.Next() {
		var i GetAttachedMessagesOnChatRoomUseNumberedPaginateRow
		if err := rows.Scan(
			&i.TAttachableItemsPkey,
			&i.AttachableItemID,
			&i.Url,
			&i.Size,
			&i.MimeTypeID,
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

const getPluralAttachableItemsOnMessage = `-- name: GetPluralAttachableItemsOnMessage :many
SELECT t_attachable_items.t_attachable_items_pkey, t_attachable_items.attachable_item_id, t_attachable_items.url, t_attachable_items.size, t_attachable_items.mime_type_id, t_images.image_id, t_images.height as image_height, t_images.width as image_width, t_files.file_id FROM t_attached_messages
LEFT JOIN t_attachable_items ON t_attached_messages.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE message_id = ANY($3::uuid[])
ORDER BY
	t_attached_messages_pkey ASC
LIMIT $1 OFFSET $2
`

type GetPluralAttachableItemsOnMessageParams struct {
	Limit      int32       `json:"limit"`
	Offset     int32       `json:"offset"`
	MessageIds []uuid.UUID `json:"message_ids"`
}

type GetPluralAttachableItemsOnMessageRow struct {
	TAttachableItemsPkey pgtype.Int8   `json:"t_attachable_items_pkey"`
	AttachableItemID     pgtype.UUID   `json:"attachable_item_id"`
	Url                  pgtype.Text   `json:"url"`
	Size                 pgtype.Float8 `json:"size"`
	MimeTypeID           pgtype.UUID   `json:"mime_type_id"`
	ImageID              pgtype.UUID   `json:"image_id"`
	ImageHeight          pgtype.Float8 `json:"image_height"`
	ImageWidth           pgtype.Float8 `json:"image_width"`
	FileID               pgtype.UUID   `json:"file_id"`
}

func (q *Queries) GetPluralAttachableItemsOnMessage(ctx context.Context, arg GetPluralAttachableItemsOnMessageParams) ([]GetPluralAttachableItemsOnMessageRow, error) {
	rows, err := q.db.Query(ctx, getPluralAttachableItemsOnMessage, arg.Limit, arg.Offset, arg.MessageIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPluralAttachableItemsOnMessageRow{}
	for rows.Next() {
		var i GetPluralAttachableItemsOnMessageRow
		if err := rows.Scan(
			&i.TAttachableItemsPkey,
			&i.AttachableItemID,
			&i.Url,
			&i.Size,
			&i.MimeTypeID,
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

const getPluralAttachedMessagesOnChatRoom = `-- name: GetPluralAttachedMessagesOnChatRoom :many
SELECT t_attachable_items.t_attachable_items_pkey, t_attachable_items.attachable_item_id, t_attachable_items.url, t_attachable_items.size, t_attachable_items.mime_type_id, t_images.image_id, t_images.height as image_height, t_images.width as image_width, t_files.file_id FROM t_attached_messages
LEFT JOIN t_attachable_items ON t_attached_messages.attachable_item_id = t_attachable_items.attachable_item_id
LEFT JOIN m_mime_types ON t_attachable_items.mime_type_id = m_mime_types.mime_type_id
LEFT JOIN t_images ON t_attachable_items.attachable_item_id = t_images.attachable_item_id
LEFT JOIN t_files ON t_attachable_items.attachable_item_id = t_files.attachable_item_id
WHERE message_id IN (
	SELECT message_id FROM t_messages WHERE chat_room_id = ANY($3::uuid[])
)
ORDER BY
	t_attached_messages_pkey ASC
LIMIT $1 OFFSET $2
`

type GetPluralAttachedMessagesOnChatRoomParams struct {
	Limit       int32       `json:"limit"`
	Offset      int32       `json:"offset"`
	ChatRoomIds []uuid.UUID `json:"chat_room_ids"`
}

type GetPluralAttachedMessagesOnChatRoomRow struct {
	TAttachableItemsPkey pgtype.Int8   `json:"t_attachable_items_pkey"`
	AttachableItemID     pgtype.UUID   `json:"attachable_item_id"`
	Url                  pgtype.Text   `json:"url"`
	Size                 pgtype.Float8 `json:"size"`
	MimeTypeID           pgtype.UUID   `json:"mime_type_id"`
	ImageID              pgtype.UUID   `json:"image_id"`
	ImageHeight          pgtype.Float8 `json:"image_height"`
	ImageWidth           pgtype.Float8 `json:"image_width"`
	FileID               pgtype.UUID   `json:"file_id"`
}

func (q *Queries) GetPluralAttachedMessagesOnChatRoom(ctx context.Context, arg GetPluralAttachedMessagesOnChatRoomParams) ([]GetPluralAttachedMessagesOnChatRoomRow, error) {
	rows, err := q.db.Query(ctx, getPluralAttachedMessagesOnChatRoom, arg.Limit, arg.Offset, arg.ChatRoomIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPluralAttachedMessagesOnChatRoomRow{}
	for rows.Next() {
		var i GetPluralAttachedMessagesOnChatRoomRow
		if err := rows.Scan(
			&i.TAttachableItemsPkey,
			&i.AttachableItemID,
			&i.Url,
			&i.Size,
			&i.MimeTypeID,
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

const pluralDeleteAttachedMessagesOnAttachableItem = `-- name: PluralDeleteAttachedMessagesOnAttachableItem :exec
DELETE FROM t_attached_messages WHERE attachable_item_id = $1 AND message_id = ANY($2::uuid[])
`

type PluralDeleteAttachedMessagesOnAttachableItemParams struct {
	AttachableItemID uuid.UUID   `json:"attachable_item_id"`
	Column2          []uuid.UUID `json:"column_2"`
}

func (q *Queries) PluralDeleteAttachedMessagesOnAttachableItem(ctx context.Context, arg PluralDeleteAttachedMessagesOnAttachableItemParams) error {
	_, err := q.db.Exec(ctx, pluralDeleteAttachedMessagesOnAttachableItem, arg.AttachableItemID, arg.Column2)
	return err
}

const pluralDeleteAttachedMessagesOnMessage = `-- name: PluralDeleteAttachedMessagesOnMessage :exec
DELETE FROM t_attached_messages WHERE message_id = $1 AND attachable_item_id = ANY($2::uuid[])
`

type PluralDeleteAttachedMessagesOnMessageParams struct {
	MessageID pgtype.UUID `json:"message_id"`
	Column2   []uuid.UUID `json:"column_2"`
}

func (q *Queries) PluralDeleteAttachedMessagesOnMessage(ctx context.Context, arg PluralDeleteAttachedMessagesOnMessageParams) error {
	_, err := q.db.Exec(ctx, pluralDeleteAttachedMessagesOnMessage, arg.MessageID, arg.Column2)
	return err
}
