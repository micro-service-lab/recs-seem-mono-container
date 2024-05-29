// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: chat_room_update_name_action.sql

package query

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const countChatRoomUpdateNameActions = `-- name: CountChatRoomUpdateNameActions :one
SELECT COUNT(*) FROM t_chat_room_update_name_actions
`

func (q *Queries) CountChatRoomUpdateNameActions(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, countChatRoomUpdateNameActions)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createChatRoomUpdateNameAction = `-- name: CreateChatRoomUpdateNameAction :one
INSERT INTO t_chat_room_update_name_actions (chat_room_action_id, updated_by, name) VALUES ($1, $2, $3) RETURNING t_chat_room_update_name_actions_pkey, chat_room_update_name_action_id, chat_room_action_id, updated_by, name
`

type CreateChatRoomUpdateNameActionParams struct {
	ChatRoomActionID uuid.UUID   `json:"chat_room_action_id"`
	UpdatedBy        pgtype.UUID `json:"updated_by"`
	Name             string      `json:"name"`
}

func (q *Queries) CreateChatRoomUpdateNameAction(ctx context.Context, arg CreateChatRoomUpdateNameActionParams) (ChatRoomUpdateNameAction, error) {
	row := q.db.QueryRow(ctx, createChatRoomUpdateNameAction, arg.ChatRoomActionID, arg.UpdatedBy, arg.Name)
	var i ChatRoomUpdateNameAction
	err := row.Scan(
		&i.TChatRoomUpdateNameActionsPkey,
		&i.ChatRoomUpdateNameActionID,
		&i.ChatRoomActionID,
		&i.UpdatedBy,
		&i.Name,
	)
	return i, err
}

type CreateChatRoomUpdateNameActionsParams struct {
	ChatRoomActionID uuid.UUID   `json:"chat_room_action_id"`
	UpdatedBy        pgtype.UUID `json:"updated_by"`
	Name             string      `json:"name"`
}

const deleteChatRoomUpdateNameAction = `-- name: DeleteChatRoomUpdateNameAction :execrows
DELETE FROM t_chat_room_update_name_actions WHERE chat_room_update_name_action_id = $1
`

