// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: chat_room_belonging.sql

package query

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const countChatRoomsOnMember = `-- name: CountChatRoomsOnMember :one
SELECT COUNT(*) FROM m_chat_room_belongings WHERE member_id = $1
AND CASE WHEN $2::boolean = true THEN
		EXISTS (SELECT 1 FROM m_chat_rooms WHERE m_chat_room_belongings.chat_room_id = m_chat_rooms.chat_room_id AND m_chat_rooms.name LIKE '%' || $3::text || '%')
	ELSE TRUE END
`

type CountChatRoomsOnMemberParams struct {
	MemberID      uuid.UUID `json:"member_id"`
	WhereLikeName bool      `json:"where_like_name"`
	SearchName    string    `json:"search_name"`
}

func (q *Queries) CountChatRoomsOnMember(ctx context.Context, arg CountChatRoomsOnMemberParams) (int64, error) {
	row := q.db.QueryRow(ctx, countChatRoomsOnMember, arg.MemberID, arg.WhereLikeName, arg.SearchName)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countMembersOnChatRoom = `-- name: CountMembersOnChatRoom :one
SELECT COUNT(*) FROM m_chat_room_belongings WHERE chat_room_id = $1
AND CASE WHEN $2::boolean = true THEN
		EXISTS (SELECT 1 FROM m_members WHERE m_chat_room_belongings.member_id = m_members.member_id AND m_members.name LIKE '%' || $3::text || '%')
	ELSE TRUE END
`

type CountMembersOnChatRoomParams struct {
	ChatRoomID    uuid.UUID `json:"chat_room_id"`
	WhereLikeName bool      `json:"where_like_name"`
	SearchName    string    `json:"search_name"`
}

func (q *Queries) CountMembersOnChatRoom(ctx context.Context, arg CountMembersOnChatRoomParams) (int64, error) {
	row := q.db.QueryRow(ctx, countMembersOnChatRoom, arg.ChatRoomID, arg.WhereLikeName, arg.SearchName)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createChatRoomBelonging = `-- name: CreateChatRoomBelonging :one
INSERT INTO m_chat_room_belongings (member_id, chat_room_id, added_at) VALUES ($1, $2, $3) RETURNING m_chat_room_belongings_pkey, member_id, chat_room_id, added_at
`

type CreateChatRoomBelongingParams struct {
	MemberID   uuid.UUID `json:"member_id"`
	ChatRoomID uuid.UUID `json:"chat_room_id"`
	AddedAt    time.Time `json:"added_at"`
}

func (q *Queries) CreateChatRoomBelonging(ctx context.Context, arg CreateChatRoomBelongingParams) (ChatRoomBelonging, error) {
	row := q.db.QueryRow(ctx, createChatRoomBelonging, arg.MemberID, arg.ChatRoomID, arg.AddedAt)
	var i ChatRoomBelonging
	err := row.Scan(
		&i.MChatRoomBelongingsPkey,
		&i.MemberID,
		&i.ChatRoomID,
		&i.AddedAt,
	)
	return i, err
}

type CreateChatRoomBelongingsParams struct {
	MemberID   uuid.UUID `json:"member_id"`
	ChatRoomID uuid.UUID `json:"chat_room_id"`
	AddedAt    time.Time `json:"added_at"`
}

const deleteChatRoomBelonging = `-- name: DeleteChatRoomBelonging :exec
DELETE FROM m_chat_room_belongings WHERE member_id = $1 AND chat_room_id = $2
`

type DeleteChatRoomBelongingParams struct {
	MemberID   uuid.UUID `json:"member_id"`
	ChatRoomID uuid.UUID `json:"chat_room_id"`
}

func (q *Queries) DeleteChatRoomBelonging(ctx context.Context, arg DeleteChatRoomBelongingParams) error {
	_, err := q.db.Exec(ctx, deleteChatRoomBelonging, arg.MemberID, arg.ChatRoomID)
	return err
}

const deleteChatRoomBelongingsOnMember = `-- name: DeleteChatRoomBelongingsOnMember :exec
DELETE FROM m_chat_room_belongings WHERE member_id = $1
`

func (q *Queries) DeleteChatRoomBelongingsOnMember(ctx context.Context, memberID uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteChatRoomBelongingsOnMember, memberID)
	return err
}

const deleteChatRoomBelongingsOnMembers = `-- name: DeleteChatRoomBelongingsOnMembers :exec
DELETE FROM m_chat_room_belongings WHERE member_id = ANY($1::uuid[])
`

func (q *Queries) DeleteChatRoomBelongingsOnMembers(ctx context.Context, dollar_1 []uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteChatRoomBelongingsOnMembers, dollar_1)
	return err
}

