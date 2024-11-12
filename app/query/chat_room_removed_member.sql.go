// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: chat_room_removed_member.sql

package query

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const countMembersOnChatRoomRemoveMemberAction = `-- name: CountMembersOnChatRoomRemoveMemberAction :one
SELECT COUNT(*) FROM t_chat_room_removed_members WHERE chat_room_remove_member_action_id = $1
`

func (q *Queries) CountMembersOnChatRoomRemoveMemberAction(ctx context.Context, chatRoomRemoveMemberActionID uuid.UUID) (int64, error) {
	row := q.db.QueryRow(ctx, countMembersOnChatRoomRemoveMemberAction, chatRoomRemoveMemberActionID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createChatRoomRemovedMember = `-- name: CreateChatRoomRemovedMember :one
INSERT INTO t_chat_room_removed_members (member_id, chat_room_remove_member_action_id) VALUES ($1, $2) RETURNING t_chat_room_removed_members_pkey, chat_room_remove_member_action_id, member_id
`

type CreateChatRoomRemovedMemberParams struct {
	MemberID                     pgtype.UUID `json:"member_id"`
	ChatRoomRemoveMemberActionID uuid.UUID   `json:"chat_room_remove_member_action_id"`
}

func (q *Queries) CreateChatRoomRemovedMember(ctx context.Context, arg CreateChatRoomRemovedMemberParams) (ChatRoomRemovedMember, error) {
	row := q.db.QueryRow(ctx, createChatRoomRemovedMember, arg.MemberID, arg.ChatRoomRemoveMemberActionID)
	var i ChatRoomRemovedMember
	err := row.Scan(&i.TChatRoomRemovedMembersPkey, &i.ChatRoomRemoveMemberActionID, &i.MemberID)
	return i, err
}

type CreateChatRoomRemovedMembersParams struct {
	MemberID                     pgtype.UUID `json:"member_id"`
	ChatRoomRemoveMemberActionID uuid.UUID   `json:"chat_room_remove_member_action_id"`
}

const deleteChatRoomRemovedMember = `-- name: DeleteChatRoomRemovedMember :execrows
DELETE FROM t_chat_room_removed_members WHERE member_id = $1 AND chat_room_remove_member_action_id = $2
`

type DeleteChatRoomRemovedMemberParams struct {
	MemberID                     pgtype.UUID `json:"member_id"`
	ChatRoomRemoveMemberActionID uuid.UUID   `json:"chat_room_remove_member_action_id"`
}

func (q *Queries) DeleteChatRoomRemovedMember(ctx context.Context, arg DeleteChatRoomRemovedMemberParams) (int64, error) {
	result, err := q.db.Exec(ctx, deleteChatRoomRemovedMember, arg.MemberID, arg.ChatRoomRemoveMemberActionID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const deleteChatRoomRemovedMembersOnChatRoomRemoveMemberAction = `-- name: DeleteChatRoomRemovedMembersOnChatRoomRemoveMemberAction :execrows
DELETE FROM t_chat_room_removed_members WHERE chat_room_remove_member_action_id = $1
`

func (q *Queries) DeleteChatRoomRemovedMembersOnChatRoomRemoveMemberAction(ctx context.Context, chatRoomRemoveMemberActionID uuid.UUID) (int64, error) {
	result, err := q.db.Exec(ctx, deleteChatRoomRemovedMembersOnChatRoomRemoveMemberAction, chatRoomRemoveMemberActionID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const deleteChatRoomRemovedMembersOnChatRoomRemoveMemberActions = `-- name: DeleteChatRoomRemovedMembersOnChatRoomRemoveMemberActions :execrows
DELETE FROM t_chat_room_removed_members WHERE chat_room_remove_member_action_id = ANY($1::uuid[])
`

func (q *Queries) DeleteChatRoomRemovedMembersOnChatRoomRemoveMemberActions(ctx context.Context, chatRoomRemoveMemberActionIds []uuid.UUID) (int64, error) {
	result, err := q.db.Exec(ctx, deleteChatRoomRemovedMembersOnChatRoomRemoveMemberActions, chatRoomRemoveMemberActionIds)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const deleteChatRoomRemovedMembersOnMember = `-- name: DeleteChatRoomRemovedMembersOnMember :execrows
DELETE FROM t_chat_room_removed_members WHERE member_id = $1
`

func (q *Queries) DeleteChatRoomRemovedMembersOnMember(ctx context.Context, memberID pgtype.UUID) (int64, error) {
	result, err := q.db.Exec(ctx, deleteChatRoomRemovedMembersOnMember, memberID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const deleteChatRoomRemovedMembersOnMembers = `-- name: DeleteChatRoomRemovedMembersOnMembers :execrows
DELETE FROM t_chat_room_removed_members WHERE member_id = ANY($1::uuid[])
`

func (q *Queries) DeleteChatRoomRemovedMembersOnMembers(ctx context.Context, memberIds []uuid.UUID) (int64, error) {
	result, err := q.db.Exec(ctx, deleteChatRoomRemovedMembersOnMembers, memberIds)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const getMembersOnChatRoomRemoveMemberAction = `-- name: GetMembersOnChatRoomRemoveMemberAction :many
SELECT t_chat_room_removed_members.t_chat_room_removed_members_pkey, t_chat_room_removed_members.chat_room_remove_member_action_id, t_chat_room_removed_members.member_id, m_members.name member_name, m_members.first_name member_first_name, m_members.last_name member_last_name, m_members.email member_email,
m_members.profile_image_id member_profile_image_id, m_members.grade_id member_grade_id, m_members.group_id member_group_id
FROM t_chat_room_removed_members
LEFT JOIN m_members ON t_chat_room_removed_members.member_id = m_members.member_id
WHERE chat_room_remove_member_action_id = $1
ORDER BY
	t_chat_room_removed_members_pkey ASC
`

type GetMembersOnChatRoomRemoveMemberActionRow struct {
	TChatRoomRemovedMembersPkey  pgtype.Int8 `json:"t_chat_room_removed_members_pkey"`
	ChatRoomRemoveMemberActionID uuid.UUID   `json:"chat_room_remove_member_action_id"`
	MemberID                     pgtype.UUID `json:"member_id"`
	MemberName                   pgtype.Text `json:"member_name"`
	MemberFirstName              pgtype.Text `json:"member_first_name"`
	MemberLastName               pgtype.Text `json:"member_last_name"`
	MemberEmail                  pgtype.Text `json:"member_email"`
	MemberProfileImageID         pgtype.UUID `json:"member_profile_image_id"`
	MemberGradeID                pgtype.UUID `json:"member_grade_id"`
	MemberGroupID                pgtype.UUID `json:"member_group_id"`
}

func (q *Queries) GetMembersOnChatRoomRemoveMemberAction(ctx context.Context, chatRoomRemoveMemberActionID uuid.UUID) ([]GetMembersOnChatRoomRemoveMemberActionRow, error) {
	rows, err := q.db.Query(ctx, getMembersOnChatRoomRemoveMemberAction, chatRoomRemoveMemberActionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetMembersOnChatRoomRemoveMemberActionRow{}
	for rows.Next() {
		var i GetMembersOnChatRoomRemoveMemberActionRow
		if err := rows.Scan(
			&i.TChatRoomRemovedMembersPkey,
			&i.ChatRoomRemoveMemberActionID,
			&i.MemberID,
			&i.MemberName,
			&i.MemberFirstName,
			&i.MemberLastName,
			&i.MemberEmail,
			&i.MemberProfileImageID,
			&i.MemberGradeID,
			&i.MemberGroupID,
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

const getMembersOnChatRoomRemoveMemberActionUseKeysetPaginate = `-- name: GetMembersOnChatRoomRemoveMemberActionUseKeysetPaginate :many
SELECT t_chat_room_removed_members.t_chat_room_removed_members_pkey, t_chat_room_removed_members.chat_room_remove_member_action_id, t_chat_room_removed_members.member_id, m_members.name member_name, m_members.first_name member_first_name, m_members.last_name member_last_name, m_members.email member_email,
m_members.profile_image_id member_profile_image_id, m_members.grade_id member_grade_id, m_members.group_id member_group_id
FROM t_chat_room_removed_members
LEFT JOIN m_members ON t_chat_room_removed_members.member_id = m_members.member_id
WHERE chat_room_remove_member_action_id = $1
AND CASE $3::text
	WHEN 'next' THEN
			t_chat_room_removed_members_pkey > $4::int
	WHEN 'prev' THEN
			t_chat_room_removed_members_pkey < $4::int
END
ORDER BY
	CASE WHEN $3::text = 'next' THEN t_chat_room_removed_members_pkey END ASC,
	CASE WHEN $3::text = 'prev' THEN t_chat_room_removed_members_pkey END DESC
LIMIT $2
`

type GetMembersOnChatRoomRemoveMemberActionUseKeysetPaginateParams struct {
	ChatRoomRemoveMemberActionID uuid.UUID `json:"chat_room_remove_member_action_id"`
	Limit                        int32     `json:"limit"`
	CursorDirection              string    `json:"cursor_direction"`
	Cursor                       int32     `json:"cursor"`
}

type GetMembersOnChatRoomRemoveMemberActionUseKeysetPaginateRow struct {
	TChatRoomRemovedMembersPkey  pgtype.Int8 `json:"t_chat_room_removed_members_pkey"`
	ChatRoomRemoveMemberActionID uuid.UUID   `json:"chat_room_remove_member_action_id"`
	MemberID                     pgtype.UUID `json:"member_id"`
	MemberName                   pgtype.Text `json:"member_name"`
	MemberFirstName              pgtype.Text `json:"member_first_name"`
	MemberLastName               pgtype.Text `json:"member_last_name"`
	MemberEmail                  pgtype.Text `json:"member_email"`
	MemberProfileImageID         pgtype.UUID `json:"member_profile_image_id"`
	MemberGradeID                pgtype.UUID `json:"member_grade_id"`
	MemberGroupID                pgtype.UUID `json:"member_group_id"`
}

func (q *Queries) GetMembersOnChatRoomRemoveMemberActionUseKeysetPaginate(ctx context.Context, arg GetMembersOnChatRoomRemoveMemberActionUseKeysetPaginateParams) ([]GetMembersOnChatRoomRemoveMemberActionUseKeysetPaginateRow, error) {
	rows, err := q.db.Query(ctx, getMembersOnChatRoomRemoveMemberActionUseKeysetPaginate,
		arg.ChatRoomRemoveMemberActionID,
		arg.Limit,
		arg.CursorDirection,
		arg.Cursor,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetMembersOnChatRoomRemoveMemberActionUseKeysetPaginateRow{}
	for rows.Next() {
		var i GetMembersOnChatRoomRemoveMemberActionUseKeysetPaginateRow
		if err := rows.Scan(
			&i.TChatRoomRemovedMembersPkey,
			&i.ChatRoomRemoveMemberActionID,
			&i.MemberID,
			&i.MemberName,
			&i.MemberFirstName,
			&i.MemberLastName,
			&i.MemberEmail,
			&i.MemberProfileImageID,
			&i.MemberGradeID,
			&i.MemberGroupID,
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

const getMembersOnChatRoomRemoveMemberActionUseNumberedPaginate = `-- name: GetMembersOnChatRoomRemoveMemberActionUseNumberedPaginate :many
SELECT t_chat_room_removed_members.t_chat_room_removed_members_pkey, t_chat_room_removed_members.chat_room_remove_member_action_id, t_chat_room_removed_members.member_id, m_members.name member_name, m_members.first_name member_first_name, m_members.last_name member_last_name, m_members.email member_email,
m_members.profile_image_id member_profile_image_id, m_members.grade_id member_grade_id, m_members.group_id member_group_id
FROM t_chat_room_removed_members
LEFT JOIN m_members ON t_chat_room_removed_members.member_id = m_members.member_id
WHERE chat_room_remove_member_action_id = $1
ORDER BY
	t_chat_room_removed_members_pkey ASC
LIMIT $2 OFFSET $3
`

type GetMembersOnChatRoomRemoveMemberActionUseNumberedPaginateParams struct {
	ChatRoomRemoveMemberActionID uuid.UUID `json:"chat_room_remove_member_action_id"`
	Limit                        int32     `json:"limit"`
	Offset                       int32     `json:"offset"`
}

type GetMembersOnChatRoomRemoveMemberActionUseNumberedPaginateRow struct {
	TChatRoomRemovedMembersPkey  pgtype.Int8 `json:"t_chat_room_removed_members_pkey"`
	ChatRoomRemoveMemberActionID uuid.UUID   `json:"chat_room_remove_member_action_id"`
	MemberID                     pgtype.UUID `json:"member_id"`
	MemberName                   pgtype.Text `json:"member_name"`
	MemberFirstName              pgtype.Text `json:"member_first_name"`
	MemberLastName               pgtype.Text `json:"member_last_name"`
	MemberEmail                  pgtype.Text `json:"member_email"`
	MemberProfileImageID         pgtype.UUID `json:"member_profile_image_id"`
	MemberGradeID                pgtype.UUID `json:"member_grade_id"`
	MemberGroupID                pgtype.UUID `json:"member_group_id"`
}

func (q *Queries) GetMembersOnChatRoomRemoveMemberActionUseNumberedPaginate(ctx context.Context, arg GetMembersOnChatRoomRemoveMemberActionUseNumberedPaginateParams) ([]GetMembersOnChatRoomRemoveMemberActionUseNumberedPaginateRow, error) {
	rows, err := q.db.Query(ctx, getMembersOnChatRoomRemoveMemberActionUseNumberedPaginate, arg.ChatRoomRemoveMemberActionID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetMembersOnChatRoomRemoveMemberActionUseNumberedPaginateRow{}
	for rows.Next() {
		var i GetMembersOnChatRoomRemoveMemberActionUseNumberedPaginateRow
		if err := rows.Scan(
			&i.TChatRoomRemovedMembersPkey,
			&i.ChatRoomRemoveMemberActionID,
			&i.MemberID,
			&i.MemberName,
			&i.MemberFirstName,
			&i.MemberLastName,
			&i.MemberEmail,
			&i.MemberProfileImageID,
			&i.MemberGradeID,
			&i.MemberGroupID,
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

const getPluralMembersOnChatRoomRemoveMemberAction = `-- name: GetPluralMembersOnChatRoomRemoveMemberAction :many
SELECT t_chat_room_removed_members.t_chat_room_removed_members_pkey, t_chat_room_removed_members.chat_room_remove_member_action_id, t_chat_room_removed_members.member_id, m_members.name member_name, m_members.first_name member_first_name, m_members.last_name member_last_name, m_members.email member_email,
m_members.profile_image_id member_profile_image_id, m_members.grade_id member_grade_id, m_members.group_id member_group_id
FROM t_chat_room_removed_members
LEFT JOIN m_members ON t_chat_room_removed_members.member_id = m_members.member_id
WHERE chat_room_remove_member_action_id = ANY($1::uuid[])
ORDER BY
	t_chat_room_removed_members_pkey ASC
`

type GetPluralMembersOnChatRoomRemoveMemberActionRow struct {
	TChatRoomRemovedMembersPkey  pgtype.Int8 `json:"t_chat_room_removed_members_pkey"`
	ChatRoomRemoveMemberActionID uuid.UUID   `json:"chat_room_remove_member_action_id"`
	MemberID                     pgtype.UUID `json:"member_id"`
	MemberName                   pgtype.Text `json:"member_name"`
	MemberFirstName              pgtype.Text `json:"member_first_name"`
	MemberLastName               pgtype.Text `json:"member_last_name"`
	MemberEmail                  pgtype.Text `json:"member_email"`
	MemberProfileImageID         pgtype.UUID `json:"member_profile_image_id"`
	MemberGradeID                pgtype.UUID `json:"member_grade_id"`
	MemberGroupID                pgtype.UUID `json:"member_group_id"`
}

func (q *Queries) GetPluralMembersOnChatRoomRemoveMemberAction(ctx context.Context, chatRoomRemoveMemberActionIds []uuid.UUID) ([]GetPluralMembersOnChatRoomRemoveMemberActionRow, error) {
	rows, err := q.db.Query(ctx, getPluralMembersOnChatRoomRemoveMemberAction, chatRoomRemoveMemberActionIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPluralMembersOnChatRoomRemoveMemberActionRow{}
	for rows.Next() {
		var i GetPluralMembersOnChatRoomRemoveMemberActionRow
		if err := rows.Scan(
			&i.TChatRoomRemovedMembersPkey,
			&i.ChatRoomRemoveMemberActionID,
			&i.MemberID,
			&i.MemberName,
			&i.MemberFirstName,
			&i.MemberLastName,
			&i.MemberEmail,
			&i.MemberProfileImageID,
			&i.MemberGradeID,
			&i.MemberGroupID,
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

const getPluralMembersOnChatRoomRemoveMemberActionUseNumberedPaginate = `-- name: GetPluralMembersOnChatRoomRemoveMemberActionUseNumberedPaginate :many
SELECT t_chat_room_removed_members.t_chat_room_removed_members_pkey, t_chat_room_removed_members.chat_room_remove_member_action_id, t_chat_room_removed_members.member_id, m_members.name member_name, m_members.first_name member_first_name, m_members.last_name member_last_name, m_members.email member_email,
m_members.profile_image_id member_profile_image_id, m_members.grade_id member_grade_id, m_members.group_id member_group_id
FROM t_chat_room_removed_members
LEFT JOIN m_members ON t_chat_room_removed_members.member_id = m_members.member_id
WHERE chat_room_remove_member_action_id = ANY($3::uuid[])
ORDER BY
	t_chat_room_removed_members_pkey ASC
LIMIT $1 OFFSET $2
`

type GetPluralMembersOnChatRoomRemoveMemberActionUseNumberedPaginateParams struct {
	Limit                         int32       `json:"limit"`
	Offset                        int32       `json:"offset"`
	ChatRoomRemoveMemberActionIds []uuid.UUID `json:"chat_room_remove_member_action_ids"`
}

type GetPluralMembersOnChatRoomRemoveMemberActionUseNumberedPaginateRow struct {
	TChatRoomRemovedMembersPkey  pgtype.Int8 `json:"t_chat_room_removed_members_pkey"`
	ChatRoomRemoveMemberActionID uuid.UUID   `json:"chat_room_remove_member_action_id"`
	MemberID                     pgtype.UUID `json:"member_id"`
	MemberName                   pgtype.Text `json:"member_name"`
	MemberFirstName              pgtype.Text `json:"member_first_name"`
	MemberLastName               pgtype.Text `json:"member_last_name"`
	MemberEmail                  pgtype.Text `json:"member_email"`
	MemberProfileImageID         pgtype.UUID `json:"member_profile_image_id"`
	MemberGradeID                pgtype.UUID `json:"member_grade_id"`
	MemberGroupID                pgtype.UUID `json:"member_group_id"`
}

func (q *Queries) GetPluralMembersOnChatRoomRemoveMemberActionUseNumberedPaginate(ctx context.Context, arg GetPluralMembersOnChatRoomRemoveMemberActionUseNumberedPaginateParams) ([]GetPluralMembersOnChatRoomRemoveMemberActionUseNumberedPaginateRow, error) {
	rows, err := q.db.Query(ctx, getPluralMembersOnChatRoomRemoveMemberActionUseNumberedPaginate, arg.Limit, arg.Offset, arg.ChatRoomRemoveMemberActionIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPluralMembersOnChatRoomRemoveMemberActionUseNumberedPaginateRow{}
	for rows.Next() {
		var i GetPluralMembersOnChatRoomRemoveMemberActionUseNumberedPaginateRow
		if err := rows.Scan(
			&i.TChatRoomRemovedMembersPkey,
			&i.ChatRoomRemoveMemberActionID,
			&i.MemberID,
			&i.MemberName,
			&i.MemberFirstName,
			&i.MemberLastName,
			&i.MemberEmail,
			&i.MemberProfileImageID,
			&i.MemberGradeID,
			&i.MemberGroupID,
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