func (q *Queries) DeleteChatRoomUpdateNameAction(ctx context.Context, chatRoomUpdateNameActionID uuid.UUID) (int64, error) {
	result, err := q.db.Exec(ctx, deleteChatRoomUpdateNameAction, chatRoomUpdateNameActionID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const getChatRoomUpdateNameActionsOnChatRoom = `-- name: GetChatRoomUpdateNameActionsOnChatRoom :many
SELECT t_chat_room_update_name_actions.t_chat_room_update_name_actions_pkey, t_chat_room_update_name_actions.chat_room_update_name_action_id, t_chat_room_update_name_actions.chat_room_action_id, t_chat_room_update_name_actions.updated_by, t_chat_room_update_name_actions.name,
m_members.name update_member_name, m_members.first_name update_member_first_name, m_members.last_name update_member_last_name, m_members.email update_member_email,
m_members.profile_image_id update_member_profile_image_id
FROM t_chat_room_update_name_actions
LEFT JOIN m_members ON t_chat_room_update_name_actions.updated_by = m_members.member_id
WHERE EXISTS (
	SELECT 1 FROM t_chat_room_actions WHERE chat_room_action_id = t_chat_room_update_name_actions.chat_room_action_id AND t_chat_room_actions.chat_room_id = $1
)
ORDER BY
	t_chat_room_update_name_actions_pkey ASC
`

type GetChatRoomUpdateNameActionsOnChatRoomRow struct {
	TChatRoomUpdateNameActionsPkey pgtype.Int8 `json:"t_chat_room_update_name_actions_pkey"`
	ChatRoomUpdateNameActionID     uuid.UUID   `json:"chat_room_update_name_action_id"`
	ChatRoomActionID               uuid.UUID   `json:"chat_room_action_id"`
	UpdatedBy                      pgtype.UUID `json:"updated_by"`
	Name                           string      `json:"name"`
	UpdateMemberName               pgtype.Text `json:"update_member_name"`
	UpdateMemberFirstName          pgtype.Text `json:"update_member_first_name"`
	UpdateMemberLastName           pgtype.Text `json:"update_member_last_name"`
	UpdateMemberEmail              pgtype.Text `json:"update_member_email"`
	UpdateMemberProfileImageID     pgtype.UUID `json:"update_member_profile_image_id"`
}

func (q *Queries) GetChatRoomUpdateNameActionsOnChatRoom(ctx context.Context, chatRoomID uuid.UUID) ([]GetChatRoomUpdateNameActionsOnChatRoomRow, error) {
	rows, err := q.db.Query(ctx, getChatRoomUpdateNameActionsOnChatRoom, chatRoomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetChatRoomUpdateNameActionsOnChatRoomRow{}
	for rows.Next() {
		var i GetChatRoomUpdateNameActionsOnChatRoomRow
		if err := rows.Scan(
			&i.TChatRoomUpdateNameActionsPkey,
			&i.ChatRoomUpdateNameActionID,
			&i.ChatRoomActionID,
			&i.UpdatedBy,
			&i.Name,
			&i.UpdateMemberName,
			&i.UpdateMemberFirstName,
			&i.UpdateMemberLastName,
			&i.UpdateMemberEmail,
			&i.UpdateMemberProfileImageID,
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

const getChatRoomUpdateNameActionsOnChatRoomUseKeysetPaginate = `-- name: GetChatRoomUpdateNameActionsOnChatRoomUseKeysetPaginate :many
SELECT t_chat_room_update_name_actions.t_chat_room_update_name_actions_pkey, t_chat_room_update_name_actions.chat_room_update_name_action_id, t_chat_room_update_name_actions.chat_room_action_id, t_chat_room_update_name_actions.updated_by, t_chat_room_update_name_actions.name,
m_members.name update_member_name, m_members.first_name update_member_first_name, m_members.last_name update_member_last_name, m_members.email update_member_email,
m_members.profile_image_id update_member_profile_image_id
FROM t_chat_room_update_name_actions
LEFT JOIN m_members ON t_chat_room_update_name_actions.updated_by = m_members.member_id
WHERE EXISTS (
	SELECT 1 FROM t_chat_room_actions WHERE chat_room_action_id = t_chat_room_update_name_actions.chat_room_action_id AND chat_room_id = $1
)
AND
	CASE $3::text
		WHEN 'next' THEN
			t_chat_room_update_name_actions_pkey > $4::int
		WHEN 'prev' THEN
			t_chat_room_update_name_actions_pkey < $4::int
	END
ORDER BY
	CASE WHEN $3::text = 'next' THEN t_chat_room_update_name_actions_pkey END ASC,
	CASE WHEN $3::text = 'prev' THEN t_chat_room_update_name_actions_pkey END DESC
LIMIT $2
`

type GetChatRoomUpdateNameActionsOnChatRoomUseKeysetPaginateParams struct {
	ChatRoomID      uuid.UUID `json:"chat_room_id"`
	Limit           int32     `json:"limit"`
	CursorDirection string    `json:"cursor_direction"`
	Cursor          int32     `json:"cursor"`
}

type GetChatRoomUpdateNameActionsOnChatRoomUseKeysetPaginateRow struct {
	TChatRoomUpdateNameActionsPkey pgtype.Int8 `json:"t_chat_room_update_name_actions_pkey"`
	ChatRoomUpdateNameActionID     uuid.UUID   `json:"chat_room_update_name_action_id"`
	ChatRoomActionID               uuid.UUID   `json:"chat_room_action_id"`
	UpdatedBy                      pgtype.UUID `json:"updated_by"`
	Name                           string      `json:"name"`
	UpdateMemberName               pgtype.Text `json:"update_member_name"`
	UpdateMemberFirstName          pgtype.Text `json:"update_member_first_name"`
	UpdateMemberLastName           pgtype.Text `json:"update_member_last_name"`
	UpdateMemberEmail              pgtype.Text `json:"update_member_email"`
	UpdateMemberProfileImageID     pgtype.UUID `json:"update_member_profile_image_id"`
}

func (q *Queries) GetChatRoomUpdateNameActionsOnChatRoomUseKeysetPaginate(ctx context.Context, arg GetChatRoomUpdateNameActionsOnChatRoomUseKeysetPaginateParams) ([]GetChatRoomUpdateNameActionsOnChatRoomUseKeysetPaginateRow, error) {
	rows, err := q.db.Query(ctx, getChatRoomUpdateNameActionsOnChatRoomUseKeysetPaginate,
		arg.ChatRoomID,
		arg.Limit,
		arg.CursorDirection,
		arg.Cursor,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetChatRoomUpdateNameActionsOnChatRoomUseKeysetPaginateRow{}
	for rows.Next() {
		var i GetChatRoomUpdateNameActionsOnChatRoomUseKeysetPaginateRow
		if err := rows.Scan(
			&i.TChatRoomUpdateNameActionsPkey,
			&i.ChatRoomUpdateNameActionID,
			&i.ChatRoomActionID,
			&i.UpdatedBy,
			&i.Name,
			&i.UpdateMemberName,
			&i.UpdateMemberFirstName,
			&i.UpdateMemberLastName,
			&i.UpdateMemberEmail,
			&i.UpdateMemberProfileImageID,
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

const getChatRoomUpdateNameActionsOnChatRoomUseNumberedPaginate = `-- name: GetChatRoomUpdateNameActionsOnChatRoomUseNumberedPaginate :many
SELECT t_chat_room_update_name_actions.t_chat_room_update_name_actions_pkey, t_chat_room_update_name_actions.chat_room_update_name_action_id, t_chat_room_update_name_actions.chat_room_action_id, t_chat_room_update_name_actions.updated_by, t_chat_room_update_name_actions.name,
m_members.name update_member_name, m_members.first_name update_member_first_name, m_members.last_name update_member_last_name, m_members.email update_member_email,
m_members.profile_image_id update_member_profile_image_id
FROM t_chat_room_update_name_actions
LEFT JOIN m_members ON t_chat_room_update_name_actions.updated_by = m_members.member_id
WHERE EXISTS (
	SELECT 1 FROM t_chat_room_actions WHERE chat_room_action_id = t_chat_room_update_name_actions.chat_room_action_id AND chat_room_id = $1
)
ORDER BY
	t_chat_room_update_name_actions_pkey ASC
LIMIT $2 OFFSET $3
`

type GetChatRoomUpdateNameActionsOnChatRoomUseNumberedPaginateParams struct {
	ChatRoomID uuid.UUID `json:"chat_room_id"`
	Limit      int32     `json:"limit"`
	Offset     int32     `json:"offset"`
}

type GetChatRoomUpdateNameActionsOnChatRoomUseNumberedPaginateRow struct {
	TChatRoomUpdateNameActionsPkey pgtype.Int8 `json:"t_chat_room_update_name_actions_pkey"`
	ChatRoomUpdateNameActionID     uuid.UUID   `json:"chat_room_update_name_action_id"`
	ChatRoomActionID               uuid.UUID   `json:"chat_room_action_id"`
	UpdatedBy                      pgtype.UUID `json:"updated_by"`
	Name                           string      `json:"name"`
	UpdateMemberName               pgtype.Text `json:"update_member_name"`
	UpdateMemberFirstName          pgtype.Text `json:"update_member_first_name"`
	UpdateMemberLastName           pgtype.Text `json:"update_member_last_name"`
	UpdateMemberEmail              pgtype.Text `json:"update_member_email"`
	UpdateMemberProfileImageID     pgtype.UUID `json:"update_member_profile_image_id"`
}

func (q *Queries) GetChatRoomUpdateNameActionsOnChatRoomUseNumberedPaginate(ctx context.Context, arg GetChatRoomUpdateNameActionsOnChatRoomUseNumberedPaginateParams) ([]GetChatRoomUpdateNameActionsOnChatRoomUseNumberedPaginateRow, error) {
	rows, err := q.db.Query(ctx, getChatRoomUpdateNameActionsOnChatRoomUseNumberedPaginate, arg.ChatRoomID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetChatRoomUpdateNameActionsOnChatRoomUseNumberedPaginateRow{}
	for rows.Next() {
		var i GetChatRoomUpdateNameActionsOnChatRoomUseNumberedPaginateRow
		if err := rows.Scan(
			&i.TChatRoomUpdateNameActionsPkey,
			&i.ChatRoomUpdateNameActionID,
			&i.ChatRoomActionID,
			&i.UpdatedBy,
			&i.Name,
			&i.UpdateMemberName,
			&i.UpdateMemberFirstName,
			&i.UpdateMemberLastName,
			&i.UpdateMemberEmail,
			&i.UpdateMemberProfileImageID,
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

const getPluralChatRoomUpdateNameActions = `-- name: GetPluralChatRoomUpdateNameActions :many
SELECT t_chat_room_update_name_actions.t_chat_room_update_name_actions_pkey, t_chat_room_update_name_actions.chat_room_update_name_action_id, t_chat_room_update_name_actions.chat_room_action_id, t_chat_room_update_name_actions.updated_by, t_chat_room_update_name_actions.name,
m_members.name update_member_name, m_members.first_name update_member_first_name, m_members.last_name update_member_last_name, m_members.email update_member_email,
m_members.profile_image_id update_member_profile_image_id
FROM t_chat_room_update_name_actions
LEFT JOIN m_members ON t_chat_room_update_name_actions.updated_by = m_members.member_id
WHERE chat_room_update_name_action_id = ANY($1::uuid[])
ORDER BY
	t_chat_room_update_name_actions_pkey ASC
`

type GetPluralChatRoomUpdateNameActionsRow struct {
	TChatRoomUpdateNameActionsPkey pgtype.Int8 `json:"t_chat_room_update_name_actions_pkey"`
	ChatRoomUpdateNameActionID     uuid.UUID   `json:"chat_room_update_name_action_id"`
	ChatRoomActionID               uuid.UUID   `json:"chat_room_action_id"`
	UpdatedBy                      pgtype.UUID `json:"updated_by"`
	Name                           string      `json:"name"`
	UpdateMemberName               pgtype.Text `json:"update_member_name"`
	UpdateMemberFirstName          pgtype.Text `json:"update_member_first_name"`
	UpdateMemberLastName           pgtype.Text `json:"update_member_last_name"`
	UpdateMemberEmail              pgtype.Text `json:"update_member_email"`
	UpdateMemberProfileImageID     pgtype.UUID `json:"update_member_profile_image_id"`
}

func (q *Queries) GetPluralChatRoomUpdateNameActions(ctx context.Context, chatRoomUpdateNameActionIds []uuid.UUID) ([]GetPluralChatRoomUpdateNameActionsRow, error) {
	rows, err := q.db.Query(ctx, getPluralChatRoomUpdateNameActions, chatRoomUpdateNameActionIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPluralChatRoomUpdateNameActionsRow{}
	for rows.Next() {
		var i GetPluralChatRoomUpdateNameActionsRow
		if err := rows.Scan(
			&i.TChatRoomUpdateNameActionsPkey,
			&i.ChatRoomUpdateNameActionID,
			&i.ChatRoomActionID,
			&i.UpdatedBy,
			&i.Name,
			&i.UpdateMemberName,
			&i.UpdateMemberFirstName,
			&i.UpdateMemberLastName,
			&i.UpdateMemberEmail,
			&i.UpdateMemberProfileImageID,
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

const getPluralChatRoomUpdateNameActionsUseNumberedPaginate = `-- name: GetPluralChatRoomUpdateNameActionsUseNumberedPaginate :many
SELECT t_chat_room_update_name_actions.t_chat_room_update_name_actions_pkey, t_chat_room_update_name_actions.chat_room_update_name_action_id, t_chat_room_update_name_actions.chat_room_action_id, t_chat_room_update_name_actions.updated_by, t_chat_room_update_name_actions.name,
m_members.name update_member_name, m_members.first_name update_member_first_name, m_members.last_name update_member_last_name, m_members.email update_member_email,
m_members.profile_image_id update_member_profile_image_id
FROM t_chat_room_update_name_actions
LEFT JOIN m_members ON t_chat_room_update_name_actions.updated_by = m_members.member_id
WHERE chat_room_update_name_action_id = ANY($3::uuid[])
ORDER BY
	t_chat_room_update_name_actions_pkey ASC
LIMIT $1 OFFSET $2
`

type GetPluralChatRoomUpdateNameActionsUseNumberedPaginateParams struct {
	Limit                       int32       `json:"limit"`
	Offset                      int32       `json:"offset"`
	ChatRoomUpdateNameActionIds []uuid.UUID `json:"chat_room_update_name_action_ids"`
}

type GetPluralChatRoomUpdateNameActionsUseNumberedPaginateRow struct {
	TChatRoomUpdateNameActionsPkey pgtype.Int8 `json:"t_chat_room_update_name_actions_pkey"`
	ChatRoomUpdateNameActionID     uuid.UUID   `json:"chat_room_update_name_action_id"`
	ChatRoomActionID               uuid.UUID   `json:"chat_room_action_id"`
	UpdatedBy                      pgtype.UUID `json:"updated_by"`
	Name                           string      `json:"name"`
	UpdateMemberName               pgtype.Text `json:"update_member_name"`
	UpdateMemberFirstName          pgtype.Text `json:"update_member_first_name"`
	UpdateMemberLastName           pgtype.Text `json:"update_member_last_name"`
	UpdateMemberEmail              pgtype.Text `json:"update_member_email"`
	UpdateMemberProfileImageID     pgtype.UUID `json:"update_member_profile_image_id"`
}

func (q *Queries) GetPluralChatRoomUpdateNameActionsUseNumberedPaginate(ctx context.Context, arg GetPluralChatRoomUpdateNameActionsUseNumberedPaginateParams) ([]GetPluralChatRoomUpdateNameActionsUseNumberedPaginateRow, error) {
	rows, err := q.db.Query(ctx, getPluralChatRoomUpdateNameActionsUseNumberedPaginate, arg.Limit, arg.Offset, arg.ChatRoomUpdateNameActionIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPluralChatRoomUpdateNameActionsUseNumberedPaginateRow{}
	for rows.Next() {
		var i GetPluralChatRoomUpdateNameActionsUseNumberedPaginateRow
		if err := rows.Scan(
			&i.TChatRoomUpdateNameActionsPkey,
			&i.ChatRoomUpdateNameActionID,
			&i.ChatRoomActionID,
			&i.UpdatedBy,
			&i.Name,
			&i.UpdateMemberName,
			&i.UpdateMemberFirstName,
			&i.UpdateMemberLastName,
			&i.UpdateMemberEmail,
			&i.UpdateMemberProfileImageID,
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

const pluralDeleteChatRoomUpdateNameActions = `-- name: PluralDeleteChatRoomUpdateNameActions :execrows
DELETE FROM t_chat_room_update_name_actions WHERE chat_room_update_name_action_id = ANY($1::uuid[])
`

func (q *Queries) PluralDeleteChatRoomUpdateNameActions(ctx context.Context, chatRoomUpdateNameActionIds []uuid.UUID) (int64, error) {
	result, err := q.db.Exec(ctx, pluralDeleteChatRoomUpdateNameActions, chatRoomUpdateNameActionIds)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}