const getChatRoomsOnMember = `-- name: GetChatRoomsOnMember :many
SELECT m_chat_room_belongings.m_chat_room_belongings_pkey, m_chat_room_belongings.member_id, m_chat_room_belongings.chat_room_id, m_chat_room_belongings.added_at, m_chat_rooms.m_chat_rooms_pkey, m_chat_rooms.chat_room_id, m_chat_rooms.name, m_chat_rooms.is_private, m_chat_rooms.cover_image_url, m_chat_rooms.owner_id, m_chat_rooms.from_organization, m_chat_rooms.created_at, m_chat_rooms.updated_at FROM m_chat_room_belongings
LEFT JOIN m_chat_rooms ON m_chat_room_belongings.chat_room_id = m_chat_rooms.chat_room_id
WHERE member_id = $1
AND CASE
	WHEN $2::boolean = true THEN m_members.name LIKE '%' || $3::text || '%'
END
ORDER BY
	CASE WHEN $4::text = 'name' THEN m_chat_rooms.name END ASC,
	CASE WHEN $4::text = 'r_name' THEN m_chat_rooms.name END DESC,
	CASE WHEN $4::text = 'old_add' THEN m_chat_room_belongings.added_at END ASC,
	CASE WHEN $4::text = 'late_add' THEN m_chat_room_belongings.added_at END DESC,
	CASE WHEN $4::text = 'old_chat' THEN
		(SELECT MAX(created_at) FROM t_messages m WHERE m.chat_room_id = m_chat_room_belongings.chat_room_id)
	END ASC,
	CASE WHEN $4::text = 'late_chat' THEN
		(SELECT MAX(created_at) FROM t_messages m WHERE m.chat_room_id = m_chat_room_belongings.chat_room_id)
	END DESC,
	m_chat_room_belongings_pkey ASC
`

type GetChatRoomsOnMemberParams struct {
	MemberID      uuid.UUID `json:"member_id"`
	WhereLikeName bool      `json:"where_like_name"`
	SearchName    string    `json:"search_name"`
	OrderMethod   string    `json:"order_method"`
}

type GetChatRoomsOnMemberRow struct {
	MChatRoomBelongingsPkey pgtype.Int8 `json:"m_chat_room_belongings_pkey"`
	MemberID                uuid.UUID   `json:"member_id"`
	ChatRoomID              uuid.UUID   `json:"chat_room_id"`
	AddedAt                 time.Time   `json:"added_at"`
	ChatRoom                ChatRoom    `json:"chat_room"`
}

