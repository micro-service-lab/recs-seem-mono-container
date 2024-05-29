// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: chat_room_create_action.sql

package query

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const countChatRoomCreateActions = `-- name: CountChatRoomCreateActions :one
SELECT COUNT(*) FROM t_chat_room_create_actions
`

func (q *Queries) CountChatRoomCreateActions(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, countChatRoomCreateActions)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createChatRoomCreateAction = `-- name: CreateChatRoomCreateAction :one
INSERT INTO t_chat_room_create_actions (chat_room_action_id, created_by, name) VALUES ($1, $2, $3) RETURNING t_chat_room_create_actions_pkey, chat_room_create_action_id, chat_room_action_id, created_by, name
`

type CreateChatRoomCreateActionParams struct {
	ChatRoomActionID uuid.UUID   `json:"chat_room_action_id"`
	CreatedBy        pgtype.UUID `json:"created_by"`
	Name             string      `json:"name"`
}

func (q *Queries) CreateChatRoomCreateAction(ctx context.Context, arg CreateChatRoomCreateActionParams) (ChatRoomCreateAction, error) {
	row := q.db.QueryRow(ctx, createChatRoomCreateAction, arg.ChatRoomActionID, arg.CreatedBy, arg.Name)
	var i ChatRoomCreateAction
	err := row.Scan(
		&i.TChatRoomCreateActionsPkey,
		&i.ChatRoomCreateActionID,
		&i.ChatRoomActionID,
		&i.CreatedBy,
		&i.Name,
	)
	return i, err
}

type CreateChatRoomCreateActionsParams struct {
	ChatRoomActionID uuid.UUID   `json:"chat_room_action_id"`
	CreatedBy        pgtype.UUID `json:"created_by"`
	Name             string      `json:"name"`
}

const deleteChatRoomCreateAction = `-- name: DeleteChatRoomCreateAction :execrows
DELETE FROM t_chat_room_create_actions WHERE chat_room_create_action_id = $1
`

