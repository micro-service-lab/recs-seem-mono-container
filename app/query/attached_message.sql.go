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

const countFilesOnChatRoom = `-- name: CountFilesOnChatRoom :one
SELECT COUNT(*) FROM t_attached_messages
WHERE message_id IN (
	SELECT message_id FROM t_messages WHERE chat_room_id = $1
)
`

func (q *Queries) CountFilesOnChatRoom(ctx context.Context, chatRoomID uuid.UUID) (int64, error) {
	row := q.db.QueryRow(ctx, countFilesOnChatRoom, chatRoomID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countFilesOnMessage = `-- name: CountFilesOnMessage :one
SELECT COUNT(*) FROM t_attached_messages WHERE message_id = $1
`

func (q *Queries) CountFilesOnMessage(ctx context.Context, messageID uuid.UUID) (int64, error) {
	row := q.db.QueryRow(ctx, countFilesOnMessage, messageID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createAttachedMessage = `-- name: CreateAttachedMessage :one
INSERT INTO t_attached_messages (message_id, attachable_item_id) VALUES ($1, $2) RETURNING t_attached_messages_pkey, attached_message_id, message_id, attachable_item_id
`

type CreateAttachedMessageParams struct {
	MessageID        uuid.UUID   `json:"message_id"`
	AttachableItemID pgtype.UUID `json:"attachable_item_id"`
}

func (q *Queries) CreateAttachedMessage(ctx context.Context, arg CreateAttachedMessageParams) (AttachedMessage, error) {
	row := q.db.QueryRow(ctx, createAttachedMessage, arg.MessageID, arg.AttachableItemID)
	var i AttachedMessage
	err := row.Scan(
		&i.TAttachedMessagesPkey,
		&i.AttachedMessageID,
		&i.MessageID,
		&i.AttachableItemID,
	)
	return i, err
}

type CreateAttachedMessagesParams struct {
	MessageID        uuid.UUID   `json:"message_id"`
	AttachableItemID pgtype.UUID `json:"attachable_item_id"`
}

const deleteAttachedMessage = `-- name: DeleteAttachedMessage :execrows
DELETE FROM t_attached_messages WHERE attached_message_id = $1
`

func (q *Queries) DeleteAttachedMessage(ctx context.Context, attachedMessageID uuid.UUID) (int64, error) {
	result, err := q.db.Exec(ctx, deleteAttachedMessage, attachedMessageID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const deleteAttachedMessagesOnMessage = `-- name: DeleteAttachedMessagesOnMessage :execrows
DELETE FROM t_attached_messages WHERE message_id = $1
`

func (q *Queries) DeleteAttachedMessagesOnMessage(ctx context.Context, messageID uuid.UUID) (int64, error) {
	result, err := q.db.Exec(ctx, deleteAttachedMessagesOnMessage, messageID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const deleteAttachedMessagesOnMessages = `-- name: DeleteAttachedMessagesOnMessages :execrows
DELETE FROM t_attached_messages WHERE message_id = ANY($1::uuid[])
`

func (q *Queries) DeleteAttachedMessagesOnMessages(ctx context.Context, dollar_1 []uuid.UUID) (int64, error) {
	result, err := q.db.Exec(ctx, deleteAttachedMessagesOnMessages, dollar_1)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const getFilesOnChatRoom = `-- name: GetFilesOnChatRoom :many
SELECT  FROM t_attached_messages
WHERE message_id IN (
	SELECT message_id FROM t_messages WHERE chat_room_id = $1
)
ORDER BY
	t_attached_messages_pkey ASC
`

type GetFilesOnChatRoomRow struct {
}

func (q *Queries) GetFilesOnChatRoom(ctx context.Context, chatRoomID uuid.UUID) ([]GetFilesOnChatRoomRow, error) {
	rows, err := q.db.Query(ctx, getFilesOnChatRoom, chatRoomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetFilesOnChatRoomRow{}
	for rows.Next() {
		var i GetFilesOnChatRoomRow
		if err := rows.Scan(); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFilesOnChatRoomUseKeysetPaginate = `-- name: GetFilesOnChatRoomUseKeysetPaginate :many
SELECT  FROM t_attached_messages
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

type GetFilesOnChatRoomUseKeysetPaginateParams struct {
	ChatRoomID      uuid.UUID `json:"chat_room_id"`
	Limit           int32     `json:"limit"`
	CursorDirection string    `json:"cursor_direction"`
	Cursor          int32     `json:"cursor"`
}

type GetFilesOnChatRoomUseKeysetPaginateRow struct {
}

func (q *Queries) GetFilesOnChatRoomUseKeysetPaginate(ctx context.Context, arg GetFilesOnChatRoomUseKeysetPaginateParams) ([]GetFilesOnChatRoomUseKeysetPaginateRow, error) {
	rows, err := q.db.Query(ctx, getFilesOnChatRoomUseKeysetPaginate,
		arg.ChatRoomID,
		arg.Limit,
		arg.CursorDirection,
		arg.Cursor,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetFilesOnChatRoomUseKeysetPaginateRow{}
	for rows.Next() {
		var i GetFilesOnChatRoomUseKeysetPaginateRow
		if err := rows.Scan(); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFilesOnChatRoomUseNumberedPaginate = `-- name: GetFilesOnChatRoomUseNumberedPaginate :many
SELECT  FROM t_attached_messages
WHERE message_id IN (
	SELECT message_id FROM t_messages WHERE chat_room_id = $1
)
ORDER BY
	t_attached_messages_pkey ASC
LIMIT $2 OFFSET $3
`

type GetFilesOnChatRoomUseNumberedPaginateParams struct {
	ChatRoomID uuid.UUID `json:"chat_room_id"`
	Limit      int32     `json:"limit"`
	Offset     int32     `json:"offset"`
}

type GetFilesOnChatRoomUseNumberedPaginateRow struct {
}

func (q *Queries) GetFilesOnChatRoomUseNumberedPaginate(ctx context.Context, arg GetFilesOnChatRoomUseNumberedPaginateParams) ([]GetFilesOnChatRoomUseNumberedPaginateRow, error) {
	rows, err := q.db.Query(ctx, getFilesOnChatRoomUseNumberedPaginate, arg.ChatRoomID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetFilesOnChatRoomUseNumberedPaginateRow{}
	for rows.Next() {
		var i GetFilesOnChatRoomUseNumberedPaginateRow
		if err := rows.Scan(); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFilesOnMessage = `-- name: GetFilesOnMessage :many
SELECT  FROM t_attached_messages
WHERE message_id = $1
ORDER BY
	t_attached_messages_pkey ASC
`

type GetFilesOnMessageRow struct {
}

func (q *Queries) GetFilesOnMessage(ctx context.Context, messageID uuid.UUID) ([]GetFilesOnMessageRow, error) {
	rows, err := q.db.Query(ctx, getFilesOnMessage, messageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetFilesOnMessageRow{}
	for rows.Next() {
		var i GetFilesOnMessageRow
		if err := rows.Scan(); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFilesOnMessageUseKeysetPaginate = `-- name: GetFilesOnMessageUseKeysetPaginate :many
SELECT  FROM t_attached_messages
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

type GetFilesOnMessageUseKeysetPaginateParams struct {
	MessageID       uuid.UUID `json:"message_id"`
	Limit           int32     `json:"limit"`
	CursorDirection string    `json:"cursor_direction"`
	Cursor          int32     `json:"cursor"`
}

type GetFilesOnMessageUseKeysetPaginateRow struct {
}

func (q *Queries) GetFilesOnMessageUseKeysetPaginate(ctx context.Context, arg GetFilesOnMessageUseKeysetPaginateParams) ([]GetFilesOnMessageUseKeysetPaginateRow, error) {
	rows, err := q.db.Query(ctx, getFilesOnMessageUseKeysetPaginate,
		arg.MessageID,
		arg.Limit,
		arg.CursorDirection,
		arg.Cursor,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetFilesOnMessageUseKeysetPaginateRow{}
	for rows.Next() {
		var i GetFilesOnMessageUseKeysetPaginateRow
		if err := rows.Scan(); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFilesOnMessageUseNumberedPaginate = `-- name: GetFilesOnMessageUseNumberedPaginate :many
SELECT  FROM t_attached_messages
WHERE message_id = $1
ORDER BY
	t_attached_messages_pkey ASC
LIMIT $2 OFFSET $3
`

type GetFilesOnMessageUseNumberedPaginateParams struct {
	MessageID uuid.UUID `json:"message_id"`
	Limit     int32     `json:"limit"`
	Offset    int32     `json:"offset"`
}

type GetFilesOnMessageUseNumberedPaginateRow struct {
}

func (q *Queries) GetFilesOnMessageUseNumberedPaginate(ctx context.Context, arg GetFilesOnMessageUseNumberedPaginateParams) ([]GetFilesOnMessageUseNumberedPaginateRow, error) {
	rows, err := q.db.Query(ctx, getFilesOnMessageUseNumberedPaginate, arg.MessageID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetFilesOnMessageUseNumberedPaginateRow{}
	for rows.Next() {
		var i GetFilesOnMessageUseNumberedPaginateRow
		if err := rows.Scan(); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPluralFilesOnChatRoom = `-- name: GetPluralFilesOnChatRoom :many
SELECT  FROM t_attached_messages
WHERE message_id IN (
	SELECT message_id FROM t_messages WHERE chat_room_id = ANY($3::uuid[])
)
ORDER BY
	t_attached_messages_pkey ASC
LIMIT $1 OFFSET $2
`

type GetPluralFilesOnChatRoomParams struct {
	Limit       int32       `json:"limit"`
	Offset      int32       `json:"offset"`
	ChatRoomIds []uuid.UUID `json:"chat_room_ids"`
}

type GetPluralFilesOnChatRoomRow struct {
}

func (q *Queries) GetPluralFilesOnChatRoom(ctx context.Context, arg GetPluralFilesOnChatRoomParams) ([]GetPluralFilesOnChatRoomRow, error) {
	rows, err := q.db.Query(ctx, getPluralFilesOnChatRoom, arg.Limit, arg.Offset, arg.ChatRoomIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPluralFilesOnChatRoomRow{}
	for rows.Next() {
		var i GetPluralFilesOnChatRoomRow
		if err := rows.Scan(); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPluralFilesOnMessage = `-- name: GetPluralFilesOnMessage :many
SELECT  FROM t_attached_messages
WHERE message_id = ANY($3::uuid[])
ORDER BY
	t_attached_messages_pkey ASC
LIMIT $1 OFFSET $2
`

type GetPluralFilesOnMessageParams struct {
	Limit      int32       `json:"limit"`
	Offset     int32       `json:"offset"`
	MessageIds []uuid.UUID `json:"message_ids"`
}

type GetPluralFilesOnMessageRow struct {
}

func (q *Queries) GetPluralFilesOnMessage(ctx context.Context, arg GetPluralFilesOnMessageParams) ([]GetPluralFilesOnMessageRow, error) {
	rows, err := q.db.Query(ctx, getPluralFilesOnMessage, arg.Limit, arg.Offset, arg.MessageIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPluralFilesOnMessageRow{}
	for rows.Next() {
		var i GetPluralFilesOnMessageRow
		if err := rows.Scan(); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