func (q *Queries) GetChatRoomsOnMember(ctx context.Context, arg GetChatRoomsOnMemberParams) ([]GetChatRoomsOnMemberRow, error) {
	rows, err := q.db.Query(ctx, getChatRoomsOnMember,
		arg.MemberID,
		arg.WhereLikeName,
		arg.SearchName,
		arg.OrderMethod,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetChatRoomsOnMemberRow{}
	for rows.Next() {
		var i GetChatRoomsOnMemberRow
		if err := rows.Scan(
			&i.MChatRoomBelongingsPkey,
			&i.MemberID,
			&i.ChatRoomID,
			&i.AddedAt,
			&i.ChatRoom.MChatRoomsPkey,
			&i.ChatRoom.ChatRoomID,
			&i.ChatRoom.Name,
			&i.ChatRoom.IsPrivate,
			&i.ChatRoom.CoverImageUrl,
			&i.ChatRoom.OwnerID,
			&i.ChatRoom.FromOrganization,
			&i.ChatRoom.CreatedAt,
			&i.ChatRoom.UpdatedAt,
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

const getChatRoomsOnMemberUseKeysetPaginate = `-- name: GetChatRoomsOnMemberUseKeysetPaginate :many
SELECT m_chat_room_belongings.m_chat_room_belongings_pkey, m_chat_room_belongings.member_id, m_chat_room_belongings.chat_room_id, m_chat_room_belongings.added_at, m_chat_rooms.m_chat_rooms_pkey, m_chat_rooms.chat_room_id, m_chat_rooms.name, m_chat_rooms.is_private, m_chat_rooms.cover_image_url, m_chat_rooms.owner_id, m_chat_rooms.from_organization, m_chat_rooms.created_at, m_chat_rooms.updated_at FROM m_chat_room_belongings
LEFT JOIN m_chat_rooms ON m_chat_room_belongings.chat_room_id = m_chat_rooms.chat_room_id
WHERE member_id = $1
AND CASE $3::text
	WHEN 'next' THEN
		CASE $4::text
			WHEN 'name' THEN m_chat_rooms.name > $5 OR (m_chat_rooms.name = $5 AND m_chat_room_belongings_pkey > $6::int)
			WHEN 'r_name' THEN m_chat_rooms.name < $5 OR (m_chat_rooms.name = $5 AND m_chat_room_belongings_pkey > $6::int)
			WHEN 'old_add' THEN m_chat_room_belongings.added_at > $7 OR (m_chat_room_belongings.added_at = $7 AND m_chat_room_belongings_pkey > $6::int)
			WHEN 'late_add' THEN m_chat_room_belongings.added_at < $7 OR (m_chat_room_belongings.added_at = $7 AND m_chat_room_belongings_pkey > $6::int)
			WHEN 'old_chat' THEN
				(SELECT MAX(created_at) FROM t_messages m WHERE m.chat_room_id = m_chat_room_belongings.chat_room_id) > $8
				OR ((SELECT MAX(created_at) FROM t_messages m WHERE m.chat_room_id = m_chat_room_belongings.chat_room_id) = $8 AND m_chat_room_belongings_pkey > $6::int)
			WHEN 'late_chat' THEN
				(SELECT MAX(created_at) FROM t_messages m WHERE m.chat_room_id = m_chat_room_belongings.chat_room_id) < $8
				OR ((SELECT MAX(created_at) FROM t_messages m WHERE m.chat_room_id = m_chat_room_belongings.chat_room_id) = $8 AND m_chat_room_belongings_pkey > $6::int)
			ELSE m_chat_room_belongings_pkey > $6::int
		END
	WHEN 'prev' THEN
		CASE $4::text
			WHEN 'name' THEN m_chat_rooms.name < $5 OR (m_chat_rooms.name = $5 AND m_chat_room_belongings_pkey < $6::int)
			WHEN 'r_name' THEN m_chat_rooms.name > $5 OR (m_chat_rooms.name = $5 AND m_chat_room_belongings_pkey < $6::int)
			WHEN 'old_add' THEN m_chat_room_belongings.added_at < $7 OR (m_chat_room_belongings.added_at = $7 AND m_chat_room_belongings_pkey < $6::int)
			WHEN 'late_add' THEN m_chat_room_belongings.added_at > $7 OR (m_chat_room_belongings.added_at = $7 AND m_chat_room_belongings_pkey < $6::int)
			WHEN 'old_chat' THEN
				(SELECT MAX(created_at) FROM t_messages m WHERE m.chat_room_id = m_chat_room_belongings.chat_room_id) < $8
				OR ((SELECT MAX(created_at) FROM t_messages m WHERE m.chat_room_id = m_chat_room_belongings.chat_room_id) = $8 AND m_chat_room_belongings_pkey < $6::int)
			WHEN 'late_chat' THEN
				(SELECT MAX(created_at) FROM t_messages m WHERE m.chat_room_id = m_chat_room_belongings.chat_room_id) > $8
				OR ((SELECT MAX(created_at) FROM t_messages m WHERE m.chat_room_id = m_chat_room_belongings.chat_room_id) = $8 AND m_chat_room_belongings_pkey < $6::int)
			ELSE m_chat_room_belongings_pkey < $6::int
		END
END
ORDER BY
	CASE WHEN $4::text = 'name' AND $3::text = 'next' THEN m_chat_rooms.name END ASC,
	CASE WHEN $4::text = 'name' AND $3::text = 'prev' THEN m_chat_rooms.name END DESC,
	CASE WHEN $4::text = 'r_name' AND $3::text = 'next' THEN m_chat_rooms.name END ASC,
	CASE WHEN $4::text = 'r_name' AND $3::text = 'prev' THEN m_chat_rooms.name END DESC,
	CASE WHEN $4::text = 'old_add' AND $3::text = 'next' THEN m_chat_room_belongings.added_at END ASC,
	CASE WHEN $4::text = 'old_add' AND $3::text = 'prev' THEN m_chat_room_belongings.added_at END DESC,
	CASE WHEN $4::text = 'late_add' AND $3::text = 'next' THEN m_chat_room_belongings.added_at END ASC,
	CASE WHEN $4::text = 'late_add' AND $3::text = 'prev' THEN m_chat_room_belongings.added_at END DESC,
	CASE WHEN $4::text = 'old_chat' AND $3::text = 'next' THEN
		(SELECT MAX(created_at) FROM t_messages m WHERE m.chat_room_id = m_chat_room_belongings.chat_room_id) END ASC,
	CASE WHEN $4::text = 'old_chat' AND $3::text = 'prev' THEN
		(SELECT MAX(created_at) FROM t_messages m WHERE m.chat_room_id = m_chat_room_belongings.chat_room_id) END DESC,
	CASE WHEN $4::text = 'late_chat' AND $3::text = 'next' THEN
		(SELECT MAX(created_at) FROM t_messages m WHERE m.chat_room_id = m_chat_room_belongings.chat_room_id) END ASC,
	CASE WHEN $4::text = 'late_chat' AND $3::text = 'prev' THEN
		(SELECT MAX(created_at) FROM t_messages m WHERE m.chat_room_id = m_chat_room_belongings.chat_room_id) END DESC,
	CASE WHEN $3::text = 'next' THEN m_chat_room_belongings_pkey END ASC,
	CASE WHEN $3::text = 'prev' THEN m_chat_room_belongings_pkey END DESC
LIMIT $2
`

type GetChatRoomsOnMemberUseKeysetPaginateParams struct {
	MemberID        uuid.UUID   `json:"member_id"`
	Limit           int32       `json:"limit"`
	CursorDirection string      `json:"cursor_direction"`
	OrderMethod     string      `json:"order_method"`
	NameCursor      pgtype.Text `json:"name_cursor"`
	Cursor          int32       `json:"cursor"`
	AddCursor       time.Time   `json:"add_cursor"`
	ChatCursor      time.Time   `json:"chat_cursor"`
}

type GetChatRoomsOnMemberUseKeysetPaginateRow struct {
	MChatRoomBelongingsPkey pgtype.Int8 `json:"m_chat_room_belongings_pkey"`
	MemberID                uuid.UUID   `json:"member_id"`
	ChatRoomID              uuid.UUID   `json:"chat_room_id"`
	AddedAt                 time.Time   `json:"added_at"`
	ChatRoom                ChatRoom    `json:"chat_room"`
}

func (q *Queries) GetChatRoomsOnMemberUseKeysetPaginate(ctx context.Context, arg GetChatRoomsOnMemberUseKeysetPaginateParams) ([]GetChatRoomsOnMemberUseKeysetPaginateRow, error) {
	rows, err := q.db.Query(ctx, getChatRoomsOnMemberUseKeysetPaginate,
		arg.MemberID,
		arg.Limit,
		arg.CursorDirection,
		arg.OrderMethod,
		arg.NameCursor,
		arg.Cursor,
		arg.AddCursor,
		arg.ChatCursor,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetChatRoomsOnMemberUseKeysetPaginateRow{}
	for rows.Next() {
		var i GetChatRoomsOnMemberUseKeysetPaginateRow
		if err := rows.Scan(
			&i.MChatRoomBelongingsPkey,
			&i.MemberID,
			&i.ChatRoomID,
			&i.AddedAt,
			&i.ChatRoom.MChatRoomsPkey,
			&i.ChatRoom.ChatRoomID,
			&i.ChatRoom.Name,
			&i.ChatRoom.IsPrivate,
			&i.ChatRoom.CoverImageUrl,
			&i.ChatRoom.OwnerID,
			&i.ChatRoom.FromOrganization,
			&i.ChatRoom.CreatedAt,
			&i.ChatRoom.UpdatedAt,
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

const getChatRoomsOnMemberUseNumberedPaginate = `-- name: GetChatRoomsOnMemberUseNumberedPaginate :many
SELECT m_chat_room_belongings.m_chat_room_belongings_pkey, m_chat_room_belongings.member_id, m_chat_room_belongings.chat_room_id, m_chat_room_belongings.added_at, m_chat_rooms.m_chat_rooms_pkey, m_chat_rooms.chat_room_id, m_chat_rooms.name, m_chat_rooms.is_private, m_chat_rooms.cover_image_url, m_chat_rooms.owner_id, m_chat_rooms.from_organization, m_chat_rooms.created_at, m_chat_rooms.updated_at FROM m_chat_room_belongings
LEFT JOIN m_chat_rooms ON m_chat_room_belongings.chat_room_id = m_chat_rooms.chat_room_id
WHERE member_id = $1
AND CASE
	WHEN $4::boolean = true THEN m_members.name LIKE '%' || $5::text || '%'
END
ORDER BY
	CASE WHEN $6::text = 'name' THEN m_chat_rooms.name END ASC,
	CASE WHEN $6::text = 'r_name' THEN m_chat_rooms.name END DESC,
	CASE WHEN $6::text = 'old_add' THEN m_chat_room_belongings.added_at END ASC,
	CASE WHEN $6::text = 'late_add' THEN m_chat_room_belongings.added_at END DESC,
	CASE WHEN $6::text = 'old_chat' THEN
		(SELECT MAX(created_at) FROM t_messages m WHERE m.chat_room_id = m_chat_room_belongings.chat_room_id)
	END ASC,
	CASE WHEN $6::text = 'late_chat' THEN
		(SELECT MAX(created_at) FROM t_messages m WHERE m.chat_room_id = m_chat_room_belongings.chat_room_id)
	END DESC,
	m_chat_room_belongings_pkey ASC
LIMIT $2 OFFSET $3
`

type GetChatRoomsOnMemberUseNumberedPaginateParams struct {
	MemberID      uuid.UUID `json:"member_id"`
	Limit         int32     `json:"limit"`
	Offset        int32     `json:"offset"`
	WhereLikeName bool      `json:"where_like_name"`
	SearchName    string    `json:"search_name"`
	OrderMethod   string    `json:"order_method"`
}

type GetChatRoomsOnMemberUseNumberedPaginateRow struct {
	MChatRoomBelongingsPkey pgtype.Int8 `json:"m_chat_room_belongings_pkey"`
	MemberID                uuid.UUID   `json:"member_id"`
	ChatRoomID              uuid.UUID   `json:"chat_room_id"`
	AddedAt                 time.Time   `json:"added_at"`
	ChatRoom                ChatRoom    `json:"chat_room"`
}

func (q *Queries) GetChatRoomsOnMemberUseNumberedPaginate(ctx context.Context, arg GetChatRoomsOnMemberUseNumberedPaginateParams) ([]GetChatRoomsOnMemberUseNumberedPaginateRow, error) {
	rows, err := q.db.Query(ctx, getChatRoomsOnMemberUseNumberedPaginate,
		arg.MemberID,
		arg.Limit,
		arg.Offset,
		arg.WhereLikeName,
		arg.SearchName,
		arg.OrderMethod,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetChatRoomsOnMemberUseNumberedPaginateRow{}
	for rows.Next() {
		var i GetChatRoomsOnMemberUseNumberedPaginateRow
		if err := rows.Scan(
			&i.MChatRoomBelongingsPkey,
			&i.MemberID,
			&i.ChatRoomID,
			&i.AddedAt,
			&i.ChatRoom.MChatRoomsPkey,
			&i.ChatRoom.ChatRoomID,
			&i.ChatRoom.Name,
			&i.ChatRoom.IsPrivate,
			&i.ChatRoom.CoverImageUrl,
			&i.ChatRoom.OwnerID,
			&i.ChatRoom.FromOrganization,
			&i.ChatRoom.CreatedAt,
			&i.ChatRoom.UpdatedAt,
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

const getMembersOnChatRoom = `-- name: GetMembersOnChatRoom :many
SELECT m_chat_room_belongings.m_chat_room_belongings_pkey, m_chat_room_belongings.member_id, m_chat_room_belongings.chat_room_id, m_chat_room_belongings.added_at, m_members.m_members_pkey, m_members.member_id, m_members.login_id, m_members.password, m_members.email, m_members.name, m_members.attend_status_id, m_members.profile_image_url, m_members.grade_id, m_members.group_id, m_members.personal_organization_id, m_members.role_id, m_members.created_at, m_members.updated_at FROM m_chat_room_belongings
LEFT JOIN m_members ON m_chat_room_belongings.member_id = m_members.member_id
WHERE chat_room_id = $1
AND CASE
	WHEN $2::boolean = true THEN m_members.name LIKE '%' || $3::text || '%'
END
ORDER BY
	CASE WHEN $4::text = 'name' THEN m_members.name END ASC,
	CASE WHEN $4::text = 'r_name' THEN m_members.name END DESC,
	CASE WHEN $4::text = 'old_add' THEN m_chat_room_belongings.added_at END ASC,
	CASE WHEN $4::text = 'late_add' THEN m_chat_room_belongings.added_at END DESC,
	m_chat_room_belongings_pkey ASC
`

type GetMembersOnChatRoomParams struct {
	ChatRoomID    uuid.UUID `json:"chat_room_id"`
	WhereLikeName bool      `json:"where_like_name"`
	SearchName    string    `json:"search_name"`
	OrderMethod   string    `json:"order_method"`
}

type GetMembersOnChatRoomRow struct {
	MChatRoomBelongingsPkey pgtype.Int8        `json:"m_chat_room_belongings_pkey"`
	MemberID                uuid.UUID          `json:"member_id"`
	ChatRoomID              uuid.UUID          `json:"chat_room_id"`
	AddedAt                 time.Time          `json:"added_at"`
	MMembersPkey            pgtype.Int8        `json:"m_members_pkey"`
	MemberID_2              pgtype.UUID        `json:"member_id_2"`
	LoginID                 pgtype.Text        `json:"login_id"`
	Password                pgtype.Text        `json:"password"`
	Email                   pgtype.Text        `json:"email"`
	Name                    pgtype.Text        `json:"name"`
	AttendStatusID          pgtype.UUID        `json:"attend_status_id"`
	ProfileImageUrl         pgtype.Text        `json:"profile_image_url"`
	GradeID                 pgtype.UUID        `json:"grade_id"`
	GroupID                 pgtype.UUID        `json:"group_id"`
	PersonalOrganizationID  pgtype.UUID        `json:"personal_organization_id"`
	RoleID                  pgtype.UUID        `json:"role_id"`
	CreatedAt               pgtype.Timestamptz `json:"created_at"`
	UpdatedAt               pgtype.Timestamptz `json:"updated_at"`
}

func (q *Queries) GetMembersOnChatRoom(ctx context.Context, arg GetMembersOnChatRoomParams) ([]GetMembersOnChatRoomRow, error) {
	rows, err := q.db.Query(ctx, getMembersOnChatRoom,
		arg.ChatRoomID,
		arg.WhereLikeName,
		arg.SearchName,
		arg.OrderMethod,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetMembersOnChatRoomRow{}
	for rows.Next() {
		var i GetMembersOnChatRoomRow
		if err := rows.Scan(
			&i.MChatRoomBelongingsPkey,
			&i.MemberID,
			&i.ChatRoomID,
			&i.AddedAt,
			&i.MMembersPkey,
			&i.MemberID_2,
			&i.LoginID,
			&i.Password,
			&i.Email,
			&i.Name,
			&i.AttendStatusID,
			&i.ProfileImageUrl,
			&i.GradeID,
			&i.GroupID,
			&i.PersonalOrganizationID,
			&i.RoleID,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const getMembersOnChatRoomUseKeysetPaginate = `-- name: GetMembersOnChatRoomUseKeysetPaginate :many
SELECT m_chat_room_belongings.m_chat_room_belongings_pkey, m_chat_room_belongings.member_id, m_chat_room_belongings.chat_room_id, m_chat_room_belongings.added_at, m_members.m_members_pkey, m_members.member_id, m_members.login_id, m_members.password, m_members.email, m_members.name, m_members.attend_status_id, m_members.profile_image_url, m_members.grade_id, m_members.group_id, m_members.personal_organization_id, m_members.role_id, m_members.created_at, m_members.updated_at FROM m_chat_room_belongings
LEFT JOIN m_members ON m_chat_room_belongings.member_id = m_members.member_id
WHERE chat_room_id = $1
AND CASE $3::text
	WHEN 'next' THEN
		CASE $4::text
			WHEN 'name' THEN m_members.name > $5 OR (m_members.name = $5 AND m_chat_room_belongings_pkey > $6::int)
			WHEN 'r_name' THEN m_members.name < $5 OR (m_members.name = $5 AND m_chat_room_belongings_pkey > $6::int)
			WHEN 'old_add' THEN m_chat_room_belongings.added_at > $7 OR (m_chat_room_belongings.added_at = $7 AND m_chat_room_belongings_pkey > $6::int)
			WHEN 'late_add' THEN m_chat_room_belongings.added_at < $7 OR (m_chat_room_belongings.added_at = $7 AND m_chat_room_belongings_pkey > $6::int)
			ELSE m_chat_room_belongings_pkey > $6::int
		END
	WHEN 'prev' THEN
		CASE $4::text
			WHEN 'name' THEN m_members.name < $5 OR (m_members.name = $5 AND m_chat_room_belongings_pkey < $6::int)
			WHEN 'r_name' THEN m_members.name > $5 OR (m_members.name = $5 AND m_chat_room_belongings_pkey < $6::int)
			WHEN 'old_add' THEN m_chat_room_belongings.added_at < $7 OR (m_chat_room_belongings.added_at = $7 AND m_chat_room_belongings_pkey < $6::int)
			WHEN 'late_add' THEN m_chat_room_belongings.added_at > $7 OR (m_chat_room_belongings.added_at = $7 AND m_chat_room_belongings_pkey < $6::int)
			ELSE m_chat_room_belongings_pkey < $6::int
		END
END
ORDER BY
	CASE WHEN $4::text = 'name' AND $3::text = 'next' THEN m_members.name END ASC,
	CASE WHEN $4::text = 'name' AND $3::text = 'prev' THEN m_members.name END DESC,
	CASE WHEN $4::text = 'r_name' AND $3::text = 'next' THEN m_members.name END ASC,
	CASE WHEN $4::text = 'r_name' AND $3::text = 'prev' THEN m_members.name END DESC,
	CASE WHEN $4::text = 'old_add' AND $3::text = 'next' THEN m_chat_room_belongings.added_at END ASC,
	CASE WHEN $4::text = 'old_add' AND $3::text = 'prev' THEN m_chat_room_belongings.added_at END DESC,
	CASE WHEN $4::text = 'late_add' AND $3::text = 'next' THEN m_chat_room_belongings.added_at END ASC,
	CASE WHEN $4::text = 'late_add' AND $3::text = 'prev' THEN m_chat_room_belongings.added_at END DESC,
	CASE WHEN $3::text = 'next' THEN m_chat_room_belongings_pkey END ASC,
	CASE WHEN $3::text = 'prev' THEN m_chat_room_belongings_pkey END DESC
LIMIT $2
`

type GetMembersOnChatRoomUseKeysetPaginateParams struct {
	ChatRoomID      uuid.UUID `json:"chat_room_id"`
	Limit           int32     `json:"limit"`
	CursorDirection string    `json:"cursor_direction"`
	OrderMethod     string    `json:"order_method"`
	NameCursor      string    `json:"name_cursor"`
	Cursor          int32     `json:"cursor"`
	AddedAtCursor   time.Time `json:"added_at_cursor"`
}

type GetMembersOnChatRoomUseKeysetPaginateRow struct {
	MChatRoomBelongingsPkey pgtype.Int8        `json:"m_chat_room_belongings_pkey"`
	MemberID                uuid.UUID          `json:"member_id"`
	ChatRoomID              uuid.UUID          `json:"chat_room_id"`
	AddedAt                 time.Time          `json:"added_at"`
	MMembersPkey            pgtype.Int8        `json:"m_members_pkey"`
	MemberID_2              pgtype.UUID        `json:"member_id_2"`
	LoginID                 pgtype.Text        `json:"login_id"`
	Password                pgtype.Text        `json:"password"`
	Email                   pgtype.Text        `json:"email"`
	Name                    pgtype.Text        `json:"name"`
	AttendStatusID          pgtype.UUID        `json:"attend_status_id"`
	ProfileImageUrl         pgtype.Text        `json:"profile_image_url"`
	GradeID                 pgtype.UUID        `json:"grade_id"`
	GroupID                 pgtype.UUID        `json:"group_id"`
	PersonalOrganizationID  pgtype.UUID        `json:"personal_organization_id"`
	RoleID                  pgtype.UUID        `json:"role_id"`
	CreatedAt               pgtype.Timestamptz `json:"created_at"`
	UpdatedAt               pgtype.Timestamptz `json:"updated_at"`
}

func (q *Queries) GetMembersOnChatRoomUseKeysetPaginate(ctx context.Context, arg GetMembersOnChatRoomUseKeysetPaginateParams) ([]GetMembersOnChatRoomUseKeysetPaginateRow, error) {
	rows, err := q.db.Query(ctx, getMembersOnChatRoomUseKeysetPaginate,
		arg.ChatRoomID,
		arg.Limit,
		arg.CursorDirection,
		arg.OrderMethod,
		arg.NameCursor,
		arg.Cursor,
		arg.AddedAtCursor,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetMembersOnChatRoomUseKeysetPaginateRow{}
	for rows.Next() {
		var i GetMembersOnChatRoomUseKeysetPaginateRow
		if err := rows.Scan(
			&i.MChatRoomBelongingsPkey,
			&i.MemberID,
			&i.ChatRoomID,
			&i.AddedAt,
			&i.MMembersPkey,
			&i.MemberID_2,
			&i.LoginID,
			&i.Password,
			&i.Email,
			&i.Name,
			&i.AttendStatusID,
			&i.ProfileImageUrl,
			&i.GradeID,
			&i.GroupID,
			&i.PersonalOrganizationID,
			&i.RoleID,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const getMembersOnChatRoomUseNumberedPaginate = `-- name: GetMembersOnChatRoomUseNumberedPaginate :many
SELECT m_chat_room_belongings.m_chat_room_belongings_pkey, m_chat_room_belongings.member_id, m_chat_room_belongings.chat_room_id, m_chat_room_belongings.added_at, m_members.m_members_pkey, m_members.member_id, m_members.login_id, m_members.password, m_members.email, m_members.name, m_members.attend_status_id, m_members.profile_image_url, m_members.grade_id, m_members.group_id, m_members.personal_organization_id, m_members.role_id, m_members.created_at, m_members.updated_at FROM m_chat_room_belongings
LEFT JOIN m_members ON m_chat_room_belongings.member_id = m_members.member_id
WHERE chat_room_id = $1
AND CASE
	WHEN $4::boolean = true THEN m_members.name LIKE '%' || $5::text || '%'
END
ORDER BY
	CASE WHEN $6::text = 'name' THEN m_members.name END ASC,
	CASE WHEN $6::text = 'r_name' THEN m_members.name END DESC,
	CASE WHEN $6::text = 'old_add' THEN m_chat_room_belongings.added_at END ASC,
	CASE WHEN $6::text = 'late_add' THEN m_chat_room_belongings.added_at END DESC,
	m_chat_room_belongings_pkey ASC
LIMIT $2 OFFSET $3
`

type GetMembersOnChatRoomUseNumberedPaginateParams struct {
	ChatRoomID    uuid.UUID `json:"chat_room_id"`
	Limit         int32     `json:"limit"`
	Offset        int32     `json:"offset"`
	WhereLikeName bool      `json:"where_like_name"`
	SearchName    string    `json:"search_name"`
	OrderMethod   string    `json:"order_method"`
}

type GetMembersOnChatRoomUseNumberedPaginateRow struct {
	MChatRoomBelongingsPkey pgtype.Int8        `json:"m_chat_room_belongings_pkey"`
	MemberID                uuid.UUID          `json:"member_id"`
	ChatRoomID              uuid.UUID          `json:"chat_room_id"`
	AddedAt                 time.Time          `json:"added_at"`
	MMembersPkey            pgtype.Int8        `json:"m_members_pkey"`
	MemberID_2              pgtype.UUID        `json:"member_id_2"`
	LoginID                 pgtype.Text        `json:"login_id"`
	Password                pgtype.Text        `json:"password"`
	Email                   pgtype.Text        `json:"email"`
	Name                    pgtype.Text        `json:"name"`
	AttendStatusID          pgtype.UUID        `json:"attend_status_id"`
	ProfileImageUrl         pgtype.Text        `json:"profile_image_url"`
	GradeID                 pgtype.UUID        `json:"grade_id"`
	GroupID                 pgtype.UUID        `json:"group_id"`
	PersonalOrganizationID  pgtype.UUID        `json:"personal_organization_id"`
	RoleID                  pgtype.UUID        `json:"role_id"`
	CreatedAt               pgtype.Timestamptz `json:"created_at"`
	UpdatedAt               pgtype.Timestamptz `json:"updated_at"`
}

func (q *Queries) GetMembersOnChatRoomUseNumberedPaginate(ctx context.Context, arg GetMembersOnChatRoomUseNumberedPaginateParams) ([]GetMembersOnChatRoomUseNumberedPaginateRow, error) {
	rows, err := q.db.Query(ctx, getMembersOnChatRoomUseNumberedPaginate,
		arg.ChatRoomID,
		arg.Limit,
		arg.Offset,
		arg.WhereLikeName,
		arg.SearchName,
		arg.OrderMethod,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetMembersOnChatRoomUseNumberedPaginateRow{}
	for rows.Next() {
		var i GetMembersOnChatRoomUseNumberedPaginateRow
		if err := rows.Scan(
			&i.MChatRoomBelongingsPkey,
			&i.MemberID,
			&i.ChatRoomID,
			&i.AddedAt,
			&i.MMembersPkey,
			&i.MemberID_2,
			&i.LoginID,
			&i.Password,
			&i.Email,
			&i.Name,
			&i.AttendStatusID,
			&i.ProfileImageUrl,
			&i.GradeID,
			&i.GroupID,
			&i.PersonalOrganizationID,
			&i.RoleID,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const getPluralChatRoomsOnMember = `-- name: GetPluralChatRoomsOnMember :many
SELECT m_chat_room_belongings.m_chat_room_belongings_pkey, m_chat_room_belongings.member_id, m_chat_room_belongings.chat_room_id, m_chat_room_belongings.added_at, m_chat_rooms.m_chat_rooms_pkey, m_chat_rooms.chat_room_id, m_chat_rooms.name, m_chat_rooms.is_private, m_chat_rooms.cover_image_url, m_chat_rooms.owner_id, m_chat_rooms.from_organization, m_chat_rooms.created_at, m_chat_rooms.updated_at FROM m_chat_room_belongings
LEFT JOIN m_chat_rooms ON m_chat_room_belongings.chat_room_id = m_chat_rooms.chat_room_id
WHERE member_id = ANY($3::uuid[])
ORDER BY
	m_chat_room_belongings_pkey ASC
LIMIT $1 OFFSET $2
`

type GetPluralChatRoomsOnMemberParams struct {
	Limit     int32       `json:"limit"`
	Offset    int32       `json:"offset"`
	MemberIds []uuid.UUID `json:"member_ids"`
}

type GetPluralChatRoomsOnMemberRow struct {
	MChatRoomBelongingsPkey pgtype.Int8 `json:"m_chat_room_belongings_pkey"`
	MemberID                uuid.UUID   `json:"member_id"`
	ChatRoomID              uuid.UUID   `json:"chat_room_id"`
	AddedAt                 time.Time   `json:"added_at"`
	ChatRoom                ChatRoom    `json:"chat_room"`
}

func (q *Queries) GetPluralChatRoomsOnMember(ctx context.Context, arg GetPluralChatRoomsOnMemberParams) ([]GetPluralChatRoomsOnMemberRow, error) {
	rows, err := q.db.Query(ctx, getPluralChatRoomsOnMember, arg.Limit, arg.Offset, arg.MemberIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPluralChatRoomsOnMemberRow{}
	for rows.Next() {
		var i GetPluralChatRoomsOnMemberRow
		if err := rows.Scan(
			&i.MChatRoomBelongingsPkey,
			&i.MemberID,
			&i.ChatRoomID,
			&i.AddedAt,
			&i.ChatRoom.MChatRoomsPkey,
			&i.ChatRoom.ChatRoomID,
			&i.ChatRoom.Name,
			&i.ChatRoom.IsPrivate,
			&i.ChatRoom.CoverImageUrl,
			&i.ChatRoom.OwnerID,
			&i.ChatRoom.FromOrganization,
			&i.ChatRoom.CreatedAt,
			&i.ChatRoom.UpdatedAt,
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

const getPluralMembersOnChatRoom = `-- name: GetPluralMembersOnChatRoom :many
SELECT m_chat_room_belongings.m_chat_room_belongings_pkey, m_chat_room_belongings.member_id, m_chat_room_belongings.chat_room_id, m_chat_room_belongings.added_at, m_members.m_members_pkey, m_members.member_id, m_members.login_id, m_members.password, m_members.email, m_members.name, m_members.attend_status_id, m_members.profile_image_url, m_members.grade_id, m_members.group_id, m_members.personal_organization_id, m_members.role_id, m_members.created_at, m_members.updated_at FROM m_chat_room_belongings
LEFT JOIN m_members ON m_chat_room_belongings.member_id = m_members.member_id
WHERE chat_room_id = ANY($3::uuid[])
ORDER BY
	m_chat_room_belongings_pkey ASC
LIMIT $1 OFFSET $2
`

type GetPluralMembersOnChatRoomParams struct {
	Limit       int32       `json:"limit"`
	Offset      int32       `json:"offset"`
	ChatRoomIds []uuid.UUID `json:"chat_room_ids"`
}

type GetPluralMembersOnChatRoomRow struct {
	MChatRoomBelongingsPkey pgtype.Int8        `json:"m_chat_room_belongings_pkey"`
	MemberID                uuid.UUID          `json:"member_id"`
	ChatRoomID              uuid.UUID          `json:"chat_room_id"`
	AddedAt                 time.Time          `json:"added_at"`
	MMembersPkey            pgtype.Int8        `json:"m_members_pkey"`
	MemberID_2              pgtype.UUID        `json:"member_id_2"`
	LoginID                 pgtype.Text        `json:"login_id"`
	Password                pgtype.Text        `json:"password"`
	Email                   pgtype.Text        `json:"email"`
	Name                    pgtype.Text        `json:"name"`
	AttendStatusID          pgtype.UUID        `json:"attend_status_id"`
	ProfileImageUrl         pgtype.Text        `json:"profile_image_url"`
	GradeID                 pgtype.UUID        `json:"grade_id"`
	GroupID                 pgtype.UUID        `json:"group_id"`
	PersonalOrganizationID  pgtype.UUID        `json:"personal_organization_id"`
	RoleID                  pgtype.UUID        `json:"role_id"`
	CreatedAt               pgtype.Timestamptz `json:"created_at"`
	UpdatedAt               pgtype.Timestamptz `json:"updated_at"`
}

func (q *Queries) GetPluralMembersOnChatRoom(ctx context.Context, arg GetPluralMembersOnChatRoomParams) ([]GetPluralMembersOnChatRoomRow, error) {
	rows, err := q.db.Query(ctx, getPluralMembersOnChatRoom, arg.Limit, arg.Offset, arg.ChatRoomIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPluralMembersOnChatRoomRow{}
	for rows.Next() {
		var i GetPluralMembersOnChatRoomRow
		if err := rows.Scan(
			&i.MChatRoomBelongingsPkey,
			&i.MemberID,
			&i.ChatRoomID,
			&i.AddedAt,
			&i.MMembersPkey,
			&i.MemberID_2,
			&i.LoginID,
			&i.Password,
			&i.Email,
			&i.Name,
			&i.AttendStatusID,
			&i.ProfileImageUrl,
			&i.GradeID,
			&i.GroupID,
			&i.PersonalOrganizationID,
			&i.RoleID,
			&i.CreatedAt,
			&i.UpdatedAt,
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