func (q *Queries) DeleteChatRoomCreateAction(ctx context.Context, chatRoomCreateActionID uuid.UUID) (int64, error) {
	result, err := q.db.Exec(ctx, deleteChatRoomCreateAction, chatRoomCreateActionID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const getChatRoomCreateActionsOnChatRoom = `-- name: GetChatRoomCreateActionsOnChatRoom :many
SELECT t_chat_room_create_actions.t_chat_room_create_actions_pkey, t_chat_room_create_actions.chat_room_create_action_id, t_chat_room_create_actions.chat_room_action_id, t_chat_room_create_actions.created_by, t_chat_room_create_actions.name,
m_members.name create_member_name, m_members.first_name create_member_first_name, m_members.last_name create_member_last_name, m_members.email create_member_email,
m_members.profile_image_id create_member_profile_image_id
FROM t_chat_room_create_actions
LEFT JOIN m_members ON t_chat_room_create_actions.created_by = m_members.member_id
WHERE EXISTS (
	SELECT 1 FROM t_chat_room_actions WHERE chat_room_action_id = t_chat_room_create_actions.chat_room_action_id AND t_chat_room_actions.chat_room_id = $1
)
ORDER BY
	t_chat_room_create_actions_pkey ASC
`

type GetChatRoomCreateActionsOnChatRoomRow struct {
	TChatRoomCreateActionsPkey pgtype.Int8 `json:"t_chat_room_create_actions_pkey"`
	ChatRoomCreateActionID     uuid.UUID   `json:"chat_room_create_action_id"`
	ChatRoomActionID           uuid.UUID   `json:"chat_room_action_id"`
	CreatedBy                  pgtype.UUID `json:"created_by"`
	Name                       string      `json:"name"`
	CreateMemberName           pgtype.Text `json:"create_member_name"`
	CreateMemberFirstName      pgtype.Text `json:"create_member_first_name"`
	CreateMemberLastName       pgtype.Text `json:"create_member_last_name"`
	CreateMemberEmail          pgtype.Text `json:"create_member_email"`
	CreateMemberProfileImageID pgtype.UUID `json:"create_member_profile_image_id"`
}

func (q *Queries) GetChatRoomCreateActionsOnChatRoom(ctx context.Context, chatRoomID uuid.UUID) ([]GetChatRoomCreateActionsOnChatRoomRow, error) {
	rows, err := q.db.Query(ctx, getChatRoomCreateActionsOnChatRoom, chatRoomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetChatRoomCreateActionsOnChatRoomRow{}
	for rows.Next() {
		var i GetChatRoomCreateActionsOnChatRoomRow
		if err := rows.Scan(
			&i.TChatRoomCreateActionsPkey,
			&i.ChatRoomCreateActionID,
			&i.ChatRoomActionID,
			&i.CreatedBy,
			&i.Name,
			&i.CreateMemberName,
			&i.CreateMemberFirstName,
			&i.CreateMemberLastName,
			&i.CreateMemberEmail,
			&i.CreateMemberProfileImageID,
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

const getChatRoomCreateActionsOnChatRoomUseKeysetPaginate = `-- name: GetChatRoomCreateActionsOnChatRoomUseKeysetPaginate :many
SELECT t_chat_room_create_actions.t_chat_room_create_actions_pkey, t_chat_room_create_actions.chat_room_create_action_id, t_chat_room_create_actions.chat_room_action_id, t_chat_room_create_actions.created_by, t_chat_room_create_actions.name,
m_members.name create_member_name, m_members.first_name create_member_first_name, m_members.last_name create_member_last_name, m_members.email create_member_email,
m_members.profile_image_id create_member_profile_image_id
FROM t_chat_room_create_actions
LEFT JOIN m_members ON t_chat_room_create_actions.created_by = m_members.member_id
WHERE EXISTS (
	SELECT 1 FROM t_chat_room_actions WHERE chat_room_action_id = t_chat_room_create_actions.chat_room_action_id AND chat_room_id = $1
)
AND
	CASE $3::text
		WHEN 'next' THEN
			t_chat_room_create_actions_pkey > $4::int
		WHEN 'prev' THEN
			t_chat_room_create_actions_pkey < $4::int
	END
ORDER BY
	CASE WHEN $3::text = 'next' THEN t_chat_room_create_actions_pkey END ASC,
	CASE WHEN $3::text = 'prev' THEN t_chat_room_create_actions_pkey END DESC
LIMIT $2
`

type GetChatRoomCreateActionsOnChatRoomUseKeysetPaginateParams struct {
	ChatRoomID      uuid.UUID `json:"chat_room_id"`
	Limit           int32     `json:"limit"`
	CursorDirection string    `json:"cursor_direction"`
	Cursor          int32     `json:"cursor"`
}

type GetChatRoomCreateActionsOnChatRoomUseKeysetPaginateRow struct {
	TChatRoomCreateActionsPkey pgtype.Int8 `json:"t_chat_room_create_actions_pkey"`
	ChatRoomCreateActionID     uuid.UUID   `json:"chat_room_create_action_id"`
	ChatRoomActionID           uuid.UUID   `json:"chat_room_action_id"`
	CreatedBy                  pgtype.UUID `json:"created_by"`
	Name                       string      `json:"name"`
	CreateMemberName           pgtype.Text `json:"create_member_name"`
	CreateMemberFirstName      pgtype.Text `json:"create_member_first_name"`
	CreateMemberLastName       pgtype.Text `json:"create_member_last_name"`
	CreateMemberEmail          pgtype.Text `json:"create_member_email"`
	CreateMemberProfileImageID pgtype.UUID `json:"create_member_profile_image_id"`
}

func (q *Queries) GetChatRoomCreateActionsOnChatRoomUseKeysetPaginate(ctx context.Context, arg GetChatRoomCreateActionsOnChatRoomUseKeysetPaginateParams) ([]GetChatRoomCreateActionsOnChatRoomUseKeysetPaginateRow, error) {
	rows, err := q.db.Query(ctx, getChatRoomCreateActionsOnChatRoomUseKeysetPaginate,
		arg.ChatRoomID,
		arg.Limit,
		arg.CursorDirection,
		arg.Cursor,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetChatRoomCreateActionsOnChatRoomUseKeysetPaginateRow{}
	for rows.Next() {
		var i GetChatRoomCreateActionsOnChatRoomUseKeysetPaginateRow
		if err := rows.Scan(
			&i.TChatRoomCreateActionsPkey,
			&i.ChatRoomCreateActionID,
			&i.ChatRoomActionID,
			&i.CreatedBy,
			&i.Name,
			&i.CreateMemberName,
			&i.CreateMemberFirstName,
			&i.CreateMemberLastName,
			&i.CreateMemberEmail,
			&i.CreateMemberProfileImageID,
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

const getChatRoomCreateActionsOnChatRoomUseNumberedPaginate = `-- name: GetChatRoomCreateActionsOnChatRoomUseNumberedPaginate :many
SELECT t_chat_room_create_actions.t_chat_room_create_actions_pkey, t_chat_room_create_actions.chat_room_create_action_id, t_chat_room_create_actions.chat_room_action_id, t_chat_room_create_actions.created_by, t_chat_room_create_actions.name,
m_members.name create_member_name, m_members.first_name create_member_first_name, m_members.last_name create_member_last_name, m_members.email create_member_email,
m_members.profile_image_id create_member_profile_image_id
FROM t_chat_room_create_actions
LEFT JOIN m_members ON t_chat_room_create_actions.created_by = m_members.member_id
WHERE EXISTS (
	SELECT 1 FROM t_chat_room_actions WHERE chat_room_action_id = t_chat_room_create_actions.chat_room_action_id AND chat_room_id = $1
)
ORDER BY
	t_chat_room_create_actions_pkey ASC
LIMIT $2 OFFSET $3
`

type GetChatRoomCreateActionsOnChatRoomUseNumberedPaginateParams struct {
	ChatRoomID uuid.UUID `json:"chat_room_id"`
	Limit      int32     `json:"limit"`
	Offset     int32     `json:"offset"`
}

type GetChatRoomCreateActionsOnChatRoomUseNumberedPaginateRow struct {
	TChatRoomCreateActionsPkey pgtype.Int8 `json:"t_chat_room_create_actions_pkey"`
	ChatRoomCreateActionID     uuid.UUID   `json:"chat_room_create_action_id"`
	ChatRoomActionID           uuid.UUID   `json:"chat_room_action_id"`
	CreatedBy                  pgtype.UUID `json:"created_by"`
	Name                       string      `json:"name"`
	CreateMemberName           pgtype.Text `json:"create_member_name"`
	CreateMemberFirstName      pgtype.Text `json:"create_member_first_name"`
	CreateMemberLastName       pgtype.Text `json:"create_member_last_name"`
	CreateMemberEmail          pgtype.Text `json:"create_member_email"`
	CreateMemberProfileImageID pgtype.UUID `json:"create_member_profile_image_id"`
}

func (q *Queries) GetChatRoomCreateActionsOnChatRoomUseNumberedPaginate(ctx context.Context, arg GetChatRoomCreateActionsOnChatRoomUseNumberedPaginateParams) ([]GetChatRoomCreateActionsOnChatRoomUseNumberedPaginateRow, error) {
	rows, err := q.db.Query(ctx, getChatRoomCreateActionsOnChatRoomUseNumberedPaginate, arg.ChatRoomID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetChatRoomCreateActionsOnChatRoomUseNumberedPaginateRow{}
	for rows.Next() {
		var i GetChatRoomCreateActionsOnChatRoomUseNumberedPaginateRow
		if err := rows.Scan(
			&i.TChatRoomCreateActionsPkey,
			&i.ChatRoomCreateActionID,
			&i.ChatRoomActionID,
			&i.CreatedBy,
			&i.Name,
			&i.CreateMemberName,
			&i.CreateMemberFirstName,
			&i.CreateMemberLastName,
			&i.CreateMemberEmail,
			&i.CreateMemberProfileImageID,
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

const getPluralChatRoomCreateActions = `-- name: GetPluralChatRoomCreateActions :many
SELECT t_chat_room_create_actions.t_chat_room_create_actions_pkey, t_chat_room_create_actions.chat_room_create_action_id, t_chat_room_create_actions.chat_room_action_id, t_chat_room_create_actions.created_by, t_chat_room_create_actions.name,
m_members.name create_member_name, m_members.first_name create_member_first_name, m_members.last_name create_member_last_name, m_members.email create_member_email,
m_members.profile_image_id create_member_profile_image_id
FROM t_chat_room_create_actions
LEFT JOIN m_members ON t_chat_room_create_actions.created_by = m_members.member_id
WHERE chat_room_create_action_id = ANY($1::uuid[])
ORDER BY
	t_chat_room_create_actions_pkey ASC
`

type GetPluralChatRoomCreateActionsRow struct {
	TChatRoomCreateActionsPkey pgtype.Int8 `json:"t_chat_room_create_actions_pkey"`
	ChatRoomCreateActionID     uuid.UUID   `json:"chat_room_create_action_id"`
	ChatRoomActionID           uuid.UUID   `json:"chat_room_action_id"`
	CreatedBy                  pgtype.UUID `json:"created_by"`
	Name                       string      `json:"name"`
	CreateMemberName           pgtype.Text `json:"create_member_name"`
	CreateMemberFirstName      pgtype.Text `json:"create_member_first_name"`
	CreateMemberLastName       pgtype.Text `json:"create_member_last_name"`
	CreateMemberEmail          pgtype.Text `json:"create_member_email"`
	CreateMemberProfileImageID pgtype.UUID `json:"create_member_profile_image_id"`
}

func (q *Queries) GetPluralChatRoomCreateActions(ctx context.Context, chatRoomCreateActionIds []uuid.UUID) ([]GetPluralChatRoomCreateActionsRow, error) {
	rows, err := q.db.Query(ctx, getPluralChatRoomCreateActions, chatRoomCreateActionIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPluralChatRoomCreateActionsRow{}
	for rows.Next() {
		var i GetPluralChatRoomCreateActionsRow
		if err := rows.Scan(
			&i.TChatRoomCreateActionsPkey,
			&i.ChatRoomCreateActionID,
			&i.ChatRoomActionID,
			&i.CreatedBy,
			&i.Name,
			&i.CreateMemberName,
			&i.CreateMemberFirstName,
			&i.CreateMemberLastName,
			&i.CreateMemberEmail,
			&i.CreateMemberProfileImageID,
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

const getPluralChatRoomCreateActionsUseNumberedPaginate = `-- name: GetPluralChatRoomCreateActionsUseNumberedPaginate :many
SELECT t_chat_room_create_actions.t_chat_room_create_actions_pkey, t_chat_room_create_actions.chat_room_create_action_id, t_chat_room_create_actions.chat_room_action_id, t_chat_room_create_actions.created_by, t_chat_room_create_actions.name,
m_members.name create_member_name, m_members.first_name create_member_first_name, m_members.last_name create_member_last_name, m_members.email create_member_email,
m_members.profile_image_id create_member_profile_image_id
FROM t_chat_room_create_actions
LEFT JOIN m_members ON t_chat_room_create_actions.created_by = m_members.member_id
WHERE chat_room_create_action_id = ANY($3::uuid[])
ORDER BY
	t_chat_room_create_actions_pkey ASC
LIMIT $1 OFFSET $2
`

type GetPluralChatRoomCreateActionsUseNumberedPaginateParams struct {
	Limit                   int32       `json:"limit"`
	Offset                  int32       `json:"offset"`
	ChatRoomCreateActionIds []uuid.UUID `json:"chat_room_create_action_ids"`
}

type GetPluralChatRoomCreateActionsUseNumberedPaginateRow struct {
	TChatRoomCreateActionsPkey pgtype.Int8 `json:"t_chat_room_create_actions_pkey"`
	ChatRoomCreateActionID     uuid.UUID   `json:"chat_room_create_action_id"`
	ChatRoomActionID           uuid.UUID   `json:"chat_room_action_id"`
	CreatedBy                  pgtype.UUID `json:"created_by"`
	Name                       string      `json:"name"`
	CreateMemberName           pgtype.Text `json:"create_member_name"`
	CreateMemberFirstName      pgtype.Text `json:"create_member_first_name"`
	CreateMemberLastName       pgtype.Text `json:"create_member_last_name"`
	CreateMemberEmail          pgtype.Text `json:"create_member_email"`
	CreateMemberProfileImageID pgtype.UUID `json:"create_member_profile_image_id"`
}

func (q *Queries) GetPluralChatRoomCreateActionsUseNumberedPaginate(ctx context.Context, arg GetPluralChatRoomCreateActionsUseNumberedPaginateParams) ([]GetPluralChatRoomCreateActionsUseNumberedPaginateRow, error) {
	rows, err := q.db.Query(ctx, getPluralChatRoomCreateActionsUseNumberedPaginate, arg.Limit, arg.Offset, arg.ChatRoomCreateActionIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPluralChatRoomCreateActionsUseNumberedPaginateRow{}
	for rows.Next() {
		var i GetPluralChatRoomCreateActionsUseNumberedPaginateRow
		if err := rows.Scan(
			&i.TChatRoomCreateActionsPkey,
			&i.ChatRoomCreateActionID,
			&i.ChatRoomActionID,
			&i.CreatedBy,
			&i.Name,
			&i.CreateMemberName,
			&i.CreateMemberFirstName,
			&i.CreateMemberLastName,
			&i.CreateMemberEmail,
			&i.CreateMemberProfileImageID,
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

const pluralDeleteChatRoomCreateActions = `-- name: PluralDeleteChatRoomCreateActions :execrows
DELETE FROM t_chat_room_create_actions WHERE chat_room_create_action_id = ANY($1::uuid[])
`

func (q *Queries) PluralDeleteChatRoomCreateActions(ctx context.Context, chatRoomCreateActionIds []uuid.UUID) (int64, error) {
	result, err := q.db.Exec(ctx, pluralDeleteChatRoomCreateActions, chatRoomCreateActionIds)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}
